package geo

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestGeoSystemSuite(t *testing.T) {
	suite.Run(t, new(GeoSystemSuite))
}

type GeoSystemSuite struct {
	suite.Suite
	gs *GeographicSystem
}

func (s *GeoSystemSuite) SetupTest() {
	s.gs = NewGeographicSystem(DefaultGeoSystemConfig)
}

func (s *GeoSystemSuite) TestGlobalPixelConverter() {
	s.Iterator(func(lat, lon float64, zoom int64) bool {
		gpx, gpy := s.gs.Projection.ToGlobalPixels(lat, lon, zoom)
		nlat, nlon := s.gs.Projection.FromGlobalPixels(gpx, gpy, zoom)
		if !s.Equal(lat, RoundToDigit(nlat, 5)) {
			return false
		}
		if !s.Equal(lon, RoundToDigit(nlon, 5)) {
			return false
		}
		tx, ty := s.gs.TileSystem.GlobalPixelsToTileXY(gpx, gpy)

		qk := s.gs.QuadKeySystem.TileXYToQuadKey(tx, ty, zoom)

		ntx, nty, err := s.gs.QuadKeySystem.QuadKeyToTileXY(qk)
		if !s.Nil(err) {
			return false
		}
		if !s.Equal(tx, ntx) {
			return false
		}
		if !s.Equal(ty, nty) {
			return false
		}
		pointQk := s.gs.CoordinatesToQuadKey(lat, lon)
		if !s.gs.QuadKeySystem.Contains(pointQk, qk) {
			s.Failf(
				"quadkey of tile do not contains quadkey of point",
				"pointQk:%s tileQk:%s", pointQk, qk)
			return false
		}
		return true
	})
}

func (s *GeoSystemSuite) Iterator(cb func(lat, lon float64, zoom int64) bool) {
	for i := MinLat / 2; i <= MaxLat/2; i++ {
		for j := MinLon / 2; j <= MaxLon/2; j++ {
			for z := s.gs.TileSystem.minZoom; z <= s.gs.TileSystem.maxZoom; z++ {
				ok := cb(i, j, z)
				if !ok {
					return
				}
			}
		}
	}
}
