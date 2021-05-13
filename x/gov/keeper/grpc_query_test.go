package keeper_test

import (
	gocontext "context"
	"fmt"
	"github.com/GeoDB-Limited/odin-core/x/common/testapp"
	"strconv"

	govtypes "github.com/GeoDB-Limited/odin-core/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

func (suite *KeeperTestSuite) TestGRPCQueryProposal() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	var (
		req         *govtypes.QueryProposalRequest
		expProposal govtypes.Proposal
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &govtypes.QueryProposalRequest{}
			},
			false,
		},
		{
			"non existing proposal request",
			func() {
				req = &govtypes.QueryProposalRequest{ProposalId: 3}
			},
			false,
		},
		{
			"zero proposal id request",
			func() {
				req = &govtypes.QueryProposalRequest{ProposalId: 0}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &govtypes.QueryProposalRequest{ProposalId: 1}
				testProposal := govtypes.NewTextProposal("Proposal", "testing proposal")
				submittedProposal, err := app.GovKeeper.SubmitProposal(ctx, testProposal)
				suite.Require().NoError(err)
				suite.Require().NotEmpty(submittedProposal)

				expProposal = submittedProposal
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			proposalRes, err := queryClient.Proposal(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expProposal.String(), proposalRes.Proposal.String())
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(proposalRes)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryProposals() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	testProposals := []govtypes.Proposal{}

	var (
		req    *govtypes.QueryProposalsRequest
		expRes *govtypes.QueryProposalsResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty state request",
			func() {
				req = &govtypes.QueryProposalsRequest{}
			},
			true,
		},
		{
			"request proposals with limit 3",
			func() {
				// create 5 test proposals
				for i := 0; i < 5; i++ {
					num := strconv.Itoa(i + 1)
					testProposal := govtypes.NewTextProposal("Proposal"+num, "testing proposal "+num)
					proposal, err := app.GovKeeper.SubmitProposal(ctx, testProposal)
					suite.Require().NotEmpty(proposal)
					suite.Require().NoError(err)
					testProposals = append(testProposals, proposal)
				}

				req = &govtypes.QueryProposalsRequest{
					Pagination: &query.PageRequest{Limit: 3},
				}

				expRes = &govtypes.QueryProposalsResponse{
					Proposals: testProposals[:3],
				}
			},
			true,
		},
		{
			"request 2nd page with limit 4",
			func() {
				req = &govtypes.QueryProposalsRequest{
					Pagination: &query.PageRequest{Offset: 3, Limit: 3},
				}

				expRes = &govtypes.QueryProposalsResponse{
					Proposals: testProposals[3:],
				}
			},
			true,
		},
		{
			"request with limit 2 and count true",
			func() {
				req = &govtypes.QueryProposalsRequest{
					Pagination: &query.PageRequest{Limit: 2, CountTotal: true},
				}

				expRes = &govtypes.QueryProposalsResponse{
					Proposals: testProposals[:2],
				}
			},
			true,
		},
		{
			"request with filter of status deposit period",
			func() {
				req = &govtypes.QueryProposalsRequest{
					ProposalStatus: govtypes.StatusDepositPeriod,
				}

				expRes = &govtypes.QueryProposalsResponse{
					Proposals: testProposals,
				}
			},
			true,
		},
		{
			"request with filter of deposit address",
			func() {
				depositCoins := sdk.NewCoins(sdk.NewCoin(govtypes.DefaultBondDenom, sdk.TokensFromConsensusPower(20)))
				deposit := govtypes.NewDeposit(testProposals[0].ProposalId, addrs[0], depositCoins)
				app.GovKeeper.SetDeposit(ctx, deposit)

				req = &govtypes.QueryProposalsRequest{
					Depositor: addrs[0].String(),
				}

				expRes = &govtypes.QueryProposalsResponse{
					Proposals: testProposals[:1],
				}
			},
			true,
		},
		{
			"request with filter of deposit address",
			func() {
				testProposals[1].Status = govtypes.StatusVotingPeriod
				app.GovKeeper.SetProposal(ctx, testProposals[1])
				suite.Require().NoError(app.GovKeeper.AddVote(ctx, testProposals[1].ProposalId, addrs[0], govtypes.OptionAbstain))

				req = &govtypes.QueryProposalsRequest{
					Voter: addrs[0].String(),
				}

				expRes = &govtypes.QueryProposalsResponse{
					Proposals: testProposals[1:2],
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			proposals, err := queryClient.Proposals(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)

				suite.Require().Len(proposals.GetProposals(), len(expRes.GetProposals()))
				for i := 0; i < len(proposals.GetProposals()); i++ {
					suite.Require().Equal(proposals.GetProposals()[i].String(), expRes.GetProposals()[i].String())
				}

			} else {
				suite.Require().Error(err)
				suite.Require().Nil(proposals)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryVote() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	var (
		req      *govtypes.QueryVoteRequest
		expRes   *govtypes.QueryVoteResponse
		proposal govtypes.Proposal
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &govtypes.QueryVoteRequest{}
			},
			false,
		},
		{
			"zero proposal id request",
			func() {
				req = &govtypes.QueryVoteRequest{
					ProposalId: 0,
					Voter:      addrs[0].String(),
				}
			},
			false,
		},
		{
			"empty voter request",
			func() {
				req = &govtypes.QueryVoteRequest{
					ProposalId: 1,
					Voter:      "",
				}
			},
			false,
		},
		{
			"non existed proposal",
			func() {
				req = &govtypes.QueryVoteRequest{
					ProposalId: 3,
					Voter:      addrs[0].String(),
				}
			},
			false,
		},
		{
			"no votes present",
			func() {
				var err error
				proposal, err = app.GovKeeper.SubmitProposal(ctx, TestProposal)
				suite.Require().NoError(err)

				req = &govtypes.QueryVoteRequest{
					ProposalId: proposal.ProposalId,
					Voter:      addrs[0].String(),
				}

				expRes = &govtypes.QueryVoteResponse{}
			},
			false,
		},
		{
			"valid request",
			func() {
				proposal.Status = govtypes.StatusVotingPeriod
				app.GovKeeper.SetProposal(ctx, proposal)
				suite.Require().NoError(app.GovKeeper.AddVote(ctx, proposal.ProposalId, addrs[0], govtypes.OptionAbstain))

				req = &govtypes.QueryVoteRequest{
					ProposalId: proposal.ProposalId,
					Voter:      addrs[0].String(),
				}

				expRes = &govtypes.QueryVoteResponse{Vote: govtypes.NewVote(proposal.ProposalId, addrs[0], govtypes.OptionAbstain)}
			},
			true,
		},
		{
			"wrong voter id request",
			func() {
				req = &govtypes.QueryVoteRequest{
					ProposalId: proposal.ProposalId,
					Voter:      addrs[1].String(),
				}

				expRes = &govtypes.QueryVoteResponse{}
			},
			false,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			vote, err := queryClient.Vote(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes, vote)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(vote)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryVotes() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	addrs := []sdk.AccAddress{testapp.Alice.Address, testapp.Bob.Address}

	var (
		req      *govtypes.QueryVotesRequest
		expRes   *govtypes.QueryVotesResponse
		proposal govtypes.Proposal
		votes    govtypes.Votes
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &govtypes.QueryVotesRequest{}
			},
			false,
		},
		{
			"zero proposal id request",
			func() {
				req = &govtypes.QueryVotesRequest{
					ProposalId: 0,
				}
			},
			false,
		},
		{
			"non existed proposals",
			func() {
				req = &govtypes.QueryVotesRequest{
					ProposalId: 2,
				}
			},
			true,
		},
		{
			"create a proposal and get votes",
			func() {
				var err error
				proposal, err = app.GovKeeper.SubmitProposal(ctx, TestProposal)
				suite.Require().NoError(err)

				req = &govtypes.QueryVotesRequest{
					ProposalId: proposal.ProposalId,
				}
			},
			true,
		},
		{
			"request after adding 2 votes",
			func() {
				proposal.Status = govtypes.StatusVotingPeriod
				app.GovKeeper.SetProposal(ctx, proposal)

				votes = []govtypes.Vote{
					{proposal.ProposalId, addrs[0].String(), govtypes.OptionAbstain},
					{proposal.ProposalId, addrs[1].String(), govtypes.OptionYes},
				}
				accAddr1, err1 := sdk.AccAddressFromBech32(votes[0].Voter)
				accAddr2, err2 := sdk.AccAddressFromBech32(votes[1].Voter)
				suite.Require().NoError(err1)
				suite.Require().NoError(err2)
				suite.Require().NoError(app.GovKeeper.AddVote(ctx, proposal.ProposalId, accAddr1, votes[0].Option))
				suite.Require().NoError(app.GovKeeper.AddVote(ctx, proposal.ProposalId, accAddr2, votes[1].Option))

				req = &govtypes.QueryVotesRequest{
					ProposalId: proposal.ProposalId,
				}

				expRes = &govtypes.QueryVotesResponse{
					Votes: votes,
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			votes, err := queryClient.Votes(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes.GetVotes(), votes.GetVotes())
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(votes)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryParams() {
	queryClient := suite.queryClient

	var (
		req    *govtypes.QueryParamsRequest
		expRes *govtypes.QueryParamsResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &govtypes.QueryParamsRequest{}
			},
			false,
		},
		{
			"deposit params request",
			func() {
				req = &govtypes.QueryParamsRequest{ParamsType: govtypes.ParamDeposit}
				expRes = &govtypes.QueryParamsResponse{
					DepositParams: govtypes.DefaultDepositParams(),
					TallyParams:   govtypes.NewTallyParams(sdk.NewDec(0), sdk.NewDec(0), sdk.NewDec(0)),
				}
			},
			true,
		},
		{
			"voting params request",
			func() {
				req = &govtypes.QueryParamsRequest{ParamsType: govtypes.ParamVoting}
				expRes = &govtypes.QueryParamsResponse{
					VotingParams: govtypes.DefaultVotingParams(),
					TallyParams:  govtypes.NewTallyParams(sdk.NewDec(0), sdk.NewDec(0), sdk.NewDec(0)),
				}
			},
			true,
		},
		{
			"tally params request",
			func() {
				req = &govtypes.QueryParamsRequest{ParamsType: govtypes.ParamTallying}
				expRes = &govtypes.QueryParamsResponse{
					TallyParams: govtypes.DefaultTallyParams(),
				}
			},
			true,
		},
		{
			"invalid request",
			func() {
				req = &govtypes.QueryParamsRequest{ParamsType: "wrongPath"}
				expRes = &govtypes.QueryParamsResponse{}
			},
			false,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			params, err := queryClient.Params(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes.GetDepositParams(), params.GetDepositParams())
				suite.Require().Equal(expRes.GetVotingParams(), params.GetVotingParams())
				suite.Require().Equal(expRes.GetTallyParams(), params.GetTallyParams())
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(params)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryDeposit() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	var (
		req      *govtypes.QueryDepositRequest
		expRes   *govtypes.QueryDepositResponse
		proposal govtypes.Proposal
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &govtypes.QueryDepositRequest{}
			},
			false,
		},
		{
			"zero proposal id request",
			func() {
				req = &govtypes.QueryDepositRequest{
					ProposalId: 0,
					Depositor:  addrs[0].String(),
				}
			},
			false,
		},
		{
			"empty deposit address request",
			func() {
				req = &govtypes.QueryDepositRequest{
					ProposalId: 1,
					Depositor:  "",
				}
			},
			false,
		},
		{
			"non existed proposal",
			func() {
				req = &govtypes.QueryDepositRequest{
					ProposalId: 2,
					Depositor:  addrs[0].String(),
				}
			},
			false,
		},
		{
			"no deposits proposal",
			func() {
				var err error
				proposal, err = app.GovKeeper.SubmitProposal(ctx, TestProposal)
				suite.Require().NoError(err)
				suite.Require().NotNil(proposal)

				req = &govtypes.QueryDepositRequest{
					ProposalId: proposal.ProposalId,
					Depositor:  addrs[0].String(),
				}
			},
			false,
		},
		{
			"valid request",
			func() {
				depositCoins := sdk.NewCoins(sdk.NewCoin(govtypes.DefaultBondDenom, sdk.TokensFromConsensusPower(20)))
				deposit := govtypes.NewDeposit(proposal.ProposalId, addrs[0], depositCoins)
				app.GovKeeper.SetDeposit(ctx, deposit)

				req = &govtypes.QueryDepositRequest{
					ProposalId: proposal.ProposalId,
					Depositor:  addrs[0].String(),
				}

				expRes = &govtypes.QueryDepositResponse{Deposit: deposit}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			deposit, err := queryClient.Deposit(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(deposit.GetDeposit(), expRes.GetDeposit())
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(expRes)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryDeposits() {
	app, ctx, queryClient, _ := suite.app, suite.ctx, suite.queryClient, suite.addrs

	var (
		req      *govtypes.QueryDepositsRequest
		expRes   *govtypes.QueryDepositsResponse
		proposal govtypes.Proposal
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &govtypes.QueryDepositsRequest{}
			},
			false,
		},
		{
			"zero proposal id request",
			func() {
				req = &govtypes.QueryDepositsRequest{
					ProposalId: 0,
				}
			},
			false,
		},
		{
			"non existed proposal",
			func() {
				req = &govtypes.QueryDepositsRequest{
					ProposalId: 2,
				}
			},
			true,
		},
		{
			"create a proposal and get deposits",
			func() {
				var err error
				proposal, err = app.GovKeeper.SubmitProposal(ctx, TestProposal)
				suite.Require().NoError(err)

				req = &govtypes.QueryDepositsRequest{
					ProposalId: proposal.ProposalId,
				}
			},
			true,
		},
		//{
		//	"get deposits with default limit",
		//	func() {
		//		depositAmount1 := sdk.NewCoins(sdk.NewCoin(govtypes.DefaultBondDenom, sdk.TokensFromConsensusPower(20)))
		//		deposit1 := govtypes.NewDeposit(proposal.ProposalId, addrs[0], depositAmount1)
		//		app.GovKeeper.SetDeposit(ctx, deposit1)
		//
		//		depositAmount2 := sdk.NewCoins(sdk.NewCoin(govtypes.DefaultBondDenom, sdk.TokensFromConsensusPower(30)))
		//		deposit2 := govtypes.NewDeposit(proposal.ProposalId, addrs[1], depositAmount2)
		//		app.GovKeeper.SetDeposit(ctx, deposit2)
		//
		//		deposits := govtypes.Deposits{deposit1, deposit2}
		//
		//		req = &govtypes.QueryDepositsRequest{
		//			ProposalId: proposal.ProposalId,
		//		}
		//
		//		expRes = &govtypes.QueryDepositsResponse{
		//			Deposits: deposits,
		//		}
		//	},
		//	true,
		//},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			deposits, err := queryClient.Deposits(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes.GetDeposits(), deposits.GetDeposits())
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(deposits)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryTally() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	addrs, _ := createValidators(suite.T(), testapp.Validators...)

	var (
		req      *govtypes.QueryTallyResultRequest
		expRes   *govtypes.QueryTallyResultResponse
		proposal govtypes.Proposal
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &govtypes.QueryTallyResultRequest{}
			},
			false,
		},
		{
			"zero proposal id request",
			func() {
				req = &govtypes.QueryTallyResultRequest{ProposalId: 0}
			},
			false,
		},
		{
			"query non existed proposal",
			func() {
				req = &govtypes.QueryTallyResultRequest{ProposalId: 1}
			},
			false,
		},
		{
			"create a proposal and get tally",
			func() {
				var err error
				proposal, err = app.GovKeeper.SubmitProposal(ctx, TestProposal)
				suite.Require().NoError(err)
				suite.Require().NotNil(proposal)

				req = &govtypes.QueryTallyResultRequest{ProposalId: proposal.ProposalId}

				expRes = &govtypes.QueryTallyResultResponse{
					Tally: govtypes.EmptyTallyResult(),
				}
			},
			true,
		},
		{
			"request tally after few votes",
			func() {
				proposal.Status = govtypes.StatusVotingPeriod
				app.GovKeeper.SetProposal(ctx, proposal)

				suite.Require().NoError(app.GovKeeper.AddVote(ctx, proposal.ProposalId, addrs[0], govtypes.OptionYes))
				suite.Require().NoError(app.GovKeeper.AddVote(ctx, proposal.ProposalId, addrs[1], govtypes.OptionYes))
				suite.Require().NoError(app.GovKeeper.AddVote(ctx, proposal.ProposalId, addrs[2], govtypes.OptionYes))

				req = &govtypes.QueryTallyResultRequest{ProposalId: proposal.ProposalId}

				expRes = &govtypes.QueryTallyResultResponse{
					Tally: govtypes.TallyResult{
						Yes: sdk.NewInt(2*10000000000 + 100000000000),
					},
				}
			},
			true,
		},
		{
			"request final tally after status changed",
			func() {
				proposal.Status = govtypes.StatusPassed
				app.GovKeeper.SetProposal(ctx, proposal)
				proposal, _ = app.GovKeeper.GetProposal(ctx, proposal.ProposalId)

				req = &govtypes.QueryTallyResultRequest{ProposalId: proposal.ProposalId}

				expRes = &govtypes.QueryTallyResultResponse{
					Tally: proposal.FinalTallyResult,
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			tally, err := queryClient.TallyResult(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes.String(), tally.String())
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(tally)
			}
		})
	}
}
