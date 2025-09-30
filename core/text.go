package core

import (
	"regexp"
	"strings"
)

type Text struct {
	Content   string
	Color     Color
	Styles    []Style
	UseMarkup bool
}

var _ Renderable = (*Text)(nil)

var tagRegex = regexp.MustCompile(`\[([a-zA-Z ]+)\](.*?)\[/[a-zA-Z ]+\]`)

func (t Text) Render() string {
	if t.UseMarkup {
		return ParseMarkup(t.Content)
	}
	if len(t.Styles) > 0 {
		return Colorize(t.Content, t.Color, t.Styles...)
	}
	return t.Content
}

func NewText(content string, color Color, styles ...Style) Text {
	return Text{Content: content, Color: color, Styles: styles, UseMarkup: false}
}

func NewMarkup(content string) Text {
	return Text{Content: content, UseMarkup: true}
}

func ParseMarkup(input string) string {
	var builder strings.Builder
	lastIdx := 0

	matches := tagRegex.FindAllStringSubmatchIndex(input, -1)
	for _, m := range matches {
		// texte avant le tag
		builder.WriteString(input[lastIdx:m[0]])

		rawStyles := input[m[2]:m[3]]
		content := input[m[4]:m[5]]

		var styleBuilder strings.Builder
		for _, style := range strings.Fields(rawStyles) {
			if val, ok := ANSIStyles[strings.ToLower(style)]; ok {
				styleBuilder.WriteString(val)
			}
		}

		builder.WriteString(styleBuilder.String())
		builder.WriteString(content)
		builder.WriteString(ANSIStyles["reset"])

		lastIdx = m[1]
	}

	builder.WriteString(input[lastIdx:])
	return builder.String()
}
