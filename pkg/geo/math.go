package geo

import "math"

func CycleRestrict(value, min, max float64) float64 {
	return value - math.Floor((value-min)/(max-min))*(max-min)
}

func Restrict(value, min, max float64) float64 {
	return math.Max(math.Min(value, max), min)
}

const RadiansToDegreesFactor = 180 / math.Pi
const DegreesToRadiansFactor = math.Pi / 180

func RadiansToDegrees(x float64) float64 {
	return x * RadiansToDegreesFactor
}
func DegreesToRadians(x float64) float64 {
	return x * DegreesToRadiansFactor
}

// RoundToDigit rounds number to N digits after dot
func RoundToDigit(v float64, n int) float64 {
	var m = math.Pow(10, float64(n))
	return math.RoundToEven(v*m) / m
}
