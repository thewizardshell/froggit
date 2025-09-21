package git

import (
	"fmt"
	"strings"
)

func GetRemotes() ([]string, error) {
	return NewGitClient("").GetRemotes()
}

func (g *GitClient) GetRemotes() ([]string, error) {
	output, err := g.runGitCommand("remote", "-v")
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var remotes []string
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 2 && !seen[parts[0]] {
			remotes = append(remotes, fmt.Sprintf("%s -> %s", parts[0], parts[1]))
			seen[parts[0]] = true
		}
	}

	return remotes, nil
}

func AddRemote(name, url string) error {
	return NewGitClient("").AddRemote(name, url)
}

func (g *GitClient) AddRemote(name, url string) error {
	_, err := g.runGitCommandCombinedOutput("remote", "add", name, url)
	return err
}

func RemoveRemote(name string) error {
	return NewGitClient("").RemoveRemote(name)
}

func (g *GitClient) RemoveRemote(name string) error {
	_, err := g.runGitCommandCombinedOutput("remote", "remove", name)
	return err
}

func Fetch() error {
	return NewGitClient("").Fetch()
}

func (g *GitClient) Fetch() error {
	return g.FetchWithConfig(false)
}

func FetchWithConfig(silent bool) error {
	return NewGitClient("").FetchWithConfig(silent)
}

func (g *GitClient) FetchWithConfig(silent bool) error {
	mu.Lock()
	if operationInProgress {
		mu.Unlock()
		if !silent {
			return fmt.Errorf("another git operation is already in progress")
		}
		return nil
	}
	operationInProgress = true
	mu.Unlock()

	defer func() {
		mu.Lock()
		operationInProgress = false
		mu.Unlock()
	}()

	output, err := g.runGitCommandCombinedOutput("fetch")
	if err != nil {
		if !silent {
			return fmt.Errorf("fetch failed: %v - %s", err, string(output))
		}
		return err
	}
	return nil
}

func Pull() error {
	return NewGitClient("").Pull()
}

func (g *GitClient) Pull() error {
	mu.Lock()
	if operationInProgress {
		mu.Unlock()
		return fmt.Errorf("another git operation is already in progress")
	}
	operationInProgress = true
	mu.Unlock()

	defer func() {
		mu.Lock()
		operationInProgress = false
		mu.Unlock()
	}()

	output, err := g.runGitCommandCombinedOutput("pull")
	if err != nil {
		return fmt.Errorf("pull failed: %v - %s", err, string(output))
	}
	return nil
}

func Push() error {
	return NewGitClient("").Push()
}

func (g *GitClient) Push() error {
	return g.PushWithBranch("")
}

func PushWithBranch(defaultBranch string) error {
	return NewGitClient("").PushWithBranch(defaultBranch)
}

func (g *GitClient) PushWithBranch(defaultBranch string) error {
	mu.Lock()
	if operationInProgress {
		mu.Unlock()
		return fmt.Errorf("another git operation is already in progress")
	}
	operationInProgress = true
	mu.Unlock()

	defer func() {
		mu.Lock()
		operationInProgress = false
		mu.Unlock()
	}()

	output, err := g.runGitCommandCombinedOutput("push")
	if err == nil {
		return nil
	}

	if strings.Contains(string(output), "set the remote as upstream") ||
		strings.Contains(string(output), "have no upstream branch") ||
		strings.Contains(string(output), "no upstream branch") {

		branch := defaultBranch
		if branch == "" {
			branchOut, branchErr := g.runGitCommand("branch", "--show-current")
			if branchErr != nil {
				return fmt.Errorf("push failed and could not determine current branch: %v", branchErr)
			}
			branch = strings.TrimSpace(string(branchOut))
		}
		if branch == "" {
			return fmt.Errorf("push failed and could not determine current branch name")
		}

		output2, err2 := g.runGitCommandCombinedOutput("push", "--set-upstream", "origin", branch)
		if err2 != nil {
			return fmt.Errorf("push failed and could not set upstream: %v - %s", err2, string(output2))
		}
		return nil
	}
	return fmt.Errorf("push failed: %v - %s", err, string(output))
}
