package utils

import (
	"log"
	"os"
	"strconv"
)

func GetEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return ""
}

func GetInt32Env(key string, defaultValue int32) int32 {
	if value, exists := os.LookupEnv(key); exists {
		parsedValue, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("Error parsing environment variable %s: %v. Using default value: %d", key, err, defaultValue)
			return defaultValue
		}
		return int32(parsedValue)
	}
	return defaultValue
}
