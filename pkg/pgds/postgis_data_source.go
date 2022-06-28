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
	"sort"
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
	ID          int64        `bun:"tile_id"`
	MinID       int64        `bun:"min_id"`
	AvgLat      float64      `bun:"cluster_lat"`
	AvgLon      float64      `bun:"cluster_lon"`
	Count       int64        `bun:"count"`
	ClusterData []*GeoObject `bun:"cluster_data"`
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
	tiles, minQk, maxQk, tileSize := p.gs.MRToTiles(mr)
	var min, max = minQk.Int64(), maxQk.Int64()
	sort.Slice(tiles, func(i, j int) bool {
		return tiles[i].QuadKey.Int64() < tiles[j].QuadKey.Int64()
	})
	objects := make([]*Cluster, 0, 100)
	subq := p.db.NewSelect().Model((*GeoObject)(nil))
	subq.ColumnExpr("COUNT(id) AS count")
	subq.ColumnExpr("MIN(id) AS min_id")
	subq.ColumnExpr("AVG(lon) AS cluster_lon")
	subq.ColumnExpr("AVG(lat) AS cluster_lat")
	subq.ColumnExpr("json_agg(row_to_json(geo_object)) as cluster_data")
	subq.ColumnExpr("width_bucket(quad_key::numeric, ?::numeric, ?::numeric,?) as tile_id", min, max, tileSize)
	subq.Where("quad_key >= ? AND quad_key <= ?", min, max)
	subq.Order("tile_id")
	subq.Group("tile_id")
	q := p.db.NewSelect()
	q.TableExpr("(?) AS cluster", subq)
	q.Join("left join geo_objects gp on gp.id = cluster.min_id")
	err := q.Scan(ctx, &objects)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	for _, object := range objects {
		if object.Count > 1 {
			tile := tiles[object.ID]
			point := p.gs.TileXYToPoint(tile.X, tile.Y, tile.Zoom)
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
	qk := p.gs.WGS84ToQuadKey(gObj.Lat, gObj.Lon)
	gObj.QuadKey = qk.Int64()
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = p.gs.QuadKeySystem.ForEachZoom(qk, func(x, y, z int64) error {
		fmt.Println(x, y, z)
		return nil
	})
	if err != nil {
		return err
	}
	_, err = tx.NewInsert().Model(d).On("conflict", "update").Exec(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
