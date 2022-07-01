package server

import (
	"context"
	"github.com/ai-zelenin/geo-host/pkg/geo"
	"net/http"
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
		r.URL.Query().Get("clusterDepth"),
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
	w.Header().Set("Content-Type", "application/javascript")
	w.Header().Set("Cache-Control", "max-age=1200")
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
