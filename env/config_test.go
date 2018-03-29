package env

import (
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	os.Setenv("TEST_KEY", "value")
	defer os.Unsetenv("TEST_KEY")

	value := Get("TEST_KEY")
	if value != "value" {
		t.Error("Unexpected value:", value)
	}

	value = Get("NonExistingKey")
	if value != "" {
		t.Error("Unexpected value:", value)
	}
}

func TestGetOrDefault(t *testing.T) {
	os.Setenv("TEST_KEY", "value")
	defer os.Unsetenv("TEST_KEY")

	value := GetOrDefault("TEST_KEY", "default1")
	if value != "value" {
		t.Error("Unexpected value:", value)
	}

	value = GetOrDefault("NonExistingKey", "default2")
	if value != "default2" {
		t.Error("Unexpected value:", value)
	}
}

func TestLookup(t *testing.T) {
	os.Setenv("TEST_KEY", "value")
	defer os.Unsetenv("TEST_KEY")
	os.Setenv("FROM_FILE", "from-env")
	defer os.Unsetenv("FROM_FILE")

	value := Lookup("FROM_FILE", "../testdata/config.txt", "missing")
	if value != "file-value" {
		t.Error("Unexpected value:", value)
	}

	value = Lookup("TEST_KEY", "../testdata/config.txt", "missing")
	if value != "value" {
		t.Error("Unexpected value:", value)
	}

	value = Lookup("DEFAULT_VALUE", "../testdata/config.txt", "default")
	if value != "default" {
		t.Error("Unexpected value:", value)
	}
}

func TestGetInt(t *testing.T) {
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")
	os.Setenv("TEST_FLOAT", "42.3")
	defer os.Unsetenv("TEST_FLOAT")

	value := GetInt("TEST_INT", "", 0)
	if value != 42 {
		t.Error("Unexpected value:", value)
	}

	value = GetInt("TEST_FLOAT", "", 1)
	if value != 1 {
		t.Error("Unexpected value:", value)
	}

	value = GetInt("TEST_DEFAULT", "", 2)
	if value != 2 {
		t.Error("Unexpected value:", value)
	}
}

func TestGetInterval(t *testing.T) {
	os.Setenv("TEST_INTERVAL", "3s200ms")
	defer os.Unsetenv("TEST_INTERVAL")
	os.Setenv("TEST_INVALID", "xabc")
	defer os.Unsetenv("TEST_INVALID")

	value := GetInterval("TEST_INTERVAL", "")
	if value.Seconds() != 3.2 {
		t.Error("Unexpected duration:", value)
	}

	value = GetInterval("TEST_INVALID", "")
	if value != DefaultInterval {
		t.Error("Unexpected duration:", value)
	}

	value = GetInterval("TEST_DEFAULT", "")
	if value != DefaultInterval {
		t.Error("Unexpected default duration:", value)
	}
}

func TestGetTimeout(t *testing.T) {
	os.Setenv("TEST_TIMEOUT", "4s150ms")
	defer os.Unsetenv("TEST_DURATION")
	os.Setenv("TEST_INVALID", "xabc")
	defer os.Unsetenv("TEST_INVALID")

	value := GetTimeout("TEST_TIMEOUT", "")
	if value.Seconds() != 4.15 {
		t.Error("Unexpected duration:", value)
	}

	value = GetTimeout("TEST_INVALID", "")
	if value != DefaultTimeout {
		t.Error("Unexpected duration:", value)
	}

	value = GetTimeout("TEST_DEFAULT", "")
	if value != DefaultTimeout {
		t.Error("Unexpected default duration:", value)
	}
}
