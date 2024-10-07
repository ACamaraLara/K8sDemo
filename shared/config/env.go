package config

import (
	"os"
	"strconv"
)

// Returns an environment string variable or a default value if not found.
func GetEnvironWithDefault(key string, defaultValue string) string {
	value, present := os.LookupEnv(key)
	if present {
		return value
	} else {
		return defaultValue
	}
}

// Returns an environment int variable or a default value if not found.
func GetEnvironIntWithDefault(key string, defaultValue int) int {
	value, present := os.LookupEnv(key)
	if present {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		} else {
			return defaultValue
		}
	} else {
		return defaultValue
	}
}
