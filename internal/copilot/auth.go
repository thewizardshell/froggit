package copilot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetConfigDir returns the configuration directory based on the XDG specification or OS defaults.
// On Unix-like systems, it defaults to ~/.config.
// On Windows, it uses the LOCALAPPDATA environment variable.

func GetConfigDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return xdg
	}
	if runtime.GOOS == "windows" {
		return os.Getenv("LOCALAPPDATA")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config")
}

// LoadAuthToken attempts to load the GitHub Copilot authentication token
// from known configuration files. It checks both "hosts.json" and "apps.json"
// for the token associated with "github.com". If found, it returns the token;
// otherwise, it returns an error indicating that the token was not found.

func LoadAuthToken() (string, error) {
	configDir := GetConfigDir()

	files := []string{"hosts.json",
		"apps.json",
	}

	for _, file := range files {
		path := filepath.Join(configDir, "github-copilot", file)
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var hosts map[string]struct {
			OAuthToken string `json:"oauth_token"`
		}
		if json.Unmarshal(data, &hosts) != nil {

			continue
		}

		for key, val := range hosts {

			if strings.HasPrefix(key, "github.com") && val.OAuthToken != "" {
				return val.OAuthToken, nil
			}
		}
	}

	return "", fmt.Errorf("copilot token not found - need copilot.vim, copilot.lua, or VS Code with Copilot")

}

// IsAvailable checks if the GitHub Copilot authentication token is available
// by attempting to load it using LoadAuthToken. It returns true if the token
// is found, and false otherwise.
func IsAvailable() bool {
	_, err := LoadAuthToken()
	return err == nil
}
