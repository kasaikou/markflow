package main

import (
	"context"
	"log/slog"
	"os"
	"path"
	"sync"

	"github.com/kasaikou/docstak/cli"
	"github.com/kasaikou/docstak/docstak"
	"github.com/kasaikou/docstak/docstak/markdown"
	"github.com/kasaikou/docstak/docstak/model"
	"github.com/kasaikou/docstak/docstak/srun"
)

func run() int {
	cwWaiter := sync.WaitGroup{}
	defer cwWaiter.Wait()
	cw, _ := cli.NewConsoleWriter(os.Stdout)
	cwWaiter.Add(1)
	go func() {
		defer cwWaiter.Done()
		cw.Route()
	}()
	defer cw.Close()

	logger := slog.New(cw.NewLoggerHandler(nil))
	ctx := docstak.WithLogger(context.Background(), logger)
	wd, _ := os.Getwd()
	environment := markdown.FileEnvironment{
		Filepath: path.Join(wd, "docstak.md"),
	}

	yml, err := os.ReadFile(environment.Filepath)
	if err != nil {
		logger.Error("cannot open file", slog.String("filepath", environment.Filepath), slog.Any("error", err))
		return -1
	}
	parsed, err := markdown.ParseMarkdown(ctx, yml)
	if err != nil {
		logger.Error("cannot parse markdown", slog.String("filepath", environment.Filepath), slog.Any("error", err))
		return -1
	}

	document, err := model.NewDocument(ctx,
		srun.NewDocumentWithPathResolver(
			srun.ResolveOption{Lang: []string{"sh"}, Command: "sh"},
		),
		markdown.NewDocFromMarkdown(environment, parsed),
	)

	if err != nil {
		logger.Error("failed to initialize document", slog.String("error", err.Error()))
		return -1
	}

	chDecoration := make(chan cli.ProcessOutputDecoration, len(cli.ProcessOutputDecorations))
	for i := range cli.ProcessOutputDecorations {
		chDecoration <- cli.ProcessOutputDecorations[i]
	}

	return docstak.ExecuteContext(ctx, document,
		docstak.ExecuteOptCalls(Cmds...),
		docstak.ExecuteOptPreProcessExec(func(ctx context.Context, task model.DocumentTask, runner *srun.ScriptRunner) (int, error) {
			decoration := <-chDecoration
			defer func() {
				chDecoration <- decoration
			}()

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

			return runner.RunContext(ctx)
		}),
	)
}
