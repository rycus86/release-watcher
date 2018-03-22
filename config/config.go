package config

import (
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var DefaultTimeout = 30 * time.Second

func Get(key string) string {
	return GetOrDefault(key, "")
}

func GetOrDefault(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	} else {
		return defaultValue
	}
}

func Lookup(key string, path string, defaultValue string) string {
	contents, err := ioutil.ReadFile(path)
	if err == nil {
		for _, line := range strings.Split(string(contents), "\n") {
			parameters := strings.SplitN(line, "=", 2)

			if parameters[0] == key {
				return strings.TrimSpace(parameters[1])
			}
		}
	}

	return GetOrDefault(key, defaultValue)
}

func GetDuration(key string, path string) time.Duration {
	if duration, err := time.ParseDuration(Lookup(key, path, "")); err == nil {
		return duration
	}

	return DefaultTimeout
}
