package github

import (
	"context"
	"fmt"
	"sort"

	gh "github.com/google/go-github/v68/github"
)

type ForkInfo struct {
	Owner         string
	Repo          string
	DefaultBranch string
	PushedAt      gh.Timestamp
	HTMLURL       string
}

func FetchForks(ctx context.Context, client *gh.Client, owner, repo string, limit int) ([]ForkInfo, *gh.Repository, error) {
	upstream, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch repository %s/%s: %w", owner, repo, err)
	}

	var allForks []*gh.Repository
	opts := &gh.RepositoryListForksOptions{
		Sort:        "newest",
		ListOptions: gh.ListOptions{PerPage: 100},
	}

	for {
		forks, resp, err := client.Repositories.ListForks(ctx, owner, repo, opts)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to list forks: %w", err)
		}
		allForks = append(allForks, forks...)
		if resp.NextPage == 0 || len(allForks) >= limit {
			break
		}
		opts.Page = resp.NextPage
	}

	// Sort by most recently pushed
	sort.Slice(allForks, func(i, j int) bool {
		return allForks[i].GetPushedAt().After(allForks[j].GetPushedAt().Time)
	})

	// Apply limit
	if len(allForks) > limit {
		allForks = allForks[:limit]
	}

	var results []ForkInfo
	for _, f := range allForks {
		branch := f.GetDefaultBranch()
		if branch == "" {
			branch = "main"
		}

		results = append(results, ForkInfo{
			Owner:         f.GetOwner().GetLogin(),
			Repo:          f.GetName(),
			DefaultBranch: branch,
			PushedAt:      f.GetPushedAt(),
			HTMLURL:       f.GetHTMLURL(),
		})
	}

	return results, upstream, nil
}
