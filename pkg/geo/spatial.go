package geo

import (
	"database/sql/driver"
	"github.com/twpayne/go-geom"
)

type SRID int

const (
	WGS84       SRID = 4326
	WebMercator SRID = 3857
)

type Data interface {
	Scan(input interface{}) error
	Value() (driver.Value, error)
}

type GEOM interface {
	ToGeom() (geom.T, error)
	FromGeom(t geom.T) error
}

type Projection interface {
	ToGlobalPixels(lat, lon float64, z int64) (gpx, gpy float64)
	FromGlobalPixels(gpx, gpy float64, z int64) (lat, lon float64)
}
