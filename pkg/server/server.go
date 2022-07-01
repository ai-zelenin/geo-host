package server

import (
	"github.com/ai-zelenin/geo-host/pkg/geo"
	"net/http"
)

type Server struct {
	cfg *Config
	ds  geo.DataSource
	gs  *geo.GeographicSystem
}

func NewServer(cfg *Config, ds geo.DataSource, gs *geo.GeographicSystem) *Server {
	return &Server{cfg: cfg, ds: ds, gs: gs}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(s.cfg.StaticDir))
	mux.Handle("/", fs)
	mux.Handle("/api/v1/yandex", NewYandexROMHandler(s.gs, s.ds))
	srv := http.Server{
		Addr:    s.cfg.ServerAddr,
		Handler: mux,
	}
	return srv.ListenAndServe()
}
