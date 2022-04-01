package src

import (
	"errors"
	"os/exec"
	"strings"
)

type HelmOverDrive struct {
	Additional_resources_folder string
	App_name                    string
	Applicaiton_folder          string
	Chart_name                  string
	Chart_version               string
	Base_folder                 string
	Env_folder                  string
	Global_file                 string
	Helm_repo                   string
	Values_file                 string
}

func (h *HelmOverDrive) CheckRequired() error {

	if h.Applicaiton_folder == "" {
		return errors.New("application_folder not defined")
	}
	if h.App_name == "" {
		return errors.New("app_name not defined")
	}
	if h.Chart_name == "" {
		return errors.New("chart_name not defined")
	}
	if h.Chart_version == "" {
		return errors.New("chart_version not defined")
	}
	if h.Base_folder == "" {
		return errors.New("base_folder not defined")
	}
	if h.Global_file == "" {
		return errors.New("global_file not defined")
	}
	if h.Helm_repo == "" {
		return errors.New("helm_repo not defined")
	}
	if h.Values_file == "" {
		return errors.New("values_file not defined")
	}

	return nil
}

func (h *HelmOverDrive) GetBaseGlobalFile() string {
	returnString := strings.Join([]string{h.Base_folder, h.Global_file}, "/")
	return returnString
}

func (h *HelmOverDrive) GetEnvGlobalFile() string {
	returnString := strings.Join([]string{h.Env_folder, h.Global_file}, "/")
	return returnString
}

func (h *HelmOverDrive) GetBaseApplicationValuesFile() string {
	returnString := strings.Join([]string{h.Base_folder, h.Applicaiton_folder, h.Values_file}, "/")
	return returnString
}

func (h *HelmOverDrive) GetEnvApplicationValuesFile() string {
	returnString := strings.Join([]string{h.Env_folder, h.Applicaiton_folder, h.Values_file}, "/")
	return returnString
}

func (h *HelmOverDrive) GetBaseApplicationAdditionalResourcesFolder() string {
	returnString := strings.Join([]string{h.Base_folder, h.Applicaiton_folder, h.Additional_resources_folder}, "/")
	return returnString
}

func (h *HelmOverDrive) GetEnvApplicationAdditionalResourcesFolder() string {
	returnString := strings.Join([]string{h.Env_folder, h.Applicaiton_folder, h.Additional_resources_folder}, "/")
	return returnString
}

func (h *HelmOverDrive) HasEnvironment() bool {
	if h.Env_folder == "" {
		return false
	}
	return true
}

// Pulls the helm chart and unpacks it, then returnes to path to it
func (h *HelmOverDrive) GetHelmChart(outDir string) error {
	cmd := exec.Command("helm", "pull", h.Chart_name, "--version", h.Chart_version, "--repo", h.Helm_repo, "--untar", "-d", outDir)

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
