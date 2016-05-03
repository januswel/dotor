package yaml

import (
	"fmt"
	"os"
	"path"
	"testing"
)

func TestShouldRead(t *testing.T) {
	const TARGET_FILE_NAME = "test.yml"

	currentDirectory, _ := os.Getwd()
	targetFilePath := path.Join(currentDirectory, TARGET_FILE_NAME)
	fmt.Println(targetFilePath)

	if _, err := ReadFromFile(targetFilePath); err != nil {
		t.Errorf("can not read %s", TARGET_FILE_NAME)
	}
}
