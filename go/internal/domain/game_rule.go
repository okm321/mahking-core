package domain

import "github.com/guregu/null/v6"

type GameRule struct {
	ID      int64
	GameID  int64
	GroupID int64
	//govalid:required
	//govalid:enum=MahjongTypeThree,MahjongTypeFour
	MahjongType MahjongType // 三麻 or 四麻
	//govalid:required
	//govalid:gte=1
	InitialPoints int // 持ち点（単位: 1,000）
	//govalid:required
	//govalid:gte=1
	ReturnPoints        int      // 返し点（単位: 1,000）
	RankingPointsFirst  int      // 一位のウマ
	RankingPointsSecond int      // 二位のウマ
	RankingPointsThird  int      // 三位のウマ
	RankingPointsFour   null.Int // 四位のウマ
	//govalid:required
	//govalid:enum=FractionalCalculationDecimal,FractionalCalculationRoundDown,FractionalCalculationRoundUp,FractionalCalculationRoundNearest,FractionalCalculationRoundFive
	FractionalCalculation FractionalCalculation // 端数計算方法
	FractionalRecipient   FractionalRecipient   // 端数を受け取る人
	UseBust               bool                  // 飛び設定
	BustPoint             null.Int              // 飛び賞のポイント
	UseChip               bool                  // チップ設定
	ChipPoint             null.Int              // チップのポイント
}

func (gr *GameRule) rankingPoint(ranking int) int {
	switch ranking {
	case 1:
		return gr.RankingPointsFirst
	case 2:
		return gr.RankingPointsSecond
	case 3:
		return gr.RankingPointsThird
	case 4:
		return int(gr.RankingPointsFour.Int64)
	default:
		return 0
	}
}

// NewGameRuleFromRule バリデーション済みの Rule からスナップショットを生成する
func NewGameRuleFromRule(rule *Rule) *GameRule {
	return &GameRule{
		GroupID:               rule.GroupID,
		MahjongType:           rule.MahjongType,
		InitialPoints:         rule.InitialPoints,
		ReturnPoints:          rule.ReturnPoints,
		RankingPointsFirst:    rule.RankingPointsFirst,
		RankingPointsSecond:   rule.RankingPointsSecond,
		RankingPointsThird:    rule.RankingPointsThird,
		RankingPointsFour:     rule.RankingPointsFour,
		FractionalCalculation: rule.FractionalCalculation,
		FractionalRecipient:   rule.FractionalRecipient,
		UseBust:               rule.UseBust,
		BustPoint:             rule.BustPoint,
		UseChip:               rule.UseChip,
		ChipPoint:             rule.ChipPoint,
	}
}
