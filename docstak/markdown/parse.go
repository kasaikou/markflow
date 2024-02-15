package markdown

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"unsafe"

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
}

type ParseResultTask struct {
	Title       string
	Description string
	Config      ParseResultTaskConfig
	Commands    []ParseResultCommand
	Internals   []ParseResultTask
}

type ParseResultTaskConfig struct {
}

type ParseResultCommand struct {
	Lang string
	Code string
}

var (
	yamlConfigRule = regexp.MustCompile(`^ya?ml:doctask.ya?ml$`)
)

func ParseMarkdown(ctx context.Context, markdown []byte) (ParseResult, error) {
	result := ParseResult{}
	node := goldmark.DefaultParser().Parse(text.NewReader(markdown))

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
			title := node.Text(markdown)
			titleStr := unsafe.String(unsafe.SliceData(title), len(title))
			if topHeaderLevel < 1 {
				topHeaderLevel = node.Level
				result.Title = titleStr

			} else { // topHeaderLevel >= 1 (title has been already defined)
				if node.Level <= topHeaderLevel {
					return result, errors.New(fmt.Sprintf("markdown has multiple document title ('%s' and '%s')", result.Title, string(node.Text(markdown))))
				}

				elemLevel := node.Level - topHeaderLevel
				if elemLevel > len(selectedElem) {
					if want := len(selectedElem) + 1; elemLevel > want {
						return result, errors.New(fmt.Sprintf("jumped header level (want: %d, have: %d) with text '%s'", want, elemLevel, string(node.Text(markdown))))
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
			desc := node.Text(markdown)
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
					return markdown[lines.At(0).Start:lines.At(lines.Len()-1).Stop]
				}
			}()
			codeStr := unsafe.String(unsafe.SliceData(code), len(code))
			lang := node.Language(markdown)
			langStr := unsafe.String(unsafe.SliceData(lang), len(lang))

			if selected == nil {

			} else { // selected != nil
				if yamlConfigRule.Match(lang) {
					if err := yaml.Unmarshal(code, &selected.Config); err != nil {
						return result, err
					}
				} else {
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
