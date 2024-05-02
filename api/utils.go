package api

import (
	"errors"
	"strings"
	"unicode"
)

// Evaluates: "1" and "true" to true
//            "0" and "false" to false
//            anything else is invalid
//
func stringToBool(s string) (bool, error) {
	lower := strings.ToLower(s)
	if lower == "1" || lower == "true" {
		return true, nil
	} else if lower == "0" || lower == "false" {
		return false, nil
	} else {
		return false, errors.New("No boolean representation of input string")
	}
}

func isPalindrome(str string) bool {
	// Lowercase and remove all non alphanumeric characters from the input
	// string
	//
	strippedStr := strings.Map(func(r rune) rune {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return -1
		}
		return unicode.ToLower(r)
	}, str)

	// Convert string to rune array
	//
	runeArr := []rune(strippedStr)

	// Iterate through the string from both sides, comparing the results until
	// they overlap
	//
	for i, j := 0, len(runeArr)-1; i < j; i, j = i+1, j-1 {
		if runeArr[i] != runeArr[j] {
			return false
		}
	}

	return true
}
