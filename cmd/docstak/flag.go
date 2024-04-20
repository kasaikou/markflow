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

import "github.com/spf13/pflag"

type parseArgResult struct {
	Verbose *bool    `json:"verbose,omitempty"`
	Quiet   *bool    `json:"quiet,omitempty"`
	Help    *bool    `json:"help,omitempty"`
	DryRun  *bool    `json:"dry_run,omitempty"`
	Cmds    []string `json:"cmds,omitempty"`
}

func parseArgs() parseArgResult {

	verbose := pflag.BoolP("verbose", "v", true, "Be verbose (default).")
	quiet := pflag.BoolP("quiet", "q", false, "Output only error message with stderr.")
	help := pflag.BoolP("help", "h", false, "Output help information.")
	dryRun := pflag.Bool("dry-run", false, "Output the operation configuration but do not execute.")

	pflag.Parse()
	cmds := pflag.Args()

	return parseArgResult{
		Verbose: verbose,
		Quiet:   quiet,
		Help:    help,
		DryRun:  dryRun,
		Cmds:    cmds,
	}
}
