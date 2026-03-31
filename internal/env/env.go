package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func init() {
	loadDotenv(".env")
}

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

// loadDotenv reads a .env file and sets environment variables.
// Already-set variables are not overwritten. Missing file is silently ignored.
func loadDotenv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		// remove surrounding quotes
		if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
			val = val[1 : len(val)-1]
		}
		// don't overwrite existing env vars
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}
