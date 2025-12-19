package app

import (
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
)

func (app *App) registerWasmModule(appOpts types.AppOptions) error {
	wasmDir := filepath.Join(DefaultNodeHome, "wasm")

	wasmConfig := wasmtypes.NodeConfig{
		ContractDebugMode:  false,
		SmartQueryGasLimit: 3000000,
		MemoryCacheSize:    100,
	}

	govModuleAddr := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	app.WasmKeeper = wasmkeeper.NewKeeper(
		app.appCodec,
		runtime.NewKVStoreService(app.WasmKey),
		app.AuthKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		distrkeeper.NewQuerier(app.DistrKeeper),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeperV2,
		app.TransferKeeper,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		wasmtypes.VMConfig{},
		wasmkeeper.BuiltInCapabilities(),
		govModuleAddr,
	)

	// Register MsgServer
	wasmtypes.RegisterMsgServer(app.MsgServiceRouter(), wasmkeeper.NewMsgServerImpl(&app.WasmKeeper))
	// Register QueryServer
	wasmtypes.RegisterQueryServer(app.GRPCQueryRouter(), wasmkeeper.Querier(&app.WasmKeeper))

	return nil
}

// RegisterWasm registers the wasm module
func RegisterWasm(cdc codec.Codec) map[string]module.AppModuleBasic {
	return map[string]module.AppModuleBasic{
		wasm.ModuleName: wasm.AppModuleBasic{},
	}
}
