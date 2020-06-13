package main

import (
	"fmt"
	"runtime"
)

var log Logger

type Logger struct {
}

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

func (l Logger) PrintColor(c Color, values ...interface{}) {
	fmt.Printf("\x1b[%dm", c)
	for _, v := range values {
		fmt.Print(v, "")
	}
	fmt.Print("\x1b[0m")
}

func (l Logger) PrintlnColor(c Color, values ...interface{}) {
	fmt.Printf("\x1b[%dm", c)
	for _, v := range values {
		fmt.Print(v, "")
	}
	fmt.Println("\x1b[0m")
}

func (l Logger) Info(v ...interface{}) {
	l.PrintColor(FgHiGreen, "INFO  | ")
	fmt.Println(v...)
}

func (l Logger) Infof(format string, v ...interface{}) {
	l.PrintColor(FgHiGreen, "INFO  | ")
	fmt.Printf(format + "\n", v...)
}

func (l Logger) Error(err error) {
	_, file, no, ok := runtime.Caller(1)
	if ok {
		l.PrintColor(FgRed, "ERROR | ")
		fmt.Println("in", file, "at line", no)
		fmt.Println("    └› ", err)
	}
}