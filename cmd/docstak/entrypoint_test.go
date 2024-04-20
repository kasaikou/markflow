package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithEmptyArgs(t *testing.T) {
	assert.NotEqual(t, 0, entrypoint(parseArgs([]string{})))
}
