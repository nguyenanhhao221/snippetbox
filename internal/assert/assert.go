package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("expected: %v, actual: %v\n", expected, actual)
	}
}

func StringContains(t *testing.T, actual, expectedSubString string) {
	t.Helper()

	if !strings.Contains(actual, expectedSubString) {
		t.Errorf("got :%q, expected to contain: %q", actual, expectedSubString)
	}
}
