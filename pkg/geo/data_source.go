package geo

import "context"

type DataSource interface {
	LoadMapView(ctx context.Context, mr *MapRequest, fc *FeatureCollection) error
	StoreGeoData(ctx context.Context, d interface{}) error
}
