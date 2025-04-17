package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inputStr string) (string, error) {
	if inputIsInvalid(inputStr) {
		return "", ErrInvalidString
	}

	if len(inputStr) == 0 {
		return "", nil
	}

	var b strings.Builder
	var savedValue string

	for i, w := 0, 0; i < len(inputStr); i += w {
		runeValue, width := utf8.DecodeRuneInString(inputStr[i:])
		stringVal := string(runeValue)
		intValue, _ := strconv.Atoi(stringVal)

		if intValue <= 0 {
			if shouldWriteChar(inputStr, i, width, stringVal) {
				b.WriteString(stringVal)
				if !isNextCharZero(inputStr, i, width) {
					savedValue = stringVal
				}
			}
		} else {
			b.WriteString(strings.Repeat(savedValue, intValue-1))
			savedValue = "digit"
		}
		w = width
	}
	return b.String(), nil
}

func shouldWriteChar(inputStr string, i, width int, currentChar string) bool {
	if currentChar == "0" {
		return false
	}

	if i+width >= len(inputStr) {
		return true
	}

	nextRune, _ := utf8.DecodeRuneInString(inputStr[i+width:])
	return nextRune != '0'
}

func isNextCharZero(inputStr string, i, width int) bool {
	if i+width >= len(inputStr) {
		return false
	}
	nextRune, _ := utf8.DecodeRuneInString(inputStr[i+width:])
	return nextRune == '0'
}

func inputIsInvalid(s string) bool {
	switch {
	case hasZeroAfterNumber(s):
		return true
	case len(s) > 0:
		_, err := strconv.Atoi(string(s[0]))
		return err == nil
	default:
		return false
	}
}

func hasZeroAfterNumber(s string) bool {
	var prevIsDigit bool
	for _, r := range s {
		switch {
		case unicode.IsDigit(r) && prevIsDigit:
			return true
		case unicode.IsDigit(r):
			prevIsDigit = true
		default:
			prevIsDigit = false
		}
	}
	return false
}
