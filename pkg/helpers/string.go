package helpers

import "unicode"

func HasUppercase(s string) bool {
	for _, c := range s {
		if unicode.IsUpper(c) {
			return true
		}
	}
	return false
}

func HasDigit(s string) bool {
	for _, c := range s {
		if c > '0' && c < '9' {
			return true
		}
	}
	return false
}
