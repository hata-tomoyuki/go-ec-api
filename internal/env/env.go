package env

import (
	"fmt"
	"os"
)

func GetString(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// MustGetString returns the value of the environment variable key.
// It panics if the variable is not set or empty.
func MustGetString(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("required environment variable %q is not set", key))
	}
	return val
}
