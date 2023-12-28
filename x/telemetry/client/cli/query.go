package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"

	telemetrytypes "github.com/ODIN-PROTOCOL/odin-core/x/telemetry/types"
)

const (
	DateFormat = "2006-01-02"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd() *cobra.Command {
	coinswapCmd := &cobra.Command{
		Use:                        telemetrytypes.ModuleName,
		Short:                      "Querying commands for the telemetry module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	coinswapCmd.AddCommand(
		GetQueryCmdTopBalances(),
		GetQueryCmdExtendedValidators(),
		/*GetQueryCmdAvgBlockSize(),
		GetQueryCmdAvgBlockTime(),
		GetQueryCmdAvgTxFee(),
		GetQueryCmdTxVolume(),
		GetQueryCmdValidatorBlocks(),*/
		GetQueryCmdTopValidators(),
		GetQueryCmdValidatorByConsAddr(),
	)
	return coinswapCmd
}

// GetQueryCmdTopBalances implements the query parameters command.
func GetQueryCmdTopBalances() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "top-balances [denom]",
		Short: "Query for top balances",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the top balances of a specific denomination.

Example:
  $ %[1]s query %[2]s top-balances [denom]
  $ %[1]s query %[2]s top-balances [denom] --limit=100 --offset=2 --desc=true
`,
				version.AppName, telemetrytypes.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			flagSet := cmd.Flags()
			pageReq, err := client.ReadPageRequest(flagSet)
			if err != nil {
				return err
			}
			desc, _ := flagSet.GetBool(flagDesc)

			queryClient := telemetrytypes.NewQueryClient(clientCtx)
			res, err := queryClient.TopBalances(cmd.Context(), &telemetrytypes.QueryTopBalancesRequest{
				Denom:      args[0],
				Pagination: pageReq,
				Desc:       desc,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "top balances")
	cmd.Flags().Bool(flagDesc, false, "desc is used in calling the data with sort by desc")

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetQueryCmdExtendedValidators() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extended-validators [status]",
		Short: "Query for extended validators",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for extended validators.

Example:
  $ %[1]s query %[2]s extended-validators [status]
  $ %[1]s query %[2]s extended-validators [status] --limit=100 --offset=2
`,
				version.AppName, telemetrytypes.ModuleName,
			),
		),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := telemetrytypes.NewQueryClient(clientCtx)
			res, err := queryClient.ExtendedValidators(cmd.Context(), &telemetrytypes.QueryExtendedValidatorsRequest{
				Status:     args[0],
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "top balances")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetQueryCmdAvgBlockSize implements the query parameters command.
/*func GetQueryCmdAvgBlockSize() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "avg-block-size [start-date] [end-date]",
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			startDate, endDate, err := ParseDateInterval(args[0], args[1])
			if err != nil {
				return sdkerrors.Wrap(err, "failed to parse date interval")
			}

			queryClient := telemetrytypes.NewQueryClient(clientCtx)
			res, err := queryClient.AvgBlockSize(cmd.Context(), &telemetrytypes.QueryAvgBlockSizeRequest{
				StartDate: startDate,
				EndDate:   endDate,
			})
			if err != nil {
				return sdkerrors.Wrap(err, "failed to query average block size")
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetQueryCmdAvgBlockTime implements the query parameters command.
func GetQueryCmdAvgBlockTime() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "avg-block-time [start-date] [end-date]",
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			startDate, endDate, err := ParseDateInterval(args[0], args[1])
			if err != nil {
				return sdkerrors.Wrap(err, "failed to parse date interval")
			}

			queryClient := telemetrytypes.NewQueryClient(clientCtx)
			res, err := queryClient.AvgBlockTime(cmd.Context(), &telemetrytypes.QueryAvgBlockTimeRequest{
				StartDate: startDate,
				EndDate:   endDate,
			})
			if err != nil {
				return sdkerrors.Wrap(err, "failed to query average block time")
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetQueryCmdAvgTxFee implements the query parameters command.
func GetQueryCmdAvgTxFee() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "avg-tx-fee [start-date] [end-date]",
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			startDate, endDate, err := ParseDateInterval(args[0], args[1])
			if err != nil {
				return sdkerrors.Wrap(err, "failed to parse date interval")
			}

			queryClient := telemetrytypes.NewQueryClient(clientCtx)
			res, err := queryClient.AvgTxFee(cmd.Context(), &telemetrytypes.QueryAvgTxFeeRequest{
				StartDate: startDate,
				EndDate:   endDate,
			})
			if err != nil {
				return sdkerrors.Wrap(err, "failed to query average tx fee")
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetQueryCmdTxVolume implements the query parameters command.
func GetQueryCmdTxVolume() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "tx-volume [start-date] [end-date]",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			startDate, endDate, err := ParseDateInterval(args[0], args[1])
			if err != nil {
				return sdkerrors.Wrap(err, "failed to parse date interval")
			}

			queryClient := telemetrytypes.NewQueryClient(clientCtx)
			res, err := queryClient.TxVolume(cmd.Context(), &telemetrytypes.QueryTxVolumeRequest{
				StartDate: startDate,
				EndDate:   endDate,
			})
			if err != nil {
				return sdkerrors.Wrap(err, "failed to query tx volume")
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}*/

// GetQueryCmdValidatorBlocks implements the query parameters command.
func GetQueryCmdValidatorBlocks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-blocks [address]",
		Short: "Query for validator blocks",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for validator blocks.

Example:
  $ %[1]s query %[2]s validator-blocks [address]
  $ %[1]s query %[2]s validator-blocks [address] --limit=100 --offset=2 --desc=true
`,
				version.AppName, telemetrytypes.ModuleName,
			),
		),
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return sdkerrors.Wrap(err, "failed to get client context")
			}

			valAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			flagSet := cmd.Flags()
			pageReq, err := client.ReadPageRequest(flagSet)
			if err != nil {
				return err
			}
			desc, _ := flagSet.GetBool(flagDesc)

			queryClient := telemetrytypes.NewQueryClient(clientCtx)
			res, err := queryClient.ValidatorBlocks(cmd.Context(), &telemetrytypes.QueryValidatorBlocksRequest{
				ValidatorAddress: valAddr.String(),
				Pagination:       pageReq,
				Desc:             desc,
			})
			if err != nil {
				return sdkerrors.Wrap(err, "failed to query validator blocks")
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "validator blocks")
	cmd.Flags().Bool(flagDesc, false, "desc is used in calling the data with sort by desc")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetQueryCmdTopValidators implements the query parameters command.
func GetQueryCmdTopValidators() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "top-validators [start-date] [end-date]",
		Short: "Query for top validators by blocks",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for top validators by blocks.

Example:
  $ %[1]s query %[2]s top-validators [start-date] [end-date]
  $ %[1]s query %[2]s top-validators [start-date] [end-date] --limit=100 --offset=2 --desc=true
`,
				version.AppName, telemetrytypes.ModuleName,
			),
		),
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return sdkerrors.Wrap(err, "failed to get client context")
			}

			startDate, endDate, err := ParseDateInterval(args[0], args[1])
			if err != nil {
				return sdkerrors.Wrap(err, "failed to parse date interval")
			}

			flagSet := cmd.Flags()
			pageReq, err := client.ReadPageRequest(flagSet)
			if err != nil {
				return err
			}
			desc, _ := flagSet.GetBool(flagDesc)

			queryClient := telemetrytypes.NewQueryClient(clientCtx)
			res, err := queryClient.TopValidators(cmd.Context(), &telemetrytypes.QueryTopValidatorsRequest{
				StartDate:  startDate,
				EndDate:    endDate,
				Pagination: pageReq,
				Desc:       desc,
			})
			if err != nil {
				return sdkerrors.Wrap(err, "failed to query top validators by blocks")
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "top balances")
	cmd.Flags().Bool(flagDesc, false, "desc is used in calling the data with sort by desc")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func ParseDateInterval(startDateArg, endDateArg string) (*time.Time, *time.Time, error) {
	var startDate, endDate *time.Time

	if startDateArg != "" {
		sd, err := time.Parse(DateFormat, startDateArg)
		if err != nil {
			return nil, nil, sdkerrors.Wrap(err, "failed to parse start date")
		}
		startDate = &sd
	}
	if endDateArg != "" {
		ed, err := time.Parse(DateFormat, endDateArg)
		if err != nil {
			return nil, nil, sdkerrors.Wrap(err, "failed to parse end date")
		}
		endDate = &ed
	}

	return startDate, endDate, nil
}

func GetQueryCmdValidatorByConsAddr() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-by-consaddr [status]",
		Short: "Query for validator by consensus address",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := telemetrytypes.NewQueryClient(clientCtx)
			res, err := queryClient.ValidatorByConsAddr(cmd.Context(), &telemetrytypes.QueryValidatorByConsAddrRequest{
				ConsensusAddress: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
