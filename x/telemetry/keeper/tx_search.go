package keeper

import (
	"github.com/cometbft/cometbft/rpc/core"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	rpctypes "github.com/cometbft/cometbft/rpc/jsonrpc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	QueryAllTxs = "tx.height >= 1"
)

// GetAccountTxsCount TODO: very bad method, it is better not to use it, or rewrite it more optimized
func (k Keeper) GetAccountTxsCount(accounts ...sdk.AccAddress) (map[string]int64, error) {
	accountsMap := make(map[string]bool)
	for _, acc := range accounts {
		accountsMap[acc.String()] = true
	}

	page := 1

	// make first call, to get total count
	txSearch, err := k.GetTxs(page, MaxCountPerPage)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get paginated transactions")
	}

	// repeat calls until we are above total count
	allTxs := make([]*coretypes.ResultTx, 0)
	for page++; page*MaxCountPerPage < txSearch.TotalCount; page++ {
		txSearch, err = k.GetTxs(page, MaxCountPerPage)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to get paginated transactions")
		}
		allTxs = append(allTxs, txSearch.Txs...)
	}

	// if something less than for a page left, get it as well
	if txSearch.TotalCount/MaxCountPerPage >= page {
		txSearch, err = k.GetTxs(page, MaxCountPerPage)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to get paginated transactions")
		}
		allTxs = append(allTxs, txSearch.Txs...)
	}

	res := make(map[string]int64)
	for _, tx := range txSearch.Txs {
		rawTx, err := k.txCfg.TxDecoder()(tx.Tx)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to decode tx")
		}
		txBuilder, err := k.txCfg.WrapTxBuilder(rawTx)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to get tx builder")
		}
		for _, signer := range txBuilder.GetTx().GetSigners() {
			if _, ok := accountsMap[signer.String()]; ok {
				res[signer.String()]++
			}
		}
	}

	return res, nil
}

func (k Keeper) GetTxs(page, maxTxsPerPage int) (*coretypes.ResultTxSearch, error) {
	resSearch, err := core.TxSearch(
		&rpctypes.Context{},
		QueryAllTxs,
		false,
		&page,
		&maxTxsPerPage,
		OrderByAsc,
	)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to find txs")
	}

	return resSearch, nil
}
