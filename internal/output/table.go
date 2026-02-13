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

		for _, fork := range cluster.Forks {
			stats := fmt.Sprintf("%s+%d%s %s-%d%s",
				colorGreen, fork.Additions, colorReset,
				colorRed, fork.Deletions, colorReset)

			fmt.Printf("  %s%-20s%s %s (%d commits ahead)\n",
				colorCyan, fork.Owner, colorReset, stats, fork.AheadBy)

			// Show first commit message
			if len(fork.CommitMessages) > 0 {
				msg := fork.CommitMessages[0]
				if len(msg) > 72 {
					msg = msg[:72] + "..."
				}
				fmt.Printf("    %s%s%s\n", colorDim, msg, colorReset)
			}

			fmt.Printf("    %s%s%s\n", colorDim, fork.HTMLURL, colorReset)
		}

		fmt.Println(strings.Repeat("─", 60))
	}
}
