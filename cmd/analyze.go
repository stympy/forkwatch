package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stympy/forkwatch/internal/analysis"
	ghclient "github.com/stympy/forkwatch/internal/github"
	"github.com/stympy/forkwatch/internal/output"
)

var (
	minAhead int
	limit    int
	jsonOut  bool
	patchOut bool
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze owner/repo",
	Short: "Analyze forks of a repository",
	Long:  `Fetches forks of the given repository, compares each to upstream, and clusters them by files changed to reveal convergent patches.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runAnalyze,
}

func init() {
	analyzeCmd.Flags().IntVar(&minAhead, "min-ahead", 1, "Minimum commits ahead to consider")
	analyzeCmd.Flags().IntVar(&limit, "limit", 100, "Max forks to analyze (sorted by most recently pushed)")
	analyzeCmd.Flags().BoolVar(&jsonOut, "json", false, "Output as JSON")
	analyzeCmd.Flags().BoolVar(&patchOut, "patch", false, "Output a unified diff suitable for git apply")
	analyzeCmd.MarkFlagsMutuallyExclusive("json", "patch")
	rootCmd.AddCommand(analyzeCmd)
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	parts := strings.SplitN(args[0], "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("repository must be in owner/repo format")
	}
	owner, repo := parts[0], parts[1]

	ctx := context.Background()

	client, err := ghclient.NewClient(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Fetching forks of %s/%s...\n", owner, repo)

	forks, upstream, err := ghclient.FetchForks(ctx, client, owner, repo, limit)
	if err != nil {
		return err
	}

	if len(forks) == 0 {
		fmt.Println("No active forks found.")
		return nil
	}

	fmt.Fprintf(os.Stderr, "Found %d active forks, comparing to upstream...\n", len(forks))

	upstreamBranch := upstream.GetDefaultBranch()
	if upstreamBranch == "" {
		upstreamBranch = "main"
	}

	var comparisons []*ghclient.ForkComparison
	for i, fork := range forks {
		fmt.Fprintf(os.Stderr, "Analyzing fork %d/%d: %s...\n", i+1, len(forks), fork.Owner)

		comp, err := ghclient.CompareFork(ctx, client, owner, repo, upstreamBranch, fork)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Warning: %v\n", err)
			continue
		}
		if comp == nil {
			continue
		}
		if comp.AheadBy < minAhead {
			continue
		}
		comparisons = append(comparisons, comp)
	}

	totalForks := upstream.GetForksCount()
	result := analysis.Cluster(comparisons, owner, repo, totalForks)

	if jsonOut {
		return output.PrintJSON(result)
	}
	if patchOut {
		output.PrintPatch(result)
		return nil
	}

	output.PrintTable(result)
	return nil
}
