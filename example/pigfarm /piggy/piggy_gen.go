package piggy

import base ".."

func Database() struct {
	Host     string `env:"DB_SERVER_IP" json:"ip"`
	Port     int    `env:"DB_PORT" json:"port"`
	User     string `env:"DB_USER" json:"user"`
	Password string `env:"DB_PW" json:"password"`
} {
	return base.GetInstance().Database()
}

func Server() struct {
	Host string `json:"ip" env:"SERVER_IP"`
	Port int    `json:"port" env:"PORT"`
} {
	return base.GetInstance().Server()
}
