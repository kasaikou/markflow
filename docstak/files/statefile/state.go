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

package statefile

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"os"
	"sort"

	"github.com/cockroachdb/errors"
	"github.com/kasaikou/docstak/docstak/model"
)

type State struct {
	Tasks map[string]StateTask `json:"tasks"`
}

type StateTask struct {
	Files map[string]StateTaskFile `json:"files,omitempty"`
}

type StateTaskFile struct {
	Rule StateTaskFileRule `json:"rule"`
	MD5  string            `json:"md5,omitempty"`
}

type StateTaskFileRule struct {
	Paths   []string `json:"paths"`
	Ignores []string `json:"ignores"`
}

func FromLocalFile(filename string) (State, error) {
	file, err := os.Open(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return State{}, nil
		}

		return State{}, errors.WithMessage(err, "cannot open state file")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var state State
	err = decoder.Decode(&state)
	if err != nil {
		return state, errors.WithMessage(err, "cannot load as toml file")
	}

	return state, nil
}

func SaveLocalFile(filename string, s State) error {
	file, err := os.Create(filename)
	if err != nil {
		return errors.WithMessage(err, "cannot create file")
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(s); err != nil {
		return errors.WithMessage(err, "cannot save state file")
	}

	return nil
}

func SetStateParsed(result State) model.NewDocumentOption {
	return func(ctx context.Context, d *model.DocumentConfig) error {

		for call, taskStates := range result.Tasks {
			config, exist := d.Document.Tasks[call]

			if exist {
				for _, file := range taskStates.Files {
					for j := range config.Skips.NotChangedPaths {
						if config.Skips.NotChangedPaths[j].IsEqualRule(file.Rule.Paths, file.Rule.Ignores) {
							config.Skips.NotChangedPaths[j].MD5 = file.MD5
						}
					}
				}
				d.Document.Tasks[call] = config
			}
		}

		return nil
	}
}

func FromDocument(ctx context.Context, d model.Document) *State {
	state := State{
		Tasks: make(map[string]StateTask),
	}

	empty := true

	for call, task := range d.Tasks {
		stateTask := StateTask{
			Files: make(map[string]StateTaskFile),
		}

		for i := range task.Skips.NotChangedPaths {
			empty = false
			paths := make([]string, 0, len(task.Skips.NotChangedPaths[i].Paths))
			for k := range task.Skips.NotChangedPaths[i].Paths {
				paths = append(paths, k)
			}

			ignores := make([]string, 0, len(task.Skips.NotChangedPaths[i].Ignores))
			for k := range task.Skips.NotChangedPaths[i].Ignores {
				ignores = append(ignores, k)
			}

			sort.Strings(paths)
			sort.Strings(ignores)

			rule := StateTaskFileRule{
				Paths:   paths,
				Ignores: ignores,
			}

			hash := md5.New()
			json.NewEncoder(hash).Encode(rule)
			key := hex.EncodeToString(hash.Sum(nil))

			if task.Skips.NotChangedPaths[i].MD5 != "" {
				stateTask.Files[key] = StateTaskFile{
					Rule: rule,
					MD5:  task.Skips.NotChangedPaths[i].MD5,
				}
			}
		}

		state.Tasks[call] = stateTask
	}

	if empty {
		return nil
	}
	return &state
}
