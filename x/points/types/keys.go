package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "points"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	// It should be synced with the gov module's name if it is ever changed.
	// See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"
)

// ParamsKey is the prefix to retrieve all Params
var ParamsKey = collections.NewPrefix("p_points")

var (
	TransactionKey      = collections.NewPrefix("transaction/value/")
	TransactionCountKey = collections.NewPrefix("transaction/count/")
)

var (
	SettlementKey      = collections.NewPrefix("settlement/value/")
	SettlementCountKey = collections.NewPrefix("settlement/count/")
)
