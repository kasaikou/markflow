package cli

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
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
			Kind:       NewStringWidth("TEST"),
			Label:      NewStringWidth(fmt.Sprintf("TEST_%d", i)),
			LabelDecoration: Decoration{
				Background: fmt.Sprintf("\033[3%dm", i%8),
			},
		})
	}
}
func TestConsoleWriterCR(t *testing.T) {

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

	for i := 0; i < 1; i++ {
		taskWg.Add(1)
		go func(ch chan<- ConsoleRecord, record ConsoleRecord) {
			defer taskWg.Done()
			for i := 0; i < 100; i++ {
				record.Text = strings.Repeat("a", rand.Intn(200))
				ch <- record
				time.Sleep(10 * time.Millisecond)
			}

		}(cw.chRecord, ConsoleRecord{
			RecordMode: RecordModeCR,
			Kind:       NewStringWidth("TEST"),
			Label:      NewStringWidth(fmt.Sprintf("TEST_%d", i)),
			LabelDecoration: Decoration{
				Background: fmt.Sprintf("\033[3%dm", i%8),
			},
		})
	}
}
