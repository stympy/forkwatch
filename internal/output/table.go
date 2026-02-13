package output

import (
	"fmt"
	"strings"

	"github.com/stympy/forkwatch/internal/analysis"
)

const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorDim    = "\033[2m"
)

func PrintTable(result *analysis.AnalysisResult) {
	fmt.Printf("\n%s%s%s/%s%s\n", colorBold, colorCyan, result.UpstreamOwner, result.UpstreamRepo, colorReset)
	fmt.Printf("%sForks: %d total, %d analyzed, %d with meaningful changes%s\n\n",
		colorDim, result.TotalForks, result.AnalyzedForks, result.ActiveForks, colorReset)

	if len(result.Clusters) == 0 {
		fmt.Println("No meaningful fork activity found.")
		return
	}

	// Show convergence clusters
	for _, cluster := range result.Clusters {
		convergenceLabel := ""
		if cluster.Convergence >= 2 {
			convergenceLabel = fmt.Sprintf(" %s%s(%d forks converge here)%s",
				colorBold, colorYellow, cluster.Convergence, colorReset)
		}

		fmt.Printf("%s%s%s%s\n", colorBold, cluster.Filename, colorReset, convergenceLabel)

		if cluster.PatchGroups != nil && len(cluster.PatchGroups.Groups) > 0 {
			printPatchGroups(cluster)
		} else {
			printForkList(cluster.Forks)
		}

		fmt.Println(strings.Repeat("─", 60))
	}
}

func printPatchGroups(cluster analysis.FileCluster) {
	for i, group := range cluster.PatchGroups.Groups {
		if len(group.Forks) > 1 {
			// Multi-fork group: show the shared diff then list owners
			if i == 0 {
				fmt.Printf("\n  %sMost common change pattern:%s\n", colorBold, colorReset)
			} else {
				fmt.Printf("\n  %sShared by %d forks:%s\n", colorDim, len(group.Forks), colorReset)
			}
			printDiff(group.Patch)
			var owners []string
			for _, f := range group.Forks {
				owners = append(owners, f.Owner)
			}
			fmt.Printf("  %s%s%s\n", colorCyan, strings.Join(owners, ", "), colorReset)
		} else {
			// Single-fork: show owner, stats, and their diff
			f := group.Forks[0]
			msg := ""
			if len(f.CommitMessages) > 0 {
				msg = f.CommitMessages[0]
				if len(msg) > 50 {
					msg = msg[:50] + "..."
				}
				msg = " — " + msg
			}
			fmt.Printf("\n  %s%s%s %s+%d%s %s-%d%s%s\n",
				colorCyan, f.Owner, colorReset,
				colorGreen, f.Additions, colorReset,
				colorRed, f.Deletions, colorReset,
				msg)
			printDiff(group.Patch)
		}
	}
}

func printDiff(patch string) {
	if patch == "" {
		return
	}
	diff := analysis.TruncateDiff(patch, 10)
	for _, line := range strings.Split(diff, "\n") {
		color := colorDim
		if strings.HasPrefix(line, "+") {
			color = colorGreen
		} else if strings.HasPrefix(line, "-") {
			color = colorRed
		}
		fmt.Printf("    %s%s%s\n", color, line, colorReset)
	}
}

func printForkList(forks []analysis.ForkSummary) {
	for _, fork := range forks {
		stats := fmt.Sprintf("%s+%d%s %s-%d%s",
			colorGreen, fork.Additions, colorReset,
			colorRed, fork.Deletions, colorReset)

		fmt.Printf("  %s%-20s%s %s (%d commits ahead)\n",
			colorCyan, fork.Owner, colorReset, stats, fork.AheadBy)

		if len(fork.CommitMessages) > 0 {
			msg := fork.CommitMessages[0]
			if len(msg) > 72 {
				msg = msg[:72] + "..."
			}
			fmt.Printf("    %s%s%s\n", colorDim, msg, colorReset)
		}

		fmt.Printf("    %s%s%s\n", colorDim, fork.HTMLURL, colorReset)
	}
}
