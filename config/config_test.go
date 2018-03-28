package config

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
