package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("log_level debug", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		logger := New("debug", buffer)
		logger.Debug("hello from debug")
		require.Equal(t, buffer.String(), "DEBUG: hello from debug\n")

		buffer.Reset()
		logger.Info("hello from info")
		require.Equal(t, buffer.String(), "INFO: hello from info\n")

		buffer.Reset()
		logger.Warn("hello from warn")
		require.Equal(t, buffer.String(), "WARN: hello from warn\n")

		buffer.Reset()
		logger.Error("hello from error")
		require.Equal(t, buffer.String(), "ERROR: hello from error\n")
	})

	t.Run("log_level info", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		logger := New("debug", buffer)
		logger.Debug("hello from debug")
		require.Equal(t, buffer.String(), "DEBUG: hello from debug\n")

		buffer.Reset()
		logger.Info("hello from info")
		require.Equal(t, buffer.String(), "INFO: hello from info\n")

		buffer.Reset()
		logger.Warn("hello from warn")
		require.Equal(t, buffer.String(), "WARN: hello from warn\n")

		buffer.Reset()
		logger.Error("hello from error")
		require.Equal(t, buffer.String(), "ERROR: hello from error\n")
	})

	t.Run("log_level warn", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		logger := New("warn", buffer)
		logger.Debug("hello from debug")
		require.Equal(t, buffer.String(), "")

		buffer.Reset()
		logger.Info("hello from info")
		require.Equal(t, buffer.String(), "")

		buffer.Reset()
		logger.Warn("hello from warn")
		require.Equal(t, buffer.String(), "WARN: hello from warn\n")

		buffer.Reset()
		logger.Error("hello from error")
		require.Equal(t, buffer.String(), "ERROR: hello from error\n")
	})

	t.Run("log_level error", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		logger := New("error", buffer)
		logger.Debug("hello from error")
		require.Equal(t, buffer.String(), "")

		buffer.Reset()
		logger.Info("hello from info")
		require.Equal(t, buffer.String(), "")

		buffer.Reset()
		logger.Warn("hello from warn")
		require.Equal(t, buffer.String(), "")

		buffer.Reset()
		logger.Error("hello from error")
		require.Equal(t, buffer.String(), "ERROR: hello from error\n")
	})
}
