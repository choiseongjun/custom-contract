package keeper_test

import (
	"context"
	"strconv"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"scontract/x/points/keeper"
	"scontract/x/points/types"
)

func createNSettlement(keeper keeper.Keeper, ctx context.Context, n int) []types.Settlement {
	items := make([]types.Settlement, n)
	for i := range items {
		iu := uint64(i)
		items[i].Id = iu
		items[i].Requester = strconv.Itoa(i)
		items[i].Amount = uint64(i)
		items[i].Status = strconv.Itoa(i)
		items[i].Timestamp = int64(i)
		_ = keeper.Settlement.Set(ctx, iu, items[i])
		_ = keeper.SettlementSeq.Set(ctx, iu)
	}
	return items
}

func TestSettlementQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNSettlement(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetSettlementRequest
		response *types.QueryGetSettlementResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetSettlementRequest{Id: msgs[0].Id},
			response: &types.QueryGetSettlementResponse{Settlement: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetSettlementRequest{Id: msgs[1].Id},
			response: &types.QueryGetSettlementResponse{Settlement: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetSettlementRequest{Id: uint64(len(msgs))},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetSettlement(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestSettlementQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNSettlement(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllSettlementRequest {
		return &types.QueryAllSettlementRequest{
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
			resp, err := qs.ListSettlement(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Settlement), step)
			require.Subset(t, msgs, resp.Settlement)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListSettlement(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Settlement), step)
			require.Subset(t, msgs, resp.Settlement)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListSettlement(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.Settlement)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListSettlement(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
