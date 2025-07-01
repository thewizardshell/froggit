package gh

import (
	"encoding/json"
	"fmt"
)

// Repository holds details about a GitHub repository.
type Repository struct {
	Name  string `json:"name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
	SshUrl      string `json:"sshUrl"`
	Description string `json:"description"`
}

// ListUserRepositories fetches and returns a slice of Repository.
func ListUserRepositories(client *GhClient) ([]Repository, error) {
	out, err := client.ListRepositories()
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}
	var repos []Repository
	if err := json.Unmarshal(out, &repos); err != nil {
		return nil, fmt.Errorf("failed to parse repositories: %w", err)
	}
	return repos, nil
}

// CloneRepository clones the given repository using 'gh repo clone'.
func CloneRepository(client *GhClient, repoFullName string, destPath string) error {
	args := []string{"repo", "clone", repoFullName}
	if destPath != "" {
		args = append(args, destPath)
	}
	out, err := client.runGhCommandCombinedOutput(args...)
	if err != nil {
		return fmt.Errorf("failed to clone repository: %v\nOutput: %s", err, string(out))
	}
	return nil
}
