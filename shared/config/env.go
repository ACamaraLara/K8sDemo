package config

import (
	"os"
	"strconv"
)

// GetEnvironWithDefault returns an environment variable as a string or a default value if not found.
func GetEnvironWithDefault(key, defaultValue string) string {
	if value, present := os.LookupEnv(key); present {
		return value
	}
	return defaultValue
}

// GetEnvironIntWithDefault returns an environment variable as an int or a default value if not found.
// If the value cannot be parsed to an int, the default value is returned.
func GetEnvironIntWithDefault(key string, defaultValue int) int {
	value := GetEnvironWithDefault(key, "")
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	return defaultValue
}
