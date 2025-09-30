package table

import (
	"fmt"
	"strings"

	"github.com/glefer/gocharm/core"
)

func ExampleTable_Render() {
	t := NewTable("Name", "Value").WithPadding(1)
	t.AddRowVar("foo", "bar")
	out := t.Render()
	// print without ANSI for example output
	fmt.Println(strings.TrimSpace(core.StripANSI(out)))
	//
}
