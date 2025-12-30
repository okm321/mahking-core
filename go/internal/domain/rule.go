package domain

import (
	"github.com/guregu/null/v6"
	pkgerror "github.com/okm321/mahking-go/pkg/error"
)

type Rule struct {
	ID                    int64
	GroupID               int64
	MahjongType           MahjongType           // 三麻 or 四麻
	InitialPoints         int                   // 持ち点（単位: 1,000）
	ReturnPoints          int                   // 返し点（単位: 1,000）
	RankingPointsFirst    int                   // 一位のウマ
	RankingPointsSecond   int                   // 二位のウマ
	RankingPointsThird    int                   // 三位のウマ
	RankingPointsFour     null.Int              // 四位のウマ
	FractionalCalculation FractionalCalculation // 端数計算方法
	UseBust               bool                  // 飛び設定
	BustPoint             null.Int              // 飛び賞のポイント
	UseChip               bool                  // チップ設定
	ChipPoint             null.Int              // チップのポイント
}

type NewRuleArgs struct {
	MahjongType           MahjongType
	InitialPoints         int
	ReturnPoints          int
	RankingPointsFirst    int
	RankingPointsSecond   int
	RankingPointsThird    int
	RankingPointsFour     null.Int
	FractionalCalculation FractionalCalculation
	UseBust               bool
	BustPoint             null.Int
	UseChip               bool
	ChipPoint             null.Int
}

const (
	MinInitialPoints = 1
	MinReturnPoints  = 1
)

func NewRule(groupID int64, args NewRuleArgs) (_ *Rule, err error) {
	r := &Rule{
		GroupID:               groupID,
		MahjongType:           args.MahjongType,
		InitialPoints:         args.InitialPoints,
		ReturnPoints:          args.ReturnPoints,
		RankingPointsFirst:    args.RankingPointsFirst,
		RankingPointsSecond:   args.RankingPointsSecond,
		RankingPointsThird:    args.RankingPointsThird,
		RankingPointsFour:     args.RankingPointsFour,
		FractionalCalculation: args.FractionalCalculation,
		UseBust:               args.UseBust,
		BustPoint:             args.BustPoint,
		UseChip:               args.UseChip,
		ChipPoint:             args.ChipPoint,
	}

	if err := r.Validate(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Rule) Validate() error {
	if r.InitialPoints < MinInitialPoints {
		return pkgerror.NewErrorf("持ち点は%d以上である必要があります: %d", MinInitialPoints, r.InitialPoints)
	}

	if r.ReturnPoints < MinReturnPoints {
		return pkgerror.NewErrorf("返し点は%d以上である必要があります: %d", MinReturnPoints, r.ReturnPoints)
	}

	if !r.MahjongType.IsValid() {
		return pkgerror.NewErrorf("無効な麻雀タイプです: %d", r.MahjongType)
	}

	switch r.MahjongType {
	case MahjongTypeThree:
		if r.RankingPointsFirst+r.RankingPointsSecond+r.RankingPointsThird != 0 {
			return pkgerror.NewErrorf(
				"ウマの合計は0である必要があります: 1位: %d, 2位: %d, 3位: %d, 合計: %d",
				r.RankingPointsFirst,
				r.RankingPointsSecond,
				r.RankingPointsThird,
				r.RankingPointsFirst+r.RankingPointsSecond+r.RankingPointsThird,
			)
		}
	case MahjongTypeFour:
		if !r.RankingPointsFour.Valid {
			return pkgerror.NewErrorf("四位のウマは必須です")
		}

		if r.RankingPointsFirst+r.RankingPointsSecond+r.RankingPointsThird+int(r.RankingPointsFour.Int64) != 0 {
			return pkgerror.NewErrorf(
				"ウマの合計は0である必要があります: 1位: %d, 2位: %d, 3位: %d, 4位: %d, 合計: %d",
				r.RankingPointsFirst,
				r.RankingPointsSecond,
				r.RankingPointsThird,
				r.RankingPointsFour.Int64,
				r.RankingPointsFirst+r.RankingPointsSecond+r.RankingPointsThird+int(r.RankingPointsFour.Int64),
			)
		}
	default:
		return pkgerror.NewErrorf("無効な麻雀タイプです: %d", r.MahjongType)
	}

	if !r.FractionalCalculation.IsValid() {
		return pkgerror.NewErrorf("無効な端数計算方法です: %d", r.FractionalCalculation)
	}

	if r.UseBust && !r.BustPoint.Valid {
		return pkgerror.NewErrorf("飛び賞のポイントは必須です")
	}

	if r.UseChip && !r.ChipPoint.Valid {
		return pkgerror.NewErrorf("チップのポイントは必須です")
	}

	return nil
}
