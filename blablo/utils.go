package blablo

import (
	"strings"
)

func TrimOrPadStringRight(s string, targetLength int) string {
	spacer := " "
	result := ""
	if len(s) > targetLength {
		result = s[:targetLength]
	} else if len(s) < targetLength {
		padding := strings.Repeat(spacer, targetLength-len(s))
		result = s + padding
	}
	return result
}
