package src

import (
	"errors"
	"strings"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
)

// Takes any amont of paths as strings
// Reads them and makes a list of all the content
// Then calls MergeYaml to merge the content together
// Returns the merged Yaml as a string
func MergeYamlFile(filepaths ...string) (string, error) {

	// Check that there are any files to merge
	if len(filepaths) <= 1 {
		return "", errors.New("At lease two file paths must be given")
	}

	var fileContents []string
	for _, path := range filepaths {

		content, err := readFile(path)
		if err != nil {
			panic(err)
		}

		fileContents = append(fileContents, content)
	}

	mergedYaml, err := MergeYaml(fileContents)
	if err != nil {
		panic(err)
	}

	return mergedYaml, nil
}

// Takes a string slice containing yaml content and merges them together
// The override priority is left to right
// The right most / last object has priority
//
// Requires at least two fiels to merge
// Returns the merged yaml as a string
func MergeYaml(yamlObjects []string) (string, error) {

	// Making an empty map to merge content of files into
	var main map[interface{}]interface{} = make(map[interface{}]interface{})

	// Check that there are any files to merge
	if len(yamlObjects) <= 1 {
		return "", errors.New("At lease two yamlObjects must be given")
	}

	// Loops over path list
	for _, override := range yamlObjects {

		var overrideYaml map[interface{}]interface{}
		yaml.Unmarshal([]byte(override), &overrideYaml)

		// Merge the override content into the main map
		mergo.Merge(&main, overrideYaml, mergo.WithOverride)
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
