package core

import (
	"regexp"
	"strings"
)

// Text represents a piece of text that can be rendered with color/styles
// or interpreted as markup when UseMarkup is true.
type Text struct {
	Content   string
	Color     Color
	Styles    []Style
	UseMarkup bool
}

var _ Renderable = (*Text)(nil)

var tagRegex = regexp.MustCompile(`\[([a-zA-Z ]+)\](.*?)\[/[a-zA-Z ]+\]`)

// Render returns the string representation of t, applying markup or
// ANSI color/styling as required.
func (t Text) Render() string {
	if t.UseMarkup {
		return ParseMarkup(t.Content)
	}
	if len(t.Styles) > 0 || t.Color != ColorReset {
		return Colorize(t.Content, t.Color, t.Styles)
	}
	return t.Content
}

// NewText constructs a Text instance with explicit color and styles.
func NewText(content string, color Color, styles ...Style) Text {
	return Text{Content: content, Color: color, Styles: styles, UseMarkup: false}
}

// NewMarkup constructs a Text instance whose content will be parsed as markup.
func NewMarkup(content string) Text {
	return Text{Content: content, UseMarkup: true}
}

// ParseMarkup converts a simple [style]...[/style] markup into ANSI
// escape sequences. Unknown style names are ignored.
func ParseMarkup(input string) string {
	var builder strings.Builder
	lastIdx := 0

	matches := tagRegex.FindAllStringSubmatchIndex(input, -1)
	for _, m := range matches {
		// text before the tag
		builder.WriteString(input[lastIdx:m[0]])

		rawStyles := input[m[2]:m[3]]
		content := input[m[4]:m[5]]

		var styleBuilder strings.Builder
		for _, style := range strings.Fields(rawStyles) {
			if val := Seq(strings.ToLower(style)); val != "" {
				styleBuilder.WriteString(val)
			}
		}

		builder.WriteString(styleBuilder.String())
		builder.WriteString(content)
		builder.WriteString(Reset())

		lastIdx = m[1]
	}

	builder.WriteString(input[lastIdx:])
	return builder.String()
}
