package geo

import (
	"fmt"
	"github.com/twpayne/go-geom"
)

func FromGeom(g geom.T) (Primitive, error) {
	var primitive Primitive
	switch g.(type) {
	case *geom.Point:
		primitive = new(GeographicPoint)
	case *geom.Polygon:
		primitive = new(GeographicPolygon)
	case *geom.GeometryCollection:
		primitive = new(GeographicCollection)
	default:
		return nil, fmt.Errorf("unexpected type %T", g)
	}
	err := primitive.FromGeom(g)
	if err != nil {
		return nil, err
	}
	return primitive, nil
}
