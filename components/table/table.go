package table

import (
	"fmt"
	"strings"

	"github.com/glefer/gocharm/core"
)

type Table struct {
	Headers []core.Renderable
	Rows    [][]core.Renderable
	Aligns  []core.Alignment
	Padding int
	Border  BorderStyle
}

func (t *Table) WithAlignments(a ...core.Alignment) *Table {
	t.Aligns = make([]core.Alignment, len(a))
	copy(t.Aligns, a)
	return t
}

func (t *Table) WithPadding(p int) *Table {
	if p < 0 {
		p = 0
	}
	t.Padding = p
	return t
}

// StringHeader is a convenience type for header values provided as strings.
type StringHeader string

// ToRenderable converts a StringHeader into a core.Renderable with bold style.
func (s StringHeader) ToRenderable() core.Renderable {
	return core.NewText(string(s), core.ColorReset, core.StyleBold)
}

type RenderableAdapter struct {
	core.Renderable
}

func (r RenderableAdapter) ToRenderable() core.Renderable {
	return r.Renderable
}
func anyToRenderable(h any) core.Renderable {
	switch v := h.(type) {
	case string:
		return StringHeader(v).ToRenderable()
	case core.Renderable:
		return RenderableAdapter{v}.ToRenderable()
	case int, int8, int16, int32, int64, float32, float64, uint, uint8, uint16, uint32, uint64:
		return core.NewText(fmt.Sprint(v), core.ColorReset, core.StyleBold)
	default:
		return core.NewText(fmt.Sprint(v), core.ColorReset, core.StyleBold)
	}
}

// NewTable creates a table with the provided headers. Headers can be strings
// or core.Renderable values; simple strings are converted automatically.
func NewTable(headers ...any) *Table {
	hdrs := make([]core.Renderable, len(headers))
	for i, h := range headers {
		hdrs[i] = anyToRenderable(h)
	}
	t := &Table{
		Headers: hdrs,
		Padding: 1,
		Border:  DefaultBorder,
	}
	// initialize aligns to default left alignment for each header
	t.Aligns = make([]core.Alignment, len(hdrs))
	return t
}

type BorderStyle struct {
	TopLeft, TopMid, TopRight          string
	MidLeft, MidMid, MidRight          string
	BottomLeft, BottomMid, BottomRight string
	Vertical                           string
	Horizontal                         string
}

// DefaultBorder is the default Unicode box-drawing border style.
var DefaultBorder = BorderStyle{
	TopLeft:     "┌",
	TopMid:      "┬",
	TopRight:    "┐",
	MidLeft:     "├",
	MidMid:      "┼",
	MidRight:    "┤",
	BottomLeft:  "└",
	BottomMid:   "┴",
	BottomRight: "┘",
	Vertical:    "│",
	Horizontal:  "─",
}

// WithBorder sets a custom border style and returns the table for chaining.
func (t *Table) WithBorder(b BorderStyle) *Table {
	t.Border = b
	return t
}

// AddRow appends a new row (cells may be nil to indicate empty cells).
func (t *Table) AddRow(cells []core.Renderable) *Table {
	if len(t.Headers) == 0 {
		return t
	}
	row := make([]core.Renderable, len(t.Headers))
	for i := 0; i < len(t.Headers); i++ {
		if i < len(cells) {
			row[i] = cells[i]
		} else {
			row[i] = nil
		}
	}
	t.Rows = append(t.Rows, row)
	return t
}

// AddRowStrict appends a row only if the number of cells matches headers.
func (t *Table) AddRowStrict(cells []core.Renderable) *Table {
	if len(t.Headers) == 0 {
		return t
	}
	if len(cells) != len(t.Headers) {
		return t
	}
	row := make([]core.Renderable, len(t.Headers))
	copy(row, cells)
	t.Rows = append(t.Rows, row)
	return t
}

// AddRowVar appends a row using variadic values. Non-renderable values are
// converted using anyToRenderable.
func (t *Table) AddRowVar(vals ...any) *Table {
	if len(t.Headers) == 0 {
		return t
	}
	row := make([]core.Renderable, len(t.Headers))
	for i := 0; i < len(t.Headers); i++ {
		if i < len(vals) {
			if vals[i] == nil {
				row[i] = nil
			} else {
				row[i] = anyToRenderable(vals[i])
			}
		} else {
			row[i] = nil
		}
	}
	t.Rows = append(t.Rows, row)
	return t
}

func (t *Table) AddRows(rows [][]core.Renderable) *Table {
	for _, r := range rows {
		t.AddRow(r)
	}
	return t
}

func (t *Table) Render() string {
	var b strings.Builder
	_ = RenderTo(&b, t)
	return b.String()
}

// renderSeparator writes a horizontal separator line using provided corner/mid chars
func renderSeparator(sb *strings.Builder, left, mid, right, horiz string, widths []int) {
	sb.WriteString(left)
	for i, w := range widths {
		sb.WriteString(strings.Repeat(horiz, w))
		if i < len(widths)-1 {
			sb.WriteString(mid)
		}
	}
	sb.WriteString(right)
	sb.WriteString("\n")
}

func renderMultiLineRow(sb *strings.Builder, cols [][]string, colWidths []int, aligns []core.Alignment, padding int, b BorderStyle) {
	// compute max height for this row
	maxH := 0
	for _, c := range cols {
		if len(c) > maxH {
			maxH = len(c)
		}
	}
	if maxH == 0 {
		maxH = 1
	}

	// for each line index, render the row line
	for ln := 0; ln < maxH; ln++ {
		sb.WriteString(b.Vertical)
		for ci, c := range cols {
			line := ""
			if ln < len(c) {
				line = c[ln]
			}
			align := core.AlignLeft
			if ci < len(aligns) {
				align = aligns[ci]
			}
			sb.WriteString(core.Pad(line, colWidths[ci], align, padding))
			sb.WriteString(b.Vertical)
		}
		sb.WriteString("\n")
	}
}
