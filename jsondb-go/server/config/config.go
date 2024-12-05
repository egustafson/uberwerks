package config

var (
	configuration Config
)

func GetConfig() Config {
	return configuration
}

type Config struct {
	Flags Flags  `yaml:"-" json:"-"`
	Port  int    `yaml:"port" json:"port"`
	DSN   string `yaml:"dsn" json:"dsn"`
}

type Flags struct {
	Verbose   bool
	DevelMode bool
}
