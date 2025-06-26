package git

import (
	"os/exec"
)

// IsGitRepository returns true if current folder is a Git repository
func IsGitRepository() bool {
	_, err := exec.Command("git", "rev-parse", "--git-dir").Output()
	return err == nil
}

// InitRepository runs 'git init' in the current directory
func InitRepository() error {
	return exec.Command("git", "init").Run()
}

