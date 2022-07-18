package geo

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type MapRequestCase struct {
	coordsStr, tileStr, zoomStr, callbackID, debugStr, clusterDepthStr string
	Result                                                             *MapRequest
	Error                                                              bool
}

var cases = []MapRequestCase{
	{
		coordsStr:  "1.1,2,3,4",
		zoomStr:    "1",
		callbackID: "asd",
		Result: &MapRequest{
			BBox: BBox{
				XMin: 1.1,
				XMax: 3,
				YMin: 2,
				YMax: 4,
			},
			Zoom:       1,
			CallbackID: "asd",
		},
	},
	{
		coordsStr: "1c,2,3,4",
		zoomStr:   "1",
		Result:    nil,
		Error:     true,
	},
	{
		tileStr:    "1,2,3,4",
		zoomStr:    "1",
		callbackID: "asd",
		Result: &MapRequest{
			TileBBox: TileBBox{
				TileXMin: 1,
				TileXMax: 3,
				TileYMin: 2,
				TileYMax: 4,
			},
			Zoom:       1,
			CallbackID: "asd",
		},
	},
	{
		tileStr:    "1,2,3,4",
		zoomStr:    "1",
		callbackID: "asd",
		Result: &MapRequest{
			TileBBox: TileBBox{
				TileXMin: 1,
				TileXMax: 3,
				TileYMin: 2,
				TileYMax: 4,
			},
			Zoom:       1,
			CallbackID: "asd",
		},
	},
}

func TestMapRequestSuite(t *testing.T) {
	suite.Run(t, new(MapRequestSuite))
}

type MapRequestSuite struct {
	suite.Suite
}

func (s *MapRequestSuite) TestParseMapRequest() {
	for i, rc := range cases {
		mr, err := ParseMapRequest(rc.coordsStr, rc.tileStr, rc.zoomStr, rc.callbackID, rc.debugStr, rc.clusterDepthStr)
		if !s.EqualValues(rc.Result, mr) {
			s.Failf("TestParseMapRequest", "fail on %d: err=[%v]", i, err)
			return
		}
		if !s.Equal(rc.Error, err != nil) {
			s.Failf("TestParseMapRequest", "fail on %d: err=[%v]", i, err)
			return
		}
	}
}
