package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func ParseString(raw json.RawMessage) string {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	var n float64
	if err := json.Unmarshal(raw, &n); err == nil {
		return fmt.Sprintf("%v", n)
	}
	return ""
}

func ParseFloat(raw json.RawMessage) (float64, error) {
	var f float64
	if err := json.Unmarshal(raw, &f); err == nil {
		return f, nil
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return strconv.ParseFloat(s, 64)
	}
	return 0, fmt.Errorf("invalid float format")
}

func ParseInt(raw json.RawMessage) (int, error) {
	var i int
	if err := json.Unmarshal(raw, &i); err == nil {
		return i, nil
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return strconv.Atoi(s)
	}
	return 0, fmt.Errorf("invalid int format")
}

func ParseMap(raw json.RawMessage) *map[string]string {
	var m map[string]string
	if err := json.Unmarshal(raw, &m); err == nil {
		return &m
	}
	return &map[string]string{}
}
