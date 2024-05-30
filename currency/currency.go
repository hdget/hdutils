package currency

import (
	"math"
)

type Utils interface {
	ToFixed(num float64, precision int) float64
}

// ToFixed 浮点数到指定小数位
func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
