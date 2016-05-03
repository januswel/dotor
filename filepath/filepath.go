package filepath

import (
	"fmt"
	goos "os"

	jwos "github.com/januswel/dotor/os"
)

// environment variable names
const (
	HOME_DIRECTORY = "HOME"
)

func ExistsPath(path string) bool {
	_, err := goos.Stat(path)
	return !goos.IsNotExist(err)
}

func GetHomeDirectory() (dir string, err error) {
	environmentVariables := jwos.GetEnvironmentVariables()
	if _, ok := environmentVariables[HOME_DIRECTORY]; !ok {
		return "", fmt.Errorf("Define the environment variable \"%s\"", HOME_DIRECTORY)
	}
	return environmentVariables[HOME_DIRECTORY], nil
}
