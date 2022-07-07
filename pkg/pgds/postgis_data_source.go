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
	Point         *geo.GeographicPoint   `bun:"point,type:geography(POINT,4326)"`
}

type Cluster struct {
	ID          int64                `bun:"tile_id"`
	MinID       int64                `bun:"min_id"`
	Count       int64                `bun:"count"`
	ClusterData []*GeoObject         `bun:"cluster_data"`
	Centroid    *geo.GeographicPoint `bun:"centroid"`
	GeoObject
}

type PostGISDataSource struct {
	gs     *geo.GeographicSystem
	DB     *bun.DB
	mapper PropertiesMapper
}

func NewPostGISDataSource(ctx context.Context, dsn string, gs *geo.GeographicSystem, mapper PropertiesMapper) (*PostGISDataSource, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	_, err := db.NewCreateTable().Model(new(GeoObject)).Table("geo_objects").IfNotExists().Exec(ctx)
	if err != nil {
		return nil, err
	}
	_, err = db.NewCreateIndex().Model(new(GeoObject)).Index("point_st_gist").Column("point").Using("SPGIST").IfNotExists().Exec(ctx)
	if err != nil {
		return nil, err
	}
	_, err = db.NewCreateIndex().Model(new(GeoObject)).Index("quad_key_btree").Column("quad_key").IfNotExists().Exec(ctx)
	if err != nil {
		return nil, err
	}
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	return &PostGISDataSource{
		gs:     gs,
		DB:     db,
		mapper: mapper,
	}, nil
}

func (p *PostGISDataSource) LoadMapView(ctx context.Context, mr *geo.MapRequest, fc *geo.FeatureCollection) error {
	tiles := p.gs.MRToTiles(mr)
	tileIDs := make([]int64, 0, len(tiles))
	for id := range tiles {
		tileIDs = append(tileIDs, id)
	}
	bitDelta := p.gs.QuadKeySystem.BitDelta(mr.Zoom)
	clusterShift := p.gs.QuadKeySystem.BitDelta(mr.Zoom + mr.ClusterDepth)
	objects := make([]*Cluster, 0, mr.TilesNumber()*(mr.ClusterDepth*4))
	subq := p.DB.NewSelect().Model((*GeoObject)(nil))
	subq.ColumnExpr("COUNT(id) AS count")
	subq.ColumnExpr("MIN(id) AS min_id")
	subq.ColumnExpr("st_centroid(st_collect(point::geometry)) as centroid")
	subq.ColumnExpr("quad_key >> ? as tile_id", clusterShift)
	subq.Where("quad_key >> ? in (?)", bitDelta, bun.In(tileIDs))
	subq.Order("tile_id")
	subq.Group("tile_id")
	q := p.DB.NewSelect()
	q.TableExpr("(?) AS cluster", subq)
	q.Join("left join geo_objects gp on gp.id = cluster.min_id")
	err := q.Scan(ctx, &objects)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	for _, object := range objects {
		// todo here we can put object into cache
		if object.Count > 1 {

			err = fc.Add(object.ID, object.Centroid, p.mapper(object))
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
	qk := p.gs.CoordinatesToQuadKey(gObj.Lat, gObj.Lon)
	gObj.QuadKey = qk.Int64()
	_, err := p.DB.NewInsert().Model(d).On("CONFLICT (id) DO UPDATE").Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func LoadMapView(ctx context.Context, db bun.DB, mr *geo.MapRequest, fc *geo.FeatureCollection) error {
	gs := geo.NewGeographicSystem(geo.DefaultGeoSystemConfig)
	tiles := gs.MRToTiles(mr)
	tileIDs := make([]int64, 0, len(tiles))
	for id := range tiles {
		tileIDs = append(tileIDs, id)
	}
	bitDelta := gs.QuadKeySystem.BitDelta(mr.Zoom)
	clusterShift := gs.QuadKeySystem.BitDelta(mr.Zoom + mr.ClusterDepth)
	objects := make([]*Cluster, 0, mr.TilesNumber()*(mr.ClusterDepth*4))
	subq := db.NewSelect().Model((*GeoObject)(nil))
	subq.ColumnExpr("COUNT(id) AS count")
	subq.ColumnExpr("MIN(id) AS min_id")
	subq.ColumnExpr("st_centroid(st_collect(point::geometry)) as centroid")
	subq.ColumnExpr("quad_key >> ? as tile_id", clusterShift)
	subq.Where("quad_key >> ? in (?)", bitDelta, bun.In(tileIDs))
	subq.Order("tile_id")
	subq.Group("tile_id")
	q := db.NewSelect()
	q.TableExpr("(?) AS cluster", subq)
	q.Join("left join geo_objects gp on gp.id = cluster.min_id")
	err := q.Scan(ctx, &objects)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	var mapper PropertiesMapper = func(obj *Cluster) map[string]interface{} {
		var iconContent string
		var balloonContent string
		if obj.Count > 1 {
			iconContent = fmt.Sprintf("%d", obj.Count)
			for _, child := range obj.ClusterData {
				balloonContent += fmt.Sprintf("%s<br>\n", child.Properties["name"])
			}
		} else {
			iconContent = obj.Properties["name"].(string)
		}
		return map[string]interface{}{
			"hintContent":    obj.ID,
			"iconContent":    iconContent,
			"balloonContent": balloonContent,
			"options": map[string]interface{}{
				"preset":    "islands#blackStretchyIcon",
				"fillColor": fmt.Sprintf("rgba(27, 125, 27, 0.2)"),
			},
		}
	}
	for _, object := range objects {
		// todo here we can put object into cache
		if object.Count > 1 {
			err = fc.Add(object.ID, object.Centroid, mapper(object))
			if err != nil {
				return err
			}
		} else {
			point := &geo.GeographicPoint{
				Latitude:  object.Lat,
				Longitude: object.Lon,
			}
			err = fc.Add(object.ID, point, mapper(object))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func StoreGeoData(ctx context.Context, db bun.DB, gObj *GeoObject) error {
	gs := geo.NewGeographicSystem(geo.DefaultGeoSystemConfig)
	qk := gs.CoordinatesToQuadKey(gObj.Lat, gObj.Lon)
	gObj.QuadKey = qk.Int64()
	_, err := db.NewInsert().Model(gObj).On("CONFLICT (id) DO UPDATE").Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
