package hw03frequencyanalysis

import (
	"slices"
	"sort"
	"strings"
)

type Pair struct {
	Key   string
	Value int
}

func Top10(input string) []string {
	result := []string{}
	if input == "" {
		return result
	}

	escapeSymbols := []rune{'\t', '\n', '\v', '\f', '\r', ' ', 0}

	stringBuilder := &strings.Builder{}

	wordsCount := make(map[string]int)

	for _, v := range input {
		if slices.Contains(escapeSymbols, v) {
			word := stringBuilder.String()
			if word == "-" {
				stringBuilder.Reset()
				continue
			}

			lowerString := strings.ToLower(word)
			lowerTrimString := strings.Trim(lowerString, "!,'\"\\.")

			wordCount, ok := wordsCount[lowerTrimString]

			if ok {
				wordsCount[lowerTrimString] = wordCount + 1
			} else {
				wordsCount[lowerTrimString] = 1
			}

			stringBuilder.Reset()
			continue
		}

		stringBuilder.WriteRune(v)
	}

	sortedMap := sortMap(wordsCount)

	keysFromSlicePairs := getKeysFromSlicePairs(sortedMap)

	return keysFromSlicePairs[1:11]
}

func sortMap(inputMap map[string]int) []Pair {
	pairs := make([]Pair, 0, len(inputMap))
	for key, value := range inputMap {
		pairs = append(pairs, Pair{Key: key, Value: value})
	}

	// Сортировка слайса по значениям
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Value == pairs[j].Value {
			return pairs[i].Key < pairs[j].Key
		}
		return pairs[i].Value > pairs[j].Value
	})

	return pairs
}

func getKeysFromSlicePairs(sortedMap []Pair) []string {
	result := []string{}
	for _, item := range sortedMap {
		result = append(result, item.Key)
	}
	return result
}
