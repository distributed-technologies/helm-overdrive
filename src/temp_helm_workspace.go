package src

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var err error

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
	dir, err := os.Stat(path)
	if err == nil {
		if dir.IsDir() {
			err = h.AddFileToTemplate(path)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *TempHelmWorkspace) AddFileToTemplate(filePath string) error {
	files, err := filepath.Glob(filePath)
	if err != nil {
		return err
	}

	for _, path := range files {
		osCmd := exec.Command("cp", "-r", path, h.getTemplatesFolderLocation())
		err = osCmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *TempHelmWorkspace) CreateHelmChart() error {
	osCmd := exec.Command("helm", "create", h.Chart_name)
	osCmd.Dir = h.Tmp_helm_dir

	err = osCmd.Run()
	if err != nil {
		return err
	}

	if err = os.RemoveAll(h.getChartsFolderLocation()); err != nil {
		return err
	}

	if err = os.RemoveAll(h.getTemplatesFolderLocation()); err != nil {
		return err
	}

	if err = os.Remove(h.getValuesFileLocation()); err != nil {
		return err
	}

	if err = os.Mkdir(h.getTemplatesFolderLocation(), 0777); err != nil {
		return err
	}

	if _, err = os.Create(h.getValuesFileLocation()); err != nil {
		return err
	}

	return nil
}

func (h *TempHelmWorkspace) TemplateChart(valueFiles ...string) (string, error) {

	args := []string{"template", h.Release_name, h.getChartFolder()}
	for _, valueFile := range valueFiles {
		args = append(args, "-f", valueFile)
	}

	osCmd := exec.Command("helm", args...)

	var out bytes.Buffer
	osCmd.Stdout = &out

	err := osCmd.Run()

	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func WriteOutputToFile(filename, output string) error {
	if err := os.WriteFile(filename, []byte(output), 0666); err != nil {
		return err
	}
	return nil
}
