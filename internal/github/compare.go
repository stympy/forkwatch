package github

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	gh "github.com/google/go-github/v68/github"
)

var botAccounts = map[string]bool{
	"dependabot[bot]":  true,
	"dependabot":       true,
	"renovate[bot]":    true,
	"renovate":         true,
	"greenkeeper[bot]": true,
	"snyk-bot":         true,
	"depfu[bot]":       true,
}

var boringFiles = map[string]bool{
	"package-lock.json": true,
	"yarn.lock":         true,
	"Gemfile.lock":      true,
	"go.sum":            true,
	"pnpm-lock.yaml":    true,
	"Cargo.lock":        true,
	"poetry.lock":       true,
	"composer.lock":     true,
}

type ForkComparison struct {
	Fork           ForkInfo
	AheadBy        int
	CommitMessages []string
	FilesChanged   []FileChange
}

type FileChange struct {
	Filename  string
	Additions int
	Deletions int
	Patch     string
}

func CompareFork(ctx context.Context, client *gh.Client, upstreamOwner, upstreamRepo, upstreamBranch string, fork ForkInfo) (*ForkComparison, error) {
	head := fmt.Sprintf("%s:%s", fork.Owner, fork.DefaultBranch)

	comparison, resp, err := client.Repositories.CompareCommits(ctx, upstreamOwner, upstreamRepo, upstreamBranch, head, nil)
	if resp != nil && resp.Rate.Remaining < 10 {
		return nil, fmt.Errorf("rate limit nearly exhausted (%d remaining, resets at %s) — stopping to avoid 403s",
			resp.Rate.Remaining, resp.Rate.Reset.Time.Format("15:04:05"))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to compare %s/%s: %w", fork.Owner, fork.Repo, err)
	}

	aheadBy := comparison.GetAheadBy()
	if aheadBy == 0 {
		return nil, nil
	}

	// Check for bot-only commits
	allBots := true
	var messages []string
	for _, c := range comparison.Commits {
		author := c.GetCommit().GetAuthor().GetName()
		if !botAccounts[author] {
			allBots = false
		}
		msg := strings.Split(c.GetCommit().GetMessage(), "\n")[0]
		messages = append(messages, msg)
	}
	if allBots && len(comparison.Commits) > 0 {
		return nil, nil
	}

	// Check for boring-only file changes
	var files []FileChange
	allBoring := true
	for _, f := range comparison.Files {
		name := f.GetFilename()
		baseName := filepath.Base(name)
		if !boringFiles[baseName] && !isCI(name) {
			allBoring = false
		}
		files = append(files, FileChange{
			Filename:  name,
			Additions: f.GetAdditions(),
			Deletions: f.GetDeletions(),
			Patch:     f.GetPatch(),
		})
	}
	if allBoring && len(files) > 0 {
		return nil, nil
	}

	return &ForkComparison{
		Fork:           fork,
		AheadBy:        aheadBy,
		CommitMessages: messages,
		FilesChanged:   files,
	}, nil
}

func isCI(path string) bool {
	return strings.HasPrefix(path, ".github/") ||
		strings.HasPrefix(path, ".circleci/") ||
		strings.HasPrefix(path, ".travis") ||
		strings.HasPrefix(path, ".gitlab-ci")
}
