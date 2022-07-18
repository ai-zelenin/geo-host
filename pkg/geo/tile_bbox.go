package geo

import (
	"fmt"
)

type TileBBox struct {
	TileXMin int64
	TileXMax int64
	TileYMin int64
	TileYMax int64
}

func NewTileBBox(tileStr string) (TileBBox, error) {
	var txmin, tymin, txmax, tymax int64
	if tileStr != "" {
		coords, err := ParseAsInt64Array(tileStr)
		if err != nil {
			return TileBBox{}, fmt.Errorf("tile bound parse error %w", err)
		}
		if len(coords) != 4 {
			return TileBBox{}, fmt.Errorf("invalid format of tile coordinates [%s]", tileStr)
		}
		txmin = coords[0]
		tymin = coords[1]
		txmax = coords[2]
		tymax = coords[3]
	}
	return TileBBox{
		TileXMin: txmin,
		TileXMax: txmax,
		TileYMin: tymin,
		TileYMax: tymax,
	}, nil
}

func (r *TileBBox) TilesNumber() int64 {
	return (1 + (r.TileXMax - r.TileXMin)) * (1 + (r.TileYMax - r.TileYMin))
}

func (r *TileBBox) IterateTiles(cb func(x, y int64) error) error {
	for i := r.TileXMin; i <= r.TileXMax; i++ {
		for j := r.TileYMin; j <= r.TileYMax; j++ {
			err := cb(i, j)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
