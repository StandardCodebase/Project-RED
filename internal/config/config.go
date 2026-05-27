package config

import (
	"encoding/json"
	"os"
)

type RemoteSync struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
}

type Config struct {
	Addr        string       `json:"addr"`
	SiteName    string       `json:"siteName"`
	DataDir     string       `json:"dataDir"`
	SourceURL   string       `json:"sourceURL"`
	SourceType  string       `json:"sourceType"`  // e.g., "tar.gz", "zip"
	AdminToken  string       `json:"adminToken"`  // NEW: Security Token
	StartupSync []RemoteSync `json:"startupSync"` // Ported from Legacy Gateway
}

func Default() Config {
	return Config{
		Addr:     ":8080",
		SiteName: "RED Engine",
		DataDir:  "./data",
	}
}

func Load(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// NEW: Allows the application to save changes back to disk
func (c *Config) Save(path string) error {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}
