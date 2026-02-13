package analysis

import (
	"fmt"
	"sort"
	"strings"
)

type PatchGroup struct {
	Patch string        // cleaned diff for display
	Full  string        // original patch text
	Forks []ForkSummary // forks with this identical patch
}

type PatchGrouping struct {
	Groups []PatchGroup // sorted largest-first
}

// GroupPatches groups forks by identical patch text. Empty patches are ungroupable
// (each gets its own single-fork group).
func GroupPatches(forks []ForkSummary) *PatchGrouping {
	grouped := make(map[string][]ForkSummary)
	var ungrouped []ForkSummary

	for _, f := range forks {
		if f.Patch == "" {
			ungrouped = append(ungrouped, f)
		} else {
			grouped[f.Patch] = append(grouped[f.Patch], f)
		}
	}

	var groups []PatchGroup
	for patch, members := range grouped {
		groups = append(groups, PatchGroup{
			Patch: CleanDiff(patch),
			Full:  patch,
			Forks: members,
		})
	}

	// Each ungrouped fork gets its own group
	for _, f := range ungrouped {
		groups = append(groups, PatchGroup{
			Forks: []ForkSummary{f},
		})
	}

	// Sort: largest groups first, then by first fork owner for stability
	sort.Slice(groups, func(i, j int) bool {
		if len(groups[i].Forks) != len(groups[j].Forks) {
			return len(groups[i].Forks) > len(groups[j].Forks)
		}
		return groups[i].Forks[0].Owner < groups[j].Forks[0].Owner
	})

	return &PatchGrouping{Groups: groups}
}

// CleanDiff strips @@ hunk headers, keeping only +/- and context lines.
func CleanDiff(patch string) string {
	var lines []string
	for _, line := range strings.Split(patch, "\n") {
		if strings.HasPrefix(line, "@@") {
			continue
		}
		lines = append(lines, line)
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

// TruncateDiff caps a diff at maxLines, appending a summary if truncated.
func TruncateDiff(diff string, maxLines int) string {
	lines := strings.Split(diff, "\n")
	if len(lines) <= maxLines {
		return diff
	}
	remaining := len(lines) - maxLines
	truncated := strings.Join(lines[:maxLines], "\n")
	return truncated + fmt.Sprintf("\n... (%d more lines)", remaining)
}
