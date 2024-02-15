package cli

import (
	"regexp"
	"strings"
	"unsafe"
)

type Decoration struct {
	Bold       string
	Italic     string
	Underline  string
	Display    string
	Foreground string
	Background string
}

var decorateRegexp = regexp.MustCompile(`(\033\[(([012348])|(2[2-48])|(([3-49]|10)[0-79])|([34]8;5;[0-9]{1,3})|([34]8;2;[0-9]{1,3};[0-9]{1,3};[0-9]{1,3})))m`)

func (d Decoration) Push(expr []byte) Decoration {
	begin := 0
	first := 0
	last := 0

	for idx := decorateRegexp.FindIndex(expr[begin:]); idx != nil; idx = decorateRegexp.FindIndex(expr[begin:]) {
		first = begin + idx[0]
		last = begin + idx[1]
		expr := unsafe.String(unsafe.SliceData(expr[first:last]), last-first)

		switch {
		case expr == "\033[0m":
			d.Bold = ""
			d.Italic = ""
			d.Underline = ""
			d.Display = ""
			d.Foreground = ""
			d.Background = ""

		case expr == "\033[1m":
			if strings.HasSuffix(d.Bold, "\033[2m") {
				d.Bold = strings.TrimSuffix(d.Bold, "\033[2m")
			} else {
				d.Bold += expr
			}

		case expr == "\033[2m":
			if strings.HasSuffix(d.Bold, "\033[1m") {
				d.Bold = strings.TrimSuffix(d.Bold, "\033[1m")
			} else {
				d.Bold += expr
			}

		case expr == "\033[3m":
			if strings.HasSuffix(d.Italic, "\033[23m") {
				d.Italic = strings.TrimSuffix(d.Italic, "\033[23m")
			} else {
				d.Italic = expr
			}

		case expr == "\033[4m":
			if strings.HasSuffix(d.Underline, "\033[24m") {
				d.Underline = strings.TrimSuffix(d.Underline, "\033[24m")
			} else {
				d.Underline = expr
			}

		case expr == "\033[8m":
			if strings.HasSuffix(d.Display, "\033[28m") {
				d.Display = strings.TrimSuffix(d.Display, "\033[28m")
			} else {
				d.Display = expr
			}

		case expr == "\033[22m":
			d.Bold = ""

		case expr == "\033[23m":
			if strings.HasSuffix(d.Italic, "\033[3m") {
				d.Italic = strings.TrimSuffix(d.Italic, "\033[3m")
			} else {
				d.Italic = expr
			}

		case expr == "\033[24m":
			if strings.HasSuffix(d.Underline, "\033[4m") {
				d.Underline = strings.TrimSuffix(d.Underline, "\033[4m")
			} else {
				d.Underline = expr
			}

		case expr == "\033[28m":
			if strings.HasSuffix(d.Display, "\033[8m") {
				d.Display = strings.TrimSuffix(d.Display, "\033[8m")
			} else {
				d.Display = expr
			}

		case expr == "\033[39m":
			d.Foreground = ""

		case expr == "\033[49m":
			d.Background = ""

		case strings.HasPrefix(expr, "\033[3") || strings.HasPrefix(expr, "\033[9"):
			d.Foreground = expr

		case strings.HasPrefix(expr, "\033[4") || strings.HasPrefix(expr, "\033[10"):
			d.Background = expr
		}

		begin = last
	}

	return d
}
