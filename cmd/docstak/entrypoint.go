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
	"log/slog"
	"os"
	"sync"

	"github.com/kasaikou/markflow/cli"
	"github.com/kasaikou/markflow/docstak"
)

func entrypoint(args parseArgResult) int {

	cwWaiter := sync.WaitGroup{}
	defer cwWaiter.Wait()
	cw, _ := cli.NewConsoleWriter(os.Stderr, cli.TerminalAutoDetect(os.Stderr))
	cwWaiter.Add(1)
	go func() {
		defer cwWaiter.Done()
		cw.Route()
	}()
	defer cw.Close()

	logger := slog.New(cw.NewLoggerHandler(nil))
	ctx := docstak.WithLogger(context.Background(), logger)

	type featureFlag struct {
		Name   string
		Enable bool
		Fn     func(context.Context, parseArgResult) int
	}

	featureFlags := []featureFlag{
		{
			Name:   "--dry-run",
			Enable: *args.DryRun,
			Fn:     func(ctx context.Context, args parseArgResult) int { return dryrun(ctx, args) },
		},
	}

	enabledFeature := []featureFlag{}
	for i := range featureFlags {
		if featureFlags[i].Enable {
			enabledFeature = append(enabledFeature, featureFlags[i])
		}
	}

	switch len(enabledFeature) {

	case 0: // Default feature (execute task).
		return run(ctx, args)

	case 1:
		return enabledFeature[0].Fn(ctx, args)

	default:
		joined := ""
		for i := range featureFlags {
			if featureFlags[i].Enable {
				if joined != "" {
					joined += ", "
				}

				joined += featureFlags[i].Name
			}
		}

		logger.Error("duplicated sub-command option", slog.String("sub-commands", joined))
		return -1
	}
}
