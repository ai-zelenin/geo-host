package pgds

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ai-zelenin/geo-host/pkg/geo"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type PropertiesMapper func(obj *Cluster) map[string]interface{}

type GeoObject struct {
	bun.BaseModel `bun:"table:geo_objects"`
	ID            int64                  `bun:"id,pk,autoincrement"`
	QuadKey       int64                  `bun:"quad_key,notnull"`
	Lat           float64                `bun:"lat,notnull"`
	Lon           float64                `bun:"lon,notnull"`
	Properties    map[string]interface{} `bun:"properties"`
}

type Cluster struct {
	ID     int64   `bun:"cluster_id"`
	MinID  int64   `bun:"min_id"`
	AvgLat float64 `bun:"cluster_lat"`
	AvgLon float64 `bun:"cluster_lon"`
	Count  int64   `bun:"count"`
	GeoObject
}

type PostGISDataSource struct {
	gs     *geo.GeographicSystem
	db     *bun.DB
	mapper PropertiesMapper
}

func NewPostGISDataSource(ctx context.Context, dsn string, gs *geo.GeographicSystem, mapper PropertiesMapper) (*PostGISDataSource, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	_, err := db.NewCreateTable().Model(new(GeoObject)).Table("geo_objects").IfNotExists().Exec(ctx)
	if err != nil {
		return nil, err
	}
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	return &PostGISDataSource{
		gs:     gs,
		db:     db,
		mapper: mapper,
	}, nil
}

func (p *PostGISDataSource) LoadMapView(ctx context.Context, mr *geo.MapRequest, fc *geo.FeatureCollection) error {
	minQk, maxQk, nextLevelMask := p.gs.MRToBoundsAndClusterMask(mr)
	objects := make([]*Cluster, 0, 100)
	subq := p.db.NewSelect().Model((*GeoObject)(nil))
	subq.ColumnExpr("COUNT(id) AS count")
	subq.ColumnExpr("MIN(id) AS min_id")
	subq.ColumnExpr("AVG(lon) AS cluster_lon")
	subq.ColumnExpr("AVG(lat) AS cluster_lat")
	subq.ColumnExpr("quad_key & ? as cluster_id", nextLevelMask.Int64())
	subq.Where("quad_key >= ? AND quad_key <= ?", minQk.Int64(), maxQk.Int64())
	subq.Group("cluster_id")

	q := p.db.NewSelect()
	q.TableExpr("(?) AS cluster", subq)
	q.Join("left join geo_objects gp on gp.id = cluster.min_id")
	err := q.Scan(ctx, &objects)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	for _, object := range objects {
		if object.Count > 1 {
			point := &geo.GeographicPoint{
				Latitude:  object.AvgLat,
				Longitude: object.AvgLon,
			}
			err = fc.Add(object.ID, point, p.mapper(object))
			if err != nil {
				return err
			}
		} else {
			point := &geo.GeographicPoint{
				Latitude:  object.Lat,
				Longitude: object.Lon,
			}
			err = fc.Add(object.ID, point, p.mapper(object))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *PostGISDataSource) StoreGeoData(ctx context.Context, d interface{}) error {
	gObj, ok := d.(*GeoObject)
	if !ok {
		return fmt.Errorf("unexpected data type %T", d)
	}
	gObj.QuadKey = p.gs.WGS84ToQuadKey(gObj.Lat, gObj.Lon).Int64()
	_, err := p.db.NewInsert().Model(d).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
