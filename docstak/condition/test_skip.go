/*
Copyright 2024 Kasai Kou

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package condition

import (
	"context"
	"log/slog"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/kasaikou/markflow/docstak"
	"github.com/kasaikou/markflow/docstak/model"
	"github.com/kasaikou/markflow/docstak/resolver"
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

	for i := range dt.Skips.NotChangedPaths {
		paths := make([]string, 0, len(dt.Skips.NotChangedPaths[i].Paths))
		for k := range dt.Skips.NotChangedPaths[i].Paths {
			paths = append(paths, k)
		}

		ignores := make([]string, 0, len(dt.Skips.NotChangedPaths[i].Ignores))
		for k := range dt.Skips.NotChangedPaths[i].Ignores {
			ignores = append(ignores, k)
		}

		container.notChangedFiles = append(container.notChangedFiles, FileNotChanged{
			Config: resolver.FileGlobConfig{
				Rootdir:    dt.Parent.Rootdir,
				Rules:      paths,
				IgnoreRule: ignores,
			},
			MD5: dt.Skips.NotChangedPaths[i].MD5,
		})
	}

	skips.container = append(skips.container, container)
	return skips
}

func (s *Skips) Test(ctx context.Context, opts TestOption) (skip bool) {
	logger := docstak.GetLogger(ctx)

	for itemIdx := range s.container {
		skip := true
		isEmpty := true

		if len(s.container[itemIdx].existFiles) > 0 {
			isEmpty = false
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

		if len(s.container[itemIdx].notChangedFiles) > 0 {
			isEmpty = false
			for ruleIdx := range s.container[itemIdx].notChangedFiles {

				enable, err := s.container[itemIdx].notChangedFiles[ruleIdx].IsEnable(ctx)

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

		if isEmpty {
			skip = false
		}

		if skip {
			return true
		}
	}

	return false
}

func (s *Skips) UpdateDocumentTask(ctx context.Context, dt *model.DocumentTask) {
	logger := docstak.GetLogger(ctx)

	for itemIdx := range s.container {
		for ruleIdx := range s.container[itemIdx].notChangedFiles {
			hash, err := s.container[itemIdx].notChangedFiles[ruleIdx].CurrentMD5(ctx)
			if err != nil {
				logger.Warn("failed to calculate md5", slog.Any("error", err))
			} else {
				logger.Info("update skip when files not changed rule's hash", slog.String("call", dt.Call), slog.String("hash", hash))
				dt.Skips.NotChangedPaths[ruleIdx].MD5 = hash
			}
		}
	}
}
