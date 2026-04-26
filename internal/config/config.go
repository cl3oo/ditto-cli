package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Token   string `json:"token"`
	BaseURL string `json:"base_url"`
}

const DefaultBaseURL = "http://api.ditto.local"

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(home, ".config", "ditto-cli")
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", err
	}
	return filepath.Join(path, "config.json"), nil
}

func LoadConfig() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}
	
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{BaseURL: DefaultBaseURL}, nil
		}
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}
	
	return &cfg, nil
}

func (c *Config) UpdateToken(token string) error {
	c.Token = token
	return c.Save()
}

func (c *Config) Save() error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(c)
}
