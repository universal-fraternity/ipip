// Package utils provides a unified processing method
package utils

import (
	"strconv"
	"strings"
)

// String2Int Convert string to int
func String2Int(s string) (int, error) {
	if s == "" || s == "*" {
		return 0, nil
	}
	return strconv.Atoi(s)
}

// String2Int32 Convert string to int32
func String2Int32(s string) (int32, error) {
	if s == "" || s == "*" {
		return 0, nil
	}
	// 使用 strconv.ParseInt 转换为 int64
	val, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(val), nil
}

// String2Int64 Convert string to int64
func String2Int64(s string) (int64, error) {
	if s == "" || s == "*" {
		return 0, nil
	}
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// String2Float64 Convert string to float64
func String2Float64(s string) (float64, error) {
	if s == "" || s == "*" {
		return 0, nil
	}

	return strconv.ParseFloat(s, 64)
}

// RefineOutput Refine output
func RefineOutput(s string) string {
	if s == "*" {
		return ""
	}
	return strings.TrimSpace(s)
}

// IsIPv4 Is it IPv4
func IsIPv4(s string) bool {
	return strings.Contains(s, ".")
}

// IsIPv6 Is it IPv6
func IsIPv6(s string) bool {
	return strings.Contains(s, ":")
}
