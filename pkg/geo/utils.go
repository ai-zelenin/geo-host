package geo

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
	"strconv"
	"strings"
)

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

func ToWKT(g Primitive) string {
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
