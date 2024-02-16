package cli

import (
	"bytes"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestConsoleWriterScanner(t *testing.T) {
	cw, _ := NewConsoleWriter(os.Stdout)
	cwWaiter := sync.WaitGroup{}
	defer cwWaiter.Wait()
	cwWaiter.Add(1)
	go func() {
		defer cwWaiter.Done()
		cw.Route()
	}()

	defer cw.Close()
	cw.NewScanner(Decoration{}, "TEST", "test_str").
		Scan(bytes.NewBufferString(strings.Repeat("012345678901234567890123456789\n", 1000)))
}
