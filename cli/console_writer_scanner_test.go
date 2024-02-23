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
	"bytes"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestAdjustLabel(t *testing.T) {
	assert.Equal(t, "abcdef/...n/opqrstu", adjustLabel("abcdef/ghijklmn/opqrstu"))
}
