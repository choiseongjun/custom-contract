package keeper

import (
	"context"

	"scontract/x/points/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) TransferPoints(ctx context.Context, msg *types.MsgTransferPoints) (*types.MsgTransferPointsResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgTransferPointsResponse{}, nil
}
