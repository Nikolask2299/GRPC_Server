package config

import (
	"flag"
	"os"
	"time"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	TokenTTL time.Duration `yaml:"token_ttl" env-default:"10m"`
	GRPC        GRPCConfig `yaml:"grpc" env-required:"true"`
}

type GRPCConfig struct {
	Port    int `yaml:"port" env-default:"8080"`
	Timeout time.Duration `yaml:"timeout" env-default:"1s"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config file not found at path")
	}
	return MustLoadByPath(path)
}


func MustLoadByPath(configPath string) *Config {
	

	if  _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist at path" + configPath)
	}
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config file: " + err.Error())
	}
	return &cfg
}




func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "Path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}