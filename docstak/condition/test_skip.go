package condition

import (
	"context"
	"log/slog"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/kasaikou/docstak/docstak"
	"github.com/kasaikou/docstak/docstak/model"
	"github.com/kasaikou/docstak/docstak/resolver"
)

type Skips struct {
	container []testContainer
}

func NewSkipsFromDocumentTask(dt *model.DocumentTask) *Skips {
	skips := &Skips{}
	container := testContainer{}
	for i := range dt.Skips.ExistPaths {
		container.existFiles = append(container.existFiles, FileIsExisted{
			Config: resolver.FileGlobConfig{
				Rootdir: dt.Parent.Rootdir,
				Rules:   []string{dt.Skips.ExistPaths[i]},
			},
		})
	}

	skips.container = append(skips.container, container)
	return skips
}

func (s *Skips) Test(ctx context.Context, opts TestOption) (skip bool) {
	logger := docstak.GetLogger(ctx)

	for itemIdx := range s.container {
		skip := true

		if len(s.container[itemIdx].existFiles) == 0 {
			skip = false
		} else {
			for ruleIdx := range s.container[itemIdx].existFiles {

				enable, err := s.container[itemIdx].existFiles[ruleIdx].IsEnable(ctx)

				if err != nil {
					skip = false
					if err != doublestar.ErrPatternNotExist {
						logger.Warn("returns error when check file is exist (check glob pattern is valid or not)", slog.Any("error", err))
					}

				} else if !enable {
					skip = false
				}
			}
		}

		if skip {
			return true
		}
	}

	return false
}
