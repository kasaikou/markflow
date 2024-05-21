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

package resolver

import (
	"context"
	"os/exec"

	"github.com/cockroachdb/errors"
	"github.com/kasaikou/markflow/docstak/model"
)

type ResolveOption struct {
	Lang    []string
	Command string
	CmdOpt  string
	Args    []string
}

func NewDocumentWithPathResolver(options ...ResolveOption) model.NewDocumentOption {
	return func(ctx context.Context, d *model.DocumentConfig) error {
		for i := range options {
			for j := range options[i].Lang {
				execPath, err := exec.LookPath(options[i].Command)
				if err != nil {
					if errors.Is(err, exec.ErrNotFound) {
						continue
					}

					return err
				}
				d.ExecPathResolver[options[i].Lang[j]] = model.ExecConfig{
					ExecPath: execPath,
					CmdOpt:   options[i].CmdOpt,
					Args:     options[i].Args,
				}
			}
		}

		return nil
	}
}
