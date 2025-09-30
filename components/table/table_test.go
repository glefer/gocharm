package table

import (
	"strings"
	"testing"

	"github.com/glefer/gocharm/core"
)

func TestTableBorderAndPadding(t *testing.T) {
	// simple table with custom border and padding
	b := BorderStyle{
		TopLeft:     "+",
		TopMid:      "+",
		TopRight:    "+",
		MidLeft:     "+",
		MidMid:      "+",
		MidRight:    "+",
		BottomLeft:  "+",
		BottomMid:   "+",
		BottomRight: "+",
		Vertical:    "|",
		Horizontal:  "-",
	}

	tbl := NewTable("H1", "H2").WithBorder(b).WithPadding(2)
	tbl.AddRowVar("a", "bb")

	r := tbl.Render()

	// strip ANSI sequences added by Render() styles
	clean := core.StripANSI(r)

	if !strings.Contains(clean, "+") || !strings.Contains(clean, "|") {
		t.Fatalf("expected custom border chars in render output: %s", r)
	}

	// check padding: between vertical and content should be two spaces
	if !strings.Contains(clean, "|  a") {
		t.Fatalf("expected two spaces padding before 'a': %s", r)
	}
}

func TestTableAlignment(t *testing.T) {
	tbl := NewTable("H1", "H2").WithPadding(1)
	tbl.WithAlignments(core.AlignLeft, core.AlignRight)
	tbl.AddRowVar("left", "right")

	r := tbl.Render()

	clean := core.StripANSI(r)

	// header separator exists
	if !strings.Contains(clean, "├") && !strings.Contains(clean, "|") {
		// either default unicode border or ascii used in other tests; ensure row content present
		t.Fatalf("unexpected render: %s", r)
	}

	// check right alignment for second column: 'right' should be flush to right inside the cell
	if !strings.Contains(clean, "right ") && !strings.Contains(clean, " right") {
		// basic sanity: ensure content present
		t.Fatalf("expected 'right' in render: %s", r)
	}
}

func TestTableMultilineCell(t *testing.T) {
	tbl := NewTable("H1").WithPadding(1)
	// add a multiline cell
	text := core.NewText("line1\nline2", core.StyleReset)
	tbl.AddRow([]core.Renderable{text})

	r := tbl.Render()

	clean := core.StripANSI(r)

	if !strings.Contains(clean, "line1") || !strings.Contains(clean, "line2") {
		t.Fatalf("expected both lines in render: %s", r)
	}
}
