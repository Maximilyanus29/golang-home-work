package filelogger

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/config"
)

var availableLogLevels = []string{"debug", "info", "warn", "error"}

type Logger struct {
	level string
	out   io.Writer
}

func New(loggerConf config.LoggerConf) *Logger {
	level := strings.ToLower(loggerConf.Level)
	if !slices.Contains(availableLogLevels, level) {
		app.Exit("unknown log_level %s", level)
	}

	logFile, err := os.OpenFile(loggerConf.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		app.Exit("Failed to open log file: %v\n", err)
	}

	return &Logger{
		level: level,
		out:   logFile,
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
