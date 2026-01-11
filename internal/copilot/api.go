package copilot

import (
	"strings"
)

// GenerateCommitMessage generates a concise git commit message based on the provided diff.
// It follows conventional commit formats and ensures the message is specific and within character limits.
// It returns the generated commit message or an error if the operation fails.
func GenerateCommitMessage(diff string) (string, error) {
	client, err := GetClient()
	if err != nil {
		return "", err
	}

	if len(diff) > 8000 {
		diff = diff[:8000] + "\n... (truncated)"
	}

	prompt := `Generate a concise git commit message for these changes.
Rules:
- Use conventional commits format: type(scope): description
- Types: feat, fix, docs, style, refactor, test, chore
- First line max 72 characters
- Be specific about what changed
- Only output the commit message, nothing else

Changes:
` + diff

	response, err := client.Chat([]Message{
		{Role: "user", Content: prompt},
	})
	if err != nil {
		return "", err
	}

	response = strings.TrimSpace(response)
	response = strings.Trim(response, "`\"'")

	return response, nil
}
