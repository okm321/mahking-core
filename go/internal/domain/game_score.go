package domain

import "github.com/guregu/null/v6"

type GameScore struct {
	ID       int64
	GameID   int64
	GroupID  int64
	MemberID int64
	//govalid:required
	//govalid:enum=SeatEast,SeatSouth,SeatWest,SeatNorth
	Seat Seat // 席
	//govalid:required
	//govalid:gte=1
	//govalid:lte=4
	Ranking  int     // 順位
	RawScore int     // 素点（100点単位, 例: 324 = 32,400点）
	Point    float64 // 計算後のポイント
	ChipCount null.Int // チップ枚数
	IsBusted  bool     // 飛びフラグ
}
