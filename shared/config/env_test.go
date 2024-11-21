package config

import (
	"os"
	"strconv"
	"testing"
)

func TestGetEnvironWithDefault(t *testing.T) {
	key := "TEST_STRING_VAR"
	defaultValue := "defaultValue"

	// Test when the environment variable is set
	expectedValue := "setValue"
	os.Setenv(key, expectedValue)
	defer os.Unsetenv(key) // Clean up after the test

	result := GetEnvironWithDefault(key, defaultValue)
	if result != expectedValue {
		t.Errorf("Expected %s, but got %s", expectedValue, result)
	}

	// Test when the environment variable is not set
	os.Unsetenv(key) // Ensure the variable is not set
	result = GetEnvironWithDefault(key, defaultValue)
	if result != defaultValue {
		t.Errorf("Expected %s, but got %s", defaultValue, result)
	}
}

func TestGetEnvironIntWithDefault(t *testing.T) {
	key := "TEST_INT_VAR"
	defaultValue := 42

	// Test when the environment variable is set to a valid integer
	expectedValue := 100
	os.Setenv(key, strconv.Itoa(expectedValue))
	defer os.Unsetenv(key) // Clean up after the test

	result := GetEnvironIntWithDefault(key, defaultValue)
	if result != expectedValue {
		t.Errorf("Expected %d, but got %d", expectedValue, result)
	}

	// Test when the environment variable is not set
	os.Unsetenv(key)
	result = GetEnvironIntWithDefault(key, defaultValue)
	if result != defaultValue {
		t.Errorf("Expected %d, but got %d", defaultValue, result)
	}

	// Test when the environment variable is set to a non-integer
	os.Setenv(key, "notAnInt")
	result = GetEnvironIntWithDefault(key, defaultValue)
	if result != defaultValue {
		t.Errorf("Expected %d (default) when non-integer value is set, but got %d", defaultValue, result)
	}
}
