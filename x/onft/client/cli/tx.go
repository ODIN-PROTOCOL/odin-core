package cli

import (
	"github.com/ODIN-PROTOCOL/odin-core/x/onft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

// NewTxCmd returns the transaction commands for the NFT module
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "NFT transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		GetCmdCreateNFTClass(),
		GetCmdTransferClassOwnership(),
		GetCmdMintNFT(),
	)

	return txCmd
}

// GetCmdCreateNFTClass implements the create NFT class command.
func GetCmdCreateNFTClass() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-class [name] [symbol] [description] [uri] [uri-hash]",
		Short: "Create a new NFT class",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateNFTClass(
				args[0], // name
				args[1], // symbol
				args[2], // description
				args[3], // uri
				args[4], // uri-hash
				nil,
				clientCtx.GetFromAddress(), // sender
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdTransferClassOwnership implements the transfer class ownership command.
func GetCmdTransferClassOwnership() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-class-ownership [class-id] [new-owner]",
		Short: "Transfer ownership of an NFT class",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgTransferClassOwnership(
				args[0],                    // class ID
				clientCtx.GetFromAddress(), // sender
				args[1],                    // new owner
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdMintNFT implements the mint NFT command.
func GetCmdMintNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [class-id] [uri] [uri-hash] [receiver]",
		Short: "Mint a new NFT in a given class",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgMintNFT(
				args[0],                    // class ID
				args[1],                    // uri
				args[2],                    // uri-hash
				clientCtx.GetFromAddress(), // sender
				args[3],                    // receiver
				nil,
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
