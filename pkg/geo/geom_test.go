package geo

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestGeomSuite(t *testing.T) {
	suite.Run(t, new(GeomSuite))
}

type GeomSuite struct {
	suite.Suite
}

func (s *GeomSuite) TestGeographicCollection() {
	polygon := &GeographicPolygon{
		SRID: WGS84,
		Points: []*GeographicPoint{
			{
				SRID:      WGS84,
				Latitude:  1,
				Longitude: 1,
			},
			{
				SRID:      WGS84,
				Latitude:  2,
				Longitude: 2,
			},
			{
				SRID:      WGS84,
				Latitude:  2,
				Longitude: 2,
			},
			{
				SRID:      WGS84,
				Latitude:  3,
				Longitude: 3,
			},
			{
				SRID:      WGS84,
				Latitude:  3,
				Longitude: 3,
			},
			{
				SRID:      WGS84,
				Latitude:  1,
				Longitude: 1,
			},
		},
	}
	point := &GeographicPoint{
		SRID:      WGS84,
		Latitude:  1,
		Longitude: 1,
	}
	gc := &GeographicCollection{
		SRID: WGS84,
		Figures: []Primitive{
			polygon,
			point,
		},
	}
	t, err := gc.ToGeom()
	if !s.Nil(err) {
		return
	}
	gc2 := &GeographicCollection{}
	err = gc2.FromGeom(t)
	if !s.Nil(err) {
		return
	}
	s.ElementsMatch(gc.Figures, gc2.Figures)
}
