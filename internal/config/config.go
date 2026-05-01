package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Token      string           `toml:"token"`
	BaseURL    string           `toml:"base_url"`
	Appearance AppearanceConfig `toml:"appearance"`
	Keys       KeyConfig        `toml:"keys"`
}

type AppearanceConfig struct {
	AccentColor     string `toml:"accent_color"`
	BackgroundColor string `toml:"background_color"`
	BorderColor     string `toml:"border_color"`
	SuccessColor    string `toml:"success_color"`
	ErrorColor      string `toml:"error_color"`
}

type KeyConfig struct {
	Quit        string `toml:"quit"`
	Back        string `toml:"back"`
	New         string `toml:"new"`
	Select      string `toml:"select"`
	Upvote      string `toml:"upvote"`
	Downvote    string `toml:"downvote"`
	Refresh     string `toml:"refresh"`
	Login       string `toml:"login"`
	Feed        string `toml:"feed"`
	Communities string `toml:"communities"`
}

const DefaultBaseURL = "http://localhost:9001/api/v1"

var ConfigPathOverride string

func GetConfigPath() (string, error) {
	if ConfigPathOverride != "" {
		return ConfigPathOverride, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(home, ".config", "ditto-cli")
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", err
	}
	return filepath.Join(path, "config.toml"), nil
}

func DefaultConfig() *Config {
	return &Config{
		BaseURL: DefaultBaseURL,
		Appearance: AppearanceConfig{
			AccentColor:     "#7D56F4",
			BackgroundColor: "#1A1B26",
			BorderColor:     "#874BFD",
			SuccessColor:    "#01BE85",
			ErrorColor:      "#FF0000",
		},
		Keys: KeyConfig{
			Quit:        "q",
			Back:        "backspace",
			New:         "n",
			Select:      "y",
			Upvote:      "a",
			Downvote:    "z",
			Refresh:     "r",
			Login:       ":L",
			Feed:        ":F",
			Communities: ":C",
		},
	}
}

func LoadConfig() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	cfg := DefaultConfig()
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// Migration check: if config.json exists, maybe we should migrate it?
		// For now, just save default.
		if err := cfg.Save(); err != nil {
			return nil, err
		}
		return cfg, nil
	}

	if _, err := toml.DecodeFile(path, cfg); err != nil {
		return nil, err
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}

	return cfg, nil
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

	// For debugging token persistence
	if os.Getenv("DEBUG_CONFIG") != "" {
		fmt.Fprintf(os.Stderr, "Saving config to: %s\n", path)
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(c)
}
