package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// FileItem represents a file in the repository, including its status and staging information.
// Name is the file path, Status indicates the Git status, Staged is true if the file is staged,
// and Selected can be used by UI clients to mark the file.
type FileItem struct {
	Name     string
	Status   string
	Staged   bool
	Selected bool
}

// IsGitRepository checks if the current directory is within a Git repository.
// It runs 'git rev-parse --git-dir' and returns true if no error occurs.
func IsGitRepository() bool {
	_, err := exec.Command("git", "rev-parse", "--git-dir").Output()
	return err == nil
}

// InitRepository initializes a new Git repository in the current directory.
// It runs 'git init' and returns any execution error.
func InitRepository() error {
	cmd := exec.Command("git", "init")
	return cmd.Run()
}

// GetModifiedFiles returns a slice of FileItem for all modified files in the working tree.
// It detects staged files via 'git diff --cached --name-status' and all changes via 'git status --porcelain'.
func GetModifiedFiles() ([]FileItem, error) {
	// Map of filenames that are staged
	stagedFiles := make(map[string]bool)
	stagedOutput, _ := exec.Command("git", "diff", "--cached", "--name-status").Output()
	if len(stagedOutput) > 0 {
		lines := strings.Split(strings.TrimSpace(string(stagedOutput)), "\n")
		for _, line := range lines {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				stagedFiles[parts[1]] = true
			}
		}
	}

	// Get status of all modified files
	output, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return nil, err
	}

	var files []FileItem
	if len(output) > 0 {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, line := range lines {
			if len(line) < 4 {
				continue
			}

			status := strings.TrimSpace(line[:2])
			filename := strings.TrimSpace(line[3:])

			// Handle filenames containing spaces
			if fields := strings.Fields(line); len(fields) >= 2 {
				filename = strings.Join(fields[1:], " ")
			}

			files = append(files, FileItem{
				Name:   filename,
				Status: status,
				Staged: stagedFiles[filename],
			})
		}
	}

	return files, nil
}

// GetBranches returns a slice of branch names and the currently checked-out branch.
// It runs 'git branch' and parses the output.
func GetBranches() ([]string, string) {
	output, err := exec.Command("git", "branch").Output()
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

// GetRemotes retrieves configured Git remotes and their URLs.
// It runs 'git remote -v' and returns unique remotes.
func GetRemotes() ([]string, error) {
	output, err := exec.Command("git", "remote", "-v").Output()
	if err != nil {
		return nil, err
	}

	var remotes []string
	seen := make(map[string]bool)
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

// Add stages the specified file using 'git add'.
func Add(filename string) error {
	return exec.Command("git", "add", filename).Run()
}

// Reset un-stages the specified file using 'git reset HEAD'.
func Reset(filename string) error {
	return exec.Command("git", "reset", "HEAD", filename).Run()
}

// Commit creates a new commit with the given message using 'git commit -m'.
func Commit(message string) error {
	return exec.Command("git", "commit", "-m", message).Run()
}

// Push sends committed changes to the remote repository using 'git push'.
func Push() error {
	return exec.Command("git", "push").Run()
}

// Checkout switches to the specified branch using 'git checkout'.
func Checkout(branch string) error {
	return exec.Command("git", "checkout", branch).Run()
}

// AddRemote adds a new remote with the given name and URL using 'git remote add'.
func AddRemote(name, url string) error {
	return exec.Command("git", "remote", "add", name, url).Run()
}

// RemoveRemote removes the specified remote using 'git remote remove'.
func RemoveRemote(name string) error {
	return exec.Command("git", "remote", "remove", name).Run()
}

// Fetch retrieves all updates from the remote repository using 'git fetch --all'.
func Fetch() error {
	output, err := exec.Command("git", "fetch", "--all").CombinedOutput()
	if err != nil {
		return fmt.Errorf("fetch failed: %v - %s", err, string(output))
	}
	return nil
}

// Pull fetches and integrates changes from the remote repository using 'git pull'.
func Pull() error {
	output, err := exec.Command("git", "pull").CombinedOutput()
	if err != nil {
		return fmt.Errorf("pull failed: %v - %s", err, string(output))
	}
	return nil
}

// CreateBranch creates and checks out a new branch using 'git checkout -b'.
func CreateBranch(name string) error {
	return exec.Command("git", "checkout", "-b", name).Run()
}

// DeleteBranch deletes the specified branch using 'git branch -d'.
func DeleteBranch(name string) error {
	return exec.Command("git", "branch", "-d", name).Run()
}

// DiscardChanges reverts changes to the specified file.
// If the file is untracked, it is removed; otherwise, changes are reset using 'git checkout --'.
func DiscardChanges(filename string) error {
	// Check if file is tracked
	if err := exec.Command("git", "ls-files", "--error-unmatch", filename).Run(); err != nil {
		// File is untracked: remove it
		if removeErr := os.Remove(filename); removeErr != nil {
			return fmt.Errorf("failed to remove untracked file: %w", removeErr)
		}
		return nil
	}

	// File is tracked: discard changes
	if err := exec.Command("git", "checkout", "--", filename).Run(); err != nil {
		return fmt.Errorf("failed to discard changes: %w", err)
	}
	return nil
}

// HasRemoteChanges checks if the local branch is behind its remote counterpart.
// It fetches updates and counts commits between HEAD and origin/branch.
func HasRemoteChanges(branch string) (bool, error) {
	if err := Fetch(); err != nil {
		return false, err
	}
	output, err := exec.Command("git", "rev-list", "--count", fmt.Sprintf("HEAD..origin/%s", branch)).Output()
	if err != nil {
		return false, err
	}
	count := strings.TrimSpace(string(output))
	return count != "0", nil
}

func LogsGraph() (string, error) {
	output, err := exec.Command("git", "log", "--graph", "--oneline", "--all").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get logs: %w", err)
	}
	return string(output), nil
}
