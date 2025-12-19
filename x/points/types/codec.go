package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRequestSettlement{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgTransferPoints{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSpendPoints{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgIssuePoints{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
