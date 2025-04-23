package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type WordAndFreq struct {
	Word string
	Freq int
}

func Top10(input string) []string {
	if len(input) == 0 {
		return nil
	}
	if !strings.Contains(input, "") {
		return nil
	}

	words := strings.Fields(input) // All words splitted by a 'space, ...,... '.

	wordCounter := map[string]int{} // code rafactoring as addressing by key is faster.
	for _, word := range words {
		if word == "" {
			continue
		}
		if wordCounter[word] > 0 {
			wordCounter[word]++
		} else {
			wordCounter[word] = 1
		}
	}

	wordFreqs := []WordAndFreq{}
	for k, v := range wordCounter {
		wordFreqs = append(wordFreqs, WordAndFreq{k, v})
	}

	sort.Slice(wordFreqs, func(i, j int) bool { // Sort slices is quick.
		if wordFreqs[i].Freq == wordFreqs[j].Freq {
			return wordFreqs[i].Word < wordFreqs[j].Word // Alphabetical if freqs are equal.
		}
		return wordFreqs[i].Freq > wordFreqs[j].Freq
	})
	result := []string{} // Creating result of 10 top.
	for i := range wordFreqs {
		result = append(result, wordFreqs[i].Word)
		if i == 9 {
			break
		}
	}

	return result
}
