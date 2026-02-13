package github

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	gh "github.com/google/go-github/v68/github"
)

func NewClient(ctx context.Context) (*gh.Client, error) {
	token, err := getGHToken()
	if err != nil {
		return nil, err
	}
	return gh.NewClient(nil).WithAuthToken(token), nil
}

func getGHToken() (string, error) {
	cmd := exec.Command("gh", "auth", "token")
	out, err := cmd.Output()
	if err != nil {
		if execErr, ok := err.(*exec.Error); ok && execErr.Err == exec.ErrNotFound {
			return "", fmt.Errorf("GitHub CLI (gh) is not installed. Install it from https://cli.github.com/")
		}
		return "", fmt.Errorf("failed to get GitHub token — run 'gh auth login' first: %w", err)
	}
	token := strings.TrimSpace(string(out))
	if token == "" {
		return "", fmt.Errorf("GitHub CLI returned empty token — run 'gh auth login' first")
	}
	return token, nil
}
