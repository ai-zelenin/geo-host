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
		TileSystem:    NewTileSystem(cfg.MinZoom, cfg.MaxZoom, cfg.TileSize),
		QuadKeySystem: NewQuadKeySystem(cfg.MinZoom, cfg.MaxZoom),
		cfg:           cfg,
	}
}

func (g *GeographicSystem) CoordinatesToQuadKey(lat, long float64) QuadKey {
	gpx, gpy := g.Projection.ToGlobalPixels(lat, long, g.cfg.MaxZoom)
	tx, ty := g.TileSystem.GlobalPixelsToTileXY(gpx, gpy)
	return g.QuadKeySystem.TileXYToQuadKey(tx, ty, g.cfg.MaxZoom)
}

func (g *GeographicSystem) DrawROMBBox(mr *MapRequest, fc *FeatureCollection) error {
	props := map[string]interface{}{
		"options": map[string]interface{}{
			"fillColor": fmt.Sprintf("rgba(27, 27, 125, 0.3)"),
		},
	}
	return fc.Add("lat-lon-polygon", mr.AsPolygon(), props)
}

func (g *GeographicSystem) DrawROMTiles(mr *MapRequest, fc *FeatureCollection) error {
	return mr.IterateTiles(func(x, y int64) error {
		tilePolygon := g.TileXYToPolygon(x, y, mr.Zoom)
		id := fmt.Sprintf("tx:%d ty:%d", x, y)
		qk := g.QuadKeySystem.TileXYToQuadKey(x, y, mr.Zoom)
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

func (g *GeographicSystem) MRToTiles(mr *MapRequest) (tiles map[int64]Tile) {
	result := make(map[int64]Tile, mr.TilesNumber())
	_ = mr.IterateTiles(func(x, y int64) error {
		qk := g.QuadKeySystem.TileXYToQuadKey(x, y, mr.Zoom)
		id := qk.Int64()
		result[id] = Tile{
			ID:      id,
			X:       x,
			Y:       y,
			Zoom:    mr.Zoom,
			QuadKey: qk,
		}
		return nil
	})
	return result
}

func (g *GeographicSystem) TileIDToCenterPoint(tileID int64) (*GeographicPoint, error) {
	qk := NewQuadKeyFromInt64(tileID)
	tx, ty, err := g.QuadKeySystem.QuadKeyToTileXY(qk)
	if err != nil {
		return nil, err
	}
	return g.TileXYToCenterPoint(tx, ty, qk.Len()), nil
}

func (g *GeographicSystem) TileXYToCenterPoint(tx, ty int64, z int64) *GeographicPoint {
	var gpx, gpy = g.TileSystem.TileXYToGlobalPixelsCenter(tx, ty)
	lat, lon := g.Projection.FromGlobalPixels(gpx, gpy, z)
	return &GeographicPoint{
		Latitude:  lat,
		Longitude: lon,
	}
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
