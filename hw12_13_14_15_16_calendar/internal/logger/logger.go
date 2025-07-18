package logger

import (
	"fmt"
	"os"
	"strings"
)

// Logger supports 'info' and 'error' levels

type Logger struct {
	level string
}

func New(level string) *Logger {
	return &Logger{level: strings.ToLower(level)}
}

func (l Logger) Info(msg string) {
	if l.level == "info" {
		fmt.Println(msg)
	}
}

func (l Logger) Error(msg string) {
	if l.level == "info" || l.level == "error" {
		fmt.Fprintln(os.Stderr, msg)
	}
}

// TODO
