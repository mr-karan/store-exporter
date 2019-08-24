package main

import (
	"strconv"
	"strings"
)

// replaceWithUnderscores replaces `-` with `_`
func replaceWithUnderscores(text string) string {
	replacer := strings.NewReplacer(" ", "_", ",", "_", "\t", "_", ",", "_", "/", "_", "\\", "_", ".", "_", "-", "_", ":", "_", "=", "_")
	return replacer.Replace(text)
}

// getStringValue converts supported data types to string
func getStringValue(i interface{}) string {
	switch v := i.(type) {
	case string:
		return v
	case []uint8:
		return string(v)
	case float64:
		return strconv.FormatFloat(v, 'f', 6, 64)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	default:
		return ""
	}
}

// getStringValue converts supported data types to float64
func getFloatValue(i interface{}) float64 {
	var value float64
	switch f := i.(type) {
	case int:
		value = float64(f)
	case int32:
		value = float64(f)
	case int64:
		value = float64(f)
	case uint:
		value = float64(f)
	case uint32:
		value = float64(f)
	case uint64:
		value = float64(f)
	case float32:
		value = float64(f)
	case float64:
		value = float64(f)
	case []uint8:
		val, err := strconv.ParseFloat(string(f), 64)
		if err != nil {
			return value
		}
		value = val
	case string:
		val, err := strconv.ParseFloat(f, 64)
		if err != nil {
			return value
		}
		value = val
	default:
		return value
	}
	return value
}
