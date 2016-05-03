package main

import (
	"fmt"
	"os"
	gofp "path/filepath"
	"runtime"

	jwfp "github.com/januswel/dotor/filepath"
	"github.com/januswel/dotor/yaml"
)

// special keys in settings
const (
	DEFAULT_KEY = "default"
	SOURCE_KEY  = "source"
	TARGET_KEY  = "target"
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

	settings, err := yaml.ReadFromFile(SETTINGS_FILE_NAME)
	if err != nil {
		panic(err)
	}

	rules, err := buildRules(settings)
	if err != nil {
		panic(err)
	}

	err = createSymbolicLinks(rules, SOURCE_PATH)
	if err != nil {
		panic(err)
	}
}

func buildRules(settings map[interface{}]interface{}) (map[string]string, error) {
	defaultRules, err := buildSpecificOsRules(DEFAULT_KEY, settings)
	if err != nil {
		return nil, err
	}

	osRules, err := buildSpecificOsRules(runtime.GOOS, settings)
	if err != nil {
		return nil, err
	}

	rules := extend(defaultRules, osRules)
	return rules, err
}

func buildSpecificOsRules(os string, settings map[interface{}]interface{}) (map[string]string, error) {
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

func extend(m1, m2 map[string]string) map[string]string {
	result := map[string]string{}

	for v, k := range m1 {
		result[k] = v
	}
	for v, k := range m2 {
		result[k] = v
	}
	return (result)
}

func createSymbolicLinks(rules map[string]string, sourceDirectoryPath string) error {
	targetDirectoryAbsolutePath, err := jwfp.GetHomeDirectory()
	if err != nil {
		return err
	}
	sourceDirectoryAbsolutePath, err := gofp.Abs(sourceDirectoryPath)
	if err != nil {
		return err
	}

	fmt.Printf("%s => %s\n", sourceDirectoryAbsolutePath, targetDirectoryAbsolutePath)

	for source, target := range rules {
		sourceAbsolutePath := gofp.Join(sourceDirectoryAbsolutePath, source)
		targetAbsolutePath := gofp.Join(targetDirectoryAbsolutePath, target)
		if !jwfp.ExistsPath(sourceAbsolutePath) {
			fmt.Printf("source file \"%s\" is not exists. skipping.\n", targetAbsolutePath)
			continue
		}
		if jwfp.ExistsPath(targetAbsolutePath) {
			fmt.Printf("target file \"%s\" is already exists. skipping.\n", targetAbsolutePath)
			continue
		}
		fmt.Printf("creating symbolic link: %s => %s\n", sourceAbsolutePath, targetAbsolutePath)
		// TODO: use os.Symlink()
	}

	return nil
}
