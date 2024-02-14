package markdown

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	if workspaceDir, exist := os.LookupEnv("WORKSPACE_DIR"); exist {
		os.Chdir(workspaceDir)
	}
}

func TestParseSample(t *testing.T) {
	wd, _ := os.Getwd()
	b, e := os.ReadFile(path.Join(wd, "example/parse-test.doctask.md"))
	if e != nil {
		panic(e)
	}

	result, err := ParseMarkdown(context.Background(), b)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "Scripts for developpers", result.Title)
}
