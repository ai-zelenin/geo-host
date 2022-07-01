package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ai-zelenin/geo-host/pkg/geo"
	"github.com/ai-zelenin/geo-host/pkg/pgds"
	"github.com/ai-zelenin/geo-host/pkg/server"
	"github.com/uptrace/bun"
	"io/ioutil"
	"log"
	"os"
)

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

func ImportPoints(db *bun.DB) {
	var points = make([]*pgds.GeoObject, 0)
	err := db.NewSelect().Model(&points).Scan(context.Background(), &points)
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
	points := make([]*pgds.GeoObject, 0)
	err = json.Unmarshal(data, &points)
	if err != nil {
		panic(err)
	}
	for _, point := range points {
		err = ds.StoreGeoData(context.Background(), point)
		if err != nil {
			panic(err)
		}
	}
}
