package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Print("Должно быть как минимум 3 агрумента")
		os.Exit(0)
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Print(err)
	}

	command := os.Args[2:]

	exitCode := RunCmd(command, env)

	os.Exit(exitCode)
}
