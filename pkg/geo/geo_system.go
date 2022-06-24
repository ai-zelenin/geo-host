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
	for i := rom.TileXMin; i <= rom.TileXMax; i++ {
		for j := rom.TileYMin; j <= rom.TileYMax; j++ {
			tilePolygon := g.TileXYToPolygon(i, j, rom.Zoom)
			id := fmt.Sprintf("tx:%d ty:%d", i, j)
			qk := g.QuadKeySystem.TileXYToQuadKey(i, j, rom.Zoom)
			minQk, maxQk := g.QuadKeySystem.QuadKeyRange(qk)
			err := fc.Add(id, tilePolygon, map[string]interface{}{
				"hintContent":  id,
				"quadKey":      qk.String(),
				"leftQuadKey":  minQk.String(),
				"rightQuadKey": maxQk.String(),
				"options": map[string]interface{}{
					"fillColor": fmt.Sprintf("rgba(27, 125, 27, 0.2)"),
				},
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *GeographicSystem) MRToQuadKeys(rom *MapRequest) (lqk, rqk QuadKey) {
	lqk = g.QuadKeySystem.TileXYToQuadKey(rom.TileXMin, rom.TileYMin, rom.Zoom)
	rqk = g.QuadKeySystem.TileXYToQuadKey(rom.TileXMax, rom.TileYMax, rom.Zoom)
	return lqk, rqk
}

func (g *GeographicSystem) MRToBoundsAndClusterMask(rom *MapRequest) (lqkMin, rqkMax, nextLevelMask QuadKey) {
	lqk := g.QuadKeySystem.TileXYToQuadKey(rom.TileXMin, rom.TileYMin, rom.Zoom)
	rqk := g.QuadKeySystem.TileXYToQuadKey(rom.TileXMax, rom.TileYMax, rom.Zoom)
	lqkMin, _ = g.QuadKeySystem.QuadKeyRange(lqk)
	_, rqkMax = g.QuadKeySystem.QuadKeyRange(rqk)
	nextLevelMask = g.QuadKeySystem.CreateMask(lqk, rom.Zoom, rom.ClusterLevel)
	return lqkMin, rqkMax, nextLevelMask
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
