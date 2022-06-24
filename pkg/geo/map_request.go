package geo

import (
	"fmt"
	"strconv"
)

type MapRequest struct {
	XMin         float64
	XMax         float64
	YMin         float64
	YMax         float64
	TileXMin     int64
	TileXMax     int64
	TileYMin     int64
	TileYMax     int64
	Zoom         int64
	CallbackID   string
	Debug        bool
	ClusterLevel int64
}

// ParseRequest from comma separated strings
func ParseRequest(coordsStr, tileStr, zoomStr, callbackID, debugStr, clusterLevelStr string) (*MapRequest, error) {
	var err error
	var xmin, ymin, xmax, ymax float64
	if coordsStr != "" {
		coords, err := ParseAsFloatArray(coordsStr)
		if err != nil {
			return nil, err
		}
		if len(coords) != 4 {
			return nil, fmt.Errorf("invalid format of lat-lon coordinates")
		}
		xmin = coords[0]
		ymin = coords[1]
		xmax = coords[2]
		ymax = coords[3]
	}
	var txmin, tymin, txmax, tymax int64
	if tileStr != "" {
		coords, err := ParseAsInt64Array(tileStr)
		if err != nil {
			return nil, err
		}
		if len(coords) != 4 {
			return nil, fmt.Errorf("invalid format of tile coordinates")
		}
		txmin = coords[0]
		tymin = coords[1]
		txmax = coords[2]
		tymax = coords[3]
	}
	zoom, err := strconv.ParseUint(zoomStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("zoom string %v", err)
	}
	var debug bool
	if debugStr != "" {
		debug, err = strconv.ParseBool(debugStr)
		if err != nil {
			return nil, err
		}
	}
	var cl int64
	if clusterLevelStr != "" {
		cl, err = strconv.ParseInt(clusterLevelStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("zoom string %v", err)
		}
	}
	return &MapRequest{
		XMin:         xmin,
		XMax:         xmax,
		YMin:         ymin,
		YMax:         ymax,
		TileXMin:     txmin,
		TileXMax:     txmax,
		TileYMin:     tymin,
		TileYMax:     tymax,
		Zoom:         int64(zoom),
		CallbackID:   callbackID,
		Debug:        debug,
		ClusterLevel: cl,
	}, nil
}

func (r *MapRequest) AsPolygon() *GeographicPolygon {
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
