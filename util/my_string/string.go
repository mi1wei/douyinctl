package my_string

import "strings"

func Contains(keywords []string, str string) bool {
	for _, keyword := range keywords {
		if strings.Contains(str, keyword) {
			return true
		}
	}
	return false
}
