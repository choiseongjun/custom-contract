package keeper

import (
	"context"
	"errors"

	"scontract/x/points/types"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListPointBalance(ctx context.Context, req *types.QueryAllPointBalanceRequest) (*types.QueryAllPointBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	pointBalances, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.PointBalance,
		req.Pagination,
		func(_ string, value types.PointBalance) (types.PointBalance, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPointBalanceResponse{PointBalance: pointBalances, Pagination: pageRes}, nil
}

func (q queryServer) GetPointBalance(ctx context.Context, req *types.QueryGetPointBalanceRequest) (*types.QueryGetPointBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.PointBalance.Get(ctx, req.Index)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetPointBalanceResponse{PointBalance: val}, nil
}
