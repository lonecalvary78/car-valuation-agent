package util

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bytedance/sonic"
)

func FromJSON(data string, dto any) error {
	err := sonic.Unmarshal([]byte(data), dto)
	if err != nil {
		return fmt.Errorf("util: failed to unmarshal json: %w", err)
	}
	return nil
}

func ToJSON(dto any) ([]byte, error) {
	marshaled, err := json.Marshal(dto)
	if err != nil {
		return nil, fmt.Errorf("util: failed to marshal json: %w", err)
	}
	return marshaled, nil
}

func ExtractJSONObject(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start < 0 || end < start {
		return s
	}
	return s[start : end+1]
}
