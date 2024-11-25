package config

import (
	"github.com/fatih/color"
	"io"
)

var disableLogs = false

func Blue(format string, a ...interface{}) {
	if disableLogs {
		return
	}
	color.Blue(format, a...)
}

func Cyan(format string, a ...interface{}) {
	if disableLogs {
		return
	}
	color.Cyan(format, a...)
}

func Red(format string, a ...interface{}) {
	if disableLogs {
		return
	}
	color.Red(format, a...)
}

func Disable() {
	disableLogs = true
}

func NewOutput(writer io.Writer) {
	color.Output = writer
}
