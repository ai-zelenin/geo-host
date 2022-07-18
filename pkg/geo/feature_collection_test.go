package geo

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestFeatureCollectionSuite(t *testing.T) {
	suite.Run(t, new(FeatureCollectionSuite))
}

type FeatureCollectionSuite struct {
	suite.Suite
}

func (s *FeatureCollectionSuite) TestFeatureCollection() {
	point := &GeographicPoint{
		SRID:      WGS84,
		Latitude:  1,
		Longitude: 1,
	}
	fc := NewFeatureCollection()
	err := fc.Add(1, point, nil)
	if !s.Nil(err) {
		return
	}
	gp, err := point.ToGeom()
	if !s.Nil(err) {
		return
	}
	err = fc.Add(2, gp, nil)
	if !s.Nil(err) {
		return
	}
	data, err := fc.MarshalToJSONP("TEST_1")
	if !s.Nil(err) {
		return
	}
	expectedResult := `TEST_1({"type":"FeatureCollection","features":[{"type":"Feature","id":"1","geometry":{"type":"Point","coordinates":[1,1]},"properties":null},{"type":"Feature","id":"2","geometry":{"type":"Point","coordinates":[1,1]},"properties":null}]})`
	if !s.Equal(expectedResult, string(data)) {
		return
	}
}
