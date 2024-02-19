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
	"fmt"
	"os"
	"regexp"
	"unsafe"

	"github.com/cockroachdb/errors"
	"github.com/kasaikou/docstak/docstak/resolver"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"gopkg.in/yaml.v3"
)

type FileEnvironment struct {
	Filepath string
}

type ParseResult struct {
	Title       string
	Description string
	Tasks       []ParseResultTask
	Config      ParseResultGlobalConfig
}

type ParseResultGlobalConfig struct {
	Root    string                    `json:"root" yaml:"root"`
	Environ ParseResultTaskConfigEnvs `json:"environ" yaml:"environ"`
}

type ParseResultTask struct {
	Title       string
	Description string
	Config      ParseResultTaskConfig
	Commands    []ParseResultCommand
	Internals   []ParseResultTask
}

type ParseResultTaskConfig struct {
	Environ  ParseResultTaskConfigEnvs     `json:"environ" yaml:"environ"`
	Requires ParseResultTaskConfigRequires `json:"requires" yaml:"requires"`
	Skips    ParseResultTaskConfigSkips    `json:"skips" yaml:"skips"`
}

type ParseResultTaskConfigEnvs struct {
	Dotenvs   []string          `json:"dotenv" yaml:"dotenv"`
	Variables map[string]string `json:"vars" yaml:"vars"`
}

type ParseResultTaskConfigSkips struct {
	Exists []string `json:"exist" yaml:"exist"`
}

type ParseResultTaskConfigRequires struct {
	Exists []string `json:"exist" yaml:"exist"`
}

type ParseResultCommand struct {
	Lang string
	Code string
}

var (
	yamlConfigRule = regexp.MustCompile(`^ya?ml:docstak.ya?ml$`)
)

type MarkdownOption struct {
	filename string
	bytes    []byte
}

func (mo MarkdownOption) Filename() string { return mo.filename }

func FromLocalFile(workingDir, searchname string) (MarkdownOption, error) {
	filename, exist := resolver.ResolveFileWithBasename(workingDir, searchname)
	if !exist {
		return MarkdownOption{}, errors.Errorf("cannot resolve file '%s' in directory '%s'", searchname, workingDir)
	}

	b, err := os.ReadFile(filename)
	if err != nil {
		return MarkdownOption{}, errors.WithMessagef(err, "cannot read file '%s'", filename)
	}

	return MarkdownOption{filename: filename, bytes: b}, nil
}

func ParseMarkdown(ctx context.Context, markdown MarkdownOption) (ParseResult, error) {
	result := ParseResult{}
	node := goldmark.DefaultParser().Parse(text.NewReader(markdown.bytes))

	if node.Kind() != ast.KindDocument {
		panic(fmt.Sprintf("unexpected node kind (want: %v, have: %v)", ast.KindDocument, node.Kind()))
	}

	topHeaderLevel := 0
	selectedElemIdx := -1
	selectedElem := []int{}
	var selected *ParseResultTask
	getSelected := func() *ParseResultTask {
		if len(selectedElem) == 0 {
			return nil
		}

		if selectedElem[0] == len(result.Tasks) {
			result.Tasks = append(result.Tasks, ParseResultTask{})
		} else {
			panic("jumped element index")
		}

		result := &result.Tasks[selectedElem[0]]
		for i := 1; i < selectedElemIdx; i++ {
			if selectedElem[i] > len(result.Internals) {
				panic("jumped element index")
			} else if selectedElem[i] == len(result.Internals) {
				result.Internals = append(result.Internals, ParseResultTask{})
			}
			result = &result.Internals[selectedElem[i]]
		}

		return result
	}

	for node := node.FirstChild(); node != nil; node = node.NextSibling() {
		switch node.Kind() {
		case ast.KindHeading:
			node := node.(*ast.Heading)
			title := node.Text(markdown.bytes)
			titleStr := unsafe.String(unsafe.SliceData(title), len(title))
			if topHeaderLevel < 1 {
				topHeaderLevel = node.Level
				result.Title = titleStr

			} else { // topHeaderLevel >= 1 (title has been already defined)
				if node.Level <= topHeaderLevel {
					return result, errors.Errorf("markdown has multiple document title ('%s' and '%s')", result.Title, string(node.Text(markdown.bytes)))
				}

				elemLevel := node.Level - topHeaderLevel
				if elemLevel > len(selectedElem) {
					if want := len(selectedElem) + 1; elemLevel > want {
						return result, errors.Errorf("jumped header level (want: %d, have: %d) with text '%s'", want, elemLevel, string(node.Text(markdown.bytes)))
					}

					selectedElem = append(selectedElem, 0)
					selectedElemIdx = elemLevel - 1

				} else { // elemLevel <= len(selectedElem)
					selectedElemIdx = elemLevel - 1
					selectedElem[selectedElemIdx]++
				}

				selected = getSelected()
				selected.Title = titleStr
			}

		case ast.KindParagraph:
			node := node.(*ast.Paragraph)
			desc := node.Text(markdown.bytes)
			descStr := unsafe.String(unsafe.SliceData(desc), len(desc))

			if descStr == "" {
				break
			}

			if selected == nil {
				if result.Description == "" {
					result.Description = descStr
				} else {
					result.Description = result.Description + "\n\n" + descStr
				}
			} else { // selected != nil
				if selected.Description == "" {
					selected.Description = descStr
				} else {
					selected.Description = selected.Description + "\n\n" + descStr
				}
			}

		case ast.KindFencedCodeBlock:
			node := node.(*ast.FencedCodeBlock)
			code := func() []byte {
				lines := node.Lines()
				if lines.Len() == 0 {
					return []byte{}
				} else {
					return markdown.bytes[lines.At(0).Start:lines.At(lines.Len()-1).Stop]
				}
			}()
			codeStr := unsafe.String(unsafe.SliceData(code), len(code))
			lang := node.Language(markdown.bytes)
			langStr := unsafe.String(unsafe.SliceData(lang), len(lang))

			if selected == nil {
				if yamlConfigRule.Match(lang) {
					if err := yaml.Unmarshal(code, &result.Config); err != nil {
						return result, errors.WithMessage(err, "failed to parse yaml format global config")
					}
				}
			} else { // selected != nil
				if yamlConfigRule.Match(lang) {
					if err := yaml.Unmarshal(code, &selected.Config); err != nil {
						return result, err
					}
				} else { // yamlConfigRule.Match(lang) == false
					selected.Commands = append(selected.Commands, ParseResultCommand{
						Lang: langStr,
						Code: codeStr,
					})
				}
			}
		}
	}
	return result, nil
}
