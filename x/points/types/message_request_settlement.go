package types

func NewMsgRequestSettlement(creator string, amount uint64) *MsgRequestSettlement {
	return &MsgRequestSettlement{
		Creator: creator,
		Amount:  amount,
	}
}
