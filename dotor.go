package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v2"
)

// yaml special keys
const (
	DEFAULT_KEY = "default"
	SOURCE_KEY  = "source"
	TARGET_KEY  = "target"
)

// environment variable names
const (
	HOME_DIRECTORY = "HOME"
)

// temporary definitions for dev
const (
	SETTINGS_FILE_NAME = "dotorrc.sample.yml"
	SOURCE_PATH        = "/Users/janus/work/dev/github/dotfiles"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, r)
		}
	}()

	settings, err := ReadSettings(SETTINGS_FILE_NAME)
	if err != nil {
		panic(err)
	}

	rules, err := BuildRules(settings)
	if err != nil {
		panic(err)
	}

	err = CreateSymbolicLinks(rules, SOURCE_PATH)
	if err != nil {
		panic(err)
	}
}

func ReadSettings(filepath string) (map[interface{}]interface{}, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	settings := make(map[interface{}]interface{})
	if err := yaml.Unmarshal(file, &settings); err != nil {
		return nil, err
	}

	return settings, nil
}

func BuildRules(settings map[interface{}]interface{}) (map[string]string, error) {
	defaultRules, err := BuildSpecificOsRules(DEFAULT_KEY, settings)
	if err != nil {
		return nil, err
	}

	osRules, err := BuildSpecificOsRules(runtime.GOOS, settings)
	if err != nil {
		return nil, err
	}

	rules := Extend(defaultRules, osRules)
	return rules, err
}

func BuildSpecificOsRules(os string, settings map[interface{}]interface{}) (map[string]string, error) {
	rules := make(map[string]string)
	if _, hasKey := settings[os]; !hasKey {
		return rules, nil
	}

	for _, value := range settings[os].([]interface{}) {
		switch value.(type) {
		case string:
			rules[value.(string)] = value.(string)
		case map[interface{}]interface{}:
			setting := value.(map[interface{}]interface{})
			_, hasSource := setting[SOURCE_KEY]
			_, hasTarget := setting[TARGET_KEY]
			if !(hasSource && hasTarget) {
				return nil, fmt.Errorf("specify key-values in format of \"%s:<source file name>\" and \"%s:<target file name>\" to change sysmbolic link names", SOURCE_KEY, TARGET_KEY)
			}
			source := setting[SOURCE_KEY].(string)
			target := setting[TARGET_KEY].(string)
			rules[source] = target
		}
	}

	return rules, nil
}

func Extend(m1, m2 map[string]string) map[string]string {
	result := map[string]string{}

	for v, k := range m1 {
		result[k] = v
	}
	for v, k := range m2 {
		result[k] = v
	}
	return (result)
}

func CreateSymbolicLinks(rules map[string]string, sourceDirectoryPath string) error {
	targetDirectoryAbsolutePath, err := GetHomeDirectory()
	if err != nil {
		return err
	}
	sourceDirectoryAbsolutePath, err := filepath.Abs(sourceDirectoryPath)
	if err != nil {
		return err
	}

	fmt.Printf("%s => %s\n", sourceDirectoryAbsolutePath, targetDirectoryAbsolutePath)

	for source, target := range rules {
		sourceAbsolutePath := filepath.Join(sourceDirectoryAbsolutePath, source)
		targetAbsolutePath := filepath.Join(targetDirectoryAbsolutePath, target)
		if !ExistsPath(sourceAbsolutePath) {
			fmt.Printf("source file \"%s\" is not exists. skipping.\n", targetAbsolutePath)
			continue
		}
		if ExistsPath(targetAbsolutePath) {
			fmt.Printf("target file \"%s\" is already exists. skipping.\n", targetAbsolutePath)
			continue
		}
		fmt.Printf("creating symbolic link: %s => %s\n", sourceAbsolutePath, targetAbsolutePath)
		// TODO: use os.Symlink()
	}

	return nil
}

func GetHomeDirectory() (dir string, err error) {
	environmentVariables := GetEnvironmentVariables()
	if _, ok := environmentVariables[HOME_DIRECTORY]; !ok {
		return "", fmt.Errorf("Define the environment variable \"%s\"", HOME_DIRECTORY)
	}
	return environmentVariables[HOME_DIRECTORY], nil
}

func GetEnvironmentVariables() map[string]string {
	environmentVariables := make(map[string]string)
	for _, item := range os.Environ() {
		splited := strings.Split(item, "=")
		environmentVariables[splited[0]] = splited[1]
	}
	return environmentVariables
}

func ExistsPath(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
