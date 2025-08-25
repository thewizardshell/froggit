package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Ui  UiConfig  `yaml:"ui"`
	Git GitConfig `yaml:"git"`
}

type UiConfig struct {
	Branding bool   `yaml:"branding"`
	Position string `yaml:"position"` // "left", "center", "right"
}

type GitConfig struct {
	AutoFetch     bool   `yaml:"autofetch"`
	DefaultBranch string `yaml:"defaultbranch"`
}

func LoadConfig(path string) (Config, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		return Config{}, err
	}

	if cfg.Ui.Position == "" {
		cfg.Ui.Position = "left"
	}
	if cfg.Git.DefaultBranch == "" {
		cfg.Git.DefaultBranch = "main"
	}

	return cfg, nil
}
