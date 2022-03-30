package src

import (
	"strings"
)

type HelmOverDrive struct {
	Applicaiton_folder string
	Application_name   string
	Root_path          string
	Base_folder        string
	Env_folder         string
	Environment        string
	Global_file        string
	Values_file        string
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
