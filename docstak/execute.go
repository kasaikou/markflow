package docstak

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/kasaikou/docstak/docstak/model"
	"github.com/kasaikou/docstak/docstak/srun"
)

type executeOptions struct {
	called []string
	onExec func(ctx context.Context, task model.DocumentTask, runner *srun.ScriptRunner) (int, error)
}

func newExecuteOptions() *executeOptions {
	return &executeOptions{
		onExec: func(ctx context.Context, task model.DocumentTask, runner *srun.ScriptRunner) (int, error) {
			return runner.RunContext(ctx)
		},
	}
}

type ExecuteOption func(eo *executeOptions) error

func ExecuteOptCalls(keys ...string) ExecuteOption {
	return func(eo *executeOptions) error {
		eo.called = append(eo.called, keys...)
		return nil
	}
}

func ExecuteOptPreProcessExec(fn func(ctx context.Context, task model.DocumentTask, runner *srun.ScriptRunner) (int, error)) ExecuteOption {
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

	for i := range option.called {
		task, exist := document.Tasks[option.called[i]]

		if !exist {
			logger.Error(fmt.Sprintf("cannot found task '%s'", option.called[i]))
		}

		for j := range task.Scripts {
			runner := srun.NewScriptRunner(task.Scripts[j].ExecPath, task.Scripts[j].Script)
			logger.Info("task start", slog.String("task", task.Call))
			exit, err := option.onExec(ctx, task, runner)
			logger.Info("task ended", slog.String("task", task.Call), slog.Int("exitCode", exit))

			if err != nil {
				logger.Error("task ended with error", slog.String("task", task.Call), slog.Any("error", err))
				return -1
			}
		}
	}

	return 0
}
