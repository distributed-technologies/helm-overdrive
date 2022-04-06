package pkg

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var err error
var osStat = os.Stat
var glob = filepath.Glob
var cmdRun = (*exec.Cmd).Run
var removeAll = os.RemoveAll
var remove = os.Remove
var mkdir = os.Mkdir
var create = os.Create
var writeFile = os.WriteFile

type TempHelmWorkspace struct {
	Chart_name   string
	Tmp_helm_dir string
	Release_name string
}

func (h *TempHelmWorkspace) getChartFolder() string {
	return strings.Join([]string{h.Tmp_helm_dir, h.Chart_name}, "/")
}

func (h *TempHelmWorkspace) getChartsFolderLocation() string {
	return strings.Join([]string{h.Tmp_helm_dir, h.Chart_name, "charts"}, "/")
}

func (h *TempHelmWorkspace) getTemplatesFolderLocation() string {
	return strings.Join([]string{h.Tmp_helm_dir, h.Chart_name, "templates"}, "/")
}

func (h *TempHelmWorkspace) getValuesFileLocation() string {
	return strings.Join([]string{h.Tmp_helm_dir, h.Chart_name, "values.yaml"}, "/")
}

func (h *TempHelmWorkspace) AddDirToTemplate(path string) error {
	dir, err := osStat(path)
	if err == nil {
		if dir.IsDir() {
			err = h.AddFileToTemplate(path)
			if err != nil {
				return err
			}
		}
	} else {
		return err
	}
	return nil
}

func (h *TempHelmWorkspace) AddFileToTemplate(filePath string) error {
	files, err := glob(filePath)
	if err != nil {
		return err
	}

	for _, path := range files {
		osCmd := exec.Command("cp", "-r", path, h.getTemplatesFolderLocation())
		err = cmdRun(osCmd)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *TempHelmWorkspace) CreateHelmChart() error {
	osCmd := exec.Command("helm", "create", h.Chart_name)
	osCmd.Dir = h.Tmp_helm_dir

	err = cmdRun(osCmd)
	if err != nil {
		return err
	}

	if err = removeAll(h.getChartsFolderLocation()); err != nil {
		return err
	}

	if err = removeAll(h.getTemplatesFolderLocation()); err != nil {
		return err
	}

	if err = remove(h.getValuesFileLocation()); err != nil {
		return err
	}

	if err = mkdir(h.getTemplatesFolderLocation(), 0777); err != nil {
		return err
	}

	if _, err = create(h.getValuesFileLocation()); err != nil {
		return err
	}

	return nil
}

func (h *TempHelmWorkspace) TemplateChart(valueFiles ...string) (string, error) {

	args := []string{"template", h.Release_name, h.getChartFolder()}
	for _, valueFile := range valueFiles {
		args = append(args, "--include-crds", "-f", valueFile)
	}

	osCmd := exec.Command("helm", args...)

	var out bytes.Buffer
	osCmd.Stdout = &out

	err := cmdRun(osCmd)

	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func WriteOutputToFile(filename, output string) error {
	if err := writeFile(filename, []byte(output), 0666); err != nil {
		return err
	}
	return nil
}
