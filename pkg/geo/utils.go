package geo

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
	"math"
	"strconv"
	"strings"
)

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

func ParseAsFloatArray(line string) ([]float64, error) {
	var result = make([]float64, 0)
	parts := strings.Split(line, ",")
	for i, part := range parts {
		val, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return nil, fmt.Errorf("%d element %v", i, err)
		}
		result = append(result, val)
	}
	return result, nil
}

func ParseAsInt64Array(line string) ([]int64, error) {
	var result = make([]int64, 0)
	parts := strings.Split(line, ",")
	for i, part := range parts {
		val, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%d element %v", i, err)
		}
		result = append(result, val)
	}
	return result, nil
}

func ToWKT(g GEOM) string {
	ge, err := g.ToGeom()
	if err != nil {
		return err.Error()
	}
	data, err := wkt.Marshal(ge)
	if err != nil {
		return err.Error()
	}
	return data
}
func FromWKT(s string) geom.T {
	g, err := wkt.Unmarshal(s)
	if err != nil {
		panic(err)
	}
	return g
}

func Bits(val int64) string {
	return strconv.FormatInt(val, 2)
}

// RoundToDigit rounds number to N digits after dot
func RoundToDigit(v float64, n int) float64 {
	var m = math.Pow(10, float64(n))
	return math.RoundToEven(v*m) / m
}

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

func Colorize(color string, s interface{}) string {
	return fmt.Sprintf("%s%v%s", color, s, colorReset)
}
