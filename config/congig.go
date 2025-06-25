package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	TG_API struct {
		Token     string `env-required:"true" yaml:"token" env:"TG_API_TOKEN"`
	} `yaml:"tg_api"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	_ = cleanenv.ReadEnv(&cfg)

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Println("Config load error:", err)
		return nil, err
	}
	return &cfg, nil
}
