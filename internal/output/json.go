package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/stympy/forkwatch/internal/analysis"
)

type jsonOutput struct {
	Repository string        `json:"repository"`
	TotalForks int           `json:"total_forks"`
	Analyzed   int           `json:"analyzed_forks"`
	Active     int           `json:"active_forks"`
	Clusters   []jsonCluster `json:"clusters"`
}

type jsonCluster struct {
	File        string           `json:"file"`
	Convergence int              `json:"convergence"`
	Forks       []jsonFork       `json:"forks"`
	PatchGroups []jsonPatchGroup `json:"patch_groups,omitempty"`
}

type jsonFork struct {
	Owner    string   `json:"owner"`
	URL      string   `json:"url"`
	AheadBy  int      `json:"ahead_by"`
	Commits  []string `json:"commit_messages"`
	Added    int      `json:"additions"`
	Deleted  int      `json:"deletions"`
	Patch    string   `json:"patch,omitempty"`
}

type jsonPatchGroup struct {
	Patch    string   `json:"patch"`
	ForkCount int     `json:"fork_count"`
	Forks    []string `json:"forks"`
}

func PrintJSON(result *analysis.AnalysisResult) error {
	out := jsonOutput{
		Repository: fmt.Sprintf("%s/%s", result.UpstreamOwner, result.UpstreamRepo),
		TotalForks: result.TotalForks,
		Analyzed:   result.AnalyzedForks,
		Active:     result.ActiveForks,
	}

	for _, c := range result.Clusters {
		jc := jsonCluster{
			File:        c.Filename,
			Convergence: c.Convergence,
		}
		for _, f := range c.Forks {
			jc.Forks = append(jc.Forks, jsonFork{
				Owner:   f.Owner,
				URL:     f.HTMLURL,
				AheadBy: f.AheadBy,
				Commits: f.CommitMessages,
				Added:   f.Additions,
				Deleted: f.Deletions,
				Patch:   f.Patch,
			})
		}
		if c.PatchGroups != nil {
			for _, g := range c.PatchGroups.Groups {
				var owners []string
				for _, f := range g.Forks {
					owners = append(owners, f.Owner)
				}
				jc.PatchGroups = append(jc.PatchGroups, jsonPatchGroup{
					Patch:     g.Full,
					ForkCount: len(g.Forks),
					Forks:     owners,
				})
			}
		}
		out.Clusters = append(out.Clusters, jc)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
