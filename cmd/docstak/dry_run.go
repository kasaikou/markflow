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

package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/kasaikou/markflow/app"
	"github.com/kasaikou/markflow/docstak"
)

func dryrun(ctx context.Context, args parseArgResult) int {
	logger := docstak.GetLogger(ctx)

	document, success := app.NewLocalDocument(ctx)
	if !success {
		return -1
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(document); err != nil {
		logger.Error("cannot encode to json", slog.Any("error", err))
		return -1
	}

	return 0
}
