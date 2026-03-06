package domain

import "math"

type FractionalCalculation int

const (
	FractionalCalculationDecimal      FractionalCalculation = iota + 1 // 小数点有効
	FractionalCalculationRoundDown                                     // 切り捨て
	FractionalCalculationRoundUp                                       // 切り上げ
	FractionalCalculationRoundNearest                                  // 四捨五入
	FractionalCalculationRoundFive                                     // 五捨六入
)

func (m FractionalCalculation) String() string {
	switch m {
	case FractionalCalculationDecimal:
		return "小数点有効"
	case FractionalCalculationRoundDown:
		return "切り捨て"
	case FractionalCalculationRoundUp:
		return "切り上げ"
	case FractionalCalculationRoundNearest:
		return "四捨五入"
	case FractionalCalculationRoundFive:
		return "五捨六入"
	default:
		return "その他"
	}
}

func (m FractionalCalculation) IsValid() bool {
	switch m {
	case FractionalCalculationDecimal,
		FractionalCalculationRoundDown,
		FractionalCalculationRoundUp,
		FractionalCalculationRoundNearest,
		FractionalCalculationRoundFive:
		return true
	default:
		return false
	}
}

// Apply は端数計算方法に基づいて値を丸める
func (m FractionalCalculation) Apply(value float64) float64 {
	switch m {
	case FractionalCalculationDecimal:
		return value
	case FractionalCalculationRoundDown:
		// 切り捨て（0に近づく方向）
		return math.Trunc(value)
	case FractionalCalculationRoundUp:
		// 切り上げ（0から離れる方向）
		return math.Copysign(math.Ceil(math.Abs(value)), value)
	case FractionalCalculationRoundNearest:
		// 四捨五入
		return math.Round(value)
	case FractionalCalculationRoundFive:
		// 五捨六入（0.5以下切り捨て、0.6以上切り上げ）
		abs := math.Abs(value)
		frac := abs - math.Floor(abs)
		if frac <= 0.5 {
			return math.Copysign(math.Floor(abs), value)
		}
		return math.Copysign(math.Ceil(abs), value)
	default:
		return value
	}
}

// IsDecimal は小数点有効かどうかを返す
func (m FractionalCalculation) IsDecimal() bool {
	return m == FractionalCalculationDecimal
}
