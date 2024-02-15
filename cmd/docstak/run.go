package main

import (
	"context"
	"log/slog"
	"os"
	"path"

	"github.com/kasaikou/docstak/docstak"
	"github.com/kasaikou/docstak/docstak/markdown"
	"github.com/kasaikou/docstak/docstak/model"
	"github.com/kasaikou/docstak/docstak/srun"
)

func run() int {

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}))
	ctx := docstak.WithLogger(context.Background(), logger)
	wd, _ := os.Getwd()
	environment := markdown.FileEnvironment{
		Filepath: path.Join(wd, "docstak.md"),
	}

	yml, err := os.ReadFile(environment.Filepath)
	if err != nil {
		logger.Error("cannot open file", slog.String("filepath", environment.Filepath), slog.String("error", err.Error()))
		return -1
	}
	parsed, err := markdown.ParseMarkdown(ctx, yml)
	if err != nil {
		logger.Error("cannot parse markdown", slog.String("filepath", environment.Filepath), slog.String("error", err.Error()))
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

	return docstak.ExecuteContext(ctx, document,
		docstak.ExecuteOptCalls(Cmds...),
	)
}
