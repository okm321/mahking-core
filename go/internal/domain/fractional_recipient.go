package domain

type FractionalRecipient int

const (
	FractionalRecipientFirst FractionalRecipient = iota + 1 // 1位の人
	FractionalRecipientLast                                 // 最下位の人
)

func (f FractionalRecipient) String() string {
	switch f {
	case FractionalRecipientFirst:
		return "1位の人"
	case FractionalRecipientLast:
		return "最下位の人"
	default:
		return "その他"
	}
}

// TargetRanking は端数を受け取る人の順位を返す
func (f FractionalRecipient) TargetRanking(playerCount int) int {
	switch f {
	case FractionalRecipientFirst:
		return 1
	case FractionalRecipientLast:
		return playerCount
	default:
		return 1
	}
}

func (f FractionalRecipient) IsValid() bool {
	switch f {
	case FractionalRecipientFirst, FractionalRecipientLast:
		return true
	default:
		return false
	}
}
