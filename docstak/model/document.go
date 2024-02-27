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
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
)

type Document struct {
	Title       string
	Description string
	Rootdir     string
	Tasks       map[string]DocumentTask
	GlobalEnvs  map[string]string
}

type DocumentConfig struct {
	Document         Document
	ExecPathResolver map[string]ExecConfig
}

type ExecConfig struct {
	ExecPath string
	CmdOpt   string
	Args     []string
}

type Condition interface {
	IsEnable(context.Context) (bool, error)
}

type DocumentTask struct {
	Parent      *Document
	Title       string
	Call        string
	Description string
	Scripts     []DocumentTaskScript
	Envs        map[string]string
	Skips       TaskSkipCondition
	Requires    TaskRequireCondition
	DependTasks []string
}

type TaskSkipCondition struct {
	ExistPaths []string
}

type TaskRequireCondition struct {
	ExistPaths []string
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
