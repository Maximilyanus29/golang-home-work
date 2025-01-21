package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	lastIndex := getLastIndexOnInput(input)
	result := strings.Builder{}
	var oldRune rune

	escapeMode := false

	for index, currentRune := range input {
		isDigitCurrentRune := isDigitRune(currentRune)
		isLastIndex := index == lastIndex

		if index == 0 {
			if isDigitCurrentRune {
				return "", ErrInvalidString
			}
			oldRune = currentRune
			continue
		}

		isDigitOldRune := isDigitRune(oldRune)

		if oldRune == '\\' && isDigitCurrentRune {
			escapeMode = true
			oldRune = currentRune
			continue
		}

		if isDigitCurrentRune && isDigitOldRune && !escapeMode {
			return "", ErrInvalidString
		}

		if !escapeMode {
			if isDigitRune(oldRune) {
				if isLastIndex {
					result.WriteRune(currentRune)
				}
				oldRune = currentRune
				continue
			}
		}

		if isDigitCurrentRune {
			if count, err := strconv.Atoi(string(currentRune)); err == nil {
				result.WriteString(strings.Repeat(string(oldRune), count))
			} else {
				return "", ErrInvalidString
			}
		} else {
			result.WriteRune(oldRune)

			if isLastIndex {
				result.WriteRune(currentRune)
			}
		}
		oldRune = currentRune
	}
	return result.String(), nil
}

func isDigitRune(v rune) bool {
	return v >= '0' && v <= '9'
}

func getLastIndexOnInput(inputString string) int {
	result := 0
	for index := range inputString {
		result = index
	}
	return result
}
