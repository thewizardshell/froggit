package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	Ui  UiConfig  `yaml:"ui"`
	Git GitConfig `yaml:"git"`
}

type UiConfig struct {
	Branding bool   `yaml:"branding"`
	Position string `yaml:"position"`
}

type GitConfig struct {
	AutoFetch     bool   `yaml:"autofetch"`
	DefaultBranch string `yaml:"defaultbranch"`
}

func LoadConfig(filename string) (Config, error) {
	exeDir := getExecutableDir()
	configPath := filepath.Join(exeDir, filename)

	f, err := os.ReadFile(configPath)
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

func getExecutableDir() string {
	ex, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(ex)
}
