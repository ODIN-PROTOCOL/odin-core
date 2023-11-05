package proof

import (
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/cometbft/cometbft/crypto/secp256k1"
	"github.com/cometbft/cometbft/crypto/tmhash"
	tmbytes "github.com/cometbft/cometbft/libs/bytes"
	"github.com/cometbft/cometbft/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// TMSignature contains all details of validator signature for performing signer recovery for ECDSA
// secp256k1 signature. Note that this struct is written specifically for signature signed on
// Tendermint's precommit data, which includes the block hash and some additional information prepended
// and appended to the block hash. The prepended part (prefix) and the appended part (suffix) are
// different for each signer (including signature size, machine clock, validator index, etc).
type TMSignature struct {
	R                tmbytes.HexBytes `json:"r"`
	S                tmbytes.HexBytes `json:"s"`
	V                uint8            `json:"v"`
	SignedDataPrefix tmbytes.HexBytes `json:"signed_data_prefix"`
	SignedDataSuffix tmbytes.HexBytes `json:"signed_data_suffix"`
}

// TMSignatureEthereum is an Ethereum version of TMSignature for solidity ABI-encoding.
type TMSignatureEthereum struct {
	R                common.Hash
	S                common.Hash
	V                uint8
	SignedDataPrefix []byte
	SignedDataSuffix []byte
}

func (signature *TMSignature) encodeToEthFormat() TMSignatureEthereum {
	return TMSignatureEthereum{
		R:                common.BytesToHash(signature.R),
		S:                common.BytesToHash(signature.S),
		V:                signature.V,
		SignedDataPrefix: signature.SignedDataPrefix,
		SignedDataSuffix: signature.SignedDataSuffix,
	}
}

func recoverETHAddress(msg, sig, signer []byte) ([]byte, uint8, error) {
	for i := uint8(0); i < 2; i++ {
		pubuc, err := crypto.SigToPub(tmhash.Sum(msg), append(sig, i))
		if err != nil {
			return nil, 0, err
		}
		pub := crypto.CompressPubkey(pubuc)
		var tmp [33]byte

		copy(tmp[:], pub)
		if string(signer) == string(secp256k1.PubKey(tmp[:]).Address()) {
			return crypto.PubkeyToAddress(*pubuc).Bytes(), 27 + i, nil
		}
	}
	return nil, 0, fmt.Errorf("No match address found")
}

// GetSignaturesAndPrefix returns a list of TMSignature from Tendermint signed header.
func GetSignaturesAndPrefix(info *types.SignedHeader) ([]TMSignature, error) {
	var addrs []string
	mapAddrs := map[string]TMSignature{}
	for i, vote := range info.Commit.Signatures {
		if !vote.ForBlock() {
			continue
		}
		msg := info.Commit.VoteSignBytes(info.ChainID, int32(i))
		lr := strings.Split(hex.EncodeToString(msg), hex.EncodeToString(info.Commit.BlockID.Hash))

		if len(lr) != 2 {
			return nil, fmt.Errorf("Split block hash failed")
		}
		addr, v, err := recoverETHAddress(msg, vote.Signature, vote.ValidatorAddress)
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, string(addr))
		mapAddrs[string(addr)] = TMSignature{
			vote.Signature[:32],
			vote.Signature[32:],
			v,
			mustDecodeString(lr[0]),
			mustDecodeString(lr[1]),
		}
	}
	if len(addrs) == 0 {
		return nil, fmt.Errorf("No valid precommit")
	}

	signatures := make([]TMSignature, len(addrs))
	sort.Strings(addrs)
	for i, addr := range addrs {
		signatures[i] = mapAddrs[addr]
	}

	return signatures, nil
}
