package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for key, EnvValue := range env {
		if EnvValue.NeedRemove {
			os.Unsetenv(key)
		}

		if EnvValue.Value != "" {
			os.Setenv(key, EnvValue.Value)
		}
	}

	args := cmd[1:]
	mainCommand := cmd[0]

	for _, arg := range args {
		if !isSafeArgument(arg) {
			fmt.Printf("Небезопасный аргумент: %s\n", arg)
			return
		}
	}

	command := exec.Command(mainCommand, args...)
	command.Env = os.Environ()
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	if err != nil {
		fmt.Print(err)
	}

	return command.ProcessState.ExitCode()
}

func isSafeArgument(arg string) bool {
	// Пример проверки: запрещаем использование символов, которые могут быть опасны
	invalidChars := ";|&<>*"
	for _, char := range invalidChars {
		if strings.ContainsRune(arg, char) {
			return false
		}
	}
	return true
}
