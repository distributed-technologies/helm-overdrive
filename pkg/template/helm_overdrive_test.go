package template

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var hod = HelmOverDrive{
	additionalResourcesFolder: "additional_resources",
	AppName:                   "hello_world",
	applicaitonFolder:         "applications/hello-world",
	baseFolder:                "src",
	envFolder:                 "test",
	globalFile:                "global.yaml",
	valuesFile:                "values.yaml",
	chartName:                 "hello-world",
	chartVersion:              "0.1.0",
	helmRepo:                  "https://helm.github.io/examples",
}

func TestCheckRequiredFailed(t *testing.T) {
	// Checking valuesFile
	tempHod := hod
	tempHod.valuesFile = ""

	expectedString := "valuesFile not defined"
	err := tempHod.CheckRequired()

	assert.EqualError(t, err, expectedString)

	// Checking helmRepo
	tempHod = hod
	tempHod.helmRepo = ""

	expectedString = "helmRepo not defined"
	err = tempHod.CheckRequired()

	assert.EqualError(t, err, expectedString)

	// Checking globalFile
	tempHod = hod
	tempHod.globalFile = ""

	expectedString = "globalFile not defined"
	err = tempHod.CheckRequired()

	assert.EqualError(t, err, expectedString)

	// Checking baseFolder
	tempHod = hod
	tempHod.baseFolder = ""

	expectedString = "baseFolder not defined"
	err = tempHod.CheckRequired()

	assert.EqualError(t, err, expectedString)

	// Checking chartVersion
	tempHod = hod
	tempHod.chartVersion = ""

	expectedString = "chartVersion not defined"
	err = tempHod.CheckRequired()

	assert.EqualError(t, err, expectedString)

	// Checking chartName
	tempHod = hod
	tempHod.chartName = ""

	expectedString = "chartName not defined"
	err = tempHod.CheckRequired()

	assert.EqualError(t, err, expectedString)

	// Checking AppName
	tempHod = hod
	tempHod.AppName = ""

	expectedString = "AppName not defined"
	err = tempHod.CheckRequired()

	assert.EqualError(t, err, expectedString)

	// Checking applicaitonFolder
	tempHod = hod
	tempHod.applicaitonFolder = ""

	expectedString = "application_folder not defined"
	err = tempHod.CheckRequired()

	assert.EqualError(t, err, expectedString)
}

func TestCheckRequiredSuccess(t *testing.T) {

	err := hod.CheckRequired()

	if err != nil {
		t.Fatalf("CheckRequired() = %t, wanted nil", err)
	}
}

func TestGetOverdrivePaths(t *testing.T) {
	thws.TmpHelmDir = t.TempDir()

	tests := []struct {
		expectedString string
		actualString   string
	}{
		{expectedString: fmt.Sprintf("%s/%s", hod.baseFolder, hod.globalFile), actualString: hod.GetBaseGlobalFile()},
		{expectedString: fmt.Sprintf("%s/%s", hod.envFolder, hod.globalFile), actualString: hod.GetEnvGlobalFile()},
		{expectedString: fmt.Sprintf("%s/%s/%s", hod.baseFolder, hod.applicaitonFolder, hod.valuesFile), actualString: hod.GetBaseApplicationValuesFile()},
		{expectedString: fmt.Sprintf("%s/%s/%s", hod.envFolder, hod.applicaitonFolder, hod.valuesFile), actualString: hod.GetEnvApplicationValuesFile()},
		{expectedString: fmt.Sprintf("%s/%s/%s", hod.baseFolder, hod.applicaitonFolder, hod.additionalResourcesFolder), actualString: hod.GetBaseApplicationAdditionalResourcesFolder()},
		{expectedString: fmt.Sprintf("%s/%s/%s", hod.envFolder, hod.applicaitonFolder, hod.additionalResourcesFolder), actualString: hod.GetEnvApplicationAdditionalResourcesFolder()},
	}

	for _, test := range tests {
		expectedString := test.expectedString
		actualString := test.actualString

		if actualString != expectedString {
			t.Fatalf("expected = %q, wanted %s, error", actualString, expectedString)
		}
	}
}

func TestHasEnvironmentSuccess(t *testing.T) {

	expected := true
	result := hod.HasEnvironment()

	if result != expected {
		t.Fatalf("HasEnvironment() = %v, wanted %v, error", result, expected)
	}
}

func TestHasEnvironmentFailed(t *testing.T) {

	tempHod := hod
	tempHod.envFolder = ""

	expected := false
	result := tempHod.HasEnvironment()

	if result != expected {
		t.Fatalf("HasEnvironment() = %v, wanted %v, error", result, expected)
	}
}

func TestGetHelmChartSuccess(t *testing.T) {
	outDir := t.TempDir()

	err := hod.GetHelmChart(outDir)

	assert.NoError(t, err)
	assert.DirExists(t, outDir+"/"+hod.chartName)
}

func TestGetHelmChartFailed(t *testing.T) {
	outDir := t.TempDir()

	tempHod := hod
	tempHod.chartName = "foo-bar"

	err := tempHod.GetHelmChart(outDir)

	assert.Error(t, err)
	assert.NoDirExists(t, outDir+"/"+tempHod.chartName)
}
