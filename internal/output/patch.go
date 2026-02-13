package output

import (
	"fmt"
	"strings"

	"github.com/stympy/forkwatch/internal/analysis"
)

// PrintPatch emits a combined unified diff suitable for `git apply`.
// It selects the most-converged-upon patch for each file cluster.
func PrintPatch(result *analysis.AnalysisResult) {
	recs := analysis.Recommend(result)
	for i, rec := range recs {
		if i > 0 {
			// blank line between file diffs
			fmt.Println()
		}
		fmt.Printf("--- a/%s\n", rec.File)
		fmt.Printf("+++ b/%s\n", rec.File)
		// The GitHub API patch already contains @@ hunk headers and
		// diff lines; print it as-is.
		fmt.Println(strings.TrimRight(rec.Patch, "\n"))
	}
}
