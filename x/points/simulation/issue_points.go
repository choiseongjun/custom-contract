package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"scontract/x/points/keeper"
	"scontract/x/points/types"
)

func SimulateMsgIssuePoints(
	ak types.AuthKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	txGen client.TxConfig,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgIssuePoints{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handle the IssuePoints simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "IssuePoints simulation not implemented"), nil, nil
	}
}
