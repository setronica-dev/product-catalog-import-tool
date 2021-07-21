package utils

func IsEmptyRow(input []string) bool {
	for _, v := range input {
		if v != "" {
			return false
		}
	}
	return true
}

func InArray(val string, array []string) (exists bool, index int) {
	exists = false
	index = -1

	for i, v := range array {
		if val == v {
			index = i
			exists = true
			return exists, index
		}
	}

	return exists, index
}
