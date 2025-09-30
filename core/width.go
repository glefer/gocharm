package core

import (
	"strings"
	"unicode"
)

// RuneWidth returns the display width for r when rendered in a terminal.
// It follows common East Asian width rules and treats control and
// non-spacing marks as zero-width.
func RuneWidth(r rune) int {
	switch {
	case r == 0, unicode.IsControl(r),
		unicode.Is(unicode.Mn, r),
		unicode.Is(unicode.Me, r),
		unicode.Is(unicode.Cf, r):
		return 0
	case r >= 0xFE00 && r <= 0xFE0F:
		return 0
	case (r >= 0x1100 && r <= 0x115F) ||
		(r >= 0x2329 && r <= 0x232A) ||
		(r >= 0x2E80 && r <= 0xA4CF) ||
		(r >= 0xAC00 && r <= 0xD7A3) ||
		(r >= 0xF900 && r <= 0xFAFF) ||
		(r >= 0xFE30 && r <= 0xFE6F) ||
		(r >= 0xFF00 && r <= 0xFF60) ||
		(r >= 0xFFE0 && r <= 0xFFE6) ||
		(r >= 0x2600 && r <= 0x26FF) ||
		(r >= 0x2700 && r <= 0x27BF) ||
		(r >= 0x1F300 && r <= 0x1FAFF):
		return 2
	}
	return 1
}

// VisibleLen returns the number of columns s occupies when printed to a
// terminal, ignoring ANSI escape sequences.
func VisibleLen(s string) int {
	clean := StripANSI(s)
	l := 0
	for _, r := range clean {
		l += RuneWidth(r)
	}
	return l
}

// Alignment controls horizontal alignment inside a padded field.
type Alignment int

const (
	AlignLeft Alignment = iota
	AlignCenter
	AlignRight
)

// Pad returns s placed into a field of the given width, aligned according to
// align and surrounded by padding spaces on both sides.
// width refers to the content width (excluding the added padding).
func Pad(s string, width int, align Alignment, padding int) string {
	vis := VisibleLen(s)
	contentWidth := width
	if contentWidth < 0 {
		contentWidth = 0
	}
	if vis >= contentWidth {
		return strings.Repeat(" ", padding) + s + strings.Repeat(" ", padding)
	}
	space := contentWidth - vis
	var inner string
	switch align {
	case AlignRight:
		inner = strings.Repeat(" ", space) + s
	case AlignCenter:
		left := space / 2
		right := space - left
		inner = strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
	default:
		inner = s + strings.Repeat(" ", space)
	}
	return strings.Repeat(" ", padding) + inner + strings.Repeat(" ", padding)
}
