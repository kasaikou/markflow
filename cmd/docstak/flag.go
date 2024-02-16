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

var (
	Verbose *bool
	Quiet   *bool
	Help    *bool
	Cmds    []string
)

func parseArgs() {
	Verbose = pflag.BoolP("verbose", "v", true, "be verbose (default)")
	Quiet = pflag.BoolP("quiet", "q", false, "be quiet")
	Help = pflag.BoolP("help", "h", false, "output help")

	pflag.Parse()
	Cmds = pflag.Args()
}
