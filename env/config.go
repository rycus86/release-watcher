package env

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var DefaultInterval = 4 * time.Hour
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

func GetInt(key string, path string, defaultValue int) int {
	value := Lookup(key, path, "")

	if asInt, err := strconv.Atoi(value); err == nil {
		return asInt
	}

	return defaultValue
}

func GetInterval(key string, path string) time.Duration {
	return LookupDuration(key, path, DefaultInterval)
}

func GetTimeout(key string, path string) time.Duration {
	return LookupDuration(key, path, DefaultTimeout)
}

func LookupDuration(key string, path string, defaultValue time.Duration) time.Duration {
	if duration, err := time.ParseDuration(Lookup(key, path, "")); err == nil {
		return duration
	}

	return defaultValue
}
