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

	return document.Document, nil
}
