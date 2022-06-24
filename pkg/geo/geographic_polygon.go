package geo

import (
	"database/sql/driver"
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
)

type GeographicPolygon struct {
	SRID
	Points []*GeographicPoint
}

func (p *GeographicPolygon) FromGeom(t geom.T) error {
	polygon, ok := t.(*geom.Polygon)
	if !ok {
		return fmt.Errorf("wrong type %T", t)
	}
	p.SRID = SRID(polygon.SRID())
	for _, coords := range polygon.Coords() {
		for _, pointsCoords := range coords {
			p.Points = append(p.Points, &GeographicPoint{
				SRID:      SRID(polygon.SRID()),
				Latitude:  pointsCoords[0],
				Longitude: pointsCoords[1],
			})
		}
	}
	return nil
}

func (p *GeographicPolygon) ToGeom() (geom.T, error) {
	sridType := p.SRID
	if sridType == 0 {
		sridType = WGS84
	}
	coords := make([]geom.Coord, len(p.Points))
	for i, point := range p.Points {
		tp, err := point.ToGeom()
		if err != nil {
			return nil, err
		}
		coords[i] = tp.FlatCoords()
	}
	lr, err := geom.NewLinearRing(geom.XY).SetCoords(coords)
	if err != nil {
		return nil, err
	}
	polygon := geom.NewPolygon(geom.XY)
	err = polygon.Push(lr)
	if err != nil {
		return nil, err
	}
	return polygon, nil
}

func (p *GeographicPolygon) Scan(input interface{}) error {
	gt, err := ewkbhex.Decode(string(input.([]byte)))
	if err != nil {
		return err
	}
	return p.FromGeom(gt)
}

func (p GeographicPolygon) Value() (driver.Value, error) {
	t, err := p.ToGeom()
	if err != nil {
		return nil, err
	}
	return ewkbhex.Encode(t, ewkbhex.NDR)
}
