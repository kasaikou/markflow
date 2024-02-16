package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/kasaikou/docstak/cli"
	"github.com/kasaikou/docstak/docstak"
	"github.com/kasaikou/docstak/docstak/markdown"
	"github.com/kasaikou/docstak/docstak/model"
	"github.com/kasaikou/docstak/docstak/resolver"
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
	po, err := markdown.FromLocalFile(wd, "docstak.md")
	if err != nil {
		logger.Error("cannot open file", slog.Any("error", err))
	}

	parsed, err := markdown.ParseMarkdown(ctx, po)
	if err != nil {
		logger.Error("cannot parse markdown", slog.String("filepath", po.Filename()), slog.Any("error", err))
		return -1
	}

	document, err := model.NewDocument(ctx,
		model.NewDocOptionRootDir(filepath.Dir(po.Filename())),
		resolver.NewDocumentWithPathResolver(
			resolver.ResolveOption{Lang: []string{"sh"}, Command: "sh"},
		),
		markdown.NewDocFromMarkdownParsing(parsed),
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
