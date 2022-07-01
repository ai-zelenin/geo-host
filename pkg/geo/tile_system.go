package geo

import (
	"math"
)

//TileSystem allow transform any coordinate to TileXY
type TileSystem struct {
	tileSize int64
	minZoom  int64
	maxZoom  int64
}

func NewTileSystem(minZoom, maxZoom int64, tileSize int64) *TileSystem {
	return &TileSystem{
		tileSize: tileSize,
		minZoom:  minZoom,
		maxZoom:  maxZoom,
	}
}

func (t *TileSystem) GlobalPixelsToTileXY(gpx, gpy float64) (tx, ty int64) {
	tx = int64(math.Floor(gpx / float64(t.tileSize)))
	ty = int64(math.Floor(gpy / float64(t.tileSize)))
	return tx, ty
}

func (t *TileSystem) TileXYToGlobalPixels(tx, ty int64) (gpx, gpy float64) {
	gpx = float64(tx * t.tileSize)
	gpy = float64(ty * t.tileSize)
	return
}
func (t *TileSystem) TileXYToGlobalPixelsCenter(tx, ty int64) (gpx, gpy float64) {
	gpx = float64(tx*t.tileSize + t.tileSize/2)
	gpy = float64(ty*t.tileSize + t.tileSize/2)
	return
}
