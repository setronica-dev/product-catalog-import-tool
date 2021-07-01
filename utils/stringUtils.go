package utils

import "strings"

func TrimAll(input string) string {
	var res string
	res = strings.Replace(
		strings.Replace(
			strings.Replace(
				input, "	", "", -1),
			"*", "", -1),
		" ", "", -1)

	return strings.ToLower(res)
}
