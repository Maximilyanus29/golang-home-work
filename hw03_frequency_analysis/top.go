package hw03frequencyanalysis

import (
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

	stringBuilder := &strings.Builder{}

	wordsCount := make(map[string]int)

	words := strings.Fields(input)

	for _, word := range words {
		if word == "-" {
			stringBuilder.Reset()
			continue
		}

		lowerString := strings.ToLower(word)
		lowerTrimString := strings.Trim(lowerString, "!,'\"\\.")

		wordsCount[lowerTrimString]++
	}

	wordsCountSortedDesc := sortMap(wordsCount)

	keysFromSlicePairs := getKeysFromSlicePairs(wordsCountSortedDesc)

	keysFromSlicePairsLen := len(keysFromSlicePairs)

	if keysFromSlicePairsLen > 10 {
		keysFromSlicePairsLen = 10
	}

	return keysFromSlicePairs[0:keysFromSlicePairsLen]
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
