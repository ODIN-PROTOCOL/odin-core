package keeper

import (
	telemetrytypes "github.com/ODIN-PROTOCOL/odin-core/x/telemetry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/rpc/core"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
	tendermint "github.com/tendermint/tendermint/types"
	"time"
)

const (
	MaxCountPerPage = 100

	AllBlocksQuery = "block.height > 1"
	OrderByAsc     = "asc"
)

func (k Keeper) GetBlocksByDates(startDate, endDate *time.Time) (map[time.Time][]*tendermint.Block, error) {
	if startDate != nil && endDate != nil {
		if ok := startDate.Before(*endDate); !ok {
			return nil, sdkerrors.Wrapf(telemetrytypes.ErrInvalidDateInterval, "invalid dates order")
		}
	}

	blocksPerDay := make(map[time.Time][]*tendermint.Block, 0)
	maxBlocksPerPage := MaxCountPerPage
	page := 1
	blocksParsed := 0

	for {
		blocks, err := core.BlockSearch(
			&rpctypes.Context{},
			AllBlocksQuery,
			&page,
			&maxBlocksPerPage,
			OrderByAsc,
		)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to find the blocks")
		}

		for _, r := range blocks.Blocks {
			blockDate := telemetrytypes.TimeToUTCDate(r.Block.Header.Time)

			if startDate != nil {
				if !blockDate.Equal(*startDate) && !blockDate.After(*startDate) {
					continue
				}
			}

			if endDate != nil {
				if !blockDate.Equal(*endDate) && !blockDate.Before(*endDate) {
					continue
				}
			}

			blocksPerDay[blockDate] = append(blocksPerDay[blockDate], r.Block)
		}

		blocksParsed += len(blocks.Blocks)
		page++

		if blocks.TotalCount == blocksParsed {
			break
		}
	}

	return blocksPerDay, nil
}

func (k Keeper) GetBlockValidators(blockHeight int64) ([]tendermint.Validator, error) {
	var validators []tendermint.Validator
	maxValidatorsPerPage := MaxCountPerPage
	page := 1

	for {
		resultValidators, err := core.Validators(&rpctypes.Context{}, &blockHeight, &page, &maxValidatorsPerPage)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to get the validators")
		}

		for _, val := range resultValidators.Validators {
			validators = append(validators, *val)
		}

		page++

		if resultValidators.Total == len(validators) {
			break
		}
	}

	return validators, nil
}

func (k Keeper) GetBlocksByValidator(ctx sdk.Context, valAddr sdk.ValAddress) ([]telemetrytypes.ValidatorBlock, error) {
	var validatorBlocks []telemetrytypes.ValidatorBlock
	maxBlocksPerPage := MaxCountPerPage
	page := 1
	blocksParsed := 0

	for {
		blocks, err := core.BlockSearch(
			&rpctypes.Context{},
			AllBlocksQuery,
			&page,
			&maxBlocksPerPage,
			OrderByAsc,
		)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to find the blocks")
		}

		for _, r := range blocks.Blocks {
			blockHeight := uint64(r.Block.Height)
			blockValidators, err := k.GetBlockValidators(r.Block.Height)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "failed to get block validators")
			}

			for _, val := range blockValidators {
				blockValAddr, err := sdk.ValAddressFromHex(val.Address.String())
				if err != nil {
					return nil, sdkerrors.Wrap(err, "failed to retrieve validator address from hex")
				}

				if blockValAddr.Equals(valAddr) {
					rewardBeforeBlock := k.distrKeeper.GetValidatorHistoricalRewards(ctx, valAddr, blockHeight-1)
					rewardAfterBlock := k.distrKeeper.GetValidatorHistoricalRewards(ctx, valAddr, blockHeight)
					reward, _ := rewardAfterBlock.CumulativeRewardRatio.
						Sub(rewardBeforeBlock.CumulativeRewardRatio).
						TruncateDecimal()

					validatorBlocks = append(validatorBlocks, telemetrytypes.ValidatorBlock{
						Height:   blockHeight,
						Time:     r.Block.Time,
						TxsCount: uint64(len(r.Block.Txs)),
						Reward:   reward,
					})

					break
				}
			}
		}

		blocksParsed += len(blocks.Blocks)
		page++

		if blocks.TotalCount == blocksParsed {
			break
		}
	}

	return validatorBlocks, nil
}
