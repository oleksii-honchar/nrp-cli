package configProcessor

import (
	"strings"
)

func checkIfStrContainsAny(str string, keywords []string) bool {
	var res bool = false
	for _, keyword := range keywords {
		if res = strings.Contains(str, keyword); res {
			break
		}
	}

	return res
}
