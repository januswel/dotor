package os

import (
	goos "os"
	"strings"
)

func GetEnvironmentVariables() map[string]string {
	environmentVariables := make(map[string]string)
	for _, item := range goos.Environ() {
		splited := strings.Split(item, "=")
		environmentVariables[splited[0]] = splited[1]
	}
	return environmentVariables
}
