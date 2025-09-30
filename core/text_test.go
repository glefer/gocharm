package core

import (
	"strings"
	"testing"
)

func TestText_Render_plain(t *testing.T) {
	text := NewText("hello", StyleReset)
	got := text.Render()
	want := "hello"
	if got != want {
		t.Errorf("plain Render: got %q, want %q", got, want)
	}
}

func TestText_Render_with_styles(t *testing.T) {
	text := NewText("world", StyleReset, StyleBold, StyleUnderline)
	got := text.Render()
	if got == "world" {
		t.Errorf("Render with styles should not be plain: got %q", got)
	}
}

func TestText_Render_markup(t *testing.T) {
	text := NewMarkup("[bold]foo[/bold]")
	got := text.Render()
	if !strings.Contains(got, "foo") {
		t.Errorf("Render markup: expected output to contain 'foo', got %q", got)
	}
}

func TestParseMarkup_multiple_tags(t *testing.T) {
	input := "[bold]foo[/bold] [underline]bar[/underline]"
	out := ParseMarkup(input)
	if !strings.Contains(out, "foo") || !strings.Contains(out, "bar") {
		t.Errorf("ParseMarkup: missing expected content, got %q", out)
	}
}

func TestParseMarkup_no_tags(t *testing.T) {
	input := "baz"
	out := ParseMarkup(input)
	if out != "baz" {
		t.Errorf("ParseMarkup: expected %q, got %q", input, out)
	}
}

func TestText_Render_empty(t *testing.T) {
	text := NewText("", StyleReset)
	got := text.Render()
	if got != "" {
		t.Errorf("Render empty: got %q, want empty string", got)
	}
}
