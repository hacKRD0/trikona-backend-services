package utils

import "strings"

// SplitStringArray takes a slice of strings and splits the first element by comma if it exists
// Returns an empty slice if input is empty
func SplitStringArray(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}
	return strings.Split(values[0], ",")
}

// Helper function to convert slice of strings to lowercase
func LowerCase(slice []string) []string {
	lower := make([]string, len(slice))
	for i, s := range slice {
		lower[i] = strings.ToLower(s)
	}
	return lower
}
