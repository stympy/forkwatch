package analysis

import (
	"sort"

	gh "github.com/stympy/forkwatch/internal/github"
)

type FileCluster struct {
	Filename    string
	Forks       []ForkSummary
	Convergence int            // number of independent forks touching this file
	PatchGroups *PatchGrouping // nil for single-fork files
}

type ForkSummary struct {
	Owner          string
	HTMLURL        string
	AheadBy        int
	CommitMessages []string
	Additions      int
	Deletions      int
	Patch          string
}

type AnalysisResult struct {
	UpstreamOwner string
	UpstreamRepo  string
	TotalForks    int
	AnalyzedForks int
	ActiveForks   int
	Clusters      []FileCluster
}

func Cluster(comparisons []*gh.ForkComparison, upstreamOwner, upstreamRepo string, totalForks int) *AnalysisResult {
	fileMap := make(map[string][]ForkSummary)

	for _, comp := range comparisons {
		// Build per-file additions/deletions for this fork
		fileStats := make(map[string]gh.FileChange)
		for _, f := range comp.FilesChanged {
			fileStats[f.Filename] = f
		}

		for _, f := range comp.FilesChanged {
			summary := ForkSummary{
				Owner:          comp.Fork.Owner,
				HTMLURL:        comp.Fork.HTMLURL,
				AheadBy:        comp.AheadBy,
				CommitMessages: comp.CommitMessages,
				Additions:      f.Additions,
				Deletions:      f.Deletions,
				Patch:          f.Patch,
			}
			fileMap[f.Filename] = append(fileMap[f.Filename], summary)
		}
	}

	var clusters []FileCluster
	for filename, forks := range fileMap {
		c := FileCluster{
			Filename:    filename,
			Forks:       forks,
			Convergence: len(forks),
		}
		if c.Convergence >= 2 {
			c.PatchGroups = GroupPatches(forks)
		}
		clusters = append(clusters, c)
	}

	// Sort: most convergent first, then alphabetically
	sort.Slice(clusters, func(i, j int) bool {
		if clusters[i].Convergence != clusters[j].Convergence {
			return clusters[i].Convergence > clusters[j].Convergence
		}
		return clusters[i].Filename < clusters[j].Filename
	})

	return &AnalysisResult{
		UpstreamOwner: upstreamOwner,
		UpstreamRepo:  upstreamRepo,
		TotalForks:    totalForks,
		AnalyzedForks: len(comparisons),
		ActiveForks:   len(comparisons),
		Clusters:      clusters,
	}
}
