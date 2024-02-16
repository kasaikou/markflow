package resolver

import (
	"context"
	"os/exec"

	"github.com/kasaikou/docstak/docstak/model"
)

type ResolveOption struct {
	Lang    []string
	Command string
}

func NewDocumentWithPathResolver(options ...ResolveOption) model.NewDocumentOption {
	return func(ctx context.Context, d *model.DocumentConfig) error {
		for i := range options {
			for j := range options[i].Lang {
				execPath, err := exec.LookPath(options[i].Command)
				if err != nil {
					if err == exec.ErrNotFound {
						continue
					}

					return err
				}
				d.ExecPathResolver[options[i].Lang[j]] = execPath
			}
		}

		return nil
	}
}
