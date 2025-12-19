package types

func NewMsgIssuePoints(creator string, recipient string, amount uint64, reason string) *MsgIssuePoints {
	return &MsgIssuePoints{
		Creator:   creator,
		Recipient: recipient,
		Amount:    amount,
		Reason:    reason,
	}
}
