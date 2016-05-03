package yaml

import (
	goyaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

func ReadFromFile(filepath string) (map[interface{}]interface{}, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	settings := make(map[interface{}]interface{})
	if err := goyaml.Unmarshal(file, &settings); err != nil {
		return nil, err
	}

	return settings, nil
}
