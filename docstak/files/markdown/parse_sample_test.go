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
