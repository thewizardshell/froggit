package git

import (
	"os/exec"
)

// GitClient holds repo info for running commands.
type GitClient struct {
	RepoPath string
}

// NewGitClient creates a new GitClient instance.
// repoPath can be empty string to use current working directory.
func NewGitClient(repoPath string) *GitClient {
	return &GitClient{
		RepoPath: repoPath,
	}
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

