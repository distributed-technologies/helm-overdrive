package pkg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var hod = HelmOverDrive{
	Additional_resources_folder: "additional_resources",
	App_name:                    "hello_world",
	Applicaiton_folder:          "applications/hello-world",
	Base_folder:                 "src",
	Env_folder:                  "test",
	Global_file:                 "global.yaml",
	Values_file:                 "values.yaml",
	Chart_name:                  "hello-world",
	Chart_version:               "0.1.0",
	Helm_repo:                   "https://helm.github.io/examples",
}

func TestCheckRequiredFailed(t *testing.T) {
	// Checking Values_file
	temp_hod := hod
	temp_hod.Values_file = ""

	expectedString := "values_file not defined"
	err := temp_hod.CheckRequired()

	assert.EqualError(t, err, expectedString)

	// Checking Helm_repo
	temp_hod = hod
	temp_hod.Helm_repo = ""

	expectedString = "helm_repo not defined"
	err = temp_hod.CheckRequired()

	assert.EqualError(t, err, expectedString)

	// Checking Global_file
	temp_hod = hod
	temp_hod.Global_file = ""

	expectedString = "global_file not defined"
	err = temp_hod.CheckRequired()

	assert.EqualError(t, err, expectedString)

	// Checking Base_folder
	temp_hod = hod
	temp_hod.Base_folder = ""

	expectedString = "base_folder not defined"
	err = temp_hod.CheckRequired()

	assert.EqualError(t, err, expectedString)

	// Checking Chart_version
	temp_hod = hod
	temp_hod.Chart_version = ""

	expectedString = "chart_version not defined"
	err = temp_hod.CheckRequired()

	assert.EqualError(t, err, expectedString)

	// Checking Chart_name
	temp_hod = hod
	temp_hod.Chart_name = ""

	expectedString = "chart_name not defined"
	err = temp_hod.CheckRequired()

	assert.EqualError(t, err, expectedString)

	// Checking App_name
	temp_hod = hod
	temp_hod.App_name = ""

	expectedString = "app_name not defined"
	err = temp_hod.CheckRequired()

	assert.EqualError(t, err, expectedString)

	// Checking Applicaiton_folder
	temp_hod = hod
	temp_hod.Applicaiton_folder = ""

	expectedString = "application_folder not defined"
	err = temp_hod.CheckRequired()

	assert.EqualError(t, err, expectedString)
}

func TestCheckRequiredSuccess(t *testing.T) {

	err := hod.CheckRequired()

	if err != nil {
		t.Fatalf("CheckRequired() = %t, wanted nil", err)
	}
}

func TestGetOverdrivePaths(t *testing.T) {
	thws.Tmp_helm_dir = t.TempDir()

	tests := []struct {
		expectedString string
		actualString   string
	}{
		{expectedString: fmt.Sprintf("%s/%s", hod.Base_folder, hod.Global_file), actualString: hod.GetBaseGlobalFile()},
		{expectedString: fmt.Sprintf("%s/%s", hod.Env_folder, hod.Global_file), actualString: hod.GetEnvGlobalFile()},
		{expectedString: fmt.Sprintf("%s/%s/%s", hod.Base_folder, hod.Applicaiton_folder, hod.Values_file), actualString: hod.GetBaseApplicationValuesFile()},
		{expectedString: fmt.Sprintf("%s/%s/%s", hod.Env_folder, hod.Applicaiton_folder, hod.Values_file), actualString: hod.GetEnvApplicationValuesFile()},
		{expectedString: fmt.Sprintf("%s/%s/%s", hod.Base_folder, hod.Applicaiton_folder, hod.Additional_resources_folder), actualString: hod.GetBaseApplicationAdditionalResourcesFolder()},
		{expectedString: fmt.Sprintf("%s/%s/%s", hod.Env_folder, hod.Applicaiton_folder, hod.Additional_resources_folder), actualString: hod.GetEnvApplicationAdditionalResourcesFolder()},
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

	temp_hod := hod
	temp_hod.Env_folder = ""

	expected := false
	result := temp_hod.HasEnvironment()

	if result != expected {
		t.Fatalf("HasEnvironment() = %v, wanted %v, error", result, expected)
	}

}

func TestGetHelmChartSuccess(t *testing.T) {
	outDir := t.TempDir()

	err := hod.GetHelmChart(outDir)

	assert.NoError(t, err)
	assert.DirExists(t, outDir+"/"+hod.Chart_name)

}

func TestGetHelmChartFailed(t *testing.T) {
	outDir := t.TempDir()

	temp_hod := hod
	temp_hod.Chart_name = "foo-bar"

	err := temp_hod.GetHelmChart(outDir)

	assert.Error(t, err)
	assert.NoDirExists(t, outDir+"/"+temp_hod.Chart_name)

}
