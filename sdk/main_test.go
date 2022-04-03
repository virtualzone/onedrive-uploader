package sdk

import (
	"runtime/debug"
	"testing"
)

func checkTestBool(t *testing.T, expected, actual bool) {
	if expected != actual {
		t.Fatalf("Expected '%t', but got '%t' at:\n%s", expected, actual, debug.Stack())
	}
}

func checkTestString(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Fatalf("Expected '%s', but got '%s' at:\n%s", expected, actual, debug.Stack())
	}
}
