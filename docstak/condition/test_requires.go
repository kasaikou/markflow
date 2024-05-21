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
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/kasaikou/markflow/docstak"
	"github.com/kasaikou/markflow/docstak/model"
	"github.com/kasaikou/markflow/docstak/resolver"
)

type Requires struct {
	container []testContainer
}

type testContainer struct {
	existFiles      []FileIsExisted
	notChangedFiles []FileNotChanged
}

type TestOption struct{}

func NewRequiresFromDocumentTask(dt *model.DocumentTask) *Requires {
	requires := &Requires{}
	container := testContainer{}
	for i := range dt.Requires.ExistPaths {
		container.existFiles = append(container.existFiles, FileIsExisted{
			Config: resolver.FileGlobConfig{
				Rootdir: dt.Parent.Rootdir,
				Rules:   []string{dt.Requires.ExistPaths[i]},
			},
		})
	}

	requires.container = append(requires.container, container)
	return requires
}

func (r *Requires) Test(ctx context.Context, opts TestOption) (sufficient bool) {
	logger := docstak.GetLogger(ctx)

	logFns := make([]func(), 0, len(r.container))

	for itemIdx := range r.container {
		valid := true
		for ruleIdx := range r.container[itemIdx].existFiles {

			enable, err := r.container[itemIdx].existFiles[ruleIdx].IsEnable(ctx)

			if err != nil {
				if err == doublestar.ErrPatternNotExist {
					logFns = append(logFns, func() {
						logger.Error("cannot found files matched with patterns",
							slog.String("pattern", strings.Join(r.container[itemIdx].existFiles[ruleIdx].Config.Rules, " | ")),
						)
					})
					valid = false

				} else {
					logger.Warn("returns error when check file is exist (check glob pattern is valid or not)", slog.Any("error", err))
					logFns = append(logFns, func() {
						logger.Error("returns error when check file is exist (check glob pattern is valid or not)", slog.Any("error", err))
					})
				}

			} else if !enable {
				logFns = append(logFns, func() {
					logger.Error("cannot found files matched with patterns",
						slog.String("pattern", strings.Join(r.container[itemIdx].existFiles[ruleIdx].Config.Rules, " | ")),
					)
				})
				valid = false
			}
		}

		if valid {
			return true
		}
	}

	for i := range logFns {
		logFns[i]()
	}

	return false
}
