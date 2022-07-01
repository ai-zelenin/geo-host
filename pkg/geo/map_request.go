package geo

import (
	"fmt"
	"strconv"
)

type MapRequest struct {
	*TileBBox
	*BBox
	Zoom         int64
	CallbackID   string
	Debug        bool
	ClusterLevel int64
}

// ParseMapRequest from comma separated strings
func ParseMapRequest(coordsStr, tileStr, zoomStr, callbackID, debugStr, clusterLevelStr string) (*MapRequest, error) {
	var err error
	bbox, err := NewBBox(coordsStr)
	if err != nil {
		return nil, err
	}
	tileBBox, err := NewTileBBox(tileStr)
	if err != nil {
		return nil, err
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
			return nil, fmt.Errorf("clusterLevelStr string %v", err)
		}
	}
	return &MapRequest{
		BBox:         bbox,
		TileBBox:     tileBBox,
		Zoom:         int64(zoom),
		CallbackID:   callbackID,
		Debug:        debug,
		ClusterLevel: cl,
	}, nil
}
