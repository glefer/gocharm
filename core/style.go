package core

import (
	"regexp"
	"strings"
)

// Color represents a named ANSI color.
type Color string

// Style represents a named ANSI style (bold, italic, etc.).
type Style string

// ANSI color constants
const (
	ColorBlack   Color = "black"
	ColorRed     Color = "red"
	ColorGreen   Color = "green"
	ColorYellow  Color = "yellow"
	ColorBlue    Color = "blue"
	ColorMagenta Color = "magenta"
	ColorCyan    Color = "cyan"
	ColorWhite   Color = "white"
	ColorReset   Color = "color_reset"
)

// ANSI style constants
const (
	StyleReset     Style = "reset"
	StyleBold      Style = "bold"
	StyleItalic    Style = "italic"
	StyleUnderline Style = "underline"
)

// ansiStyles maps color/style names to ANSI escape sequences. It's internal
// to avoid consumers relying on the concrete map and to allow future changes.
var ansiStyles = map[string]string{
	string(StyleReset):     "\033[0m",
	string(StyleBold):      "\033[1m",
	string(StyleItalic):    "\033[3m",
	string(StyleUnderline): "\033[4m",
	string(ColorBlack):     "\033[30m",
	string(ColorRed):       "\033[31m",
	string(ColorGreen):     "\033[32m",
	string(ColorYellow):    "\033[33m",
	string(ColorBlue):      "\033[34m",
	string(ColorMagenta):   "\033[35m",
	string(ColorCyan):      "\033[36m",
	string(ColorWhite):     "\033[37m",
	string(ColorReset):     "\033[39m",
}

// AddANSIStyle allows to add or override a custom ANSI style or color sequence.
// Example: AddANSIStyle("orange", "\033[38;5;208m")
func AddANSIStyle(name, sequence string) {
	if name != "" && sequence != "" {
		ansiStyles[name] = sequence
	}
}

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// Colorize wraps text with the ANSI sequences for color and styles, and resets styles at the end.
// If color or styles are unknown, they are ignored. If styles is nil or empty, only color is applied.
//
// Example:
//
//	Colorize("Hello", ColorRed, []Style{StyleBold})
func Colorize(text string, color Color, styles []Style) string {
	var output strings.Builder
	if val, ok := ansiStyles[string(color)]; ok && color != "" {
		output.WriteString(val)
	}
	if len(styles) > 0 {
		for _, style := range styles {
			if val, ok := ansiStyles[string(style)]; ok && style != "" {
				output.WriteString(val)
			}
		}
	}
	reset, ok := ansiStyles[string(StyleReset)]
	if !ok {
		return text
	}
	output.WriteString(text)
	output.WriteString(reset)
	return output.String()
}

// StripANSI returns s with ANSI escape sequences removed.
func StripANSI(s string) string {
	return ansiRegexp.ReplaceAllString(s, "")
}

// ExtractLeadingPrefix returns the contiguous ANSI escape sequences
// found at the start of s. If there are none, an empty string is returned.
func ExtractLeadingPrefix(s string) string {
	matches := ansiRegexp.FindAllStringIndex(s, -1)
	if len(matches) == 0 {
		return ""
	}
	end := 0
	for _, m := range matches {
		if m[0] == end {
			end = m[1]
		} else if m[0] == 0 && end == 0 {
			end = m[1]
		} else {
			break
		}
	}
	if end == 0 {
		return ""
	}
	return s[:end]
}

// EnsurePrefixPerLine applies a leading ANSI prefix (if any) to each line of
// the provided string and appends a reset sequence when needed. If a line already
// contains a leading ANSI prefix, it is preserved and not duplicated. For empty lines,
// returns an empty string (no prefix/reset). Handles edge cases robustly.
func EnsurePrefixPerLine(s string) []string {
	prefix := ExtractLeadingPrefix(s)
	lines := strings.Split(s, "\n")
	res := make([]string, len(lines))
	for i, ln := range lines {
		linePrefix := ExtractLeadingPrefix(ln)
		hasReset := strings.HasSuffix(ln, ansiStyles[string(StyleReset)])
		switch {
		case ln == "":
			res[i] = ""
		case linePrefix != "" && hasReset:
			res[i] = ln
		case linePrefix != "" && !hasReset:
			res[i] = ln + ansiStyles[string(StyleReset)]
		case strings.HasSuffix(ln, ansiStyles[string(StyleReset)]):
			res[i] = prefix + ln
		default:
			res[i] = prefix + ln + ansiStyles[string(StyleReset)]
		}
	}
	return res
}

// ColorizeLines applies Colorize to each line of the input slice.
// Useful for colorizing multi-line output with consistent style.
func ColorizeLines(lines []string, color Color, styles []Style) []string {
	res := make([]string, len(lines))
	for i, ln := range lines {
		if ln == "" {
			res[i] = ""
		} else {
			res[i] = Colorize(ln, color, styles)
		}
	}
	return res
}

// Seq returns the ANSI sequence for the given style or color name.
// It returns an empty string if the name is unknown or empty.
func Seq(name string) string {
	if name == "" {
		return ""
	}
	if v, ok := ansiStyles[name]; ok {
		return v
	}
	return ""
}

// Reset returns the ANSI reset sequence.
func Reset() string {
	return ansiStyles[string(StyleReset)]
}

// ---
// Notes:
// - User input cannot inject unexpected sequences because only known keys from ansiStyles are used.
// - Edge cases (empty lines, already colorized lines, etc.) are handled in EnsurePrefixPerLine and ColorizeLines.
// - For heavy usage, consider benchmarking EnsurePrefixPerLine.
// - See unit tests for robustness on tricky cases.
