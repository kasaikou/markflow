package markdown

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	if workspaceDir, exist := os.LookupEnv("DOCSTAK_TEST_WORKSPACE_DIR"); exist {
		os.Chdir(workspaceDir)
	}
}

func TestParseSample(t *testing.T) {
	wd, _ := os.Getwd()
	po, err := FromLocalFile(wd, "docstak.md")
	if !assert.NoError(t, err) {
		return
	}

	result, err := ParseMarkdown(context.Background(), po)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "Scripts for github.com/kasaikou/docstak developpers", result.Title)
}
