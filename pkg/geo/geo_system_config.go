package geo

var DefaultGeoSystemConfig = &Config{
	ProjectionType: WGS84,
	MinZoom:        MinZoom,
	MaxZoom:        MaxZoom,
	TileSize:       TileSize,
}

type Config struct {
	TileSize       int64
	MinZoom        int64
	MaxZoom        int64
	ProjectionType SRID
}
