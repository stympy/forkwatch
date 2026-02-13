package analysis

// Recommendation is the top-voted change for a single file.
type Recommendation struct {
	File          string
	Patch         string // raw patch from GitHub API (with @@ headers)
	Convergence   int    // total forks touching this file
	AgreedBy      int    // forks with this exact patch
	Forks         []string
	CommitMessage string // representative first-line commit message
}

// Recommend returns the most-converged-upon patch for each convergent
// cluster (convergence >= 2). The result is ordered by convergence
// descending, matching the cluster sort order.
func Recommend(result *AnalysisResult) []Recommendation {
	var recs []Recommendation
	for _, c := range result.Clusters {
		if c.Convergence < 2 || c.PatchGroups == nil || len(c.PatchGroups.Groups) == 0 {
			continue
		}
		top := c.PatchGroups.Groups[0]
		if len(top.Forks) < 2 || top.Full == "" {
			continue
		}
		var owners []string
		var msg string
		for _, f := range top.Forks {
			owners = append(owners, f.Owner)
			if msg == "" && len(f.CommitMessages) > 0 {
				msg = f.CommitMessages[0]
			}
		}
		recs = append(recs, Recommendation{
			File:          c.Filename,
			Patch:         top.Full,
			Convergence:   c.Convergence,
			AgreedBy:      len(top.Forks),
			Forks:         owners,
			CommitMessage: msg,
		})
	}
	return recs
}
