package cli

import (
	"errors"
	"log/slog"
	"os"
	"sync"
	"testing"
)

func TestConsoleWriterLoggerHandler(t *testing.T) {
	cw, _ := NewConsoleWriter(os.Stdout)
	cwWaiter := sync.WaitGroup{}
	cwWaiter.Add(1)
	go func() {
		defer cwWaiter.Done()
		cw.Route()
	}()

	defer cw.Close()
	logger := slog.New(cw.NewLoggerHandler(nil))

	for i := 0; i < 100; i++ {
		logger.Info("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", slog.String("type", "text"))
		logger.Error("errorerrorerrorerror", slog.Any("error", errors.New("error occered")))
	}
}
