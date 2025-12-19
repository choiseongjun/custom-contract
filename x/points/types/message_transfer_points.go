package types

func NewMsgTransferPoints(creator string, recipient string, amount uint64) *MsgTransferPoints {
	return &MsgTransferPoints{
		Creator:   creator,
		Recipient: recipient,
		Amount:    amount,
	}
}
