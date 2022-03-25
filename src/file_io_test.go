package src

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	content, err := readFile(tfName)
	if err != nil {
		t.Fatalf("Read yaml test file: %v", err)
	}

	// Check that the content read from the file and the content unmarshaled is the same
	assert.NotNil(t, content)
	assert.Equal(t, s1, content)
}

/*
Test that if no file given readYaml throws an error
*/
func TestReadYamlFileNoPath(t *testing.T) {
	expectedErrorString := "No path were given"

	_, err := readFile("")

	assert.EqualError(t, err, expectedErrorString)
}
