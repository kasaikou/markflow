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

package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kasaikou/docstak/app"
	"github.com/kasaikou/docstak/cli"
	"github.com/kasaikou/docstak/docstak"
	"github.com/kasaikou/docstak/docstak/condition"
	"github.com/kasaikou/docstak/docstak/model"
	"github.com/kasaikou/docstak/docstak/srun"
)

func run(ctx context.Context, args parseArgResult) int {
	cwWaiter := sync.WaitGroup{}
	defer cwWaiter.Wait()
	cw, _ := cli.NewConsoleWriter(os.Stdout, cli.TerminalAutoDetect(os.Stdout))
	cwWaiter.Add(1)
	go func() {
		defer cwWaiter.Done()
		cw.Route()
	}()
	defer cw.Close()

	logger := slog.New(cw.NewLoggerHandler(nil))
	ctx = docstak.WithLogger(ctx, logger)
	document, success := app.NewLocalDocument(ctx)
	if !success {
		return -1
	}

	chDecoration := make(chan cli.ProcessOutputDecoration, len(cli.ProcessOutputDecorations))
	for i := range cli.ProcessOutputDecorations {
		chDecoration <- cli.ProcessOutputDecorations[i]
	}

	ctx, cancel := context.WithCancel(ctx)

	sigWaiter := sync.WaitGroup{}
	sigWaiter.Add(1)
	go func() {
		defer sigWaiter.Done()
		sig := make(chan os.Signal, 1)
		defer close(sig)

		signal.Notify(sig, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)
		select {
		case <-ctx.Done():
			return
		case signal := <-sig:
			logger.Warn("signal received, make graceful shutdown started", slog.String("signal", signal.String()))
			cancel()
			return
		}
	}()
	defer sigWaiter.Wait()
	defer cancel()

	exit := docstak.ExecuteContext(ctx, document.Document,
		docstak.ExecuteOptCalls(args.Cmds...),
		docstak.ExecuteOptProcessExec(func(ctx context.Context, task model.DocumentTask, runner *srun.ScriptRunner) (int, error) {
			decoration := <-chDecoration
			defer func() {
				chDecoration <- decoration
			}()

			skip := condition.NewSkipsFromDocumentTask(&task)
			isSkip := skip.Test(ctx, condition.TestOption{})
			if isSkip {
				logger.Info("task execute is not required", slog.String("task", task.Call))
				return 0, nil
			}

			sufficient := condition.NewRequiresFromDocumentTask(&task).Test(ctx, condition.TestOption{})
			if !sufficient {
				logger.Error("task's require rules are insufficient", slog.String("task", task.Call))
				return -1, nil
			}

			stdOutScanner := cw.NewScanner(decoration.Stdout, "STDOUT", task.Title)
			stdout, _ := runner.Stdout()
			stderrScanner := cw.NewScanner(decoration.Stderr, "ERROUT", task.Title)
			stderr, _ := runner.Stderr()

			wg := sync.WaitGroup{}
			defer wg.Wait()

			wg.Add(1)
			go func() {
				defer wg.Done()
				stdOutScanner.Scan(stdout)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				stderrScanner.Scan(stderr)
			}()

			logger.Info("task start", slog.String("task", task.Call))
			exit, err := runner.RunContext(ctx)
			logger.Info("task ended", slog.String("task", task.Call), slog.Int("exitCode", exit))

			if err != nil {
				return exit, err
			}

			skip.UpdateDocumentTask(ctx, &task)
			document.Document.Tasks[task.Call] = task

			return exit, nil
		}),
	)

	document.SaveState(ctx)
	return exit
}
