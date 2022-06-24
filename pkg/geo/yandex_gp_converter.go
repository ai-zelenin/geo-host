package geo

import "math"

const (
	// E Eccentricity of Earth ellipsoid
	E                = 0.0818191908426
	MajorEarthRadius = 6378137.0
	SubRadius        = 1 / MajorEarthRadius
	Equator          = 2 * math.Pi * MajorEarthRadius
	HalfEquator      = Equator / 2
	SubEquator       = 1 / Equator
	HalfPi           = math.Pi / 2
	Rad2Deg          = 180 / math.Pi

	MaxLat = 90.0
	MinLat = -90.0
	MaxLon = 180.0
	MinLon = -180.0

	EGS3857MinLat = -85.0840
	EGS3857MaxLat = 85.0840

	// precalculated coefficients for fast inverse Mercator calculations
	e2 = E * E
	e4 = e2 * e2
	e6 = e4 * e4
	e8 = e4 * e4
	d2 = e2/2 + 5*e4/24 + e6/12 + 13*e8/360
	d4 = 7*e4/48 + 29*e6/240 + 811*e8/11520
	d6 = 7*e6/120 + 81*e8/1120
	d8 = 4279 * e8 / 161280
)

//YandexGPConverter
// https://yandex.ru/dev/maps/tiles/doc/dg/concepts/about-tiles.html#get-tile-number
// https://yastatic.net/s3/front-maps-static/maps-front-jsapi-v2-1/2.1.79-41/build/debug/full-f444800dea2e3f74c30f6b13a88f1a6a7b7eb78a.js
type YandexGPConverter struct {
	E float64
}

func NewYandexGPConverter(e float64) *YandexGPConverter {
	return &YandexGPConverter{
		E: e,
	}
}

func (g YandexGPConverter) ToGlobalPixels(lat, lon float64, zoom int64) (gpx, gpy float64) {
	var z = float64(zoom)
	var ro = math.Pow(2, z+8) / 2
	// Longitude to gpx
	gpx = ro * (1 + lon/180)

	// Latitude to gpy
	// epsilon needed for prevent gpy=Inf case
	// for latitude -90 teta=0 and gpy = Inf
	var epsilon = 1e-10
	lat = Restrict(lat, MinLat+epsilon, MaxLat-epsilon)
	var beta = (math.Pi * lat) / 180
	var eSinBeta = g.E * math.Sin(beta)
	var fi = (1 - eSinBeta) / (1 + eSinBeta)
	var teta = math.Tan((math.Pi/4)+(beta/2)) * math.Pow(fi, g.E/2)
	gpy = ro * (1 - (math.Log(teta) / math.Pi))
	return gpx, gpy
}

func (g YandexGPConverter) FromGlobalPixels(gpx, gpy float64, zoom int64) (lat, lon float64) {
	var z = float64(zoom)
	var f = math.Pow(2, z+8)
	var ro = f / 2

	// GPX to longitude
	lon = (180*gpx)/ro - 180

	// GPY to latitude
	var y = HalfEquator - gpy/(f*SubEquator)
	var phi = HalfPi - 2*math.Atan(1/math.Exp(y*SubRadius))
	phi = phi + d2*math.Sin(2*phi) + d4*math.Sin(4*phi) + d6*math.Sin(6*phi) + d8*math.Sin(8*phi)
	lat = phi * Rad2Deg
	return RoundToDigit(lat, 7), RoundToDigit(lon, 7)
}

//GlobalPixelsToWGS84WithGD
// GD - https://en.wikipedia.org/wiki/Gudermannian_function
func GlobalPixelsToWGS84WithGD(gpx, gpy float64, zoom int64) (lat, lon float64) {
	var z = float64(zoom)
	var f = math.Pow(2, z+8)
	var ro = f / 2

	// Longitude to gpx
	//var lonInRad = CycleRestrict(lon * DegreesToRadiansFactor, -math.Pi, math.Pi)
	//gpx = MajorEarthRadius * lonInRad
	// GPX to longitude
	lon = (180*gpx)/ro - 180
	//lon = RadiansToDegrees(CycleRestrict(gpx*SubRadius, -math.Pi, math.Pi))

	// GPY to latitude
	var y = HalfEquator - gpy/(f*SubEquator)
	var ts = math.Exp(-y / MajorEarthRadius)
	var phi = HalfPi - 2*math.Atan(ts)
	for {
		con := E * math.Sin(phi)
		dPhi := HalfPi - 2*math.Atan(ts*math.Pow((1-con)/(1+con), E)) - phi
		phi += dPhi
		if math.Abs(dPhi) < 0.00000001 {
			break
		}
	}
	lat = phi * Rad2Deg
	return lat, lon
}
