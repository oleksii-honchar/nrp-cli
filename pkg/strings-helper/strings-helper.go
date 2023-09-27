package stringsHelper

import (
	"strings"
)

func CheckIfStrInArray(str string, arr []string) bool {
	var res bool = false
	for _, s := range arr {
		if s == str {
			res = true
			break
		}
	}

	return res
}

func CheckIfStrContainsAny(str string, keywords []string) bool {
	var res bool = false
	for _, keyword := range keywords {
		if res = strings.Contains(str, keyword); res {
			break
		}
	}

	return res
}
