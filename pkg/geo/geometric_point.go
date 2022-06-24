package geo

import (
	"database/sql/driver"
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
)

type GeometricPoint struct {
	X float64
	Y float64
}

func (p *GeometricPoint) ToGeom() geom.T {
	return geom.NewPoint(geom.XY).MustSetCoords([]float64{p.X, p.Y})
}

func (p *GeometricPoint) Scan(input interface{}) error {
	gt, err := ewkbhex.Decode(string(input.([]byte)))
	if err != nil {
		return err
	}
	coords := gt.FlatCoords()
	if len(coords) != 2 {
		return fmt.Errorf("unexpected coordinates format %v", coords)
	}
	p.X = coords[0]
	p.Y = coords[1]
	return nil
}

func (p GeometricPoint) Value() (driver.Value, error) {
	geomPoint, err := geom.NewPoint(geom.XY).SetCoords([]float64{p.X, p.Y})
	if err != nil {
		return nil, err
	}
	ewkbhexGeom, err := ewkbhex.Encode(geomPoint, ewkbhex.NDR)
	if err != nil {
		return nil, err
	}
	return ewkbhexGeom, nil
}

func (p *GeometricPoint) GetX() float64 {
	return p.X
}

func (p *GeometricPoint) GetY() float64 {
	return p.Y
}
