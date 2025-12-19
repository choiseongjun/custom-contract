package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"scontract/x/points/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema         collections.Schema
	Params         collections.Item[types.Params]
	PointBalance   collections.Map[string, types.PointBalance]
	TransactionSeq collections.Sequence
	Transaction    collections.Map[uint64, types.Transaction]
	SettlementSeq  collections.Sequence
	Settlement     collections.Map[uint64, types.Settlement]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,

) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		Params:       collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		PointBalance: collections.NewMap(sb, types.PointBalanceKey, "pointBalance", collections.StringKey, codec.CollValue[types.PointBalance](cdc)), Transaction: collections.NewMap(sb, types.TransactionKey, "transaction", collections.Uint64Key, codec.CollValue[types.Transaction](cdc)),
		TransactionSeq: collections.NewSequence(sb, types.TransactionCountKey, "transactionSequence"),
		Settlement:     collections.NewMap(sb, types.SettlementKey, "settlement", collections.Uint64Key, codec.CollValue[types.Settlement](cdc)),
		SettlementSeq:  collections.NewSequence(sb, types.SettlementCountKey, "settlementSequence"),
	}
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}
