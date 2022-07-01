package geo

import (
	"database/sql/driver"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
)

type AbstractGeographic struct {
	SRID
	geom.T
}

func (g *AbstractGeographic) Scan(input interface{}) error {
	gt, err := ewkbhex.Decode(string(input.([]byte)))
	if err != nil {
		return err
	}
	g.SRID = SRID(gt.SRID())
	g.T = gt
	return nil
}

func (g AbstractGeographic) Value() (driver.Value, error) {
	return ewkbhex.Encode(g.T, ewkbhex.NDR)
}
