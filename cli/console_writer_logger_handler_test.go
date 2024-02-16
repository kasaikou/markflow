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
