package cli

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"
)

var LoggerDebugDecoration = Decoration{Faint: DC_FAINT}
var LoggerInfoDecoration = Decoration{Faint: DC_FAINT, Foreground: FG_GREEN}
var LoggerWarnDecoration = Decoration{Faint: DC_FAINT, Foreground: FG_YELLOW}
var LoggerErrorDecoration = Decoration{Foreground: FG_RED}

var stringReplacer = strings.NewReplacer("\n", " ", "\r", "")

type consoleWriterLoggerHandler struct {
	pool  sync.Pool
	level slog.Leveler
	ch    chan<- ConsoleRecord
}

func (h *consoleWriterLoggerHandler) Enabled(ctx context.Context, lv slog.Level) bool {
	return lv >= h.level.Level()
}

func appendAttr(buffer []byte, attr slog.Attr, isLast bool, keyPrefix string) []byte {

	if len(keyPrefix) > 0 {
		buffer = append(buffer, keyPrefix...)
		buffer = append(buffer, ':')
	}
	buffer = append(buffer, attr.Key...)

	buffer = append(buffer, ": "...)

	switch attr.Value.Kind() {
	case slog.KindString:
		buffer = append(buffer, '"')
		buffer = append(buffer, stringReplacer.Replace(attr.Value.String())...)
		buffer = append(buffer, '"')

	case slog.KindTime:
		buffer = append(buffer, attr.Value.Time().Format(time.RFC3339Nano)...)

	case slog.KindAny:
		switch value := attr.Value.Any().(type) {
		case error:
			buffer = append(buffer, '"')
			buffer = append(buffer, stringReplacer.Replace(value.Error())...)
			buffer = append(buffer, '"')
		default:
			buffer = append(buffer, stringReplacer.Replace(fmt.Sprint(attr.Value))...)
		}

	case slog.KindGroup:
		attrs := attr.Value.Group()
		if len(attrs) == 0 {
			buffer = append(buffer, "<empty group>"...)
		}

		buffer = append(buffer, '{')
		for i := range attrs {
			buffer = appendAttr(buffer, attrs[i], i+1 >= len(attrs), attr.Key)
		}
		buffer = append(buffer, '}')

	default:
		buffer = append(buffer, stringReplacer.Replace(fmt.Sprint(attr.Value))...)
	}

	if !isLast {
		buffer = append(buffer, ", "...)
	}

	return buffer
}

func (h *consoleWriterLoggerHandler) Handle(ctx context.Context, record slog.Record) error {
	buffer := h.pool.Get().(*[]byte)
	(*buffer) = (*buffer)[:0]
	(*buffer) = append((*buffer), record.Message...)
	(*buffer) = append((*buffer), " ("...)

	numAttrs := record.NumAttrs()
	attrsIdx := 0

	record.Attrs(func(attr slog.Attr) bool {
		attrsIdx++
		(*buffer) = appendAttr(*buffer, attr, attrsIdx >= numAttrs, "")
		return true
	})

	(*buffer) = append((*buffer), ")"...)

	var decoration Decoration
	switch record.Level {
	case slog.LevelDebug:
		decoration = LoggerDebugDecoration
	case slog.LevelInfo:
		decoration = LoggerInfoDecoration
	case slog.LevelWarn:
		decoration = LoggerWarnDecoration
	case slog.LevelError:
		decoration = LoggerErrorDecoration
	default:
		panic("unknown logging level")
	}

	h.ch <- ConsoleRecord{
		sender:          h,
		RecordMode:      RecordModeLF,
		LabelDecoration: decoration,
		Kind:            "DOCSTAK",
		Label:           record.Level.String(),
		Text:            string(*buffer),
		TextDecoration:  decoration,
	}

	return nil
}

type ConsoleWriterLoggerHandler struct {
	handler   *consoleWriterLoggerHandler
	attrs     []slog.Attr
	withGroup string
}

func (cw *ConsoleWriter) NewLoggerHandler(level slog.Leveler) *ConsoleWriterLoggerHandler {

	if level == nil {
		level = slog.LevelInfo
	}

	handler := &consoleWriterLoggerHandler{
		pool: sync.Pool{
			New: func() any {
				buffer := make([]byte, 0, 1024)
				return &buffer
			},
		},
		ch:    cw.chRecord,
		level: level,
	}

	return &ConsoleWriterLoggerHandler{handler: handler}
}

func (h ConsoleWriterLoggerHandler) Enabled(ctx context.Context, lv slog.Level) bool {
	return h.handler.Enabled(ctx, lv)
}

func (h ConsoleWriterLoggerHandler) Handle(ctx context.Context, record slog.Record) error {
	record.AddAttrs(h.attrs...)
	return h.handler.Handle(ctx, record)
}

func (h ConsoleWriterLoggerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) > 0 {
		if h.withGroup != "" {
			for i := range attrs {
				h.attrs = append(h.attrs, slog.Group(h.withGroup, attrs[i]))
			}
		} else {
			h.attrs = append(h.attrs, attrs...)
		}
	}

	return h
}

func (h ConsoleWriterLoggerHandler) WithGroup(name string) slog.Handler {
	if name != "" {
		h.withGroup = name
	}

	return h
}
