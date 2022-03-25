package src

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
)

// Reads all the files given, and merges them toghether.
// The last file has priority
//
// Requires at least two fiels to merge
// Returns the merged yaml as a string
func MergeYaml(filePaths ...string) (string, error) {

	// Making an empty map to merge content of files into
	var main map[interface{}]interface{} = make(map[interface{}]interface{})

	// Check that there are any files to merge
	if len(filePaths) <= 1 {
		return "", errors.New("At lease two paths must be given")
	}

	// Loops over path list
	for _, path := range filePaths {

		// Read yaml file
		override, err := readYamlFile(path)
		if err != nil {
			panic(err)
		}

		// Merge the override content into the main map
		mergo.Merge(&main, override, mergo.WithOverride)
	}

	if len(main) == 0 {
		return "", errors.New("Nothing was merged")
	}

	// Converts the yaml struct into a byte array
	bs, err := yaml.Marshal(main)
	if err != nil {
		panic(err)
	}

	return formatYamlString(string(bs)), nil
}

// Reads a yaml file and returns a map with the content
func readYamlFile(filePath string) (map[interface{}]interface{}, error) {
	var yamlFileContent map[interface{}]interface{}

	if filePath == "" {
		return nil, errors.New("No path were given")
	}

	// Reads content of the fil
	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	// Decodes the byte array a yaml struct
	if err := yaml.Unmarshal(bs, &yamlFileContent); err != nil {
		panic(err)
	}

	return yamlFileContent, nil
}

func formatYamlString(s string) string {
	var returnString string

	// Remove all tabs from yaml string
	// Tabs are not allowed in yaml
	returnString = strings.ReplaceAll(s, "\t", "")

	// Remove new line on the left side
	// Yaml cannot begin with new line
	returnString = strings.TrimLeft(returnString, "\n")

	// Replace quadruple spaces with double spaces
	// This is for the indents of maps
	returnString = strings.ReplaceAll(returnString, "    ", "  ")

	return returnString
}
