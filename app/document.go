package app

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/kasaikou/docstak/docstak"
	"github.com/kasaikou/docstak/docstak/markdown"
	"github.com/kasaikou/docstak/docstak/model"
	"github.com/kasaikou/docstak/docstak/resolver"
)

var LanguageCmdPairs = []resolver.ResolveOption{
	{Lang: []string{"sh", "shell"}, Command: "sh", CmdOpt: "-c"},
	{Lang: []string{"bash"}, Command: "bash", CmdOpt: "-c"},
	{Lang: []string{"powershell", "posh"}, Command: "powershell", CmdOpt: "-Command"},
	{Lang: []string{"py", "python"}, Command: "python", CmdOpt: "-c"},
	{Lang: []string{"js", "javascript"}, Command: "node", CmdOpt: "-e"},
}

type documentOpt struct {
	resolveFilename string
	workingDir      string
}

type DocumentOpt func(d *documentOpt)

func NewDocument(ctx context.Context, opts ...DocumentOpt) (document model.Document, success bool) {

	logger := docstak.GetLogger(ctx)
	wd, _ := os.Getwd()
	documentOpt := documentOpt{
		resolveFilename: "docstak.md",
		workingDir:      wd,
	}

	for i := range opts {
		opts[i](&documentOpt)
	}

	po, err := markdown.FromLocalFile(wd, "docstak.md")
	if err != nil {
		logger.Error("cannot open file", slog.Any("error", err))
		return document, false
	}

	parsed, err := markdown.ParseMarkdown(ctx, po)
	if err != nil {
		logger.Error("cannot parse markdown", slog.String("filepath", po.Filename()), slog.Any("error", err))
		return document, false
	}

	document, err = model.NewDocument(ctx,
		model.NewDocOptionRootDir(filepath.Dir(po.Filename())),
		resolver.NewDocumentWithPathResolver(LanguageCmdPairs...),
		markdown.NewDocFromMarkdownParsing(parsed),
	)

	if err != nil {
		logger.Error("failed to initialize document", slog.String("error", err.Error()))
		return document, false
	}

	return document, true
}
