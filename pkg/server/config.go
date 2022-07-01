package server

type Config struct {
	ServerAddr string `json:"server_addr" yaml:"server_addr"`
	StaticDir  string `json:"static_dir" yaml:"static_dir"`
}
