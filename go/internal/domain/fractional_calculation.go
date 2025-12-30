package domain

type FractionalCalculation int

const (
	FractionalCalculationRoundUp           FractionalCalculation = iota + 1 // 切り上げ
	FractionalCalculationRoundDown                                          // 切り捨て
	FractionalCalculationRoundNearest                                       // 四捨五入
	FractionalCalculationRoundUpBelowTen                                    // 10点未満切り上げ
	FractionalCalculationRoundDownBelowTen                                  // 10点未満切り捨て
)

func (m FractionalCalculation) String() string {
	switch m {
	case FractionalCalculationRoundUp:
		return "切り上げ"
	case FractionalCalculationRoundDown:
		return "切り捨て"
	case FractionalCalculationRoundNearest:
		return "四捨五入"
	case FractionalCalculationRoundUpBelowTen:
		return "10点未満切り上げ"
	case FractionalCalculationRoundDownBelowTen:
		return "10点未満切り捨て"
	default:
		return "その他"
	}
}

func (m FractionalCalculation) IsValid() bool {
	switch m {
	case FractionalCalculationRoundUp,
		FractionalCalculationRoundDown,
		FractionalCalculationRoundNearest,
		FractionalCalculationRoundUpBelowTen,
		FractionalCalculationRoundDownBelowTen:
		return true
	default:
		return false
	}
}
