package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)
	dirEntry, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntry {
		fileName := entry.Name()
		if strings.ContainsRune(fileName, '=') {
			continue
		}

		file, err := os.Open(dir + "/" + fileName)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		fi, err := file.Stat()
		if err != nil {
			return nil, err
		}

		if fi.Size() == 0 {
			env[fileName] = EnvValue{
				NeedRemove: true,
			}
			continue
		}

		firstLine, err := getFirstLine(file)
		if err != nil {
			env[fileName] = EnvValue{
				NeedRemove: true,
			}
			continue
		}

		trimValue := strings.TrimRight(firstLine, "\t ")
		totalValue := strings.ReplaceAll(trimValue, "\x00", "\n")

		env[fileName] = EnvValue{
			Value:      totalValue,
			NeedRemove: true,
		}
	}

	// Place your code here
	return env, nil
}

func getFirstLine(file *os.File) (string, error) {
	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		return scanner.Text(), nil
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", errors.New("something wrong")
}
