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
