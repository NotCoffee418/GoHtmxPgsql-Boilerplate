package utils

// SliceContains checks if a slice contains an item
func SliceContains[T comparable](list []T, item T) bool {
	for _, listItem := range list {
		if listItem == item {
			return true
		}
	}
	return false
}
