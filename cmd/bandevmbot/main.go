package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/GeoDB-Limited/odin-core/cmd/bandevmbot/generated"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

var (
	logger log.Logger
)

func getValidators(nodeURI string) []generated.BridgeValidatorWithPower {
	node, err := rpchttp.New(nodeURI, "/websocket")
	if err != nil {
		panic(err)
	}
	validators, err := node.Validators(nil, nil, nil, nil)
	if err != nil {
		panic(err)
	}

	vals := make([]generated.BridgeValidatorWithPower, len(validators.Validators))

	for idx, validator := range validators.Validators {
		pubKeyBytes, ok := validator.PubKey.(secp256k1.PubKey)
		if !ok {
			panic("fail to cast pubkey")
		}

		pubkey, err := crypto.DecompressPubkey(pubKeyBytes[:])
		if err != nil {
			panic(err)
		}
		vals[idx] = generated.BridgeValidatorWithPower{
			Addr:  crypto.PubkeyToAddress(*pubkey),
			Power: big.NewInt(validator.VotingPower),
		}
	}
	return vals
}

func updateValidators(rpcURI string, address string, node string, privateKey string, gasPrice uint64) {
	vals := getValidators(node)
	contractAddress := common.HexToAddress(address)
	evmClient, err := ethclient.Dial(rpcURI)
	if err != nil {
		panic(err)
	}

	bridge, err := generated.NewCacheBridge(contractAddress, evmClient)
	if err != nil {
		panic(err)
	}

	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		panic(err)
	}
	publicKey := pk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic(err)
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := evmClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		panic(err)
	}
	chainId, err := evmClient.NetworkID(context.TODO())
	if err != nil {
		panic(err)
	}
	opts, err := bind.NewKeyedTransactorWithChainID(pk, chainId)
	if err != nil {
		panic(err)
	}
	opts.From = fromAddress
	opts.Nonce = big.NewInt(int64(nonce))
	opts.GasLimit = uint64(2000000000)
	opts.GasPrice = new(big.Int).SetUint64(gasPrice)

	_, err = bridge.UpdateValidatorPowers(opts, vals)
	if err != nil {
		panic(err)
	}
}

const (
	flagRPCUri          = "rpc-uri"
	flagContractAddress = "contract-address"
	flagNodeUri         = "node-uri"
	flagPrivKey         = "priv-key"
	flagGasPrice        = "gas-price"
	flagPollInterval    = "poll-interval"
)

func main() {
	cmd := &cobra.Command{
		Use:   "(--rpc-uri [rpc-uri]) (--contract-address [contract-address]) (--node-uri [node-uri]) (--priv-key [priv-key]) (--gas-price [gas-price]) (--poll-interval [poll-interval])",
		Short: "Periodically update validator set to the destination EVM blockchain",
		Args:  cobra.ExactArgs(0),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Periodically update validator set to the destination EVM blockchain
Example:
$ --rpc-uri https://kovan.infura.io/v3/d3301689638b40dabad8395bf00d3945 --contract-address 0x0d8152D22a05A3Cf2cE1c5bEfCc2F8658f75a59d --node-uri http://d3n-debug.bandprotocol.com:26657 --priv-key AA0C65C16D4B8511C58122966F94192F6963D0EB7896435430BCDFF56E9F13B9 --gas-price 1000000 --poll-interval 24
`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {

			rpcURI, err := cmd.Flags().GetString(flagRPCUri)
			if err != nil {
				return err
			}
			contractAddress, err := cmd.Flags().GetString(flagContractAddress)
			if err != nil {
				return err
			}
			nodeURI, err := cmd.Flags().GetString(flagNodeUri)
			if err != nil {
				return err
			}
			privateKey, err := cmd.Flags().GetString(flagPrivKey)
			if err != nil {
				return err
			}
			rawGasPrice, err := cmd.Flags().GetString(flagGasPrice)
			if err != nil {
				return err
			}
			gasPrice, err := strconv.ParseInt(rawGasPrice, 10, 64)
			if err != nil {
				return err
			}
			rawInterval, err := cmd.Flags().GetString(flagPollInterval)
			if err != nil {
				return err
			}
			interval, err := strconv.ParseInt(rawInterval, 10, 64)
			if err != nil {
				return err
			}

			for {
				updateValidators(rpcURI, contractAddress, nodeURI, privateKey, uint64(gasPrice))
				fmt.Println("finish round")
				time.Sleep(time.Duration(interval) * time.Hour)
			}
		},
	}
	cmd.Flags().String(flagRPCUri, "", "RPC URI")
	cmd.Flags().String(flagContractAddress, "", "Address of contract")
	cmd.Flags().String(flagNodeUri, "", "Node URI")
	cmd.Flags().String(flagPrivKey, "", "Private key")
	cmd.Flags().String(flagGasPrice, "", "Gas Price")
	cmd.Flags().String(flagPollInterval, "", "Interval of update validatos (Hours)")

	err := cmd.Execute()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed executing: %s, exiting...", err))
		os.Exit(1)

	}
}
