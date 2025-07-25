package git

import (
	"os/exec"
	"strings"
)

// GitClient holds repo info for running commands.
type GitClient struct {
	RepoPath string
}

// NewGitClient creates a new GitClient instance.
// repoPath can be empty string to use current working directory.
func NewGitClient(repoPath string) *GitClient {
	if repoPath == "" {
		// Find the git repository root
		if root, err := findGitRoot(); err == nil {
			repoPath = root
		}
	}
	return &GitClient{
		RepoPath: repoPath,
	}
}

// findGitRoot finds the root directory of the git repository
func findGitRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func (g *GitClient) runGitCommand(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	if g.RepoPath != "" {
		cmd.Dir = g.RepoPath
	}
	return cmd.Output()
}

func (g *GitClient) runGitCommandCombinedOutput(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	if g.RepoPath != "" {
		cmd.Dir = g.RepoPath
	}
	return cmd.CombinedOutput()
}
