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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecorate(t *testing.T) {
	tests := []struct {
		From   Decoration
		Expr   string
		Expect Decoration
	}{
		{
			Expr:   "\033[3mHello world",
			Expect: Decoration{Italic: "\033[3m"},
		},
		{
			Expr:   "\033[3mHello world\033[0m",
			Expect: Decoration{},
		},
	}

	for i := range tests {
		assert.Equal(t, tests[i].Expect, tests[i].From.Push([]byte(tests[i].Expr)))
	}
}
