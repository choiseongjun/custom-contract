package types_test

import (
	"testing"

	"scontract/x/points/types"

	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc:     "valid genesis state",
			genState: &types.GenesisState{PointBalanceMap: []types.PointBalance{{Index: "0"}, {Index: "1"}}, TransactionList: []types.Transaction{{Id: 0}, {Id: 1}}, TransactionCount: 2, SettlementList: []types.Settlement{{Id: 0}, {Id: 1}}, SettlementCount: 2}, valid: true,
		}, {
			desc: "duplicated pointBalance",
			genState: &types.GenesisState{
				PointBalanceMap: []types.PointBalance{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
				TransactionList: []types.Transaction{{Id: 0}, {Id: 1}}, TransactionCount: 2,
				SettlementList: []types.Settlement{{Id: 0}, {Id: 1}}, SettlementCount: 2}, valid: false,
		}, {
			desc: "duplicated transaction",
			genState: &types.GenesisState{
				TransactionList: []types.Transaction{
					{
						Id: 0,
					},
					{
						Id: 0,
					},
				},
				SettlementList: []types.Settlement{{Id: 0}, {Id: 1}}, SettlementCount: 2,
			}, valid: false,
		}, {
			desc: "invalid transaction count",
			genState: &types.GenesisState{
				TransactionList: []types.Transaction{
					{
						Id: 1,
					},
				},
				TransactionCount: 0,
				SettlementList:   []types.Settlement{{Id: 0}, {Id: 1}}, SettlementCount: 2,
			}, valid: false,
		}, {
			desc: "duplicated settlement",
			genState: &types.GenesisState{
				SettlementList: []types.Settlement{
					{
						Id: 0,
					},
					{
						Id: 0,
					},
				},
			},
			valid: false,
		}, {
			desc: "invalid settlement count",
			genState: &types.GenesisState{
				SettlementList: []types.Settlement{
					{
						Id: 1,
					},
				},
				SettlementCount: 0,
			},
			valid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
