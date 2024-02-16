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
	"regexp"
	"strings"
	"unsafe"
)

type Decoration struct {
	Bold       string
	Faint      string
	Italic     string
	Underline  string
	Display    string
	Foreground string
	Background string
}

const (
	DC_RESET          = "\033[0m"
	DC_BOLD           = "\033[1m"
	DC_FAINT          = "\033[2m"
	DC_ITALIC         = "\033[3m"
	DC_UNDERLINE      = "\033[4m"
	DC_HIDE           = "\033[8m"
	DC_NORM_INTENSITY = "\033[22m"
	DC_NOT_ITALIC     = "\033[23m"
	DC_NOT_UNDERLINE  = "\033[24m"
	DC_REVEAL         = "\033[28m"
)

const (
	FG_BLACK   = "\033[30m"
	FG_RED     = "\033[31m"
	FG_GREEN   = "\033[32m"
	FG_YELLOW  = "\033[33m"
	FG_BLUE    = "\033[34m"
	FG_MAGENTA = "\033[35m"
	FG_CYAN    = "\033[36m"
	FG_WHITE   = "\033[37m"
	FG_RESET   = "\033[39m"
	BG_BLACK   = "\033[40m"
	BG_RED     = "\033[41m"
	BG_GREEN   = "\033[42m"
	BG_YELLOW  = "\033[43m"
	BG_BLUE    = "\033[44m"
	BG_MAGENTA = "\033[45m"
	BG_CYAN    = "\033[46m"
	BG_WHITE   = "\033[47m"
	BG_RESET   = "\033[49m"
)

var decorateRegexp = regexp.MustCompile(`(\033\[(([0-48])|(2[2-48])|(([3-49]|10)[0-79])|([34]8;5;[0-9]{1,3})|([34]8;2;[0-9]{1,3};[0-9]{1,3};[0-9]{1,3})))m`)

func (d Decoration) Push(expr []byte) Decoration {
	begin := 0
	first := 0
	last := 0

	for idx := decorateRegexp.FindIndex(expr[begin:]); idx != nil; idx = decorateRegexp.FindIndex(expr[begin:]) {
		first = begin + idx[0]
		last = begin + idx[1]
		expr := unsafe.String(unsafe.SliceData(expr[first:last]), last-first)
		d = d.update(expr)

		begin = last
	}

	return d
}

func (d Decoration) PushString(expr string) Decoration {
	begin := 0
	first := 0
	last := 0

	for idx := decorateRegexp.FindStringIndex(expr[begin:]); idx != nil; idx = decorateRegexp.FindStringIndex(expr[begin:]) {
		first = begin + idx[0]
		last = begin + idx[1]
		expr := expr[first:last]
		d = d.update(expr)

		begin = last
	}

	return d
}

func (d Decoration) update(expr string) Decoration {
	switch {
	case expr == DC_RESET:
		d.Bold = ""
		d.Faint = ""
		d.Italic = ""
		d.Underline = ""
		d.Display = ""
		d.Foreground = ""
		d.Background = ""

	case expr == DC_BOLD:
		d.Bold = expr

	case expr == DC_FAINT:
		d.Faint = expr

	case expr == DC_ITALIC:
		d.Italic = expr

	case expr == DC_UNDERLINE:
		d.Underline = expr

	case expr == DC_HIDE:
		d.Display = expr

	case expr == DC_NORM_INTENSITY:
		d.Bold = ""
		d.Faint = ""

	case expr == DC_NOT_ITALIC:
		d.Italic = expr

	case expr == DC_NOT_UNDERLINE:
		d.Underline = expr

	case expr == DC_REVEAL:
		d.Display = expr

	case expr == FG_RESET:
		d.Foreground = ""

	case expr == BG_RESET:
		d.Background = ""

	case strings.HasPrefix(expr, "\033[3") || strings.HasPrefix(expr, "\033[9"):
		d.Foreground = expr

	case strings.HasPrefix(expr, "\033[4") || strings.HasPrefix(expr, "\033[10"):
		d.Background = expr
	}

	return d
}

func (d *Decoration) AppendBytes(src []byte) []byte {

	src = append(src, "\033[0m"...)

	if len(d.Bold) > 0 {
		src = append(src, d.Bold...)
	}

	if len(d.Faint) > 0 {
		src = append(src, d.Faint...)
	}

	if len(d.Italic) > 0 {
		src = append(src, d.Italic...)
	}

	if len(d.Underline) > 0 {
		src = append(src, d.Italic...)
	}

	if len(d.Display) > 0 {
		src = append(src, d.Display...)
	}

	if len(d.Foreground) > 0 {
		src = append(src, d.Foreground...)
	}

	if len(d.Background) > 0 {
		src = append(src, d.Background...)
	}

	return src
}
