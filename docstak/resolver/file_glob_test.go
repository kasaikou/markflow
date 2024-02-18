package resolver

import (
	"fmt"
	"os"
	"testing"
)

func init() {
	if workspaceDir, exist := os.LookupEnv("DOCSTAK_TEST_WORKSPACE_DIR"); exist {
		os.Chdir(workspaceDir)
	}
}

func TestGlobRules(t *testing.T) {

	wd, _ := os.Getwd()
	fmt.Println(ResolveFileGlob(FileGlobConfig{
		Rootdir:    wd,
		Rules:      []string{"**/*.go"},
		IgnoreRule: []string{"**/*_test.*"},
	}))
}
