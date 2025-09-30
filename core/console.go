package core

import (
	"fmt"
	"io"
	"os"
)

type Mode int

// Console gère l'affichage en mode texte ou markup sur un writer donné.
type Console struct {
	out  io.Writer
	mode Mode
}

const (
	ModePlain Mode = iota
	ModeMarkup
)

func NewConsole(opts ...interface{}) *Console {
	var out io.Writer = os.Stdout
	mode := ModeMarkup
	for _, opt := range opts {
		switch v := opt.(type) {
		case io.Writer:
			out = v
		case Mode:
			mode = v
		}
	}
	return &Console{out: out, mode: mode}
}

func (c *Console) SetMode(mode Mode) {
	c.mode = mode
}

func (c *Console) Println(v interface{}) {
	c.printInternal(v, true)
}

func (c *Console) Print(v interface{}) {
	c.printInternal(v, false)
}

func (c *Console) printInternal(v interface{}, newline bool) {
	var r Renderable
	switch val := v.(type) {
	case string:
		if c.mode == ModeMarkup {
			r = NewMarkup(val)
		} else {
			r = NewText(val, StyleReset)
		}
	case Renderable:
		r = val
	default:
		r = NewText(fmt.Sprint(val), StyleReset)
	}
	output := r.Render()
	if newline {
		fmt.Fprintln(c.out, output)
	} else {
		fmt.Fprint(c.out, output)
	}
}
