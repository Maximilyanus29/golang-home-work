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
	input = `qwe\\5`
	// var escapedRune rune

	var result string
	runes := []rune(input)
	length := len(runes)

	if isDigitRune(runes[0]) {
		return "", ErrInvalidString
	}

	for i := 0; i < length; i++ {
		r := runes[i]
		count := 1

		if i+1 < len(runes) {
			r2 := runes[i+1]

			if isDigitRune(r2) {
				if isDigitRune(r) {
					return "", ErrInvalidString
				}

				if r == '\\' {
					if i+2 < len(runes) {
						r3 := runes[i+2]

						if isDigitRune(r3) {
							count, _ = strconv.Atoi(string(r3))
							result += strings.Repeat(string(r2), count)
							i += 2
						}
					} else {
						result += string(r2)
						i++
					}

					continue
				}

				count, _ = strconv.Atoi(string(r2))
				result += strings.Repeat(string(r), count)
				i++
				continue
			}
		}

		result += string(r)
	}

	return result, nil
}

func isDigitRune(v rune) bool {
	if v >= '0' && v <= '9' {
		return true
	}
	return false
}
