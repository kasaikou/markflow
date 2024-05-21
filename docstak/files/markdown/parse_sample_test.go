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
	"encoding/json"
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

	expect := ParseResult{
		Title: "Scripts for github.com/kasaikou/markflow developpers",
		Tasks: []ParseResultTask{
			{
				Title:        "Scripts for github.com/kasaikou/markflow developpers",
				HeadingLevel: 1,
			},
			{
				Title:        "hello_world",
				HeadingLevel: 2,
				Description:  "Echo \"Hello World\"",
				Commands: []ParseResultCommand{{
					Lang: "sh",
					Code: "echo \"Hello World, docstak!\"\n",
				}},
			},
			{
				Title:        "download",
				HeadingLevel: 2,
				Description:  "Download dependencies",
				Config: ParseResultTaskConfig{
					Requires: ParseResultTaskConfigRequires{
						File: ParseResultTaskConfigFiles{
							Exists: []string{"go.mod", "go.sum"},
						},
					},
					Skips: ParseResultTaskConfigSkips{
						File: ParseResultTaskConfigFiles{
							NotChangeds: []string{"go.sum"},
						},
					},
				},
				Commands: []ParseResultCommand{{
					Lang: "sh",
					Code: "go mod download\n",
				}},
			},
			{
				Title:        "test",
				HeadingLevel: 2,
				Description:  "Run go test",
				Config: ParseResultTaskConfig{
					Previous: []string{"download"},
				},
				Commands: []ParseResultCommand{{
					Lang: "sh",
					Code: "DOCSTAK_TEST_WORKSPACE_DIR=$(pwd) go test ./...\n",
				}},
			},
			{
				Title:        "fmt",
				HeadingLevel: 2,
				Description:  "Format source codes",
				Commands: []ParseResultCommand{{
					Lang: "sh",
					Code: "go fmt ./...\n",
				}},
			},
			{
				Title:        "ci",
				HeadingLevel: 2,
				Description:  "Running on GitHub Actions, local, and so on.",
				Config: ParseResultTaskConfig{
					Previous: []string{"ci/fmt", "ci/depends", "ci/coverage-test"},
				},
			},
			{
				Title:        "ci/depends",
				HeadingLevel: 3,
				Config: ParseResultTaskConfig{
					Skips: ParseResultTaskConfigSkips{
						File: ParseResultTaskConfigFiles{
							NotChangeds: []string{"**.go", "go.sum", "go.mod"},
						},
					},
				},
				Commands: []ParseResultCommand{{
					Lang: "sh",
					Code: "go mod tidy &&\ngit diff --no-patch --exit-code go.sum\n",
				}},
			},
			{
				Title:        "ci/fmt",
				HeadingLevel: 3,
				Commands: []ParseResultCommand{{
					Lang: "sh",
					Code: "gofmt -l .\n",
				}},
			},
			{
				Title:        "ci/coverage-test",
				HeadingLevel: 3,
				Config: ParseResultTaskConfig{
					Previous: []string{"ci/coverage-test/go", "download"},
				},
			},
			{
				Title:        "ci/coverage-test/go",
				HeadingLevel: 4,
				Config: ParseResultTaskConfig{
					Skips: ParseResultTaskConfigSkips{
						File: ParseResultTaskConfigFiles{
							NotChangeds: []string{"**.go", "go.sum", "go.mod"},
						},
					},
				},
				Commands: []ParseResultCommand{{
					Lang: "sh",
					Code: "go test -coverprofile=coverage.txt -covermode=atomic ./...\n",
				}},
			},
		},
	}

	// Convert to JSON and compare and verify.
	resultJson, _ := json.MarshalIndent(result, "", "  ")
	expectJson, _ := json.MarshalIndent(expect, "", "  ")
	assert.Equal(t, string(expectJson), string(resultJson))

}
