package logger

import (
	"fmt"
	"io"
	"log"
	"slices"
	"strings"
)

var availableLogLevels = []string{"debug", "info", "warn", "error"}

type Logger struct {
	level string
	out   io.Writer
}

func New(level string, out io.Writer) *Logger {
	level = strings.ToLower(level)
	if !slices.Contains(availableLogLevels, level) {
		log.Fatalf("unknown log level: %s, using level 'info'", level)
	}
	return &Logger{
		level: level,
		out:   out,
	}
}

func (l Logger) Debug(msg string) {
	if slices.Contains([]string{"debug"}, l.level) {
		l.out.Write([]byte(fmt.Sprintf("DEBUG: %s\n", msg)))
	}
}

func (l Logger) Info(msg string) {
	if slices.Contains([]string{"debug", "info"}, l.level) {
		l.out.Write([]byte(fmt.Sprintf("INFO: %s\n", msg)))
	}
}

func (l Logger) Warn(msg string) {
	if slices.Contains([]string{"debug", "info", "warn"}, l.level) {
		l.out.Write([]byte(fmt.Sprintf("WARN: %s\n", msg)))
	}
}

func (l Logger) Error(msg string) {
	if slices.Contains([]string{"debug", "info", "warn", "error"}, l.level) {
		l.out.Write([]byte(fmt.Sprintf("ERROR: %s\n", msg)))
	}
}
