package geo

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
)

type FeatureCollection struct {
	geojson.FeatureCollection
}

func NewFeatureCollection() *FeatureCollection {
	return &FeatureCollection{FeatureCollection: geojson.FeatureCollection{
		Features: make([]*geojson.Feature, 0),
	}}
}

func (f *FeatureCollection) Add(id interface{}, geo interface{}, properties map[string]interface{}) error {
	var feature *geojson.Feature
	switch g := geo.(type) {
	case Primitive:
		gt, er := g.ToGeom()
		if er != nil {
			return er
		}
		feature = &geojson.Feature{
			ID:         fmt.Sprintf("%v", id),
			Geometry:   gt,
			Properties: properties,
		}
	case geom.T:
		feature = &geojson.Feature{
			ID:         fmt.Sprintf("%v", id),
			Geometry:   g,
			Properties: properties,
		}
	default:
		return fmt.Errorf("unexpected feature type %T", geo)
	}
	if feature == nil {
		return fmt.Errorf("bad type %T", geo)
	}
	f.Features = append(f.Features, feature)
	return nil
}

func (f *FeatureCollection) MarshalToJSONP(callbackID string) ([]byte, error) {
	data, err := f.MarshalJSON()
	if err != nil {
		return nil, err
	}
	callbackIDBytes := []byte(callbackID + "(")
	data = append(callbackIDBytes, data...)
	data = append(data, []byte(")")...)
	return data, nil
}
