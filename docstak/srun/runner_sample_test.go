package srun

import (
	"context"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunnerSample(t *testing.T) {
	runner := NewScriptRunner("/bin/sh", "echo 'hello world' &&\n echo 'good bye'\n")
	stdout, _ := runner.Stdout()
	stderr, _ := runner.Stderr()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		io.Copy(os.Stdout, stdout)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		io.Copy(os.Stderr, stderr)
	}()

	_, err := runner.RunContext(context.Background())
	assert.NoError(t, err)
	wg.Wait()
}
