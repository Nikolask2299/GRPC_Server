package config

import "time"

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
	

}