package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
)

// GetQueryCmd returns the cli query commands for the minting module.
func GetQueryCmd() *cobra.Command {
	mintingQueryCmd := &cobra.Command{
		Use:                        minttypes.ModuleName,
		Short:                      "Querying commands for the minting module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	mintingQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryInflation(),
		GetCmdQueryAnnualProvisions(),
		//GetCmdQueryIntegrationAddress(),
		GetCmdQueryTreasuryPool(),
		GetCmdQueryCurrentMintVolume(),
	)

	return mintingQueryCmd
}

// GetCmdQueryParams implements a command to return the current minting
// parameters.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current minting parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := minttypes.NewQueryClient(clientCtx)

			params := &minttypes.QueryParamsRequest{}
			res, err := queryClient.Params(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryInflation implements a command to return the current minting
// inflation value.
func GetCmdQueryInflation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inflation",
		Short: "Query the current minting inflation value",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := minttypes.NewQueryClient(clientCtx)

			params := &minttypes.QueryInflationRequest{}
			res, err := queryClient.Inflation(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintString(fmt.Sprintf("%s\n", res.Inflation))
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryAnnualProvisions implements a command to return the current minting
// annual provisions value.
func GetCmdQueryAnnualProvisions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "annual-provisions",
		Short: "Query the current minting annual provisions value",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := minttypes.NewQueryClient(clientCtx)

			params := &minttypes.QueryAnnualProvisionsRequest{}
			res, err := queryClient.AnnualProvisions(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintString(fmt.Sprintf("%s\n", res.AnnualProvisions))
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryIntegrationAddress returns the command for fetching integration address
//func GetCmdQueryIntegrationAddress() *cobra.Command {
//	cmd := &cobra.Command{
//		Use:   "integration-address [network-name]",
//		Short: "Query current integration address by network name",
//		Args:  cobra.ExactArgs(1),
//		RunE: func(cmd *cobra.Command, args []string) error {
//			clientCtx, err := client.GetClientQueryContext(cmd)
//			if err != nil {
//				return err
//			}
//
//			queryClient := minttypes.NewQueryClient(clientCtx)
//
//			res, err := queryClient.IntegrationAddress(context.Background(), &minttypes.QueryIntegrationAddressRequest{
//				NetworkName: args[0],
//			})
//			if err != nil {
//				return err
//			}
//
//			return clientCtx.PrintString(fmt.Sprintf("%s\n", res.IntegrationAddress))
//		},
//	}
//
//	flags.AddQueryFlagsToCmd(cmd)
//
//	return cmd
//}

// GetCmdQueryTreasuryPool returns the command for fetching treasury pool info
func GetCmdQueryTreasuryPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "treasury-pool",
		Short: "Query the amount of coins in the treasury pool",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := minttypes.NewQueryClient(clientCtx)

			params := &minttypes.QueryTreasuryPoolRequest{}
			res, err := queryClient.TreasuryPool(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintString(fmt.Sprintf("%s\n", res.TreasuryPool))
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryCurrentMintVolume returns the command for getting minted coins volume
func GetCmdQueryCurrentMintVolume() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current-mint-volume",
		Short: "Query the amount of minted coins",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := minttypes.NewQueryClient(clientCtx)

			params := &minttypes.QueryCurrentMintVolumeRequest{}
			res, err := queryClient.CurrentMintVolume(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintString(fmt.Sprintf("%s\n", res.CurrentMintVolume))
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
