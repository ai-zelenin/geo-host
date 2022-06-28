package geo

import "fmt"

type GeographicSystem struct {
	TileSystem    *TileSystem
	QuadKeySystem *QuadKeySystem
	Projection    Projection
	cfg           *Config
}

func NewGeographicSystem(cfg *Config) *GeographicSystem {
	var projection Projection
	switch cfg.ProjectionType {
	case WebMercator:
		projection = NewYandexGPConverter(0)
	default:
		projection = NewYandexGPConverter(E)
	}
	return &GeographicSystem{
		Projection:    projection,
		TileSystem:    NewTileSystem(cfg.MinZoom, cfg.MaxZoom, float64(cfg.TileSize)),
		QuadKeySystem: NewQuadKeySystem(cfg.MinZoom, cfg.MaxZoom),
		cfg:           cfg,
	}
}

func (g *GeographicSystem) WGS84ToQuadKey(lat, long float64) QuadKey {
	gpx, gpy := g.Projection.ToGlobalPixels(lat, long, g.cfg.MaxZoom)
	tx, ty := g.TileSystem.GlobalPixelsToTileXY(gpx, gpy)
	return g.QuadKeySystem.TileXYToQuadKey(tx, ty, g.cfg.MaxZoom)
}

func (g *GeographicSystem) DrawROMBBox(rom *MapRequest, fc *FeatureCollection) error {
	props := map[string]interface{}{
		"options": map[string]interface{}{
			"fillColor": fmt.Sprintf("rgba(27, 27, 125, 0.3)"),
		},
	}
	return fc.Add("lat-lon-polygon", rom.AsPolygon(), props)
}

func (g *GeographicSystem) DrawROMTiles(rom *MapRequest, fc *FeatureCollection) error {
	return rom.IterateTiles(func(x, y int64) error {
		tilePolygon := g.TileXYToPolygon(x, y, rom.Zoom)
		id := fmt.Sprintf("tx:%d ty:%d", x, y)
		qk := g.QuadKeySystem.TileXYToQuadKey(x, y, rom.Zoom)
		minQk, maxQk := g.QuadKeySystem.QuadKeyRange(qk)
		return fc.Add(id, tilePolygon, map[string]interface{}{
			"hintContent":  id,
			"quadKey":      qk.String(),
			"leftQuadKey":  minQk.String(),
			"rightQuadKey": maxQk.String(),
			"options": map[string]interface{}{
				"fillColor": fmt.Sprintf("rgba(27, 125, 27, 0.2)"),
			},
		})

	})
}

func (g *GeographicSystem) MRToTiles(rom *MapRequest) (tiles []Tile, lqkMin, rqkMax QuadKey, tileBucketBitSize int64) {
	result := make([]Tile, 0, rom.TilesNumber())
	_ = rom.IterateTiles(func(x, y int64) error {
		qk := g.QuadKeySystem.TileXYToQuadKey(x, y, rom.Zoom)
		result = append(result, Tile{
			X:       x,
			Y:       y,
			Zoom:    rom.Zoom,
			QuadKey: qk,
		})
		return nil
	})
	lqkMin, _ = g.QuadKeySystem.QuadKeyRange(result[0].QuadKey)
	_, rqkMax = g.QuadKeySystem.QuadKeyRange(result[len(result)-1].QuadKey)
	return result, lqkMin, rqkMax, g.QuadKeySystem.Base10Delta(rom.Zoom)
}

func (g *GeographicSystem) TileXYToPoint(tx, ty int64, z int64) *GeographicPoint {
	var gpx, gpy = g.TileSystem.TileXYToGlobalPixels(tx, ty)
	lat, lon := g.Projection.FromGlobalPixels(gpx, gpy, z)
	return &GeographicPoint{
		Latitude:  lat,
		Longitude: lon,
	}
}

func (g *GeographicSystem) TileXYToPolygon(tx, ty int64, zoom int64) *GeographicPolygon {
	return &GeographicPolygon{
		Points: []*GeographicPoint{
			g.TileXYToPoint(tx, ty, zoom),
			g.TileXYToPoint(tx, ty+1, zoom),

			g.TileXYToPoint(tx, ty+1, zoom),
			g.TileXYToPoint(tx+1, ty+1, zoom),

			g.TileXYToPoint(tx+1, ty+1, zoom),
			g.TileXYToPoint(tx+1, ty, zoom),

			g.TileXYToPoint(tx+1, ty, zoom),
			g.TileXYToPoint(tx, ty, zoom),
		},
	}
}
