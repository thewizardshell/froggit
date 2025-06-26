package git

import (
	"fmt"
	"strings"
)

// Add stages a file
func Add(filename string) error {
	return NewGitClient("").Add(filename)
}

func (g *GitClient) Add(filename string) error {
	_, err := g.runGitCommandCombinedOutput("add", filename)
	return err
}

// Reset unstages a file
func Reset(filename string) error {
	return NewGitClient("").Reset(filename)
}

func (g *GitClient) Reset(filename string) error {
	_, err := g.runGitCommandCombinedOutput("reset", "HEAD", filename)
	return err
}

// Commit creates a commit
func Commit(message string) error {
	return NewGitClient("").Commit(message)
}

func (g *GitClient) Commit(message string) error {
	_, err := g.runGitCommandCombinedOutput("commit", "-m", message)
	return err
}

// HasRemoteChanges checks if local branch is behind remote
func HasRemoteChanges(branch string) (bool, error) {
	return NewGitClient("").HasRemoteChanges(branch)
}

func (g *GitClient) HasRemoteChanges(branch string) (bool, error) {
	if err := g.Fetch(); err != nil {
		return false, err
	}

	output, err := g.runGitCommand("rev-list", "--count", fmt.Sprintf("HEAD..origin/%s", branch))
	if err != nil {
		return false, err
	}

	count := strings.TrimSpace(string(output))
	return count != "0", nil
}

// LogsGraph returns a graph view of logs
func LogsGraph() (string, error) {
	return NewGitClient("").LogsGraph()
}

func (g *GitClient) LogsGraph() (string, error) {
	output, err := g.runGitCommand("log", "--graph", "--oneline", "--all")
	if err != nil {
		return "", fmt.Errorf("failed to get logs: %w", err)
	}
	return string(output), nil
}

