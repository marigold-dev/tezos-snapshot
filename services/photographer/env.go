package main

import (
	"os"
	"strconv"
)

func getEnvInt(key string, defaultValue int) int {
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
