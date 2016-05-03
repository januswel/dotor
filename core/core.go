package core

import (
	"fmt"
	"os"
	gofp "path/filepath"
	"runtime"
	"sort"

	jwfp "github.com/januswel/dotor/filepath"
	"github.com/januswel/dotor/yaml"
)

// special keys in settings
const (
	DEFAULT_KEY = "default"
	SOURCE_KEY  = "source"
	TARGET_KEY  = "target"
)

type Rule struct {
	source string
	target string
}
type Rules []Rule

func (rules Rules) Len() int {
	return len(rules)
}
func (rules Rules) Swap(i, j int) {
	rules[i], rules[j] = rules[j], rules[i]
}
func (rules Rules) Less(i, j int) bool {
	if rules[i].source == rules[j].source {
		return (rules[i].target < rules[j].target)
	} else {
		return (rules[i].source < rules[j].source)
	}
}

func Execute(settingsFileName, sourcePath string) error {
	settings, err := yaml.ReadFromFile(settingsFileName)
	if err != nil {
		return err
	}

	rules, err := buildRules(settings)
	if err != nil {
		return err
	}

	err = createSymbolicLinks(rules, sourcePath)
	if err != nil {
		return err
	}

	return nil
}

func buildRules(settings map[interface{}]interface{}) (Rules, error) {
	defaultRules, err := buildSpecificOsRules(DEFAULT_KEY, settings)
	if err != nil {
		return nil, err
	}

	osRules, err := buildSpecificOsRules(runtime.GOOS, settings)
	if err != nil {
		return nil, err
	}

	rules := extend(defaultRules, osRules)

	result := Rules{}
	for source, target := range rules {
		rule := Rule{source, target}
		result = append(result, rule)
	}

	sort.Sort(result)

	return result, err
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

func createSymbolicLinks(rules Rules, sourceDirectoryPath string) error {
	targetDirectoryAbsolutePath, err := jwfp.GetHomeDirectory()
	if err != nil {
		return err
	}
	sourceDirectoryAbsolutePath, err := gofp.Abs(sourceDirectoryPath)
	if err != nil {
		return err
	}

	fmt.Printf("%s => %s\n", sourceDirectoryAbsolutePath, targetDirectoryAbsolutePath)

	for _, rule := range rules {
		sourceAbsolutePath := gofp.Join(sourceDirectoryAbsolutePath, rule.source)
		targetAbsolutePath := gofp.Join(targetDirectoryAbsolutePath, rule.target)
		if !jwfp.ExistsPath(sourceAbsolutePath) {
			fmt.Printf("[skipped] \"%s\" source file is not exists.\n", sourceAbsolutePath)
			continue
		}
		if jwfp.ExistsPath(targetAbsolutePath) {
			fmt.Printf("[skipped] \"%s\" target file is already exists.\n", targetAbsolutePath)
			continue
		}

		err := os.Symlink(sourceAbsolutePath, targetAbsolutePath)
		if err != nil {
			panic(err)
		}

		fmt.Printf("[success] \"%s\" => \"%s\"\n", sourceAbsolutePath, targetAbsolutePath)
	}

	return nil
}
