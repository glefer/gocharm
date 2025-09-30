package core

import (
	"bytes"
	"testing"
)

type mockRenderable struct{ out string }

func (m mockRenderable) Render() string { return m.out }

func assertConsoleOutput(t *testing.T, mode Mode, input interface{}, want string, newline bool) {
	t.Helper()
	buf := &bytes.Buffer{}
	console := NewConsole(buf, mode)
	if newline {
		console.Println(input)
	} else {
		console.Print(input)
	}
	if got := buf.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestConsole_PrintAndPrintln(t *testing.T) {
	cases := []struct {
		name    string
		mode    Mode
		input   interface{}
		want    string
		newline bool
	}{
		{"Println plain string", ModePlain, "hello", NewText("hello", StyleReset).Render() + "\n", true},
		{"Println markup string", ModeMarkup, "hello", NewMarkup("hello").Render() + "\n", true},
		{"Print plain string", ModePlain, "world", NewText("world", StyleReset).Render(), false},
		{"Print markup string", ModeMarkup, "world", NewMarkup("world").Render(), false},
		{"Println Renderable", ModePlain, mockRenderable{"foo"}, "foo\n", true},
		{"Print int", ModePlain, 42, NewText("42", StyleReset).Render(), false},
		{"Println empty string", ModePlain, "", NewText("", StyleReset).Render() + "\n", true},
		{"Println nil", ModePlain, nil, NewText("<nil>", StyleReset).Render() + "\n", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assertConsoleOutput(t, c.mode, c.input, c.want, c.newline)
		})
	}
}

func TestConsole_SetMode(t *testing.T) {
	buf := &bytes.Buffer{}
	console := NewConsole(buf, ModePlain)
	console.SetMode(ModeMarkup)
	console.Println("test")
	want := NewMarkup("test").Render() + "\n"
	if got := buf.String(); got != want {
		t.Errorf("SetMode failed: got %q, want %q", got, want)
	}
}

func TestConsole_Render(t *testing.T) {
	buf := &bytes.Buffer{}
	console := NewConsole(buf, ModePlain)
	r := mockRenderable{"rendered output"}
	console.Render(r)
	want := "rendered output\n"
	if got := buf.String(); got != want {
		t.Errorf("Render failed: got %q, want %q", got, want)
	}
}

func TestConsole_ModeSwitching(t *testing.T) {
	buf := &bytes.Buffer{}
	console := NewConsole(buf, ModePlain)
	console.Println("plain")
	console.SetMode(ModeMarkup)
	console.Println("markup")
	want := NewText("plain", StyleReset).Render() + "\n" + NewMarkup("markup").Render() + "\n"
	if got := buf.String(); got != want {
		t.Errorf("ModeSwitching failed: got %q, want %q", got, want)
	}
}
