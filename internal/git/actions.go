package git

import (
	"fmt"
	"strings"
)

func Add(filename string) error {
	return NewGitClient("").Add(filename)
}

func (g *GitClient) Add(filename string) error {
	_, err := g.runGitCommandCombinedOutput("add", filename)
	return err
}

func Reset(filename string) error {
	return NewGitClient("").Reset(filename)
}

func (g *GitClient) Reset(filename string) error {
	_, err := g.runGitCommandCombinedOutput("reset", "HEAD", filename)
	return err
}

func Commit(message string) error {
	return NewGitClient("").Commit(message)
}

func (g *GitClient) Commit(message string) error {
	output, err := g.runGitCommandCombinedOutput("commit", "-m", message)
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

func Merge(branch string) error {
	return NewGitClient("").Merge(branch)
}

func (g *GitClient) Merge(branch string) error {
	_, err := g.runGitCommandCombinedOutput("merge", branch)
	return err
}

func Rebase(branch string) error {
	return NewGitClient("").Rebase(branch)
}

func (g *GitClient) Rebase(branch string) error {
	_, err := g.runGitCommandCombinedOutput("rebase", branch)
	return err
}

func SaveStash(message string) error {
	return NewGitClient("").SaveStash(message)
}

func (g *GitClient) SaveStash(message string) error {
	_, err := g.runGitCommandCombinedOutput("stash", "push", "-m", message)
	return err
}

func StashPop() error {
	return NewGitClient("").StashPop()
}

func (g *GitClient) StashPop() error {
	_, err := g.runGitCommandCombinedOutput("stash", "pop")
	return err
}

func StashList() (string, error) {
	return NewGitClient("").StashList()
}

func (g *GitClient) StashList() (string, error) {
	output, err := g.runGitCommand("stash", "list")
	if err != nil {
		return "", fmt.Errorf("failed to list stashes: %w", err)
	}
	return string(output), nil
}

func HasRemoteChanges(branch string) (bool, error) {
	return NewGitClient("").HasRemoteChanges(branch)
}

func HasRemoteChangesWithFetch(branch string, doFetch bool) (bool, error) {
	return NewGitClient("").HasRemoteChangesWithFetch(branch, doFetch)
}

func (g *GitClient) HasRemoteChanges(branch string) (bool, error) {
	return g.HasRemoteChangesWithFetch(branch, true)
}

func (g *GitClient) HasRemoteChangesWithFetch(branch string, doFetch bool) (bool, error) {
	if doFetch {
		if err := g.FetchWithConfig(true); err != nil {
			return false, err
		}
	}

	output, err := g.runGitCommand("rev-list", "--count", fmt.Sprintf("HEAD..origin/%s", branch))
	if err != nil {
		return false, err
	}

	count := strings.TrimSpace(string(output))
	return count != "0", nil
}

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

func GetConflictFiles() ([]string, error) {
	output, err := NewGitClient("").runGitCommand("diff", "--name-only", "--diff-filter=U")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var files []string
	for _, l := range lines {
		if l != "" {
			files = append(files, l)
		}
	}
	return files, nil
}

func MergeContinue() error {
	_, err := NewGitClient("").runGitCommandCombinedOutput("merge", "--continue")
	return err
}

func MergeAbort() error {
	_, err := NewGitClient("").runGitCommandCombinedOutput("merge", "--abort")
	return err
}

func RebaseContinue() error {
	_, err := NewGitClient("").runGitCommandCombinedOutput("rebase", "--continue")
	return err
}

func RebaseAbort() error {
	_, err := NewGitClient("").runGitCommandCombinedOutput("rebase", "--abort")
	return err
}

func StashApply(stashRef string) error {
	return NewGitClient("").StashApply(stashRef)
}

func (g *GitClient) StashApply(stashRef string) error {
	_, err := g.runGitCommandCombinedOutput("stash", "apply", stashRef)
	return err
}

func StashDrop(stashRef string) error {
	return NewGitClient("").StashDrop(stashRef)
}

func (g *GitClient) StashDrop(stashRef string) error {
	_, err := g.runGitCommandCombinedOutput("stash", "drop", stashRef)
	return err
}

func StashShow(stashRef string) (string, error) {
	return NewGitClient("").StashShow(stashRef)
}

func (g *GitClient) StashShow(stashRef string) (string, error) {
	output, err := g.runGitCommand("stash", "show", "-p", stashRef)
	if err != nil {
		return "", fmt.Errorf("failed to show stash: %w", err)
	}
	return string(output), nil
}

func GetStashRef(stashLine string) string {
	stashLine = strings.TrimSpace(stashLine)
	if stashLine == "" {
		return "stash@{0}"
	}
	parts := strings.Split(stashLine, ":")
	if len(parts) > 0 {
		ref := strings.TrimSpace(parts[0])
		if ref == "" {
			return "stash@{0}"
		}
		return ref
	}
	return "stash@{0}"
}

func HasCommitsToush() (bool, error) {
	return NewGitClient("").HasCommitsToush()
}

func GetStagedDiff() (string, error) {
	return NewGitClient("").GetStagedDiff()
}

func (g *GitClient) GetStagedDiff() (string, error) {
	output, err := g.runGitCommand("diff", "--cached")
	if err != nil {
		return "", fmt.Errorf("failed to get staged diff: %w", err)
	}
	return string(output), nil
}

func (g *GitClient) HasCommitsToush() (bool, error) {
	output, err := g.runGitCommand("rev-list", "--count", "HEAD")
	if err != nil {
		return false, nil
	}

	localCommits := strings.TrimSpace(string(output))
	if localCommits == "0" {
		return false, nil
	}

	output, err = g.runGitCommand("rev-list", "--count", "HEAD", "^@{u}")
	if err != nil {
		return localCommits != "0", nil
	}

	ahead := strings.TrimSpace(string(output))
	return ahead != "0", nil
}
