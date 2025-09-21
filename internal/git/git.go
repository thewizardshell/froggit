package git

import (
	"os/exec"
)

func IsGitRepository() bool {
	_, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	return err == nil
}

func InitRepository() error {
	return exec.Command("git", "init").Run()
}
