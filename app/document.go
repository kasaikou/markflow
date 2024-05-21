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

package app

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/kasaikou/markflow/docstak"
	"github.com/kasaikou/markflow/docstak/files/markdown"
	"github.com/kasaikou/markflow/docstak/files/statefile"
	"github.com/kasaikou/markflow/docstak/model"
	"github.com/kasaikou/markflow/docstak/resolver"
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

type LocalDocument struct {
	MarkdownFilename string
	StateFilename    string
	Document         model.Document
}

func NewLocalDocument(ctx context.Context, opts ...DocumentOpt) (document LocalDocument, success bool) {

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

	stateFilename := filepath.Join(filepath.Dir(po.Filename()), ".docstak_state.json")
	state, err := statefile.FromLocalFile(stateFilename)
	if err != nil {
		logger.Error("cannot parse state", slog.String("filepath", po.Filename()), slog.Any("error", err))
		return document, false
	}

	doc, err := model.NewDocument(ctx,
		model.NewDocOptionRootDir(filepath.Dir(po.Filename())),
		resolver.NewDocumentWithPathResolver(LanguageCmdPairs...),
		markdown.NewDocFromMarkdownParsing(parsed),
		statefile.SetStateParsed(state),
	)

	if err != nil {
		logger.Error("failed to initialize document", slog.String("error", err.Error()))
		return document, false
	}

	return LocalDocument{
		MarkdownFilename: po.Filename(),
		StateFilename:    stateFilename,
		Document:         doc,
	}, true
}

func (local *LocalDocument) SaveState(ctx context.Context) {
	logger := docstak.GetLogger(ctx)
	state := statefile.FromDocument(ctx, local.Document)
	if state != nil {
		if err := statefile.SaveLocalFile(local.StateFilename, *state); err != nil {
			logger.Error("cannot save statefile", slog.String("filepath", local.StateFilename), slog.Any("error", err))
		}
	}
}
