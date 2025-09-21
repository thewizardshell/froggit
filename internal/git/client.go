package git

import (
	"os/exec"
	"strings"
)

type GitClient struct {
	RepoPath string
}

func NewGitClient(repoPath string) *GitClient {
	if repoPath == "" {
		if root, err := findGitRoot(); err == nil {
			repoPath = root
		}
	}
	return &GitClient{
		RepoPath: repoPath,
	}
}

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
