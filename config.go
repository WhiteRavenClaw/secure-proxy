package main

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type SessionConfig struct {
	CookieDomain string `yaml:"cookieDomain"`
	CookieName   string `yaml:"cookieName"`
	TTLSeconds   int    `yaml:"ttlSeconds"`
}

type UserConfig struct {
	Username         string   `yaml:"username"`
	TOTPSecret       string   `yaml:"totpSecret"`
	AvailableDomains []string `yaml:"availableDomains"`
}

type UpstreamConfig struct {
	Host        string `yaml:"host"`
	Destination string `yaml:"destination"`
}

type AppConfig struct {
	Sessions  SessionConfig    `yaml:"sessions"`
	Users     []UserConfig     `yaml:"users"`
	Upstreams []UpstreamConfig `yaml:"upstreams"`
}

func LoadConfig() *AppConfig {
	var cfg AppConfig
	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		panic(fmt.Errorf("не удалось загрузить конфиг: %w", err))
	}
	return &cfg
}
