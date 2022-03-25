package src

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// -------------------- formatYamlString Tests --------------------------

/*
Test that formatYamlString removes all tabs
*/
func TestFormatYamlStringRemoveTabs(t *testing.T) {
	s1 := formatYamlString("\tname: Foo")

	expectedString := "name: Foo"

	assert.Equal(t, expectedString, s1)
}

/*
Test that formatYamlString removes beginngin newlines
*/
func TestFormatYamlStringRemoveBeginningNewLine(t *testing.T) {
	s1 := formatYamlString("\nname: Foo")

	expectedString := "name: Foo"

	assert.Equal(t, expectedString, s1)
}

/*
Test that formatYamlString replaces 4 spaces with 2 spaces
This is because the merger replaces the indented map content with 4 spaces

f.eks.
map:
  content: 1
  content2: 2
  content3: 3

Is made into:
map:
    content: 1
    content2: 2
    content3: 3

*/
func TestFormatYamlStringReplce4SpacesWith2Spaces(t *testing.T) {
	s1 := formatYamlString("name: Foo\nmarks:\n    eng: 4")

	expectedString := "name: Foo\nmarks:\n  eng: 4"

	assert.Equal(t, expectedString, s1)
}

// -------------------- MergeYaml Tests --------------------------

/*
Test whether the error returned when no path given fits the expected error
*/
func TestMergeYamlNoYamlObjects(t *testing.T) {
	_, err := MergeYaml([]string{""})

	expectedErrorString := "At lease two yamlObjects must be given"

	assert.EqualError(t, err, expectedErrorString)
}

/*
Test whether the error returned when only a single path given fits the expected error
*/
func TestMergeYamlSingleYamlObject(t *testing.T) {
	file1 := genTestYamlFile(``)

	expectedErrorString := "At lease two yamlObjects must be given"

	_, err := MergeYaml([]string{file1})

	assert.EqualError(t, err, expectedErrorString)
}

/*
Test whether the error returned when there is no content in the files given,
fits the expected error
*/
func TestMergeYamlNoContent(t *testing.T) {
	file1 := genTestYamlFile(``)
	file2 := genTestYamlFile(``)

	expectedErrorString := "Nothing was merged"

	_, err := MergeYaml([]string{file1, file2})

	assert.EqualError(t, err, expectedErrorString)
}

/*
Test whether the merging the map returns the result of the last merged map
*/
func TestMergeYamlMapValue(t *testing.T) {
	s1 := formatYamlString(`
	name: Foo
	`)
	s2 := formatYamlString(`
	name: Bar
	`)
	expectedYaml := formatYamlString(`
	name: Bar
	`)

	result, _ := MergeYaml([]string{s1, s2})

	assert.Equal(t, expectedYaml, result)
}

/*
Test whether the merging the map returns the result of the last merged map
without setting the other values to nil
*/
func TestMergeYamlMapValueOverride(t *testing.T) {
	s1 := formatYamlString(`
	name: Foo
	`)
	s2 := formatYamlString(`
	name: Bar
	age: 23
	`)
	expectedYaml := formatYamlString(`
	age: 23
	name: Bar
	`)

	result, _ := MergeYaml([]string{s1, s2})

	assert.Equal(t, expectedYaml, result)
}

/*
Test that marging a map inside a map, merges the values
*/
func TestMergeYamlMapValueInMap(t *testing.T) {
	s1 := formatYamlString(`
	name: Foo
	marks:
	  math: 12
	`)

	s2 := formatYamlString(`
	name: Bar
	marks:
	  eng: 7
	`)

	expectedYaml := formatYamlString(`
	marks:
	  eng: 7
	  math: 12
	name: Bar
	`)

	result, _ := MergeYaml([]string{s1, s2})

	assert.Equal(t, expectedYaml, result)
}

/*
Test that merging slices overrids the entire slice instead of appending to it
*/
func TestMergeYamlOverrideSlice(t *testing.T) {
	s1 := formatYamlString(`
	name: Foo
	sports:
	  - football
	  - tennis
	`)

	s2 := formatYamlString(`
	sports:
	  - Boxing
	`)

	expectedYaml := formatYamlString(`
	name: Foo
	sports:
	  - Boxing
	`)

	result, _ := MergeYaml([]string{s1, s2})

	assert.Equal(t, expectedYaml, result)
}

// -------------------- MergeYamlFile Tests --------------------------

/*
Test whether the error returned when no path given fits the expected error
*/
func TestMergeYamlFileNoYamlObjects(t *testing.T) {
	_, err := MergeYamlFile()

	expectedErrorString := "At lease two file paths must be given"

	assert.EqualError(t, err, expectedErrorString)
}

/*
Test whether the error returned when only a single path given fits the expected error
*/
func TestMergeYamlFileSingleYamlObject(t *testing.T) {
	file1 := genTestYamlFile(``)

	expectedErrorString := "At lease two file paths must be given"

	_, err := MergeYamlFile(file1)

	assert.EqualError(t, err, expectedErrorString)
}

/*
Test whether the merging the map returns the result of the last merged map
*/
func TestMergeYamlFileMapValue(t *testing.T) {
	s1 := formatYamlString(`
	name: Foo
	`)
	s2 := formatYamlString(`
	name: Bar
	`)
	expectedYaml := formatYamlString(`
	name: Bar
	`)

	file1 := genTestYamlFile(s1)
	file2 := genTestYamlFile(s2)

	result, _ := MergeYamlFile(file1, file2)

	assert.Equal(t, expectedYaml, result)
}

// -------------------- Helper functions --------------------------

/*
Generates a temp file and writes a string to it

Returnes the files name (which contains the path)
*/
func genTestYamlFile(s string) string {
	// Generate temp file
	tf, err := os.CreateTemp("", "")
	if err != nil {
		fmt.Printf("create test file: %v", err)
	}

	// Write yaml data to temp file
	err = ioutil.WriteFile(tf.Name(), []byte(s), 0644)
	if err != nil {
		panic("Unable to write data into the file")
	}

	// Closes the file to modification
	if err := tf.Close(); err != nil {
		fmt.Printf("close test file: %v", err)
	}

	// Returns the filename
	return tf.Name()
}
