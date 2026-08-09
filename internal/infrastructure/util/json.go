package util

import (
	"encoding/json"
	"strings"

	"github.com/bytedance/sonic"
)

func FromJSON(data string, dto any) error {
	return sonic.Unmarshal([]byte(data), dto)
}

func ToJSON(dto any) ([]byte, error) {
	return json.Marshal(dto)
}

func ExtractJSONObject(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start < 0 || end < start {
		return s
	}
	return s[start : end+1]
}
