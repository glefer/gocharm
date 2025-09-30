package core

type Color string
type Style string

const (
	ColorRed       Color = "red"
	ColorGreen     Color = "green"
	ColorYellow    Color = "yellow"
	ColorBlue      Color = "blue"
	ColorBlack     Color = "black"
	ColorMagenta   Color = "magenta"
	ColorCyan      Color = "cyan"
	ColorWhite     Color = "white"
	StyleReset     Color = "reset"
	StyleBold      Style = "bold"
	StyleItalic    Style = "italic"
	StyleUnderline Style = "underline"
)

var ANSIStyles = map[string]string{
	string(StyleReset):     "\033[0m",
	string(StyleBold):      "\033[1m",
	string(ColorRed):       "\033[31m",
	string(ColorGreen):     "\033[32m",
	string(ColorYellow):    "\033[33m",
	string(ColorBlue):      "\033[34m",
	string(ColorBlack):     "\033[30m",
	string(ColorMagenta):   "\033[35m",
	string(ColorCyan):      "\033[36m",
	string(ColorWhite):     "\033[37m",
	string(StyleItalic):    "\033[3m",
	string(StyleUnderline): "\033[4m",
}

func Colorize(text string, color Color, styles ...Style) string {
	output := ""

	if val, ok := ANSIStyles[string(color)]; ok {
		output += val
	}

	for _, style := range styles {
		if val, ok := ANSIStyles[string(style)]; ok {
			output += val
		}
	}
	return output + text + ANSIStyles[string(StyleReset)]
}
