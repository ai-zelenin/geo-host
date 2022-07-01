package geo

type Tile struct {
	ID         int64
	X          int64
	Y          int64
	Zoom       int64
	QuadKey    QuadKey
	MinQuadKey QuadKey
	MaxQuadKey QuadKey
}
