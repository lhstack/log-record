package utils

import (
	"fmt"
	"os"
	"strconv"
)

func EnvGetOrDefault[T comparable](key string, defaultValue T, apply func(value string) T) T {
	v, ok := os.LookupEnv(key)
	if ok {
		return apply(v)
	}
	return defaultValue
}

func EnvGetOrDefaultStringValue(key string, defaultValue string) string {
	v, ok := os.LookupEnv(key)
	if ok {
		return v
	}
	return defaultValue
}

func EnvGetOrDefaultIntValue(key string, defaultValue int) int {
	return EnvGetOrDefault(key, defaultValue, func(v string) int {
		intV, err := strconv.Atoi(v)
		if err != nil {
			fmt.Println(err)
			return defaultValue
		}
		return intV
	})
}
