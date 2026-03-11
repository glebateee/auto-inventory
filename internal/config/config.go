package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"dev"`
	GRPCConfig gRPCConfig `yaml:"grpc_config"`
}

type gRPCConfig struct {
	Host string `yaml:"host" env-default:"localhost"`
	Port int    `yaml:"port" env-default:"8877"`
}

var (
	emptyPath = ""
)

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == emptyPath {
		panic("config path not set")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found: " + path)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}
	return &cfg
}

func fetchConfigPath() string {
	var path string
	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == emptyPath {
		path = os.Getenv("CONFIG_PATH")
	}
	return path
}
