package keeper

import (
	"context"

	"scontract/x/points/types"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SpendPoints(ctx context.Context, msg *types.MsgSpendPoints) (*types.MsgSpendPointsResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// 1. 잔액 조회
	balance, err := k.PointBalance.Get(ctx, msg.Creator)
	if err != nil {
		if errorsmod.IsOf(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(types.ErrInsufficientFunds, "balance not found")
		}
		return nil, err
	}

	// 2. 잔액 확인 (부족하면 에러)
	if balance.Balance < msg.Amount {
		return nil, errorsmod.Wrapf(types.ErrInsufficientFunds, "balance is %d but needed %d", balance.Balance, msg.Amount)
	}

	// 3. 잔액 차감
	balance.Balance -= msg.Amount
	if err := k.PointBalance.Set(ctx, msg.Creator, balance); err != nil {
		return nil, err
	}

	// 4. 거래 기록 추가
	id, err := k.TransactionSeq.Next(ctx)
	if err != nil {
		return nil, err
	}
	
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	tx := types.Transaction{
		Id:        id,
		Sender:    msg.Creator,
		Recipient: "MERCHANT", // or BURN address
		Amount:    msg.Amount,
		TxType:    "spend",
		Timestamp: sdkCtx.BlockTime().Unix(),
	}
	if err := k.Transaction.Set(ctx, id, tx); err != nil {
		return nil, err
	}

	return &types.MsgSpendPointsResponse{}, nil
}
