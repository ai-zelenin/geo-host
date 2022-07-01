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
	ClusterDepth int64
}

// ParseMapRequest from comma separated strings
func ParseMapRequest(coordsStr, tileStr, zoomStr, callbackID, debugStr, clusterDepthStr string) (*MapRequest, error) {
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
		return nil, fmt.Errorf("zoom parse error [%v]", err)
	}
	var debug bool
	if debugStr != "" {
		debug, err = strconv.ParseBool(debugStr)
		if err != nil {
			return nil, fmt.Errorf("debug parse error [%v]", err)
		}
	}
	var cl int64
	if clusterDepthStr != "" {
		cl, err = strconv.ParseInt(clusterDepthStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("clusterDepth parse error [%v]", err)
		}
		if cl > 4 {
			return nil, fmt.Errorf("clusterDepth cannot be greater then 4")
		}
	}
	return &MapRequest{
		BBox:         bbox,
		TileBBox:     tileBBox,
		Zoom:         int64(zoom),
		CallbackID:   callbackID,
		Debug:        debug,
		ClusterDepth: cl,
	}, nil
}
