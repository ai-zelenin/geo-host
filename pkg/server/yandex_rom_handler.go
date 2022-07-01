package server

import (
	"context"
	"fmt"
	"github.com/ai-zelenin/geo-host/pkg/geo"
	"net/http"
	"time"
)

type YandexROMHandler struct {
	gs *geo.GeographicSystem
	ds geo.DataSource
}

func NewYandexROMHandler(gs *geo.GeographicSystem, ds geo.DataSource) *YandexROMHandler {
	return &YandexROMHandler{gs: gs, ds: ds}
}

func (y *YandexROMHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mr, err := geo.ParseMapRequest(
		r.URL.Query().Get("bbox"),
		r.URL.Query().Get("tiles"),
		r.URL.Query().Get("zoom"),
		r.URL.Query().Get("callback"),
		r.URL.Query().Get("debug"),
		r.URL.Query().Get("clusterLevel"),
	)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	fc, err := y.handleMapRequest(r.Context(), mr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data, err := fc.MarshalToJSONP(mr.CallbackID)
	if err != nil {
		panic(err)
	}
	now := time.Now()
	w.Header().Set("Content-Type", "application/javascript")
	w.Header().Set("Cache-Control", "max-age=1200")
	w.Header().Set("Last-Modified", time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Format(http.TimeFormat))
	w.Header().Set("Etag", fmt.Sprintf("%d-%d-%d-%d", mr.TileXMin, mr.TileXMax, mr.TileYMin, mr.TileYMax))
	_, _ = w.Write(data)
}

func (y *YandexROMHandler) handleMapRequest(ctx context.Context, mr *geo.MapRequest) (*geo.FeatureCollection, error) {
	fc := geo.NewFeatureCollection()
	if mr.Debug {
		err := y.gs.DrawROMTiles(mr, fc)
		if err != nil {
			return nil, err
		}
	}
	err := y.ds.LoadMapView(ctx, mr, fc)
	if err != nil {
		return nil, err
	}
	return fc, nil
}
