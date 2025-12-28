package colors

import (
	"os"

	"golang.org/x/term"
)

var (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[38;5;51m"
)

func Red(val string) string {
	return red + val + reset
}

func Green(val string) string {
	return green + val + reset
}

func Yellow(val string) string {
	return yellow + val + reset
}

func Cyan(val string) string {
	return cyan + val + reset
}

func init() {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		reset = ""
		red = ""
		green = ""
		yellow = ""
		cyan = ""
	}
}
