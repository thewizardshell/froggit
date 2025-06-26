package git

import (
	"strings"
)

// GetBranches public function
func GetBranches() ([]string, string) {
	return NewGitClient("").GetBranches()
}

// GetBranches method
func (g *GitClient) GetBranches() ([]string, string) {
	output, err := g.runGitCommand("branch")
	if err != nil {
		return nil, ""
	}

	var branches []string
	var current string
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "* ") {
			current = line[2:]
			branches = append(branches, current)
		} else {
			branches = append(branches, line)
		}
	}

	return branches, current
}

func CreateBranch(name string) error {
	return NewGitClient("").CreateBranch(name)
}

func (g *GitClient) CreateBranch(name string) error {
	_, err := g.runGitCommandCombinedOutput("checkout", "-b", name)
	return err
}

func DeleteBranch(name string) error {
	return NewGitClient("").DeleteBranch(name)
}

func (g *GitClient) DeleteBranch(name string) error {
	_, err := g.runGitCommandCombinedOutput("branch", "-d", name)
	return err
}

func Checkout(branch string) error {
	return NewGitClient("").Checkout(branch)
}

func (g *GitClient) Checkout(branch string) error {
	_, err := g.runGitCommandCombinedOutput("checkout", branch)
	return err
}

