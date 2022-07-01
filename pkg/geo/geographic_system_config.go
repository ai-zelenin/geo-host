package geo

const (
	DefaultTileSize = 256
	DefaultMinZoom  = 0
	DefaultMaxZoom  = 23
)

var DefaultGeoSystemConfig = &Config{
	ProjectionType: WGS84,
	MinZoom:        DefaultMinZoom,
	MaxZoom:        DefaultMaxZoom,
	TileSize:       DefaultTileSize,
}

type Config struct {
	TileSize       int64 `json:"tile_size" yaml:"tile_size"`
	MinZoom        int64 `json:"min_zoom" yaml:"min_zoom"`
	MaxZoom        int64 `json:"max_zoom" yaml:"max_zoom"`
	ProjectionType SRID  `json:"projection_type" yaml:"projection_type"`
}
