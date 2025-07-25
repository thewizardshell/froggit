package git

import (
	"os/exec"
)

// IsGitRepository returns true if current folder is within a Git repository
func IsGitRepository() bool {
	_, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	return err == nil
}

// InitRepository runs 'git init' in the current directory
func InitRepository() error {
	return exec.Command("git", "init").Run()
}
