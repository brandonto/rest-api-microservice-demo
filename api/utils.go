package api

import (
    "errors"
    "strings"
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

func isPalindrome(s string) bool {
    // TODO
    return true
}
