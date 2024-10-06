package helper

import (
	"encoding/json"
)

func ConvertStructToJsonOr(s interface{}, def string) string {
	jsonData, err := json.Marshal(s)
	if err != nil {
		return def
	}
	return string(jsonData)
}
