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

func (q queryServer) ListTransaction(ctx context.Context, req *types.QueryAllTransactionRequest) (*types.QueryAllTransactionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	transactions, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Transaction,
		req.Pagination,
		func(_ uint64, value types.Transaction) (types.Transaction, error) {
			return value, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTransactionResponse{Transaction: transactions, Pagination: pageRes}, nil
}

func (q queryServer) GetTransaction(ctx context.Context, req *types.QueryGetTransactionRequest) (*types.QueryGetTransactionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	transaction, err := q.k.Transaction.Get(ctx, req.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, sdkerrors.ErrKeyNotFound
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetTransactionResponse{Transaction: transaction}, nil
}
