package proof

import (
	"context"
	"encoding/json"
	"github.com/ODIN-PROTOCOL/odin-core/pkg/obi"
	"net/http"
	"strconv"

	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/ethereum/go-ethereum/accounts/abi"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

type JsonMultiProof struct {
	BlockHeight          uint64            `json:"block_height"`
	OracleDataMultiProof []OracleDataProof `json:"oracle_data_multi_proof"`
	BlockRelayProof      BlockRelayProof   `json:"block_relay_proof"`
}

type MultiProof struct {
	JsonProof     JsonMultiProof   `json:"json_proof"`
	EVMProofBytes tmbytes.HexBytes `json:"evm_proof_bytes"`
}

func GetMutiProofHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestIDs := r.URL.Query()["id"]
		ctx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		height := &ctx.Height
		if ctx.Height == 0 {
			height = nil
		}

		commit, err := ctx.Client.Commit(context.Background(), height)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		signatures, err := GetSignaturesAndPrefix(&commit.SignedHeader)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		blockRelay := BlockRelayProof{
			BlockHeaderMerkleParts: GetBlockHeaderMerkleParts(commit.Header),
			Signatures:             signatures,
		}

		oracleDataBytesList := make([][]byte, len(requestIDs))
		oracleDataList := make([]OracleDataProof, len(requestIDs))

		for idx, requestID := range requestIDs {
			intRequestID, err := strconv.ParseUint(requestID, 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			requestID := oracletypes.RequestID(intRequestID)

			// Extract multiStoreEp in the first iteration only, since multiStoreEp is the same for all requests.
			value, iavlEp, multiStoreEp, err := getProofsByKey(
				ctx,
				oracletypes.ResultStoreKey(requestID),
				rpcclient.ABCIQueryOptions{Height: commit.Height - 1, Prove: true},
				idx == 0,
			)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			var rs oracletypes.Result
			obi.MustDecode(value, &rs)

			oracleData := OracleDataProof{
				Result:      rs,
				Version:     decodeIAVLLeafPrefix(iavlEp.Leaf.Prefix),
				MerklePaths: GetMerklePaths(iavlEp),
			}
			oracleDataBytes, err := oracleData.encodeToEthData(uint64(commit.Height))
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			// Append oracle data proof to list
			oracleDataBytesList[idx] = oracleDataBytes
			oracleDataList[idx] = oracleData

			if idx == 0 {
				blockRelay.MultiStoreProof = GetMultiStoreProof(multiStoreEp)
			}
		}

		blockRelayBytes, err := blockRelay.encodeToEthData()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Calculate byte for MultiProofbytes
		var relayAndVerifyArguments abi.Arguments
		format := `[{"type":"bytes"},{"type":"bytes[]"}]`
		err = json.Unmarshal([]byte(format), &relayAndVerifyArguments)
		if err != nil {
			panic(err)
		}

		evmProofBytes, err := relayAndVerifyArguments.Pack(blockRelayBytes, oracleDataBytesList)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, ctx, MultiProof{
			JsonProof: JsonMultiProof{
				BlockHeight:          uint64(commit.Height),
				OracleDataMultiProof: oracleDataList,
				BlockRelayProof:      blockRelay,
			},
			EVMProofBytes: evmProofBytes,
		})
	}
}
