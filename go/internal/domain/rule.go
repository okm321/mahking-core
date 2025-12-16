package domain

type Rule struct {
	ID                    int64
	GroupID               int64
	MahjongType           MahjongType // 三麻 or 四麻
	InitialPoints         int         // 持ち点（単位: 1,000）
	ReturnPoints          int         // 返し点（単位: 1,000）
	RankingPointsFirst    int         // 一位のウマ
	RankingPointsSecond   int         // 二位のウマ
	RankingPointsThird    int         // 三位のウマ
	RankingPointsFour     int         // 四位のウマ
	FractionalCalculation int         // 1: 切り上げ, 2: 切り捨て, 3: 四捨五入, 4: 10点未満切り上げ, 5: 10点未満切り捨て
	UseBust               bool        // 飛び設定
	BustPoint             int         // 飛び賞のポイント
	UseChip               bool        // チップ設定
	ChipPoint             int         // チップのポイント
}

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
