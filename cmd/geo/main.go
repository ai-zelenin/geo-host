package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ai-zelenin/geo-host/pkg/geo"
	"github.com/ai-zelenin/geo-host/pkg/pgds"
	"github.com/ai-zelenin/geo-host/pkg/server"
	"github.com/go-pg/pg/v10"
	"io/ioutil"
	"log"
	"os"
)

type GlobalPoints struct {
	ID       int64
	Name     string
	Location *geo.GeographicPoint
	Lat      float64
	Lon      float64
	Qk10     int64  `pg:"qk_10"`
	Qk4      string `pg:"qk_4"`
}

func main() {
	ctx := context.Background()
	dsn := "postgres://postgres:postgis@localhost:5432/postgres?sslmode=disable"
	gs := geo.NewGeographicSystem(geo.DefaultGeoSystemConfig)
	ds, err := pgds.NewPostGISDataSource(ctx, dsn, gs, func(obj *pgds.Cluster) map[string]interface{} {
		var iconContent string
		var balloonContent string
		if obj.Count > 1 {
			iconContent = fmt.Sprintf("%d", obj.Count)
			for _, child := range obj.ClusterData {
				balloonContent += fmt.Sprintf("%s<br>\n", child.Properties["name"])
			}
			//"options": map[string]interface{}{
			//				"fillColor": fmt.Sprintf("rgba(27, 125, 27, 0.2)"),
			//			},

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
	})
	if err != nil {
		log.Fatal(err)
	}
	ExportPoints(ds, gs)
	cfg := &server.Config{
		ServerAddr: ":8080",
		StaticDir:  "./front",
	}
	srv := server.NewServer(cfg, ds, gs)
	err = srv.Start()
	if err != nil {
		log.Fatal(err)
	}

}

func ImportPoints(db *pg.DB) {
	var points = make([]*GlobalPoints, 0)
	err := db.Model(&points).Select()
	if err != nil {
		panic(err)
	}
	data, err := json.MarshalIndent(points, "", "\t")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("metro.json", data, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func ExportPoints(ds geo.DataSource, gs *geo.GeographicSystem) {
	data, err := ioutil.ReadFile("metro.json")
	if err != nil {
		panic(err)
	}
	points := make([]*GlobalPoints, 0)
	err = json.Unmarshal(data, &points)
	if err != nil {
		panic(err)
	}
	for _, point := range points {
		err = ds.StoreGeoData(context.Background(), &pgds.GeoObject{
			ID:    point.ID,
			Lat:   point.Location.Latitude,
			Lon:   point.Location.Longitude,
			Point: point.Location,
			Properties: map[string]interface{}{
				"name": point.Name,
				"qk4":  gs.WGS84ToQuadKey(point.Location.Latitude, point.Location.Longitude).String(),
			},
		})
		if err != nil {
			panic(err)
		}
	}
}
