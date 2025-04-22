package hw03frequencyanalysis

import (
	// "fmt".
	"sort"
	"strings"
)

type WordAndFreq struct {
	Word string
	Freq int
}

// Add new element as word and 1 repeat or
// if it exists, update its qty of repeats aka frequency.
func updateFreq(wordFreqs []WordAndFreq, targetWord string) []WordAndFreq {
	if len(wordFreqs) == 0 {
		wordFreqs = append(wordFreqs, WordAndFreq{targetWord, 1})
	} else {
		foundFlag := false
		for i := range wordFreqs {
			if wordFreqs[i].Word == targetWord {
				wordFreqs[i].Freq++
				foundFlag = true
				break
			}
		}
		if !foundFlag {
			wordFreqs = append(wordFreqs, WordAndFreq{targetWord, 1})
		}
	}

	return wordFreqs
}

func Top10(input string) []string {
	if len(input) == 0 {
		return nil
	}
	if !strings.Contains(input, " ") {
		return nil
	}

	words := strings.Fields(input)
	wordFreqs := []WordAndFreq{}
	for _, word := range words {
		if word == "" {
			continue
		}
		wordFreqs = updateFreq(wordFreqs, word)
	}

	sort.Slice(wordFreqs, func(i, j int) bool {
		if wordFreqs[i].Freq == wordFreqs[j].Freq {
			return wordFreqs[i].Word < wordFreqs[j].Word // Alphabetical if freqs are equal
		}
		return wordFreqs[i].Freq > wordFreqs[j].Freq
	})
	result := []string{} // Creating result
	for i := range wordFreqs {
		result = append(result, wordFreqs[i].Word)
		if i == 9 {
			break
		}
	}

	return result
}
