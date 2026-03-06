package domain

type Seat int

const (
	SeatEast  Seat = iota + 1 // 東
	SeatSouth                 // 南
	SeatWest                  // 西
	SeatNorth                 // 北
)

func (s Seat) String() string {
	switch s {
	case SeatEast:
		return "東"
	case SeatSouth:
		return "南"
	case SeatWest:
		return "西"
	case SeatNorth:
		return "北"
	default:
		return "その他"
	}
}

func (s Seat) IsValid() bool {
	switch s {
	case SeatEast, SeatSouth, SeatWest, SeatNorth:
		return true
	default:
		return false
	}
}
