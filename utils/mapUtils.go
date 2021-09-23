package utils

func IsEmptyMap(row map[string]interface{}) bool {
	for _, value := range row {
		if value != nil && value != "" {
			return false
		}
	}
	return true
}
