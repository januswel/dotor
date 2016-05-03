package os

import (
	"testing"
)

func TestShouldExistPATH(t *testing.T) {
	actuals := GetEnvironmentVariables()

	if _, hasPath := actuals["PATH"]; !hasPath {
		t.Error("not found environment variable PATH")
	}
}
