package logger

import (
	"fmt"
	"slices"
	"strings"
)

type Logger struct {
	Level string
}

func New(level string) *Logger {
	return &Logger{
		Level: level,
	}
}

func (l Logger) Debug(msg ...string) {
	if slices.Contains([]string{"DEBUG"}, l.Level) {
		fmt.Println(format(msg, l.Level))

	}

}

func (l Logger) Info(msg ...string) {
	if slices.Contains([]string{"INFO", "DEBUG"}, l.Level) {
		fmt.Println(format(msg, l.Level))

	}
}

func (l Logger) Error(msg ...string) {
	if slices.Contains([]string{"INFO", "DEBUG", "ERROR"}, l.Level) {
		fmt.Println(format(msg, l.Level))

	}
}

func (l Logger) Warning(msg ...string) {
	if slices.Contains([]string{"INFO", "DEBUG", "ERROR", "WARNING"}, l.Level) {
		fmt.Println(format(msg, l.Level))
	}
}

func format(msg []string, level string) string {
	var builder strings.Builder
	builder.WriteString(level)
	builder.WriteRune(':')
	builder.WriteRune(' ')
	for _, v := range msg {
		builder.WriteString(v)
		builder.WriteRune(' ')
	}
	return builder.String()
}
