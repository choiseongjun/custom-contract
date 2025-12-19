package keeper

import (
	"context"

	"scontract/x/points/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) SpendPoints(ctx context.Context, msg *types.MsgSpendPoints) (*types.MsgSpendPointsResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgSpendPointsResponse{}, nil
}
