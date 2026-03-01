package in

import (
	"github.com/guregu/null/v6"
	"github.com/okm321/mahking-go/internal/domain"
)

type CreateGroupWithRule struct {
	//govalid:maxlength=100
	Name string `json:"name"` // グループ名
	//govalid:required
	MemberNames []string `json:"member_names"` // メンバー名
	//govalid:required
	Rules rules `json:"rules"` // ルール設定
}

type rules struct {
	//govalid:required
	//govalid:enum=domain.MahjongTypeThree,domain.MahjongTypeFour
	MahjongType domain.MahjongType `json:"mahjong_type"` // 三麻 or 四麻
	//govalid:required
	//govalid:gte=1
	InitialPoints int `json:"initial_points"` // 持ち点（単位: 1,000）
	//govalid:required
	//govalid:gte=1
	ReturnPoints int `json:"return_points"` // 返し点（単位: 1,000）
	//govalid:required
	RankingPointsFirst int `json:"ranking_points_first"` // 一位のウマ
	//govalid:required
	RankingPointsSecond int `json:"ranking_points_second"` // 二位のウマ
	//govalid:required
	RankingPointsThird int `json:"ranking_points_third"` // 三位のウマ
	//govalid:cel=this.MahjongType != 2 || value.Valid
	RankingPointsFour null.Int `json:"ranking_points_four"` // 四位のウマ
	//govalid:required
	FractionalCalculation domain.FractionalCalculation `json:"fractional_calculation"` // 1: 切り上げ, 2: 切り捨て, 3: 四捨五入, 4: 10点未満切り上げ, 5: 10点未満切り捨て
	//govalid:required
	UseBust   bool     `json:"use_bust"`   // 飛び設定
	BustPoint null.Int `json:"bust_point"` // 飛び賞のポイント
	//govalid:required
	UseChip   bool     `json:"use_chip"`   // チップ設定
	ChipPoint null.Int `json:"chip_point"` // チップのポイント
}
