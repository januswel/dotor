package filepath

import (
	"runtime"
	"testing"
)

func TestShouldExistHomeDirectoryOnMacosx(t *testing.T) {
	if runtime.GOOS != "darwin" {
		return
	}

	if _, err := GetHomeDirectory(); err != nil {
		t.Error("environment variable HOME does not exist")
	}
}

func TestShouldExistHomeDirectoryOnLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		return
	}

	if _, err := GetHomeDirectory(); err != nil {
		t.Error("environment variable HOME does not exist")
	}
}
