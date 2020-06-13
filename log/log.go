package log

import (
	"fmt"
	"runtime"
)

type Color int

const Reset Color = 0

const (
	FgBlack Color = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Reset iota
const (
	FgHiBlack Color = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

const EscapeCharacter = "\x1b"

func PrintColor(c Color, values ...interface{}) {
	fmt.Printf("\x1b[%dm", c)
	for _, v := range values {
		fmt.Print(v, "")
	}
	fmt.Print("\x1b[0m")
}

func PrintlnColor(c Color, values ...interface{}) {
	fmt.Printf("\x1b[%dm", c)
	for _, v := range values {
		fmt.Print(v, "")
	}
	fmt.Println("\x1b[0m")
}

func Info(v ...interface{}) {
	PrintColor(FgHiGreen, "INFO  | ")
	fmt.Println(v...)
}

func Infof(format string, v ...interface{}) {
	PrintColor(FgHiGreen, "INFO  | ")
	fmt.Printf(format + "\n", v...)
}

func Error(err error) {
	_, file, no, ok := runtime.Caller(1)
	if ok {
		PrintColor(FgRed, "ERROR | ")
		fmt.Println("in", file, "at line", no)
		fmt.Println("    └› ", err)
	}
}