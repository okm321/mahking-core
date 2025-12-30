package domain

type MahjongType int

const (
	MahjongTypeThree MahjongType = iota + 1
	MahjongTypeFour
)

func (m MahjongType) String() string {
	switch m {
	case MahjongTypeThree:
		return "三麻"
	case MahjongTypeFour:
		return "四麻"
	default:
		return "その他"
	}
}

func (m MahjongType) IsValid() bool {
	switch m {
	case MahjongTypeThree, MahjongTypeFour:
		return true
	default:
		return false
	}
}

func (m MahjongType) RequiredMemberCount() int {
	switch m {
	case MahjongTypeThree:
		return 3
	case MahjongTypeFour:
		return 4
	default:
		return 0
	}
}
