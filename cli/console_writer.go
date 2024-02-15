package cli

import (
	"bytes"
	"io"

	"github.com/mattn/go-runewidth"
)

type RecordMode byte

const (
	RecordModeNA = '\000'
	RecordModeLF = '\n'
	RecordModeCR = '\r'
)

type StringWidth struct {
	text  string
	width int
}

func NewStringWidth(text string) StringWidth {
	return StringWidth{
		text:  text,
		width: runewidth.StringWidth(text),
	}
}

func (s *StringWidth) String() string { return s.text }
func (s *StringWidth) Width() int     { return s.width }

func firstLineWithWidthIndex(text string, width int) (idx int) {
	countWidth := 0
	for i, r := range text {
		w := runewidth.RuneWidth(r)
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
	Kind            StringWidth // Kind.Width() <= 7
	Label           StringWidth // Label.Width() <= 19
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

	prefixBeginAt := len(src)
	src = cr.LabelDecoration.AppendBytes(src)
	src = append(src, cr.Kind.String()...)
	src = append(src, appendBytesSpaces[:KindWidth-cr.Kind.Width()]...)
	src = append(src, cr.Label.String()...)
	src = append(src, appendBytesSpaces[:LabelWidth-cr.Label.Width()]...)
	prefixEndAt := len(src)
	prefix := src[prefixBeginAt:prefixEndAt]
	src = cr.TextDecoration.AppendBytes(src)

	for {
		if idx := firstLineWithWidthIndex(text, width-PrefixWidth); idx == len(text) {
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
}

type ConsoleWriter struct {
	_            struct{}
	defaultWidth int
	chRecord     chan ConsoleRecord
	dest         io.Writer
}

type LoggerOption func(*ConsoleWriter) error

func NewConsoleWriter(dest io.Writer, options ...LoggerOption) (*ConsoleWriter, error) {

	logger := ConsoleWriter{
		dest:         dest,
		defaultWidth: 128,
		chRecord:     make(chan ConsoleRecord),
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

		width := cw.defaultWidth
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
						buffer = append(buffer, '\r')
					case RecordModeLF:
						buffer = append(buffer, '\n')
					}

					dest.Write(buffer)
					return
				}

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
						prev.TextDecoration = Decoration{}
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
				}

				if record.RecordMode == RecordModeCR {
					prevCountLF = bytes.Count(buffer, lfbytes)
				}
				prev = record
			}
		}
	}(cw.dest, cw.chRecord)

}
