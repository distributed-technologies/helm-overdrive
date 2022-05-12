package template

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var thws = TempHelmWorkspace{
	chartName: "foo",
	// TmpHelmDir, is set on every test to have different tmp folders for each test
	ReleaseName: "bar",
}

func TestGetChartFolder(t *testing.T) {
	thws.TmpHelmDir = t.TempDir()

	expectedString := fmt.Sprintf("%s/%s", thws.TmpHelmDir, thws.chartName)
	actualString := thws.getChartFolder()

	if actualString != expectedString {
		t.Fatalf("getChartFolder() = %q, wanted %s, error", actualString, expectedString)
	}
}

func TestGetWorkspacePaths(t *testing.T) {
	thws.TmpHelmDir = t.TempDir()

	tests := []struct {
		expectedString string
		actualString   string
	}{
		{expectedString: fmt.Sprintf("%s/%s/%s", thws.TmpHelmDir, thws.chartName, "charts"), actualString: thws.getChartsFolderLocation()},
		{expectedString: fmt.Sprintf("%s/%s/%s", thws.TmpHelmDir, thws.chartName, "templates"), actualString: thws.getTemplatesFolderLocation()},
		{expectedString: fmt.Sprintf("%s/%s/%s", thws.TmpHelmDir, thws.chartName, "values.yaml"), actualString: thws.getValuesFileLocation()},
	}

	for _, test := range tests {
		expectedString := test.expectedString
		actualString := test.actualString

		if actualString != expectedString {
			t.Fatalf("expected = %q, wanted %s, error", actualString, expectedString)
		}
	}
}

func TestCreateHelmChartSuccess(t *testing.T) {
	thws.TmpHelmDir = t.TempDir()

	err = thws.CreateHelmChart()
	if err != nil {
		t.Fatalf("CreateHelmChart threw err %q", err)
	}

	assert.FileExists(t, thws.getValuesFileLocation())
	assert.DirExists(t, thws.getTemplatesFolderLocation())
	assert.DirExists(t, thws.getChartFolder())
}

func TestCreateHelmChartFailed(t *testing.T) {
	thws.TmpHelmDir = t.TempDir()

	create = func(name string) (*os.File, error) {
		return nil, errors.New("")
	}

	err = thws.CreateHelmChart()
	if err == nil {
		t.Fatalf("CreateHelmChart threw err %q", err)
	}

	mkdir = func(name string, perm os.FileMode) error {
		return errors.New("")
	}

	err = thws.CreateHelmChart()
	if err == nil {
		t.Fatalf("CreateHelmChart threw err %q", err)
	}

	remove = func(name string) error {
		return errors.New("")
	}

	err = thws.CreateHelmChart()
	if err == nil {
		t.Fatalf("CreateHelmChart threw err %q", err)
	}

	removeAll = func(path string) error {
		return errors.New("")
	}

	err = thws.CreateHelmChart()
	if err == nil {
		t.Fatalf("CreateHelmChart threw err %q", err)
	}

	cmdRun = func(c *exec.Cmd) error {
		return errors.New("")
	}

	err = thws.CreateHelmChart()
	if err == nil {
		t.Fatalf("CreateHelmChart threw err %q", err)
	}

	t.Cleanup(setStubsToDefault)
}

func TestAddFileToTempalteSuccess(t *testing.T) {
	thws.TmpHelmDir = t.TempDir()

	// Generate files for test
	tmp_dir := t.TempDir()
	file1Name := "cm.yaml"
	file1 := makeTestFiles(t, file1Name, tmp_dir).Name()

	// Create test helm chart to move files into
	err := thws.CreateHelmChart()
	if err != nil {
		t.Fatalf("CreateHelmChart threw err %q", err)
	}

	thws.AddFileToTemplate(file1)

	expectedFile1 := fmt.Sprintf("%s/%s", thws.getTemplatesFolderLocation(), filepath.Base(file1))

	assert.FileExists(t, expectedFile1)
}

func TestAddFileToTempalteFailed(t *testing.T) {
	thws.TmpHelmDir = t.TempDir()

	// Generate files for test
	tmp_dir := t.TempDir()
	file1Name := "cm.yaml"
	file1 := makeTestFiles(t, file1Name, tmp_dir).Name()

	// Create test helm chart to move files into
	err = thws.CreateHelmChart()
	if err != nil {
		t.Fatalf("CreateHelmChart threw err %q", err)
	}

	// Making (*exec.Cmd).run return an error
	cmdRun = func(c *exec.Cmd) error {
		return errors.New("")
	}
	t.Cleanup(setStubsToDefault)

	err = thws.AddFileToTemplate(file1)

	if err == nil {
		t.Errorf("err = %q, wanted err to not be nil", err)
	}

	// Making filepath.Glob return an error
	glob = func(pattern string) (matches []string, err error) {
		return nil, errors.New("")
	}
	t.Cleanup(setStubsToDefault)

	err := thws.AddFileToTemplate("test")

	if err == nil {
		t.Errorf("err = %q, wanted err to not be nil", err)
	}
}

func TestAddDirToTempalteSuccess(t *testing.T) {
	thws.TmpHelmDir = t.TempDir()

	// Generate files for test
	tmp_dir := t.TempDir()
	file1Name := "cm.yaml"
	file2Name := "sec.yaml"
	file1 := filepath.Base(makeTestFiles(t, file1Name, tmp_dir).Name())
	file2 := filepath.Base(makeTestFiles(t, file2Name, tmp_dir).Name())

	// Create test helm chart to move files into
	err := thws.CreateHelmChart()
	if err != nil {
		t.Fatalf("CreateHelmChart threw err %q", err)
	}
	thws.AddDirToTemplate(tmp_dir)

	expectedDir := fmt.Sprintf("%s/%s", thws.getTemplatesFolderLocation(), filepath.Base(tmp_dir))
	expectedFile1 := fmt.Sprintf("%s/%s/%s", thws.getTemplatesFolderLocation(), filepath.Base(tmp_dir), file1)
	expectedFile2 := fmt.Sprintf("%s/%s/%s", thws.getTemplatesFolderLocation(), filepath.Base(tmp_dir), file2)

	assert.DirExists(t, expectedDir)
	assert.FileExists(t, expectedFile1)
	assert.FileExists(t, expectedFile2)
}

func TestAddDirToTempalteOsStatFail(t *testing.T) {
	thws.TmpHelmDir = t.TempDir()

	// Generate files for test
	tmp_dir := t.TempDir()

	// Create test helm chart to move files into
	err := thws.CreateHelmChart()
	if err != nil {
		t.Fatalf("CreateHelmChart threw err %q", err)
	}

	osStat = func(name string) (os.FileInfo, error) {
		return nil, errors.New("")
	}

	err = thws.AddDirToTemplate(tmp_dir)

	if err == nil {
		t.Fatalf("err = %q, wanted err", err)
	}
}

func TestAddDirToTempalteAddFileToTmplateFail(t *testing.T) {
	thws.TmpHelmDir = t.TempDir()

	// Generate files for test
	tmp_dir := t.TempDir()

	// Create test helm chart to move files into
	err := thws.CreateHelmChart()
	if err != nil {
		t.Fatalf("CreateHelmChart threw err %q", err)
	}

	// Set glob to return error so that AddDirToTemplates returns an error
	glob = func(pattern string) (matches []string, err error) {
		return nil, errors.New("glob failed")
	}
	t.Cleanup(setStubsToDefault)

	err = thws.AddDirToTemplate(tmp_dir)

	assert.Error(t, err)
	if err == nil {
		t.Fatalf("err = %q, wanted err", err)
	}
}

func TestWriteOutputToFileSuccess(t *testing.T) {
	tmp_dir := t.TempDir()
	file, err := os.CreateTemp(tmp_dir, "tmp")
	if err != nil {
		t.Fatalf("CreateTemp failed")
	}

	fileContent := "test"
	WriteOutputToFile(file.Name(), fileContent)

	output, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatalf("readFile failed with %q", err)
	}

	if string(output) != fileContent {
		t.Fatalf("output = %q, wanted %s", string(output), fileContent)
	}
}

func TestWriteOutputToFileFailed(t *testing.T) {

	writeFile = func(name string, data []byte, perm os.FileMode) error {
		return errors.New("")
	}
	t.Cleanup(setStubsToDefault)

	err = WriteOutputToFile("foo", "bar")

	assert.Error(t, err)
	if err == nil {
		t.Fatalf("err = %q, wanted err == \"\"", err)
	}
}

func TestTemplateChartSuccess(t *testing.T) {
	_thws := thws
	_thws.TmpHelmDir = t.TempDir()
	tmp_dir := t.TempDir()

	t.Cleanup(setStubsToDefault)

	file, err := os.CreateTemp(tmp_dir, "tmp*.yaml")
	if err != nil {
		t.Fatalf("CreateTemp failed")
	}

	valueFile, err := os.CreateTemp(tmp_dir, "value*.yaml")
	if err != nil {
		t.Fatalf("CreateTemp failed")
	}

	fileContent := "test: test"
	WriteOutputToFile(file.Name(), fileContent)

	err = _thws.CreateHelmChart()
	if err != nil {
		t.Fatal(err)
	}

	err = _thws.AddFileToTemplate(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	helmOutput, err := _thws.TemplateChart(valueFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, helmOutput, fileContent)
}

func TestTemplateChartFailure(t *testing.T) {
	_thws := thws
	_thws.TmpHelmDir = t.TempDir()

	cmdRun = func(c *exec.Cmd) error {
		return errors.New("")
	}

	_, err = _thws.TemplateChart()

	assert.Error(t, err)
}

func TestCheckFolderExistsSuccess(t *testing.T) {
	dir := t.TempDir()
	expectedBool := true

	result := CheckFolderExists(dir)

	assert.Equal(t, expectedBool, result)
}

func TestCheckFolderExistsfailure(t *testing.T) {
	tests := []struct {
		expectedBool bool
		actualBool   bool
	}{
		{expectedBool: true, actualBool: CheckFolderExists(t.TempDir())},
		{expectedBool: false, actualBool: CheckFolderExists(makeTestFiles(t, "tmp", t.TempDir()).Name())},
		{expectedBool: false, actualBool: CheckFolderExists("")},
	}

	for _, test := range tests {
		expectedBool := test.expectedBool
		actualBool := test.actualBool

		assert.Equal(t, expectedBool, actualBool)
	}
}

func TestCleanUpSuccess(t *testing.T) {
	tmpDir := t.TempDir()

	assert.DirExists(t, tmpDir)

	CleanUp(tmpDir)

	assert.NoDirExists(t, tmpDir)
}

func TestCleanUpFailure(t *testing.T) {
	removeAll = func(path string) error {
		return errors.New("")
	}
	tmpDir := t.TempDir()

	assert.DirExists(t, tmpDir)

	err := CleanUp(tmpDir)

	assert.Error(t, err)
}

func makeTestFiles(t *testing.T, filename, tmp_dir string) *os.File {
	file, err := os.CreateTemp(tmp_dir, filename)
	if err != nil {
		t.Fatal(err)
	}
	return file
}

func setStubsToDefault() {
	osStat = os.Stat
	glob = filepath.Glob
	cmdRun = (*exec.Cmd).Run
	removeAll = os.RemoveAll
	remove = os.Remove
	mkdir = os.Mkdir
	create = os.Create
	writeFile = os.WriteFile
}
