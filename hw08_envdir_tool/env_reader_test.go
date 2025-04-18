package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		env, err := ReadDir("testdata/env")
		require.NoError(t, err)

		val, ok := env["FOO"]

		require.True(t, ok)
		require.Len(t, env, 5)
		require.Equal(t, val.Value, `   foo
with new line`)
	})

	t.Run("error", func(t *testing.T) {
		env, err := ReadDir("testdata/envv")
		require.Error(t, err)

		val, ok := env["FOO"]

		require.False(t, ok)
		require.Len(t, env, 0)
		require.Equal(t, val.Value, "")
	})
}
