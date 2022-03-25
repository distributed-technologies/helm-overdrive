package src

import (
	"errors"
	"io/ioutil"
)

// Reads a yaml file and returns a map with the content
func readFile(filePath string) (string, error) {
	if filePath == "" {
		return "", errors.New("No path were given")
	}

	// Reads content of the fil
	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	return formatYamlString(string(bs)), nil
}
