package core

import "testing"

func TestStripANSI(t *testing.T) {
	red := "\x1b[31m"
	bold := "\x1b[1m"
	reset := "\x1b[0m"
	in := red + bold + "Hello" + reset + " World"
	got := StripANSI(in)
	want := "Hello World"
	if got != want {
		t.Fatalf("StripANSI: got %q, want %q", got, want)
	}
}

func TestExtractLeadingPrefix(t *testing.T) {
	red := "\x1b[31m"
	bold := "\x1b[1m"
	reset := "\x1b[0m"

	cases := []struct {
		name string
		in   string
		want string
	}{
		{"no_ansi", "Hello", ""},
		{"ansi_at_start", red + bold + "Hello" + reset, red + bold},
		{"ansi_not_at_start", "Hello " + red + "world" + reset, ""},
		{"multiple_contiguous", red + bold + reset + "X", red + bold + reset},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ExtractLeadingPrefix(c.in)
			if got != c.want {
				t.Fatalf("ExtractLeadingPrefix(%q): got %q, want %q", c.in, got, c.want)
			}
		})
	}
}
