package core

import (
	"fmt"
	"io"
	"os"
)

// Mode represents the output mode used by a Console.
type Mode int

// Console manages rendering of values to an io.Writer.
// It supports plain and markup modes.
type Console struct {
	out  io.Writer
	mode Mode
}

const (
	ModePlain Mode = iota
	ModeMarkup
)

// NewConsole creates a new Console. By default it writes to os.Stdout
// and uses ModeMarkup. Options may include an io.Writer and/or a Mode.
// ConsoleOption configures a Console.
type ConsoleOption func(*Console)

// WithWriter sets the writer used by the Console.
func WithWriter(w io.Writer) ConsoleOption {
	return func(c *Console) { c.out = w }
}

// WithMode sets the output mode of the Console.
func WithMode(m Mode) ConsoleOption {
	return func(c *Console) { c.mode = m }
}

// NewConsole creates a new Console. It accepts either functional options
// (ConsoleOption) or the legacy variadic signature (io.Writer and/or Mode).
// If both styles are mixed, functional options are applied first, then
// legacy options are applied.
func NewConsole(opts ...interface{}) *Console {
	c := &Console{out: os.Stdout, mode: ModeMarkup}

	// first pass: functional options
	for _, o := range opts {
		if fn, ok := o.(ConsoleOption); ok {
			fn(c)
		}
	}

	// second pass: legacy-style options
	for _, opt := range opts {
		switch v := opt.(type) {
		case io.Writer:
			c.out = v
		case Mode:
			c.mode = v
		}
	}
	return c
}

// SetMode changes the Console output mode.
func (c *Console) SetMode(mode Mode) {
	c.mode = mode
}

// Println renders v and writes it followed by a newline.
func (c *Console) Println(v interface{}) {
	c.printInternal(v, true)
}

// Print renders v and writes it without adding a newline.
func (c *Console) Print(v interface{}) {
	c.printInternal(v, false)
}

// printInternal performs rendering and writes to the underlying writer.
func (c *Console) printInternal(v interface{}, newline bool) {
	var r Renderable
	switch val := v.(type) {
	case string:
		if c.mode == ModeMarkup {
			r = NewMarkup(val)
		} else {
			r = NewText(val, ColorReset)
		}
	case Renderable:
		r = val
	default:
		r = NewText(fmt.Sprint(val), ColorReset)
	}
	output := r.Render()
	if newline {
		if _, err := fmt.Fprintln(c.out, output); err != nil {
			_ = err
		}
	} else {
		if _, err := fmt.Fprint(c.out, output); err != nil {
			_ = err
		}
	}
}

func (c *Console) Render(renderable Renderable) {
	c.Println(renderable)
}
