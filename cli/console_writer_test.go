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
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConsoleWriterLF(t *testing.T) {

	wg := sync.WaitGroup{}
	defer wg.Wait()

	cw, _ := NewConsoleWriter(os.Stdout)
	defer cw.Close()

	wg.Add(1)
	go func() {
		defer wg.Done()
		cw.Route()
	}()

	taskWg := sync.WaitGroup{}
	defer taskWg.Wait()

	for i := 0; i < 100; i++ {
		taskWg.Add(1)
		go func(ch chan<- ConsoleRecord, record ConsoleRecord) {
			defer taskWg.Done()
			for i := 0; i < 100; i++ {
				record.Text = strings.Repeat("a", rand.Intn(200))
				ch <- record
			}

		}(cw.chRecord, ConsoleRecord{
			RecordMode: RecordModeLF,
			Kind:       "TEST",
			Label:      fmt.Sprintf("TEST_%d", i),
			LabelDecoration: Decoration{
				Background: fmt.Sprintf("\033[3%dm", i%8),
			},
		})
	}
}
func TestConsoleWriterCR(t *testing.T) {

	wg := sync.WaitGroup{}
	defer wg.Wait()

	cw, _ := NewConsoleWriter(os.Stdout, LimitedWidth(50))
	defer cw.Close()

	wg.Add(1)
	go func() {
		defer wg.Done()
		cw.Route()
	}()

	taskWg := sync.WaitGroup{}
	defer taskWg.Wait()

	for i := 0; i < 1; i++ {
		taskWg.Add(1)
		go func(ch chan<- ConsoleRecord, record ConsoleRecord) {
			defer taskWg.Done()
			for i := 0; i < 100; i++ {
				record.Text = strings.Repeat("\t0123456789", rand.Intn(20))
				ch <- record
				time.Sleep(10 * time.Millisecond)
			}

		}(cw.chRecord, ConsoleRecord{
			sender:     &taskWg,
			RecordMode: RecordModeLF,
			Kind:       "TEST",
			Label:      fmt.Sprintf("TEST_%d", i),
			LabelDecoration: Decoration{
				Background: fmt.Sprintf("\033[3%dm", i%8),
			},
		})
	}
}

func TestFirstLineWithWidthIndex(t *testing.T) {

	tests := []struct {
		Text          string
		Width, Except int
	}{
		{
			Text:   "\taaa",
			Width:  16,
			Except: 4,
		},
		{
			Text:   "aaaaaaaaaa\taaaaaaaaaaaaa",
			Width:  20,
			Except: 15,
		},
		{
			Text:   "aaaaaaaaaaaaaaaaaaaaaaaaaa\taaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\taaaaaaaaaaaaaaaaa",
			Width:  79,
			Except: 68,
		},
	}

	for _, test := range tests {
		idx := firstLineWithWidthIndex(test.Text, test.Width, 0)
		assert.Equal(t, test.Except, idx, fmt.Sprintf("text: '%s'", test.Text[:idx]))
	}
}
