package util

import (
	"os"
	"strconv"
	"strings"
)

func GetEnvInt(key string, defaultValue int) int {
	stringValue := os.Getenv(key)
	if stringValue == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(stringValue)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func GetEnvBool(key string, defaultValue bool) bool {
	stringValue := os.Getenv(key)
	if stringValue == "" {
		return defaultValue
	}

	return strings.Contains(strings.ToLower(stringValue), "true")
}
