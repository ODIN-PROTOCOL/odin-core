package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ODIN-PROTOCOL/odin-core/x/globalfee/types"
)

func (s *IntegrationTestSuite) TestUpdateParams() {
	testCases := []struct {
		name      string
		request   *types.MsgUpdateParams
		expectErr string
	}{
		{
			name: "set invalid authority",
			request: &types.MsgUpdateParams{
				Authority: "foo",
			},
			expectErr: "invalid authority",
		},
		{
			name: "set full valid params",
			request: &types.MsgUpdateParams{
				Authority: s.globalfeeKeeper.GetAuthority(),
				Params: types.Params{
					MinimumGasPrices: sdk.NewDecCoins(
						sdk.NewDecCoin("ALX", sdk.NewInt(1)),
						sdk.NewDecCoinFromDec("BLX", sdk.NewDecWithPrec(1, 3)),
					),
				},
			},
			expectErr: "",
		},
		{
			name: "set empty coin",
			request: &types.MsgUpdateParams{
				Authority: s.globalfeeKeeper.GetAuthority(),
				Params: types.Params{
					MinimumGasPrices: sdk.DecCoins(nil),
				},
			},
			expectErr: "",
		},
		{
			name: "set invalid denom",
			request: &types.MsgUpdateParams{
				Authority: s.globalfeeKeeper.GetAuthority(),
				Params: types.Params{
					MinimumGasPrices: []sdk.DecCoin{
						{
							Denom:  "1AAAA",
							Amount: sdk.NewDecFromInt(sdk.NewInt(1)),
						},
					},
				},
			},
			expectErr: "invalid denom",
		},
		{
			name: "set negative value",
			request: &types.MsgUpdateParams{
				Authority: s.globalfeeKeeper.GetAuthority(),
				Params: types.Params{
					MinimumGasPrices: []sdk.DecCoin{
						{
							Denom:  "AAAA",
							Amount: sdk.NewDecFromInt(sdk.NewInt(-1)),
						},
					},
				},
			},
			expectErr: "is not positive",
		},
		{
			name: "set duplicated denom",
			request: &types.MsgUpdateParams{
				Authority: s.globalfeeKeeper.GetAuthority(),
				Params: types.Params{
					MinimumGasPrices: []sdk.DecCoin{
						{
							Denom:  "AAAA",
							Amount: sdk.NewDecFromInt(sdk.NewInt(1)),
						},
						{
							Denom:  "AAAA",
							Amount: sdk.NewDecFromInt(sdk.NewInt(2)),
						},
					},
				},
			},
			expectErr: "duplicate denomination",
		},
		{
			name: "set unsorted denom",
			request: &types.MsgUpdateParams{
				Authority: s.globalfeeKeeper.GetAuthority(),
				Params: types.Params{
					MinimumGasPrices: []sdk.DecCoin{
						{
							Denom:  "BBBB",
							Amount: sdk.NewDecFromInt(sdk.NewInt(1)),
						},
						{
							Denom:  "AAAA",
							Amount: sdk.NewDecFromInt(sdk.NewInt(2)),
						},
					},
				},
			},
			expectErr: "is not sorted",
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			_, err := s.msgServer.UpdateParams(s.ctx, tc.request)
			if tc.expectErr != "" {
				s.Require().ErrorContains(err, tc.expectErr)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
