package pigfarm

import (
	"github.com/ironpark/genpig"
)

func init() {
	genpig.SetConfigPaths("$HOME", ".", "./config")
	genpig.SetConfigNames("myconfig")
}

//go:generate genpig -struct Config
type Config struct {
	Database struct {
		Host     string `env:"DB_SERVER_IP" json:"ip"`
		Port     int    `env:"DB_PORT" json:"port"`
		User     string `env:"DB_USER" json:"user"`
		Password string `env:"DB_PW" json:"password"`
	} `json:"database"`
	Server struct {
		Host string `json:"ip" env:"SERVER_IP"`
		Port int    `json:"port" env:"PORT"`
	} `json:"server"`
}
