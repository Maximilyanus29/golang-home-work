package logger

import (
	"log"
	"slices"
	"strings"
)

var availableLogLevels = []string{"debug", "info", "warn", "error"}

type Logger struct {
	level string
}

func New(level string) *Logger {
	level = strings.ToLower(level)
	if !slices.Contains(availableLogLevels, level) {
		log.Printf("unknown log level: %s, using level 'info'", level)
		return &Logger{
			level: "info",
		}
	}
	return &Logger{
		level: level,
	}
}

func (l Logger) Info(msg string) {
	log.Printf("INFO: %s", msg)
}

func (l Logger) Error(msg string) {
	log.Printf("ERROR: %s", msg)
}

func (l Logger) Warn(msg string) {
	log.Printf("WARN: %s", msg)
}

func (l Logger) Debug(msg string) {
	log.Printf("DEBUG: %s", msg)
}
