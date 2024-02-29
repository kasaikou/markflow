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
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-runewidth"
	"golang.org/x/term"
)

type RecordMode byte

const (
	RecordModeNA RecordMode = '\000'
	RecordModeLF RecordMode = '\n'
	RecordModeCR RecordMode = '\r'
)

func firstLineWithWidthIndex(text string, width int, prefix int) (idx int) {
	countWidth := 0

	for i, r := range text {
		var w int
		switch r {
		case '\t':
			w = 8 - ((countWidth + prefix) % 8)
		default:
			w = runewidth.RuneWidth(r)
		}

		if countWidth+w > width {
			return i
		}
		countWidth += w
	}

	return len(text)
}

type ConsoleRecord struct {
	sender          any
	RecordMode      RecordMode
	LabelDecoration Decoration
	Kind            string // Kind.Width() <= 7
	Label           string // Label.Width() <= 19
	Text            string
	TextDecoration  Decoration
}

const RecordKindWidthLimit = 7
const RecordLabelWidthLimit = 19

var appendBytesSpaces = bytes.Repeat([]byte{' '}, 24)

func (cr *ConsoleRecord) AppendBytes(src []byte, width int) []byte {

	const (
		KindWidth   = RecordKindWidthLimit + 1
		LabelWidth  = RecordLabelWidthLimit + 1
		PrefixWidth = KindWidth + LabelWidth
	)

	text := cr.Text
	decoration := cr.TextDecoration
	if width > 0 {

		prefixBeginAt := len(src)
		src = cr.LabelDecoration.AppendBytes(src)
		src = append(src, cr.Kind...)
		src = append(src, appendBytesSpaces[:KindWidth-len(cr.Kind)]...)
		src = append(src, cr.Label...)
		src = append(src, appendBytesSpaces[:LabelWidth-len(cr.Label)]...)
		prefixEndAt := len(src)

		prefix := src[prefixBeginAt:prefixEndAt]
		src = cr.TextDecoration.AppendBytes(src)

		for {
			if idx := firstLineWithWidthIndex(text, width-PrefixWidth, PrefixWidth); idx == len(text) {
				src = append(src, text...)
				return src
			} else {
				src = append(src, text[:idx]...)
				decoration = decoration.PushString((text[:idx]))
				text = text[idx:]
				src = append(src, '\n')
				src = append(src, prefix...)
				src = decoration.AppendBytes(src)
			}
		}

	} else {
		src = cr.LabelDecoration.AppendBytes(src)
		src = append(src, cr.Kind...)
		src = append(src, appendBytesSpaces[:KindWidth-len(cr.Kind)]...)
		src = append(src, cr.Label...)
		src = append(src, appendBytesSpaces[:LabelWidth-len(cr.Label)]...)

		src = cr.TextDecoration.AppendBytes(src)
		src = append(src, text...)
		return src
	}
}

type ConsoleWriter struct {
	_        struct{}
	chRecord chan ConsoleRecord
	dest     io.Writer
	getWidth func() int
}

type LoggerOption func(*ConsoleWriter) error

func LimitedWidth(width int) LoggerOption {
	return func(cw *ConsoleWriter) error {
		cw.getWidth = func() int {
			return width
		}
		return nil
	}
}

func UnlimitedWidth() LoggerOption {
	return func(cw *ConsoleWriter) error {
		cw.getWidth = func() int { return 0 }
		return nil
	}
}

func TerminalWidth() LoggerOption {
	fd := os.Stdout.Fd()
	return func(cw *ConsoleWriter) error {
		cw.getWidth = func() int {
			width, _, err := term.GetSize(int(fd))
			if err != nil {
				panic(err)
			}

			return width
		}
		return nil
	}
}

func TerminalAutoDetect() LoggerOption {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Println("terminal mode")
		return TerminalWidth()
	} else {
		return UnlimitedWidth()
	}
}

func NewConsoleWriter(dest io.Writer, options ...LoggerOption) (*ConsoleWriter, error) {

	logger := ConsoleWriter{
		dest:     dest,
		chRecord: make(chan ConsoleRecord),
		getWidth: func() int { return 0 },
	}

	for i := range options {
		if err := options[i](&logger); err != nil {
			return nil, err
		}
	}

	return &logger, nil
}

func (cw *ConsoleWriter) Close() error {
	close(cw.chRecord)
	return nil
}

func (cw *ConsoleWriter) Route() {
	func(dest io.Writer, chRecord <-chan ConsoleRecord) {

		lfbytes := []byte{'\n'}

		prevCountLF := 0
		prev := ConsoleRecord{RecordMode: RecordModeNA}
		buffer := make([]byte, 0, 1024)

		for {
			buffer := buffer[:0]

			select {
			case record, exist := <-chRecord:
				if !exist {
					switch prev.RecordMode {
					case RecordModeCR:
						buffer = append(buffer, "\r\033[0m"...)
					case RecordModeLF:
						buffer = append(buffer, "\n\033[0m"...)
					}

					dest.Write(buffer)
					return
				}

				width := cw.getWidth()

				switch prev.RecordMode {
				case RecordModeNA:
					buffer = record.AppendBytes(buffer, width)

				case RecordModeCR:
					if record.sender == prev.sender {
						buffer = append(buffer, '\r')
						for i := 0; i < prevCountLF; i++ {
							buffer = append(buffer, "\033[1A\033[K"...)
						}
						buffer = append(buffer, "\033[K"...)
						buffer = record.AppendBytes(buffer, width)

					} else { // prev.sender != record.sender
						buffer = append(buffer, '\r')
						for i := 0; i < prevCountLF; i++ {
							buffer = append(buffer, "\033[1A\033[K"...)
						}
						buffer = append(buffer, "\033[K"...)
						prev.TextDecoration = Decoration{} // plain text
						prev.Text = "(strip output with CR by docstak)"
						buffer = prev.AppendBytes(buffer, width)
						buffer = append(buffer, '\n')
						buffer = record.AppendBytes(buffer, width)

					}

				case RecordModeLF:
					if record.RecordMode == RecordModeCR {
						if record.sender == prev.sender {
							buffer = append(buffer, '\n')
							buffer = record.AppendBytes(buffer, width)
						}

					} else { // record.RecordMode == RecordModeLF
						buffer = append(buffer, '\n')
						buffer = record.AppendBytes(buffer, width)

					}
				}

				if len(buffer) > 0 {
					dest.Write(buffer)

					if record.RecordMode == RecordModeCR {
						prevCountLF = bytes.Count(buffer, lfbytes)
					}
					prev = record
				}
			}
		}
	}(cw.dest, cw.chRecord)

}
