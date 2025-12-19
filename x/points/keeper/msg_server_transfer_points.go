package keeper

import (
	"context"

	"scontract/x/points/types"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) TransferPoints(ctx context.Context, msg *types.MsgTransferPoints) (*types.MsgTransferPointsResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// 1. 보내는 사람 잔액 조회
	senderBalance, err := k.PointBalance.Get(ctx, msg.Creator)
	if err != nil {
		if errorsmod.IsOf(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(types.ErrInsufficientFunds, "sender balance not found")
		}
		return nil, err
	}

	// 2. 잔액 확인
	if senderBalance.Balance < msg.Amount {
		return nil, errorsmod.Wrapf(types.ErrInsufficientFunds, "insufficient funds")
	}

	// 3. 받는 사람 잔액 조회 (없으면 0)
	recipientBalance, err := k.PointBalance.Get(ctx, msg.Recipient)
	if err != nil {
		if !errorsmod.IsOf(err, collections.ErrNotFound) {
			return nil, err
		}
		recipientBalance = types.PointBalance{Address: msg.Recipient, Balance: 0}
	}

	// 4. 잔액 이동
	senderBalance.Balance -= msg.Amount
	recipientBalance.Balance += msg.Amount

	// 5. 저장
	if err := k.PointBalance.Set(ctx, msg.Creator, senderBalance); err != nil {
		return nil, err
	}
	if err := k.PointBalance.Set(ctx, msg.Recipient, recipientBalance); err != nil {
		return nil, err
	}

	// 6. 거래 기록
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
		TxType:    "transfer",
		Timestamp: sdkCtx.BlockTime().Unix(),
	}
	if err := k.Transaction.Set(ctx, id, tx); err != nil {
		return nil, err
	}

	return &types.MsgTransferPointsResponse{}, nil
}
