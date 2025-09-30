package table

import (
	"io"
	"strings"

	"github.com/glefer/gocharm/core"
)

// RenderTo writes the table representation to w. It returns any write error
// encountered while emitting the rendered table.
func RenderTo(w io.Writer, t *Table) error {
	if len(t.Headers) == 0 {
		_, _ = w.Write([]byte(""))
		return nil
	}

	// Pre-render headers and cells and split into lines
	hdrStrs := make([]string, len(t.Headers))
	for i, header := range t.Headers {
		hdrStrs[i] = header.Render()
	}

	hdrLines := make([][]string, len(t.Headers))
	for i := range hdrStrs {
		hdrLines[i] = core.EnsurePrefixPerLine(hdrStrs[i])
	}

	rowsLines := make([][][]string, len(t.Rows))
	for ri, row := range t.Rows {
		rowsLines[ri] = make([][]string, len(t.Headers))
		for ci := 0; ci < len(t.Headers); ci++ {
			if ci < len(row) && row[ci] != nil {
				rowsLines[ri][ci] = core.EnsurePrefixPerLine(row[ci].Render())
			} else {
				rowsLines[ri][ci] = []string{""}
			}
		}
	}

	// compute max content width per column (visible length, without padding) across all lines
	colWidths := make([]int, len(t.Headers))
	for i := range t.Headers {
		for _, line := range hdrLines[i] {
			if core.VisibleLen(line) > colWidths[i] {
				colWidths[i] = core.VisibleLen(line)
			}
		}
	}
	for _, row := range rowsLines {
		for i, cellLines := range row {
			for _, line := range cellLines {
				if core.VisibleLen(line) > colWidths[i] {
					colWidths[i] = core.VisibleLen(line)
				}
			}
		}
	}

	// compute box widths = content width + 2*padding (used for borders)
	cellBoxWidths := make([]int, len(colWidths))
	for i, w := range colWidths {
		cellBoxWidths[i] = w + 2*t.Padding
	}

	// estimate builder size
	totalWidth := 1 // left border
	for _, w := range cellBoxWidths {
		totalWidth += w
		totalWidth += 1 // vertical separator or right border
	}
	height := 4 + len(t.Rows)
	sb := strings.Builder{}
	sb.Grow(totalWidth * height)
	b := t.Border

	// top border
	renderSeparator(&sb, b.TopLeft, b.TopMid, b.TopRight, b.Horizontal, cellBoxWidths)

	// header block
	renderMultiLineRow(&sb, hdrLines, colWidths, t.Aligns, t.Padding, b)

	// header separator
	renderSeparator(&sb, b.MidLeft, b.MidMid, b.MidRight, b.Horizontal, cellBoxWidths)

	// body rows
	for _, row := range rowsLines {
		renderMultiLineRow(&sb, row, colWidths, t.Aligns, t.Padding, b)
	}

	// bottom border
	renderSeparator(&sb, b.BottomLeft, b.BottomMid, b.BottomRight, b.Horizontal, cellBoxWidths)

	_, err := io.WriteString(w, sb.String())
	return err
}
