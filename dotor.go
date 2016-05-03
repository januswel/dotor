package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"gopkg.in/yaml.v2"
)

// yaml special keys
const (
	DEFAULT_KEY = "default"
	SOURCE_KEY  = "source"
	TARGET_KEY  = "target"
)

// temporary definitions for dev
const (
	SETTINGS_FILE_NAME = "dotorrc.sample.yml"
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

	fmt.Println(rules)
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
