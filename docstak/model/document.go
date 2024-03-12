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

package model

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
)

type Document struct {
	Title       string                  `json:"title,omitempty"`
	Description string                  `json:"description,omitempty"`
	Rootdir     string                  `json:"rootdir"`
	Tasks       map[string]DocumentTask `json:"tasks,omitempty"`
	GlobalEnvs  map[string]string       `json:"global_envs,omitempty"`
}

type DocumentConfig struct {
	Document         Document
	ExecPathResolver map[string]ExecConfig
}

type ExecConfig struct {
	ExecPath string   `json:"exec_path"`
	CmdOpt   string   `json:"cmd_opt,omitempty"`
	Args     []string `json:"args,omitempty"`
}

type Condition interface {
	IsEnable(context.Context) (bool, error)
}

type DocumentTask struct {
	Parent      *Document            `json:"-"`
	Title       string               `json:"omitempty"`
	Call        string               `json:"call"`
	Description string               `json:"description,omitempty"`
	Scripts     []DocumentTaskScript `json:"scripts"`
	Envs        map[string]string    `json:"envs,omitempty"`
	Skips       TaskSkipCondition    `json:"skips,omitempty"`
	Requires    TaskRequireCondition `json:"requires,omitempty"`
	DependTasks []string             `json:"depend_tasks,omitempty"`
}

type TaskSkipCondition struct {
	ExistPaths      []string                      `json:"exist_paths,omitempty"`
	NotChangedPaths []TaskFileNotChangedCondition `json:"not_changed_paths,omitempty"`
}

type TaskRequireCondition struct {
	ExistPaths []string `json:"exist_paths,omitempty"`
}

type setString map[string]struct{}

type TaskFileNotChangedCondition struct {
	Paths   setString `json:"paths,omitempty"`
	Ignores setString `json:"ignores,omitempty"`
	MD5     string    `json:"md5,omitempty"`
}

func (t setString) MarshalJSON() ([]byte, error) {

	keys := make([]string, 0, len(t))
	for k := range t {
		keys = append(keys, k)
	}

	return json.Marshal(keys)
}

func (cond *TaskFileNotChangedCondition) IsEqualRule(paths []string, ignores []string) bool {
	if len(cond.Paths) != len(paths) || len(cond.Ignores) != len(ignores) {
		return false
	}

	for i := range paths {
		if _, exist := cond.Paths[paths[i]]; !exist {
			return false
		}
	}

	for i := range ignores {
		if _, exist := cond.Paths[ignores[i]]; !exist {
			return false
		}
	}

	return true
}

type DocumentTaskScript struct {
	Config ExecConfig
	Script string
}

type NewDocumentOption func(ctx context.Context, d *DocumentConfig) error

func NewDocOptionRootDir(dirname string) NewDocumentOption {
	if !filepath.IsAbs(dirname) {
		panic("dirname must be absolute path")
	}

	return func(ctx context.Context, d *DocumentConfig) error {
		d.Document.Rootdir = dirname
		return nil
	}
}

func NewDocument(ctx context.Context, options ...NewDocumentOption) (Document, error) {
	document := DocumentConfig{
		ExecPathResolver: map[string]ExecConfig{},
		Document: Document{
			Tasks: map[string]DocumentTask{},
		},
	}

	for i := range options {
		if err := options[i](ctx, &document); err != nil {
			return document.Document, err
		}
	}

	if err := validateIsTaskDependencyCirculated(document.Document); err != nil {
		return document.Document, err
	}

	return document.Document, nil
}

func validateIsTaskDependencyCirculated(document Document) error {

	mpRead := map[string]map[string]struct{}{}
	history := make([]string, 0, len(document.Tasks))
	for task := range document.Tasks {
		history := history[:0]
		if err := validateIsTaskDependencyCirculatedInternal(task, history, mpRead, document); err != nil {
			return err
		}
	}
	return nil
}

var ErrCirculatedDependency = errors.New("task dependency is circulated")

func validateIsTaskDependencyCirculatedInternal(task string, mpHistory []string, mpDepends map[string]map[string]struct{}, document Document) error {

	if _, exist := mpDepends[task]; exist {
		return nil
	}

	errCirculated := func(task string) error {
		return errors.WithDetail(ErrCirculatedDependency, strings.Join(append(mpHistory, task+" (circulated)"), " -> "))
	}

	depends := map[string]struct{}{}

	mpHistory = append(mpHistory, task)
	t := document.Tasks[task]
	for i := range t.DependTasks {
		depends[t.DependTasks[i]] = struct{}{}

		for j := range mpHistory {
			if t.DependTasks[i] == mpHistory[j] {
				return errCirculated(t.DependTasks[i])
			}
		}

		if indirectDepends, cached := mpDepends[t.DependTasks[i]]; cached {
			for dependTask := range indirectDepends {
				for j := range mpHistory {
					if dependTask == mpHistory[j] {
						return errCirculated("..." + dependTask)
					}
				}
				depends[dependTask] = struct{}{}
			}

		} else {
			err := validateIsTaskDependencyCirculatedInternal(t.DependTasks[i], mpHistory, mpDepends, document)
			if err != nil {
				return err
			}

			indirectDepends := mpDepends[t.DependTasks[i]]
			for dependTask := range indirectDepends {
				depends[dependTask] = struct{}{}
			}
		}
	}

	mpDepends[task] = depends
	return nil
}
