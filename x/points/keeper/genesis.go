package keeper

import (
	"context"

	"scontract/x/points/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	for _, elem := range genState.PointBalanceMap {
		if err := k.PointBalance.Set(ctx, elem.Index, elem); err != nil {
			return err
		}
	}
	for _, elem := range genState.TransactionList {
		if err := k.Transaction.Set(ctx, elem.Id, elem); err != nil {
			return err
		}
	}

	if err := k.TransactionSeq.Set(ctx, genState.TransactionCount); err != nil {
		return err
	}
	for _, elem := range genState.SettlementList {
		if err := k.Settlement.Set(ctx, elem.Id, elem); err != nil {
			return err
		}
	}

	if err := k.SettlementSeq.Set(ctx, genState.SettlementCount); err != nil {
		return err
	}

	return k.Params.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	if err := k.PointBalance.Walk(ctx, nil, func(_ string, val types.PointBalance) (stop bool, err error) {
		genesis.PointBalanceMap = append(genesis.PointBalanceMap, val)
		return false, nil
	}); err != nil {
		return nil, err
	}
	err = k.Transaction.Walk(ctx, nil, func(key uint64, elem types.Transaction) (bool, error) {
		genesis.TransactionList = append(genesis.TransactionList, elem)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	genesis.TransactionCount, err = k.TransactionSeq.Peek(ctx)
	if err != nil {
		return nil, err
	}
	err = k.Settlement.Walk(ctx, nil, func(key uint64, elem types.Settlement) (bool, error) {
		genesis.SettlementList = append(genesis.SettlementList, elem)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	genesis.SettlementCount, err = k.SettlementSeq.Peek(ctx)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}
