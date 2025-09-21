package git

import (
	"fmt"
  "os/exec"	
	"os"
	"strings"
)

type FileItem struct {
	Name     string
	Status   string
	Staged   bool
	Selected bool
}

func GetModifiedFiles() ([]FileItem, error) {
	return NewGitClient("").GetModifiedFiles()
}

func (g *GitClient) GetModifiedFiles() ([]FileItem, error) {
	stagedFiles := make(map[string]bool)

	stagedOutput, _ := g.runGitCommand("diff", "--cached", "--name-status")
	if len(stagedOutput) > 0 {
		lines := strings.Split(strings.TrimSpace(string(stagedOutput)), "\n")
		for _, line := range lines {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				stagedFiles[parts[1]] = true
			}
		}
	}

	output, err := g.runGitCommand("status", "--porcelain")
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

func DiscardChanges(filename string) error {
	return NewGitClient("").DiscardChanges(filename)
}

func (g *GitClient) DiscardChanges(filename string) error {
	cmd := exec.Command("git", "ls-files", "--error-unmatch", filename)
	if g.RepoPath != "" {
		cmd.Dir = g.RepoPath
	}
	if err := cmd.Run(); err != nil {
		if err := os.Remove(filename); err != nil {
			return fmt.Errorf("failed to remove untracked file: %w", err)
		}
		return nil
	}

	_, err := g.runGitCommandCombinedOutput("checkout", "--", filename)
	return err
}
