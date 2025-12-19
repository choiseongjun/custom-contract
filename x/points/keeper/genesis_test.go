package keeper_test

import (
	"testing"

	"scontract/x/points/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params:          types.DefaultParams(),
		PointBalanceMap: []types.PointBalance{{Index: "0"}, {Index: "1"}}, TransactionList: []types.Transaction{{Id: 0}, {Id: 1}},
		TransactionCount: 2,
		SettlementList:   []types.Settlement{{Id: 0}, {Id: 1}},
		SettlementCount:  2,
	}
	f := initFixture(t)
	err := f.keeper.InitGenesis(f.ctx, genesisState)
	require.NoError(t, err)
	got, err := f.keeper.ExportGenesis(f.ctx)
	require.NoError(t, err)
	require.NotNil(t, got)

	require.EqualExportedValues(t, genesisState.Params, got.Params)
	require.EqualExportedValues(t, genesisState.PointBalanceMap, got.PointBalanceMap)
	require.EqualExportedValues(t, genesisState.TransactionList, got.TransactionList)
	require.Equal(t, genesisState.TransactionCount, got.TransactionCount)
	require.EqualExportedValues(t, genesisState.SettlementList, got.SettlementList)
	require.Equal(t, genesisState.SettlementCount, got.SettlementCount)

}
