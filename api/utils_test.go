package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPalindrome(t *testing.T) {
	assert.False(t, isPalindrome("test"), "incorrect result")
	assert.False(t, isPalindrome("racecard"), "incorrect result")

	assert.True(t, isPalindrome(" "), "incorrect result")
	assert.True(t, isPalindrome("?"), "incorrect result")
	assert.True(t, isPalindrome("??"), "incorrect result")
	assert.True(t, isPalindrome("racecar"), "incorrect result")
	assert.True(t, isPalindrome("rac.ecar"), "incorrect result")
	assert.True(t, isPalindrome("r!ac.ecar"), "incorrect result")
	assert.True(t, isPalindrome("r!ac.ec  ar"), "incorrect result")
}
