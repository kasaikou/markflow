package resolver

import (
	"os"
	"slices"

	"github.com/bmatcuk/doublestar/v4"
)

type FileGlobConfig struct {
	Rootdir    string
	Rules      []string
	IgnoreRule []string
}

func ResolveFileGlob(config FileGlobConfig) ([]string, error) {

	fileSystem := os.DirFS(config.Rootdir)
	candidates := []string{}
	for i := range config.Rules {
		matched, err := doublestar.Glob(fileSystem, config.Rules[i], doublestar.WithFilesOnly())
		if err != nil {
			return nil, err
		}

		candidates = append(candidates, matched...)
	}

	slices.Sort(candidates)
	candidates = slices.Compact(candidates)

	results := make([]string, 0, len(candidates))
	for i := range candidates {
		matched := true
		for j := range config.IgnoreRule {
			ignored, err := doublestar.PathMatch(config.IgnoreRule[j], candidates[i])
			if err != nil {
				return nil, err
			} else if ignored {
				matched = false
				break
			}
		}

		if matched {
			results = append(results, candidates[i])
		}
	}

	return results, nil
}
