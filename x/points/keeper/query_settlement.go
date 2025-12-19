package keeper

import (
	"context"
	"errors"

	"scontract/x/points/types"

	"cosmossdk.io/collections"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListSettlement(ctx context.Context, req *types.QueryAllSettlementRequest) (*types.QueryAllSettlementResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	settlements, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Settlement,
		req.Pagination,
		func(_ uint64, value types.Settlement) (types.Settlement, error) {
			return value, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllSettlementResponse{Settlement: settlements, Pagination: pageRes}, nil
}

func (q queryServer) GetSettlement(ctx context.Context, req *types.QueryGetSettlementRequest) (*types.QueryGetSettlementResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	settlement, err := q.k.Settlement.Get(ctx, req.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, sdkerrors.ErrKeyNotFound
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetSettlementResponse{Settlement: settlement}, nil
}
