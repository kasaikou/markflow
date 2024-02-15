package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlag(t *testing.T) {
	os.Args = append(os.Args, "-v", "-q", "fmt", "test")
	parseArgs()

	assert.Equal(t, true, *Verbose)
	assert.Equal(t, true, *Quiet)
	assert.Equal(t, false, *Help)
	assert.Contains(t, Cmds, "fmt")
	assert.Contains(t, Cmds, "test")
}
