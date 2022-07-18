package geo

import (
	"database/sql/driver"
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
)

type GeographicPoint struct {
	SRID
	Latitude  float64
	Longitude float64
}

func (p *GeographicPoint) FromGeom(t geom.T) error {
	point, ok := t.(*geom.Point)
	if !ok {
		return fmt.Errorf("wrong type %T", t)
	}
	p.SRID = SRID(point.SRID())
	p.Latitude = point.X()
	p.Longitude = point.Y()
	return nil
}

func (p *GeographicPoint) ToGeom() (geom.T, error) {
	srid := DefaultSRID(p.SRID)
	return geom.NewPoint(geom.XY).SetSRID(int(srid)).SetCoords([]float64{p.Latitude, p.Longitude})
}

func (p *GeographicPoint) Scan(input interface{}) error {
	gt, err := ewkbhex.Decode(string(input.([]byte)))
	if err != nil {
		return err
	}
	return p.FromGeom(gt)
}

func (p GeographicPoint) Value() (driver.Value, error) {
	t, err := p.ToGeom()
	if err != nil {
		return nil, err
	}
	return ewkbhex.Encode(t, ewkbhex.NDR)
}
