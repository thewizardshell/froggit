package gh

import (
	"os/exec"
)

type GhClient struct {
}

func NewGhClient() *GhClient {
	return &GhClient{}
}

func (g *GhClient) runGhCommand(args ...string) ([]byte, error) {
	cmd := exec.Command("gh", args...)
	return cmd.Output()
}

func (g *GhClient) runGhCommandCombinedOutput(args ...string) ([]byte, error) {
	cmd := exec.Command("gh", args...)
	return cmd.CombinedOutput()
}

// ListRepositories returns a list of repositories for the authenticated user.
// It uses 'gh repo list' with --json for easier parsing.
func (g *GhClient) ListRepositories() ([]byte, error) {
	return g.runGhCommand("repo", "list", "--json", "name,owner,sshUrl,description", "--limit", "100")
}
