package points

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	pointssimulation "scontract/x/points/simulation"
	"scontract/x/points/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	pointsGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&pointsGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgIssuePoints          = "op_weight_msg_points"
		defaultWeightMsgIssuePoints int = 100
	)

	var weightMsgIssuePoints int
	simState.AppParams.GetOrGenerate(opWeightMsgIssuePoints, &weightMsgIssuePoints, nil,
		func(_ *rand.Rand) {
			weightMsgIssuePoints = defaultWeightMsgIssuePoints
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgIssuePoints,
		pointssimulation.SimulateMsgIssuePoints(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgSpendPoints          = "op_weight_msg_points"
		defaultWeightMsgSpendPoints int = 100
	)

	var weightMsgSpendPoints int
	simState.AppParams.GetOrGenerate(opWeightMsgSpendPoints, &weightMsgSpendPoints, nil,
		func(_ *rand.Rand) {
			weightMsgSpendPoints = defaultWeightMsgSpendPoints
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSpendPoints,
		pointssimulation.SimulateMsgSpendPoints(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgTransferPoints          = "op_weight_msg_points"
		defaultWeightMsgTransferPoints int = 100
	)

	var weightMsgTransferPoints int
	simState.AppParams.GetOrGenerate(opWeightMsgTransferPoints, &weightMsgTransferPoints, nil,
		func(_ *rand.Rand) {
			weightMsgTransferPoints = defaultWeightMsgTransferPoints
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgTransferPoints,
		pointssimulation.SimulateMsgTransferPoints(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgRequestSettlement          = "op_weight_msg_points"
		defaultWeightMsgRequestSettlement int = 100
	)

	var weightMsgRequestSettlement int
	simState.AppParams.GetOrGenerate(opWeightMsgRequestSettlement, &weightMsgRequestSettlement, nil,
		func(_ *rand.Rand) {
			weightMsgRequestSettlement = defaultWeightMsgRequestSettlement
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRequestSettlement,
		pointssimulation.SimulateMsgRequestSettlement(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
