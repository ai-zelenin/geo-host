package geo

import "fmt"

type BBox struct {
	XMin float64
	XMax float64
	YMin float64
	YMax float64
}

func NewBBox(coordsStr string) (*BBox, error) {
	var xmin, ymin, xmax, ymax float64
	if coordsStr != "" {
		coords, err := ParseAsFloatArray(coordsStr)
		if err != nil {
			return nil, err
		}
		if len(coords) != 4 {
			return nil, fmt.Errorf("invalid format of bbox coordinates [%s]", coordsStr)
		}
		xmin = coords[0]
		ymin = coords[1]
		xmax = coords[2]
		ymax = coords[3]
	}
	return &BBox{
		XMin: xmin,
		XMax: xmax,
		YMin: ymin,
		YMax: ymax,
	}, nil
}

func (r *BBox) AsPolygon() *GeographicPolygon {
	return &GeographicPolygon{
		Points: []*GeographicPoint{
			{
				Latitude:  r.XMin,
				Longitude: r.YMin,
			},
			{
				Latitude:  r.XMin,
				Longitude: r.YMax,
			},

			{
				Latitude:  r.XMin,
				Longitude: r.YMax,
			},
			{
				Latitude:  r.XMax,
				Longitude: r.YMax,
			},

			{
				Latitude:  r.XMax,
				Longitude: r.YMax,
			},
			{
				Latitude:  r.XMax,
				Longitude: r.YMin,
			},

			{
				Latitude:  r.XMax,
				Longitude: r.YMin,
			},
			{
				Latitude:  r.XMin,
				Longitude: r.YMin,
			},
		},
	}
}
