package geo

import (
	"math"
)

const (
	TileSize = 256
	MinZoom  = 0
	MaxZoom  = 23
)

//TileSystem allow transform any coordinate to TileXY
type TileSystem struct {
	tileSize float64
	minZoom  int64
	maxZoom  int64
}

func NewTileSystem(minZoom, maxZoom int64, tileSize float64) *TileSystem {
	return &TileSystem{
		tileSize: tileSize,
		minZoom:  minZoom,
		maxZoom:  maxZoom,
	}
}

func (t *TileSystem) GlobalPixelsToTileXY(gpx, gpy float64) (tx, ty int64) {
	tx = int64(math.Floor(gpx / t.tileSize))
	ty = int64(math.Floor(gpy / t.tileSize))
	return tx, ty
}

func (t *TileSystem) TileXYToGlobalPixels(tx, ty int64) (gpx, gpy float64) {
	gpx = float64(tx * int64(t.tileSize))
	gpy = float64(ty * int64(t.tileSize))
	return
}
