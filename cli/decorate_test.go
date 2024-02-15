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
