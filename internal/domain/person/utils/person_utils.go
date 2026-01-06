package person

import (
	"strings"
	"unicode"
)

func OnlyDigits(input string) string {
	var b strings.Builder

	for _, r := range input {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}

	return b.String()
}
