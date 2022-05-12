package template

import (
	"errors"
	"fmt"
	"os/exec"
)

type HelmOverDrive struct {
	additionalResourcesFolder string
	AppName                   string
	applicaitonFolder         string
	chartName                 string
	chartVersion              string
	baseFolder                string
	envFolder                 string
	globalFile                string
	helmRepo                  string
	valuesFile                string
}

func (h *HelmOverDrive) CheckRequired() error {

	if h.applicaitonFolder == "" {
		return errors.New("application_folder not defined")
	}
	if h.AppName == "" {
		return errors.New("AppName not defined")
	}
	if h.chartName == "" {
		return errors.New("chartName not defined")
	}
	if h.chartVersion == "" {
		return errors.New("chartVersion not defined")
	}
	if h.baseFolder == "" {
		return errors.New("baseFolder not defined")
	}
	if h.globalFile == "" {
		return errors.New("globalFile not defined")
	}
	if h.helmRepo == "" {
		return errors.New("helmRepo not defined")
	}
	if h.valuesFile == "" {
		return errors.New("valuesFile not defined")
	}

	return nil
}

func (h *HelmOverDrive) GetBaseGlobalFile() string {
	return fmt.Sprintf("%s/%s", h.baseFolder, h.globalFile)
}

func (h *HelmOverDrive) GetEnvGlobalFile() string {
	return fmt.Sprintf("%s/%s", h.envFolder, h.globalFile)
}

func (h *HelmOverDrive) GetBaseApplicationValuesFile() string {
	return fmt.Sprintf("%s/%s/%s", h.baseFolder, h.applicaitonFolder, h.valuesFile)
}

func (h *HelmOverDrive) GetEnvApplicationValuesFile() string {
	return fmt.Sprintf("%s/%s/%s", h.envFolder, h.applicaitonFolder, h.valuesFile)
}

func (h *HelmOverDrive) GetBaseApplicationAdditionalResourcesFolder() string {
	return fmt.Sprintf("%s/%s/%s", h.baseFolder, h.applicaitonFolder, h.additionalResourcesFolder)
}

func (h *HelmOverDrive) GetEnvApplicationAdditionalResourcesFolder() string {
	return fmt.Sprintf("%s/%s/%s", h.envFolder, h.applicaitonFolder, h.additionalResourcesFolder)
}

func (h *HelmOverDrive) HasEnvironment() bool {
	return h.envFolder != ""
}

// Pulls the helm chart and unpacks it, then returnes to path to it
func (h *HelmOverDrive) GetHelmChart(outDir string) error {
	cmd := exec.Command("helm", "pull", h.chartName, "--version", h.chartVersion, "--repo", h.helmRepo, "--untar", "-d", outDir)

	return cmd.Run()
}
