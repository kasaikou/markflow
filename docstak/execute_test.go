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

package docstak_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/kasaikou/markflow/docstak"
	"github.com/kasaikou/markflow/docstak/model"
)

func TestExecute(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}))
	ctx := docstak.WithLogger(context.Background(), logger)

	document := model.Document{
		Tasks: map[string]model.DocumentTask{
			"echo": {
				Title: "echo",
				Call:  "echo",
				Scripts: []model.DocumentTaskScript{
					{
						Config: model.ExecConfig{
							ExecPath: "sh",
							CmdOpt:   "-c",
						},
						Script: "echo 'hello world'",
					},
				},
			},
		},
	}

	docstak.ExecuteContext(ctx, document, docstak.ExecuteOptCalls("echo"))
}
func TestExecuteParallel(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}))
	ctx := docstak.WithLogger(context.Background(), logger)

	document := model.Document{
		Tasks: map[string]model.DocumentTask{
			"echo-parallel-1": {
				Title: "echo-parallel-1",
				Call:  "echo-parallel-1",
				Scripts: []model.DocumentTaskScript{
					{
						Config: model.ExecConfig{
							ExecPath: "sh",
							CmdOpt:   "-c",
						},
						Script: "echo 'hello world'",
					},
				},
			},
			"echo-parallel-2": {
				Title: "echo-parallel-2",
				Call:  "echo-parallel-2",
				Scripts: []model.DocumentTaskScript{
					{
						Config: model.ExecConfig{
							ExecPath: "sh",
							CmdOpt:   "-c",
						},
						Script: "echo 'hello world'",
					},
				},
			},
			"echo-parallel-3": {
				Title: "echo-parallel-3",
				Call:  "echo-parallel-3",
				Scripts: []model.DocumentTaskScript{
					{
						Config: model.ExecConfig{
							ExecPath: "sh",
							CmdOpt:   "-c",
						},
						Script: "echo 'hello world'",
					},
				},
			},
			"echo-parallel-4": {
				Title: "echo-parallel-4",
				Call:  "echo-parallel-4",
				Scripts: []model.DocumentTaskScript{
					{
						Config: model.ExecConfig{
							ExecPath: "sh",
							CmdOpt:   "-c",
						},
						Script: "echo 'hello world'",
					},
				},
			},
		},
	}

	docstak.ExecuteContext(ctx, document, docstak.ExecuteOptCalls("echo-parallel-1", "echo-parallel-2", "echo-parallel-3", "echo-parallel-4"))
}
func TestExecuteDepends(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}))
	ctx := docstak.WithLogger(context.Background(), logger)

	document := model.Document{
		Tasks: map[string]model.DocumentTask{
			"echo-parallel-1": {
				Title: "echo-parallel-1",
				Call:  "echo-parallel-1",
				Scripts: []model.DocumentTaskScript{
					{
						Config: model.ExecConfig{
							ExecPath: "sh",
							CmdOpt:   "-c",
						},
						Script: "echo 'hello world'",
					},
				},
			},
			"echo-parallel-2": {
				Title: "echo-parallel-2",
				Call:  "echo-parallel-2",
				Scripts: []model.DocumentTaskScript{
					{
						Config: model.ExecConfig{
							ExecPath: "sh",
							CmdOpt:   "-c",
						},
						Script: "echo 'hello world'",
					},
				},
				DependTasks: []string{"echo-parallel-1"},
			},
			"echo-parallel-3": {
				Title: "echo-parallel-3",
				Call:  "echo-parallel-3",
				Scripts: []model.DocumentTaskScript{
					{
						Config: model.ExecConfig{
							ExecPath: "sh",
							CmdOpt:   "-c",
						},
						Script: "echo 'hello world'",
					},
				},
				DependTasks: []string{"echo-parallel-2"},
			},
			"echo-parallel-4": {
				Title: "echo-parallel-4",
				Call:  "echo-parallel-4",
				Scripts: []model.DocumentTaskScript{
					{
						Config: model.ExecConfig{
							ExecPath: "sh",
							CmdOpt:   "-c",
						},
						Script: "echo 'hello world'",
					},
				},
				DependTasks: []string{"echo-parallel-3"},
			},
		},
	}

	docstak.ExecuteContext(ctx, document, docstak.ExecuteOptCalls("echo-parallel-4"))
}
