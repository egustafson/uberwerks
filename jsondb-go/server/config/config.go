package config

type Config struct {
	Port int    `json:"port"`
	DSN  string `json:"dsn"`
}
