package core

import "os"

type Config struct {
	CMCAPIKey string
	Port      string
}

func Load() *Config {
	return &Config{
		CMCAPIKey: os.Getenv("CMC_API_KEY"),
		Port:      getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}