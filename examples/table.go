package main

import (
	"fmt"

	"github.com/glefer/gocharm/components/table"
	"github.com/glefer/gocharm/core"
)

func main() {
	// Example 1: Simple table
	fmt.Println("=== Simple Table ===")
	t := table.NewTable("Character", "Level", "Game")
	t.AddRowVar("Mario", 12, "Super Mario Bros")
	t.AddRowVar("Link", 25, "Zelda: Ocarina of Time")
	t.AddRowVar("Master Chief", 42, "Halo")
	fmt.Println(t.Render())

	// Example 2: Table with custom padding and border
	fmt.Println("=== Table with Custom Padding and Border ===")
	t2 := table.NewTable("Item", "Price", "Stock").WithPadding(2)
	t2.AddRowVar("Potion", "50G", 12)
	t2.AddRowVar("Phoenix Down", "500G", 3)
	fmt.Println(t2.Render())

	// Example 3: Table with ANSI color in cell
	fmt.Println("=== Table with ANSI Color in Cell ===")
	t3 := table.NewTable("Message", "Status")
	t3.AddRowVar(
		"Boss defeated", core.Colorize("OK", core.ColorGreen, nil),
	)
	t3.AddRowVar(
		"Game Over", core.Colorize("FAIL", core.ColorRed, []core.Style{core.StyleBold}),
	)
	fmt.Println(t3.Render())

	// Example 4: Table with styled headers
	fmt.Println("=== Table with Styled Headers ===")
	t4 := table.NewTable(
		core.Colorize("ID", core.ColorCyan, []core.Style{core.StyleBold}),
		core.Colorize("Player", core.ColorYellow, []core.Style{core.StyleUnderline}),
		core.Colorize("Class", core.ColorMagenta, nil),
	)
	t4.AddRowVar(1, "Link", core.Colorize("Hero", core.ColorRed, []core.Style{core.StyleBold}))
	t4.AddRowVar(2, core.Colorize("Kirby", core.ColorGreen, nil), "Copycat")
	fmt.Println(t4.Render())

	// Example 5: Table with mixed styles in cells
	fmt.Println("=== Table with Mixed Styles in Cells ===")
	t5 := table.NewTable("Quest", "Status")
	t5.AddRowVar(
		core.Colorize("Defeat Sephiroth", core.ColorBlue, []core.Style{core.StyleItalic}),
		core.Colorize("DONE", core.ColorGreen, []core.Style{core.StyleBold}),
	)
	t5.AddRowVar(
		core.Colorize("Collect 100 rings", core.ColorYellow, nil),
		core.Colorize("PENDING", core.ColorRed, []core.Style{core.StyleUnderline}),
	)
	fmt.Println(t5.Render())
}
