package filepath

import (
	"testing"
)

func TestFailUnexistPath(t *testing.T) {
	const UNEXIST_PATH = "3kf9l/q93mk/103dh/vuhyo/9lei9"
	if ExistsPath(UNEXIST_PATH) {
		t.Errorf("you have really %s ?", UNEXIST_PATH)
	}
}
