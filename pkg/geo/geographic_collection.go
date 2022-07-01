package geo

import (
	"database/sql/driver"
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
)

type GeographicCollection struct {
	SRID
	Figures []Primitive
}

func (p *GeographicCollection) FromGeom(t geom.T) error {
	gc, ok := t.(*geom.GeometryCollection)
	if !ok {
		return fmt.Errorf("wrong type %T", t)
	}
	p.SRID = SRID(gc.SRID())
	p.Figures = make([]Primitive, gc.NumGeoms())
	for i, gt := range gc.Geoms() {
		figure, err := FromGeom(gt)
		if err != nil {
			return err
		}
		p.Figures[i] = figure
	}
	return nil
}

func (p *GeographicCollection) ToGeom() (geom.T, error) {
	sridType := p.SRID
	if sridType == 0 {
		sridType = WGS84
	}
	gc := geom.NewGeometryCollection()
	gc.SetSRID(int(p.SRID))
	for _, figure := range p.Figures {
		gt, err := figure.ToGeom()
		if err != nil {
			return nil, err
		}
		err = gc.Push(gt)
		if err != nil {
			return nil, err
		}
	}
	return gc, nil
}

func (p *GeographicCollection) Scan(input interface{}) error {
	gt, err := ewkbhex.Decode(string(input.([]byte)))
	if err != nil {
		return err
	}
	return p.FromGeom(gt)
}

func (p GeographicCollection) Value() (driver.Value, error) {
	t, err := p.ToGeom()
	if err != nil {
		return nil, err
	}
	return ewkbhex.Encode(t, ewkbhex.NDR)
}
