package domain

import (
	"testing"

	"github.com/guregu/null/v6"
	"github.com/stretchr/testify/assert"
)

func testScore(ranking int, rawScore int) *GameScore {
	return &GameScore{
		Ranking:  ranking,
		RawScore: rawScore,
	}
}

func TestGame_calculatePoints(t *testing.T) {
	// 四麻共通: 持ち25,000 / 返し30,000 / ウマ +20,+10,-10,-20
	fourPlayerRule := func(fc FractionalCalculation, fr FractionalRecipient) *GameRule {
		return &GameRule{
			MahjongType:           MahjongTypeFour,
			InitialPoints:         25,
			ReturnPoints:          30,
			RankingPointsFirst:    20,
			RankingPointsSecond:   10,
			RankingPointsThird:    -10,
			RankingPointsFour:     null.IntFrom(-20),
			FractionalCalculation: fc,
			FractionalRecipient:   fr,
		}
	}

	tests := []struct {
		name           string
		game           *Game
		expectedPoints map[int]float64 // ranking -> expected point
	}{
		{
			// 端数がそのまま残る
			name: "四麻/小数点有効",
			game: &Game{
				GameRule: fourPlayerRule(FractionalCalculationDecimal, 0),
				GameScores: []*GameScore{
					testScore(1, 453), // diff=15.3 → 15.3+20+20(oka)=55.3
					testScore(2, 287), // diff=-1.3 → -1.3+10=8.7
					testScore(3, 156), // diff=-14.4 → -14.4-10=-24.4
					testScore(4, 104), // diff=-19.6 → -19.6-20=-39.6
				},
			},
			expectedPoints: map[int]float64{1: 55.3, 2: 8.7, 3: -24.4, 4: -39.6},
		},
		{
			// Trunc(15.3)=15 等。合計=1 → 1位に-1
			name: "四麻/切り捨て/端数→1位",
			game: &Game{
				GameRule: fourPlayerRule(FractionalCalculationRoundDown, FractionalRecipientFirst),
				GameScores: []*GameScore{
					testScore(1, 453),
					testScore(2, 287),
					testScore(3, 156),
					testScore(4, 104),
				},
			},
			expectedPoints: map[int]float64{1: 54, 2: 9, 3: -24, 4: -39},
		},
		{
			// 切り捨ての端数を最下位が受け取る
			name: "四麻/切り捨て/端数→最下位",
			game: &Game{
				GameRule: fourPlayerRule(FractionalCalculationRoundDown, FractionalRecipientLast),
				GameScores: []*GameScore{
					testScore(1, 453),
					testScore(2, 287),
					testScore(3, 156),
					testScore(4, 104),
				},
			},
			expectedPoints: map[int]float64{1: 55, 2: 9, 3: -24, 4: -40},
		},
		{
			// Ceil(15.3)=16, Ceil(1.3)→-2 等。合計=-1 → 1位に+1
			name: "四麻/切り上げ/端数→1位",
			game: &Game{
				GameRule: fourPlayerRule(FractionalCalculationRoundUp, FractionalRecipientFirst),
				GameScores: []*GameScore{
					testScore(1, 453),
					testScore(2, 287),
					testScore(3, 156),
					testScore(4, 104),
				},
			},
			expectedPoints: map[int]float64{1: 57, 2: 8, 3: -25, 4: -40},
		},
		{
			// Round(15.5)=16（0.5は0から離れる方向）。合計=-1 → 1位に+1
			name: "四麻/四捨五入/端数→1位",
			game: &Game{
				GameRule: fourPlayerRule(FractionalCalculationRoundNearest, FractionalRecipientFirst),
				GameScores: []*GameScore{
					testScore(1, 455), // diff=15.5 → Round=16
					testScore(2, 285), // diff=-1.5 → Round=-2
					testScore(3, 155), // diff=-14.5 → Round=-15
					testScore(4, 105), // diff=-19.5 → Round=-20
				},
			},
			expectedPoints: map[int]float64{1: 57, 2: 8, 3: -25, 4: -40},
		},
		{
			// 0.5は切り捨て（四捨五入との違い）。合計=1 → 1位に-1
			name: "四麻/五捨六入/端数→1位",
			game: &Game{
				GameRule: fourPlayerRule(FractionalCalculationRoundFive, FractionalRecipientFirst),
				GameScores: []*GameScore{
					testScore(1, 455), // diff=15.5 → 15
					testScore(2, 285), // diff=-1.5 → -1
					testScore(3, 155), // diff=-14.5 → -14
					testScore(4, 105), // diff=-19.5 → -19
				},
			},
			expectedPoints: map[int]float64{1: 54, 2: 9, 3: -24, 4: -39},
		},
		{
			// オカ=(25-25)*4=0
			name: "四麻/オカなし",
			game: &Game{
				GameRule: &GameRule{
					MahjongType:           MahjongTypeFour,
					InitialPoints:         25,
					ReturnPoints:          25,
					RankingPointsFirst:    20,
					RankingPointsSecond:   10,
					RankingPointsThird:    -10,
					RankingPointsFour:     null.IntFrom(-20),
					FractionalCalculation: FractionalCalculationDecimal,
				},
				GameScores: []*GameScore{
					testScore(1, 400), // diff=15 → 15+20=35
					testScore(2, 280), // diff=3 → 3+10=13
					testScore(3, 200), // diff=-5 → -5-10=-15
					testScore(4, 120), // diff=-13 → -13-20=-33
				},
			},
			expectedPoints: map[int]float64{1: 35, 2: 13, 3: -15, 4: -33},
		},
		{
			// 端数が出ない素点 → 調整不要
			name: "四麻/切り捨て/端数なし",
			game: &Game{
				GameRule: fourPlayerRule(FractionalCalculationRoundDown, FractionalRecipientFirst),
				GameScores: []*GameScore{
					testScore(1, 450), // diff=15.0
					testScore(2, 300), // diff=0.0
					testScore(3, 150), // diff=-15.0
					testScore(4, 100), // diff=-20.0
				},
			},
			expectedPoints: map[int]float64{1: 55, 2: 10, 3: -25, 4: -40},
		},
		{
			// TODO(human): 三麻のテストケースを追加してください
			// 三麻共通: 持ち35,000 / 返し40,000 / ウマ +15,0,-15
			// ヒント:
			//   - 素点の合計 = InitialPoints * 10 * 3 = 1050
			//   - diff = (RawScore - 400) / 10
			//   - オカ = (40-35) * 3 = 15
			//   - 端数が出る素点を選び、丸め+端数調整が正しく動くことを確認
			name: "三麻/切り捨て/端数→1位",
			game: &Game{
				GameRule: &GameRule{
					MahjongType:           MahjongTypeThree,
					InitialPoints:         35,
					ReturnPoints:          40,
					RankingPointsFirst:    15,
					RankingPointsSecond:   0,
					RankingPointsThird:    -15,
					FractionalCalculation: FractionalCalculationRoundDown,
					FractionalRecipient:   FractionalRecipientFirst,
				},
				GameScores: []*GameScore{
					// TODO(human): 素点と期待値を埋めてください
				},
			},
			expectedPoints: map[int]float64{
				// TODO(human): 期待値を埋めてください
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.game.GameScores) == 0 {
				t.Skip("テストデータ未実装")
			}
			tt.game.calculatePoints()
			for _, score := range tt.game.GameScores {
				expected := tt.expectedPoints[score.Ranking]
				assert.InDelta(t, expected, score.Point, 0.0001, "ranking=%d", score.Ranking)
			}
		})
	}
}
