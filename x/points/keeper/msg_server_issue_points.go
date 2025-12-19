package keeper

import (
	"context"

	"scontract/x/points/types"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) IssuePoints(ctx context.Context, msg *types.MsgIssuePoints) (*types.MsgIssuePointsResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// 1. PointBalance 가져오기 (없으면 생성)
	balance, err := k.PointBalance.Get(ctx, msg.Recipient)
	if err != nil && !errorsmod.IsOf(err, collections.ErrNotFound) {
		return nil, err
	}
	
	// 2. 잔액 증가
	newBalance := balance.Balance + msg.Amount
	
	// 3. 저장
	newPointBalance := types.PointBalance{
		Address: msg.Recipient,
		Balance: newBalance,
	}
	if err := k.PointBalance.Set(ctx, msg.Recipient, newPointBalance); err != nil {
		return nil, err
	}

	// 4. 거래 기록 (Transaction) 추가
	id, err := k.TransactionSeq.Next(ctx)
	if err != nil {
		return nil, err
	}
	
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	tx := types.Transaction{
		Id:        id,
		Sender:    msg.Creator,
		Recipient: msg.Recipient,
		Amount:    msg.Amount,
		TxType:    "issue",
		Timestamp: sdkCtx.BlockTime().Unix(),
	}
	if err := k.Transaction.Set(ctx, id, tx); err != nil {
		return nil, err
	}

	return &types.MsgIssuePointsResponse{}, nil
}
