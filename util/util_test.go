package util

import (
	"testing"
)

func TestGetString(t *testing.T) {
	t.Setenv("APP_ENV", "111")

	val := GetString("APP_ENV", "default")

	if val != "111" {
		t.Fatalf(
			"expected %q but got %q",
			"111",
			val,
		)
	}
}

func TestGetString_DefaultValue(t *testing.T) {
	val := GetString("UNKNOWN_ENV", "default")

	if val != "default" {
		t.Fatalf(
			"expected %q but got %q",
			"default",
			val,
		)
	}
}
