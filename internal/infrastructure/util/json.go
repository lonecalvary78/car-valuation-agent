package util

import (
	"encoding/json"

	"github.com/bytedance/sonic"
)

func FromJSON(data string, dto any) error {
	return sonic.Unmarshal([]byte(data), dto)
}

func ToJSON(dto any) ([]byte, error) {
	return json.Marshal(dto)
}
