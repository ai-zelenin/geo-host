package geo

import (
	"database/sql/driver"
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
)

type GeographicMultiPolygon struct {
	SRID
	Polygons []*GeographicPolygon
}

func (p *GeographicMultiPolygon) FromGeom(t geom.T) error {
	multiPolygon, ok := t.(*geom.MultiPolygon)
	if !ok {
		return fmt.Errorf("wrong type %T", t)
	}
	p.SRID = SRID(multiPolygon.SRID())
	for i := 0; i < multiPolygon.NumPolygons(); i++ {
		polygon := multiPolygon.Polygon(i)
		newGeographicPolygon := new(GeographicPolygon)
		err := newGeographicPolygon.FromGeom(polygon)
		if err != nil {
			return err
		}
		p.Polygons = append(p.Polygons, newGeographicPolygon)
	}
	return nil
}

func (p *GeographicMultiPolygon) ToGeom() (geom.T, error) {
	sridType := p.SRID
	if sridType == 0 {
		sridType = WGS84
	}
	mp := geom.NewMultiPolygon(geom.XY)
	for _, polygon := range p.Polygons {
		tp, err := polygon.ToGeom()
		if err != nil {
			return nil, err
		}
		err = mp.Push(geom.NewPolygonFlat(geom.XY, tp.FlatCoords(), tp.Ends()))
		if err != nil {
			return nil, err
		}
	}
	return mp, nil
}

func (p *GeographicMultiPolygon) Scan(input interface{}) error {
	gt, err := ewkbhex.Decode(string(input.([]byte)))
	if err != nil {
		return err
	}
	return p.FromGeom(gt)
}

func (p GeographicMultiPolygon) Value() (driver.Value, error) {
	t, err := p.ToGeom()
	if err != nil {
		return nil, err
	}
	return ewkbhex.Encode(t, ewkbhex.NDR)
}
