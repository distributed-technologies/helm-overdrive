package src

import (
	"os/exec"
	"strings"
)

type HelmOverDrive struct {
	Additional_resources_folder string
	Applicaiton_folder          string
	Application_name            string
	Chart_name                  string
	Chart_version               string
	Base_folder                 string
	Env_folder                  string
	Environment                 string
	Global_file                 string
	Helm_repo                   string
	Values_file                 string
	Root_path                   string
}

func (h *HelmOverDrive) GetBaseGlobalFile() string {
	returnString := strings.Join([]string{h.Root_path, h.Base_folder, h.Global_file}, "/")
	return returnString
}

func (h *HelmOverDrive) GetEnvGlobalFile() string {
	returnString := strings.Join([]string{h.Root_path, h.Env_folder, h.Environment, h.Global_file}, "/")
	return returnString
}

func (h *HelmOverDrive) GetBaseApplicationValuesFile() string {
	returnString := strings.Join([]string{h.Root_path, h.Base_folder, h.Applicaiton_folder, h.Application_name, h.Values_file}, "/")
	return returnString
}

func (h *HelmOverDrive) GetEnvApplicationValuesFile() string {
	returnString := strings.Join([]string{h.Root_path, h.Env_folder, h.Environment, h.Applicaiton_folder, h.Application_name, h.Values_file}, "/")
	return returnString
}

func (h *HelmOverDrive) GetBaseApplicationAdditionalResourcesFolder() string {
	returnString := strings.Join([]string{h.Root_path, h.Base_folder, h.Applicaiton_folder, h.Application_name, h.Additional_resources_folder}, "/")
	return returnString
}

func (h *HelmOverDrive) GetEnvApplicationAdditionalResourcesFolder() string {
	returnString := strings.Join([]string{h.Root_path, h.Env_folder, h.Environment, h.Applicaiton_folder, h.Application_name, h.Additional_resources_folder}, "/")
	return returnString
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
