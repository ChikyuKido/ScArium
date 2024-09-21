package helper

import "strings"

func ConvertStringArrayToDatabaseString(arr []string) string {
	return strings.Join(arr, ",")
}
func ConvertDatabaseStringToStringArray(text string) []string {
	return strings.Split(text, ",")
}
