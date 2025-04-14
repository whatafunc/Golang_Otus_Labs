package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inputStr string) (string, error) {
	inputStr = strings.ReplaceAll(inputStr, "🙃", "")
	inputStrLen := len(inputStr)
	if inputStrLen > 0 {
		savedValue := ""
		var b strings.Builder
		for i, w := 0, 0; i < inputStrLen; i += w {
			runeValue, width := utf8.DecodeRuneInString(inputStr[i:])

			stringVal := ""
			if width > 1 {
				width = 1
			} else {
				stringVal = string(runeValue)
			}
			fmt.Println(" checkCur = ", stringVal)
			intValue, _ := strconv.Atoi(stringVal)
			if intValue > 0 {
				if i == 0 {
					return "", ErrInvalidString
				}
				fmt.Fprintf(&b, "%s", strings.Repeat(savedValue, intValue-1))
				savedValue = "digit"
			} else {
				doWriteFlag := true
				if inputStrLen > i+1 {
					runeValueNext, _ := utf8.DecodeRuneInString(inputStr[i+1:])
					stringValNext := string(runeValueNext)

					if stringValNext == "0" {
						doWriteFlag = false
					}
				}
				if stringVal == "0" {
					doWriteFlag = false
				}
				if doWriteFlag {
					fmt.Fprintf(&b, "%s", stringVal)
					savedValue = stringVal
				}
			}
			w = width
		}
		return b.String(), nil
	} else {
		return "", nil
	}
}
