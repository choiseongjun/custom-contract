package types

func NewMsgSpendPoints(creator string, amount uint64, description string) *MsgSpendPoints {
	return &MsgSpendPoints{
		Creator:     creator,
		Amount:      amount,
		Description: description,
	}
}
