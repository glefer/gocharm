package core

import (
	"testing"
)

func TestColorize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		color    Color
		styles   []Style
		expected string
	}{
		{
			name:     "Color only",
			input:    "hello",
			color:    ColorRed,
			styles:   nil,
			expected: "\033[31mhello\033[0m",
		},
		{
			name:     "Color and style",
			input:    "world",
			color:    ColorGreen,
			styles:   []Style{StyleBold},
			expected: "\033[32m\033[1mworld\033[0m",
		},
		{
			name:     "Multiple styles",
			input:    "test",
			color:    ColorBlue,
			styles:   []Style{StyleBold, StyleUnderline},
			expected: "\033[34m\033[1m\033[4mtest\033[0m",
		},
		{
			name:     "No color, one style",
			input:    "plain",
			color:    "",
			styles:   []Style{StyleBold},
			expected: "\033[1mplain\033[0m",
		},
		{
			name:     "Unknown color and style",
			input:    "unknown",
			color:    "unknowncolor",
			styles:   []Style{"unknownstyle"},
			expected: "unknown\033[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{tt.input}
			if tt.color != "" {
				args = append(args, string(tt.color))
			}
			for _, s := range tt.styles {
				args = append(args, string(s))
			}
			result := Colorize(tt.input, tt.color, tt.styles...)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestANSIStylesMapContainsAllConstants(t *testing.T) {
	colors := []Color{ColorRed, ColorGreen, ColorYellow, ColorBlue, ColorBlack, ColorMagenta, ColorCyan, ColorWhite}
	styles := []Style{StyleBold, StyleItalic, StyleUnderline}
	for _, c := range colors {
		if _, ok := ANSIStyles[string(c)]; !ok {
			t.Errorf("ANSIStyles missing color: %s", c)
		}
	}
	for _, s := range styles {
		if _, ok := ANSIStyles[string(s)]; !ok {
			t.Errorf("ANSIStyles missing style: %s", s)
		}
	}
	if _, ok := ANSIStyles[string(StyleReset)]; !ok {
		t.Errorf("ANSIStyles missing reset style")
	}
}
