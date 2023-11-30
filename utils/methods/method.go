package methods

import "math"

func UpperBound(originSize int, maxNum int) float64 {
	divNum := float64(originSize) / float64(maxNum)
	sqrtNum := math.Sqrt(float64(divNum))
	scale := 1 / sqrtNum
	return scale
}
