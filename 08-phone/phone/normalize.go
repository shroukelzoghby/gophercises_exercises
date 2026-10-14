package phone

import "unicode"

// Normalize strips all non-digit characters from a phone number,
// leaving only digits in ########## form.
func Normalize(number string) string {
	var digits []rune
	for _, r := range number {
		if unicode.IsDigit(r) {
			digits = append(digits, r)
		}
	}
	return string(digits)
}
