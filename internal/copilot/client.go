package copilot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type CopilotToken struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	Endpoints struct {
		API string `json:"api"`
	} `json:"endpoints"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Client struct {
	httpClient   *http.Client
	oauthToken   string
	copilotToken *CopilotToken
	mu           sync.RWMutex
}

var (
	instance *Client
	once     sync.Once
	initErr  error
)

func GetClient() (*Client, error) {
	once.Do(func() {
		oauth, err := LoadAuthToken()
		if err != nil {
			initErr = err
			return
		}
		instance = &Client{
			httpClient: &http.Client{Timeout: 60 * time.Second},
			oauthToken: oauth,
		}
	})
	return instance, initErr
}

func (c *Client) getValidToken() (string, error) {
	c.mu.RLock()
	if c.copilotToken != nil && time.Now().Unix() < c.copilotToken.ExpiresAt-300 {
		token := c.copilotToken.Token
		c.mu.RUnlock()
		return token, nil
	}
	c.mu.RUnlock()
	return c.refreshToken()
}

func (c *Client) refreshToken() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.copilotToken != nil && time.Now().Unix() < c.copilotToken.ExpiresAt-300 {
		return c.copilotToken.Token, nil
	}

	req, _ := http.NewRequest("GET", "https://api.github.com/copilot_internal/v2/token", nil)
	req.Header = http.Header{
		"Authorization": []string{"Token " + c.oauthToken},
		"Accept":        []string{"application/json"},
		"User-Agent":    []string{"GithubCopilot/1.155.0"},
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("token refresh failed: %d", resp.StatusCode)
	}

	var token CopilotToken
	if err := json.Unmarshal(bodyBytes, &token); err != nil {
		return "", fmt.Errorf("failed to decode token: %w", err)
	}

	c.copilotToken = &token
	return token.Token, nil
}

func (c *Client) Chat(messages []Message) (string, error) {
	token, err := c.getValidToken()
	if err != nil {
		return "", err
	}

	baseUrl := "https://api.githubcopilot.com"
	if c.copilotToken.Endpoints.API != "" {
		baseUrl = c.copilotToken.Endpoints.API
	}

	body := map[string]any{
		"model":      "gpt-4o",
		"messages":   messages,
		"stream":     false,
		"max_tokens": 500,
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", baseUrl+"/chat/completions", bytes.NewBuffer(jsonBody))
	req.Header = http.Header{
		"Authorization":  []string{"Bearer " + token},
		"Content-Type":   []string{"application/json"},
		"User-Agent":     []string{"GithubCopilot/1.155.0"},
		"Editor-Version": []string{"vscode/1.95.3"},
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("copilot error: %d", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("empty response from copilot")
}
