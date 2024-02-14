package docstak_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/kasaikou/docstak/docstak"
	"github.com/kasaikou/docstak/docstak/model"
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
						ExecPath: "bash",
						Script:   "echo 'hello world'",
					},
				},
			},
		},
	}

	docstak.ExecuteContext(ctx, document, docstak.ExecuteOptCalls("echo"))
}
