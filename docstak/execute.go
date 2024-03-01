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

package docstak

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"sync"

	"github.com/kasaikou/docstak/docstak/model"
	"github.com/kasaikou/docstak/docstak/srun"
)

type executeOptions struct {
	called    []string
	onExec    func(ctx context.Context, task model.DocumentTask, runner *srun.ScriptRunner) (int, error)
	numWorker int
}

func newExecuteOptions() *executeOptions {
	return &executeOptions{
		onExec: func(ctx context.Context, task model.DocumentTask, runner *srun.ScriptRunner) (int, error) {
			return runner.RunContext(ctx)
		},
		numWorker: runtime.NumCPU(),
	}
}

type ExecuteOption func(eo *executeOptions) error

func ExecuteOptCalls(keys ...string) ExecuteOption {
	return func(eo *executeOptions) error {
		eo.called = append(eo.called, keys...)
		return nil
	}
}

func ExecuteOptProcessExec(fn func(ctx context.Context, task model.DocumentTask, runner *srun.ScriptRunner) (int, error)) ExecuteOption {
	return func(eo *executeOptions) error {
		eo.onExec = fn
		return nil
	}
}

func ExecuteContext(ctx context.Context, document model.Document, options ...ExecuteOption) int {

	logger := GetLogger(ctx)
	option := newExecuteOptions()

	for i := range options {
		if err := options[i](option); err != nil {
			logger.Error("failed to load execute options", slog.String("error", err.Error()))
		}
	}

	execTasks := map[string]struct{}{}
	called := option.called
	for i := 0; i < len(called); i++ {
		task, exist := document.Tasks[called[i]]

		if !exist {
			logger.Error(fmt.Sprintf("cannot found task '%s'", called[i]))
		}

		if _, exist := execTasks[called[i]]; exist {
			continue
		}

		execTasks[called[i]] = struct{}{}
		called = append(called, task.DependTasks...)
	}

	tasks := make([]string, 0, len(execTasks))
	for task := range execTasks {
		tasks = append(tasks, task)
	}

	return executeTasks(ctx, document, option, tasks)
}

func executeTasks(ctx context.Context, document model.Document, option *executeOptions, executeTasks []string) int {
	wg := sync.WaitGroup{}

	type taskResp struct {
		Call string
		Exit int
	}

	chTaskResp := make(chan taskResp)
	defer close(chTaskResp)
	var cancel func()
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()
	taskChs := make([]chan taskResp, 0, len(executeTasks))

	for i := range executeTasks {
		ch := make(chan taskResp, len(executeTasks))
		defer close(ch)
		task := document.Tasks[executeTasks[i]]
		wg.Add(1)
		go func(ctx context.Context, task model.DocumentTask, chEnded <-chan taskResp, chRes chan<- taskResp) {
			defer wg.Done()

			depends := map[string]struct{}{}
			for i := range task.DependTasks {
				depends[task.DependTasks[i]] = struct{}{}
			}

			<-chEnded

			for len(depends) > 0 {
				select {
				case <-ctx.Done():
					return
				case res := <-chEnded:
					delete(depends, res.Call)
				}
			}

			ch := make(chan taskResp)
			wg := sync.WaitGroup{}
			for j := range task.Scripts {
				wg.Add(1)
				go func(ctx context.Context, task model.DocumentTask, script model.DocumentTaskScript, chRes chan<- taskResp) {
					defer wg.Done()
					exit := executeTask(ctx, task, script, option)

					if ctx.Err() == nil {
						chRes <- taskResp{
							Call: task.Call,
							Exit: exit,
						}
					}

				}(ctx, task, task.Scripts[j], ch)
			}
			defer wg.Wait()

			ended := 0

			if len(task.Scripts) == 0 {
				chRes <- taskResp{
					Call: task.Call,
					Exit: 0,
				}
			}

			for ended < len(task.Scripts) {
				select {
				case <-ctx.Done():
					return
				case result := <-ch:
					ended++
					if ended >= len(task.Scripts) {
						chRes <- result
					} else if result.Exit != 0 {
						chRes <- result
					}
				}
			}

		}(ctx, task, ch, chTaskResp)

		taskChs = append(taskChs, ch)
	}
	defer wg.Wait()
	ended := 0

	for i := range taskChs {
		taskChs[i] <- taskResp{}
	}

	for {
		select {
		case <-ctx.Done():
			return -1
		case res := <-chTaskResp:
			if res.Exit != 0 {
				cancel()
				return res.Exit
			}

			ended++
			if ended >= len(executeTasks) {
				return 0
			}

			for i := range taskChs {
				taskChs[i] <- res
			}
		}
	}
}

func executeTask(ctx context.Context, task model.DocumentTask, script model.DocumentTaskScript, option *executeOptions) int {
	logger := GetLogger(ctx)

	runner := srun.NewScriptRunner(script.Config.ExecPath, script.Config.CmdOpt, script.Script, script.Config.Args...)

	for key, value := range task.Envs {
		runner.SetEnv(key, value)
	}

	environ := os.Environ()
	for i := range environ {
		runner.SetEnviron(environ[i])
	}

	exit, err := option.onExec(ctx, task, runner)

	if err != nil {
		logger.Error("task ended with error", slog.String("task", task.Call), slog.Any("error", err))
		return -1
	}

	return exit
}
