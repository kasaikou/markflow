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

package markdown

import (
	"context"
	"log/slog"
	"os"
	"path"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/kasaikou/docstak/docstak"
	"github.com/kasaikou/docstak/docstak/environ"
	"github.com/kasaikou/docstak/docstak/model"
)

func setDocumentTask(ctx context.Context, document *model.DocumentConfig, result ParseResultTask) error {
	name := result.Title

	if _, exist := document.Document.Tasks[name]; exist {
		return errors.Errorf("duplicated task: '%s'", name)
	}

	config := model.DocumentTask{
		Parent:      &document.Document,
		Title:       result.Title,
		Call:        name,
		Description: result.Description,
		Envs:        make(map[string]string),
		DependTasks: result.Config.Previous,
	}

	// Read dotenv files.
	for i := range result.Config.Environ.Dotenvs {
		if !path.IsAbs(result.Config.Environ.Dotenvs[i]) {
			result.Config.Environ.Dotenvs[i] = path.Join(document.Document.Rootdir, result.Config.Environ.Dotenvs[i])
		}
		err := environ.LoadDotenv(result.Config.Environ.Dotenvs[i], func(key, value string) {
			config.Envs[key] = value
		})
		if err != nil {
			if os.IsNotExist(err) {
				docstak.GetLogger(ctx).Warn("dotenv file not found", slog.String("filename", result.Config.Environ.Dotenvs[i]))
			} else {
				return err
			}
		}
	}

	for key, value := range result.Config.Environ.Variables {
		config.Envs[key] = value
	}

	config.Requires.ExistPaths = result.Config.Requires.Exists

	for i := range result.Commands {
		execConfig, exist := document.ExecPathResolver[result.Commands[i].Lang]
		if !exist {
			return errors.Errorf("cannot resolve execute path in defined script language '%s'", result.Commands[i].Lang)
		}

		config.Scripts = append(config.Scripts, model.DocumentTaskScript{
			Config: execConfig,
			Script: result.Commands[i].Code,
		})
	}

	document.Document.Tasks[name] = config
	return nil
}

func NewDocFromMarkdownParsing(result ParseResult) model.NewDocumentOption {
	return func(ctx context.Context, document *model.DocumentConfig) error {
		document.Document.Title = result.Title
		document.Document.Description = result.Description

		// Update root directory.
		// This parameter is optional.
		if result.Config.Root != "" {
			if filepath.IsAbs(result.Config.Root) {
				document.Document.Rootdir = result.Config.Root
			} else {
				filepath.Join(document.Document.Rootdir, result.Config.Root)
			}
		}

		// Read dotenv files.
		for i := range result.Config.Environ.Dotenvs {
			if !path.IsAbs(result.Config.Environ.Dotenvs[i]) {
				result.Config.Environ.Dotenvs[i] = path.Join(document.Document.Rootdir, result.Config.Environ.Dotenvs[i])
			}
			err := environ.LoadDotenv(result.Config.Environ.Dotenvs[i], func(key, value string) {
				document.Document.GlobalEnvs[key] = value
			})
			if err != nil {
				if os.IsNotExist(err) {
					docstak.GetLogger(ctx).Warn("dotenv file not found", slog.String("filename", result.Config.Environ.Dotenvs[i]))
				} else {
					return err
				}
			}
		}

		// Set environment variables.
		// It's higher priority than dotenv files.
		for key, value := range result.Config.Environ.Variables {
			document.Document.GlobalEnvs[key] = value
		}

		for i := range result.Tasks {
			if err := setDocumentTask(ctx, document, result.Tasks[i]); err != nil {
				return err
			}
		}

		return nil
	}
}
