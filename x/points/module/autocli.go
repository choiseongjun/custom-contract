package points

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"scontract/x/points/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ListPointBalance",
					Use:       "list-point-balance",
					Short:     "List all point-balance",
				},
				{
					RpcMethod:      "GetPointBalance",
					Use:            "get-point-balance [id]",
					Short:          "Gets a point-balance",
					Alias:          []string{"show-point-balance"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
				{
					RpcMethod: "ListTransaction",
					Use:       "list-transaction",
					Short:     "List all transaction",
				},
				{
					RpcMethod:      "GetTransaction",
					Use:            "get-transaction [id]",
					Short:          "Gets a transaction by id",
					Alias:          []string{"show-transaction"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod: "ListSettlement",
					Use:       "list-settlement",
					Short:     "List all settlement",
				},
				{
					RpcMethod:      "GetSettlement",
					Use:            "get-settlement [id]",
					Short:          "Gets a settlement by id",
					Alias:          []string{"show-settlement"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "IssuePoints",
					Use:            "issue-points [recipient] [amount] [reason]",
					Short:          "Send a issue-points tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "recipient"}, {ProtoField: "amount"}, {ProtoField: "reason"}},
				},
				{
					RpcMethod:      "SpendPoints",
					Use:            "spend-points [amount] [description]",
					Short:          "Send a spend-points tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}, {ProtoField: "description"}},
				},
				{
					RpcMethod:      "TransferPoints",
					Use:            "transfer-points [recipient] [amount]",
					Short:          "Send a transfer-points tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "recipient"}, {ProtoField: "amount"}},
				},
				{
					RpcMethod:      "RequestSettlement",
					Use:            "request-settlement [amount]",
					Short:          "Send a request-settlement tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
