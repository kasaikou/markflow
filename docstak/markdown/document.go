package markdown

import (
	"context"
	"errors"
	"fmt"

	"github.com/kasaikou/docstak/docstak/model"
)

func setDocumentTask(ctx context.Context, document *model.DocumentConfig, result ParseResultTask) error {
	name := result.Title

	if _, exist := document.Document.Tasks[name]; exist {
		return errors.New(fmt.Sprintf("duplicated task: '%s'", name))
	}

	config := model.DocumentTask{
		Title:       result.Title,
		Call:        name,
		Description: result.Description,
	}

	for i := range result.Commands {
		path, exist := document.ExecPathResolver[result.Commands[i].Lang]
		if !exist {
			return errors.New(fmt.Sprintf("cannot resolve execute path in defined script language '%s'", result.Commands[i].Lang))
		}

		config.Scripts = append(config.Scripts, model.DocumentTaskScript{
			ExecPath: path,
			Script:   result.Commands[i].Code,
		})
	}

	document.Document.Tasks[name] = config
	return nil
}

func NewDocFromMarkdown(environment FileEnvironment, result ParseResult) model.NewDocumentOption {
	// dirpath := path.Dir(environment.Filepath)

	return func(ctx context.Context, document *model.DocumentConfig) error {
		document.Document.Title = result.Title
		document.Document.Description = result.Description

		for i := range result.Tasks {
			if err := setDocumentTask(ctx, document, result.Tasks[i]); err != nil {
				return err
			}
		}

		return nil
	}
}
