package keeper

import (
	"context"

	"scontract/x/points/types"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RequestSettlement(ctx context.Context, msg *types.MsgRequestSettlement) (*types.MsgRequestSettlementResponse, error) {
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

	// 2. 잔액 확인
	if balance.Balance < msg.Amount {
		return nil, errorsmod.Wrapf(types.ErrInsufficientFunds, "balance is %d but needed %d", balance.Balance, msg.Amount)
	}

	// 3. 잔액 차감 (정산 요청 시 포인트 차감)
	balance.Balance -= msg.Amount
	if err := k.PointBalance.Set(ctx, msg.Creator, balance); err != nil {
		return nil, err
	}

	// 4. 정산 ID 생성
	id, err := k.SettlementSeq.Next(ctx)
	if err != nil {
		return nil, err
	}
	
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// 5. 정산 기록 저장
	settlement := types.Settlement{
		Id:        id,
		Requester: msg.Creator,
		Amount:    msg.Amount,
		Status:    "pending",
		Timestamp: sdkCtx.BlockTime().Unix(),
	}
	
	if err := k.Settlement.Set(ctx, id, settlement); err != nil {
		return nil, err
	}

	return &types.MsgRequestSettlementResponse{}, nil
}
