package gosearch

import (
	"fmt"
	"strings"
)

const (
	RenderName = 1 << iota
	RenderImportPath
	RenderLicense
	RenderSynopsis
	RenderHomeSite

	notRenderTitle
	notRenderEmpty
)

const RenderUndefined = 0
const RenderAll = RenderName | RenderImportPath | RenderLicense | RenderSynopsis | RenderHomeSite

// Formatter print a pretty Package description
type Formatter interface {
	Format(p Package) string
}

type FormatterOption func(f *ShortFormatter)

type ShortFormatter struct {
	renderFlag int
}

func (s *ShortFormatter) Format(p Package) string {
	if hasBitMask(s.renderFlag, RenderName|RenderImportPath) {
		return fmt.Sprintf("%-16s %-s", p.Name, p.ImportPath)
	}
	if hasBitMask(s.renderFlag, RenderImportPath) {
		return fmt.Sprintf(p.ImportPath)
	}
	return fmt.Sprintf(p.Name)
}

func RenderParts(flags ...int) FormatterOption {
	flag := RenderUndefined
	for _, cf := range flags {
		flag |= cf
	}
	return func(f *ShortFormatter) {
		f.renderFlag |= flag
	}
}

func NewShortFormatter(options ...FormatterOption) *ShortFormatter {
	f := ShortFormatter{}
	for _, option := range options {
		option(&f)
	}
	return &f
}

type DefaultFormatter struct {
	*ShortFormatter
}

func (f *DefaultFormatter) Format(p Package) string {
	flag := f.renderFlag
	if flag == RenderUndefined {
		flag = RenderAll
	}
	isNotRenderEmpty := !hasBitMask(flag, notRenderEmpty)
	title := func() func(s string) string {
		isNotRenderTitle := hasBitMask(f.renderFlag, notRenderTitle)
		if isNotRenderTitle {
			return func(s string) string {
				return ""
			}
		} else {
			return func(s string) string {
				return s + ": "
			}
		}
	}()
	var b strings.Builder
	if hasBitMask(flag, RenderName) && (isNotRenderEmpty || p.Name != "") {
		b.WriteString(title("Name") + p.Name + "\n")
	}
	if hasBitMask(flag, RenderImportPath) && (isNotRenderEmpty || p.ImportPath != "") {
		b.WriteString(title("ImportPath") + p.ImportPath + "\n")
	}
	if hasBitMask(flag, RenderLicense) && (isNotRenderEmpty || p.License != "") {
		b.WriteString(title("License") + p.License + "\n")
	}
	if hasBitMask(flag, RenderSynopsis) && (isNotRenderEmpty || p.Synopsis != "") {
		b.WriteString(title("Synopsis") + p.Synopsis + "\n")
	}
	if hasBitMask(flag, RenderHomeSite) && (isNotRenderEmpty || p.HomeSite != "") {
		b.WriteString(title("HomeSite") + p.HomeSite + "\n")
	}
	return b.String()
}

func hasBitMask(flag, mask int) bool {
	return flag&mask == mask
}

func NotRenderTitle() FormatterOption {
	return RenderParts(notRenderTitle)
}

func NotRenderEmpty() FormatterOption {
	return RenderParts(notRenderEmpty)
}

func NewDefaultFormatter(options ...FormatterOption) *DefaultFormatter {
	return &DefaultFormatter{ShortFormatter: NewShortFormatter(options...)}
}
