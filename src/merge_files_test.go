package src

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
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
func TestMergeYamlNoPaths(t *testing.T) {
	_, err := MergeYaml()

	expectedErrorString := "At lease two paths must be given"

	assert.EqualError(t, err, expectedErrorString)
}

/*
Test whether the error returned when only a single path given fits the expected error
*/
func TestMergeYamlSinglePath(t *testing.T) {
	file1 := genTestYamlFile(``)

	expectedErrorString := "At lease two paths must be given"

	_, err := MergeYaml(file1)

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

	_, err := MergeYaml(file1, file2)

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

	base := genTestYamlFile(s1)
	override := genTestYamlFile(s2)

	result, _ := MergeYaml(base, override)

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

	base := genTestYamlFile(s1)
	override := genTestYamlFile(s2)

	result, _ := MergeYaml(base, override)

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

	base := genTestYamlFile(s1)
	override := genTestYamlFile(s2)

	result, _ := MergeYaml(base, override)

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

	base := genTestYamlFile(s1)
	override := genTestYamlFile(s2)

	result, _ := MergeYaml(base, override)

	assert.Equal(t, expectedYaml, result)
}

// -------------------- ReadYaml Tests --------------------------

/*
Test that we can read a file and get a expected result
*/
func TestReadYamlFileWithFile(t *testing.T) {
	s1 := formatYamlString(`
	name: Foo
	sports:
	  - football
	  - tennis
	`)

	tfName := genTestYamlFile(s1)

	// Use the readYamlFile to read the content of the test file
	content, err := readYamlFile(tfName)
	if err != nil {
		t.Fatalf("Read yaml test file: %v", err)
	}

	// Create the variable to hold the yaml content
	var expectedContent map[interface{}]interface{}

	// Unmarshal the content into maps
	err = yaml.Unmarshal([]byte(s1), &expectedContent)
	if err != nil {
		t.Fatalf("Generate expected content: %v", err)
	}

	// Check that the content read from the file and the content unmarshaled is the same
	assert.NotNil(t, content)
	assert.Equal(t, content, expectedContent)
}

/*
Test that if no file given readYaml throws an error
*/
func TestReadYamlFileNoPath(t *testing.T) {

	expectedErrorString := "No path were given"

	content, err := readYamlFile("")

	assert.Nil(t, content)
	assert.EqualError(t, err, expectedErrorString)
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
