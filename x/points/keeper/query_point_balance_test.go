package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"scontract/x/points/keeper"
	"scontract/x/points/types"
)

func createNPointBalance(keeper keeper.Keeper, ctx context.Context, n int) []types.PointBalance {
	items := make([]types.PointBalance, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)
		items[i].Address = strconv.Itoa(i)
		items[i].Balance = uint64(i)
		_ = keeper.PointBalance.Set(ctx, items[i].Index, items[i])
	}
	return items
}

func TestPointBalanceQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNPointBalance(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetPointBalanceRequest
		response *types.QueryGetPointBalanceResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetPointBalanceRequest{
				Index: msgs[0].Index,
			},
			response: &types.QueryGetPointBalanceResponse{PointBalance: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetPointBalanceRequest{
				Index: msgs[1].Index,
			},
			response: &types.QueryGetPointBalanceResponse{PointBalance: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetPointBalanceRequest{
				Index: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetPointBalance(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestPointBalanceQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNPointBalance(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllPointBalanceRequest {
		return &types.QueryAllPointBalanceRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListPointBalance(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.PointBalance), step)
			require.Subset(t, msgs, resp.PointBalance)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListPointBalance(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.PointBalance), step)
			require.Subset(t, msgs, resp.PointBalance)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListPointBalance(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.PointBalance)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListPointBalance(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
