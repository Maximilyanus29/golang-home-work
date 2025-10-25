package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		env := Environment{}
		os.Setenv("FOO", "123")

		env["FOO"] = EnvValue{
			Value:      "FOO",
			NeedRemove: true,
		}

		returnCode := RunCmd([]string{"echo", "$FOO"}, Environment{})
		require.Equal(t, returnCode, 0)
	})

	t.Run("negative", func(t *testing.T) {
		env := Environment{}
		os.Setenv("FOO", "123")

		env["FOO"] = EnvValue{
			Value:      "FOO",
			NeedRemove: true,
		}

		returnCode := RunCmd([]string{"ech", "$FOO"}, Environment{})
		require.NotEqual(t, returnCode, 0)
	})
}
