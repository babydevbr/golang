// Package env provide env vars
package env

import "os"

// Get a value or returns a default value.
func Get(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}

	return defaultValue
}
