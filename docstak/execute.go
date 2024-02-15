package docstak

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/kasaikou/docstak/docstak/model"
	"github.com/kasaikou/docstak/docstak/srun"
)

type executeOptions struct {
	called []string
}

func newExecuteOptions() *executeOptions {
	return &executeOptions{}
}

type ExecuteOption func(eo *executeOptions) error

func ExecuteOptCalls(keys ...string) ExecuteOption {
	return func(eo *executeOptions) error {
		eo.called = append(eo.called, keys...)
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
			stdout, _ := runner.Stdout()
			stderr, _ := runner.Stderr()

			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				io.Copy(os.Stdout, stdout)
			}()
			wg.Add(1)
			go func() {
				defer wg.Done()
				io.Copy(os.Stderr, stderr)
			}()
			exit, err := runner.RunContext(ctx)
			wg.Wait()

			if err != nil {
				logger.Error(fmt.Sprintf("error in task '%s'", option.called[i]), slog.String("error", err.Error()))
				return -1
			} else if exit != 0 {
				logger.Error(fmt.Sprintf("exit code is not %d in task '%s'", exit, option.called[i]))
				return exit
			}
		}
	}

	return 0
}
