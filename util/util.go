package util

import (
	"os"
)

// GetString is a function which returns a value by looking up the environment varibale
// Returns the fallback value if the evenronment varibale key doesn't exist
func GetString(env, fallback string) string {
	if value, ok := os.LookupEnv(env); ok {
		return value
	}
	return fallback
}
