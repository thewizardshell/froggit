package config

import (
	"os"
	"testing"
)

func TestLoadConfig_DefaultsAndParsing(t *testing.T) {
	// create temp file
	content := "ui:\n  branding: true\n  position: \"\"\ngit:\n  autofetch: true\n"
	f, err := os.CreateTemp("", "cfg-*.yml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()

	cfg, err := LoadConfig(f.Name())
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if cfg.Ui.Branding != true {
		t.Fatalf("expected Ui.Branding true, got %v", cfg.Ui.Branding)
	}
	// position should default to left when empty
	if cfg.Ui.Position != "left" {
		t.Fatalf("expected Ui.Position 'left', got %q", cfg.Ui.Position)
	}
	// default branch should be main
	if cfg.Git.DefaultBranch != "main" {
		t.Fatalf("expected DefaultBranch 'main', got %q", cfg.Git.DefaultBranch)
	}
}
