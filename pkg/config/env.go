package config

import (
	"os"
)

type Config struct {
	TTSAPIURL string
	TTSAPIKey string
	Port      string
	TTSEndpoint string
	APIVersion  string
}

func LoadConfig() *Config {
	return &Config{
		TTSAPIURL: os.Getenv("TTS_API_URL"),
		TTSAPIKey: os.Getenv("TTS_API_KEY"),
		Port:      getEnvOrDefault("PORT", "8080"),
		TTSEndpoint: getEnvOrDefault("TTS_ENDPOINT", "/tts"),
		APIVersion:  getEnvOrDefault("API_VERSION", "v1"),
	}
}

func getEnvOrDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
} 