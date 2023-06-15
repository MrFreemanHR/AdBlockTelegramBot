package config

import (
	"encoding/json"
	"os"
)

var CurrentConfig Config

type Config struct {
	// Logger
	VerbosityLevel uint `json:"verbosity_level"`
	// SQLite DSN
	SQLiteDSN string `json:"sqlite_dsn"`
	// Telegram
	Token string `json:"token"`
	// tdlib telegram
	APIid            string `json:"api_id"`
	APIHash          string `json:"api_hash"`
	DatabasePassword string `json:"database_password"`
	// Locales
	LocalesFolder string `json:"locales_folder"`
	DefaultLocale string `json:"default_locale"`
}

func ParseConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
