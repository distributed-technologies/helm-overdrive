package template

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/distributed-technologies/helm-overdrive/pkg/logging"
)

var (
	err       error
	osStat    = os.Stat
	glob      = filepath.Glob
	cmdRun    = (*exec.Cmd).Run
	removeAll = os.RemoveAll
	remove    = os.Remove
	mkdir     = os.Mkdir
	create    = os.Create
	writeFile = os.WriteFile
)

type TempHelmWorkspace struct {
	chartName   string
	TmpHelmDir  string
	ReleaseName string
}

func (h *TempHelmWorkspace) getChartFolder() string {
	return fmt.Sprintf("%s/%s", h.TmpHelmDir, h.chartName)
}

func (h *TempHelmWorkspace) getChartsFolderLocation() string {
	return fmt.Sprintf("%s/%s/%s", h.TmpHelmDir, h.chartName, "charts")
}

func (h *TempHelmWorkspace) getTemplatesFolderLocation() string {
	return fmt.Sprintf("%s/%s/%s", h.TmpHelmDir, h.chartName, "templates")
}

func (h *TempHelmWorkspace) getValuesFileLocation() string {
	return fmt.Sprintf("%s/%s/%s", h.TmpHelmDir, h.chartName, "values.yaml")
}

func (h *TempHelmWorkspace) AddDirToTemplate(path string) error {
	dir, err := osStat(path)
	if err != nil {
		return err
	}
	if dir.IsDir() {
		return h.AddFileToTemplate(path)
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

		if err := cmdRun(osCmd); err != nil {
			return err
		}
	}

	return nil
}

func (h *TempHelmWorkspace) CreateHelmChart() error {
	osCmd := exec.Command("helm", "create", h.chartName)
	osCmd.Dir = h.TmpHelmDir

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

	args := []string{"template", h.ReleaseName, h.getChartFolder(), "--include-crds"}
	for _, valueFile := range valueFiles {
		args = append(args, "-f", valueFile)
	}

	osCmd := exec.Command("helm", args...)

	var out bytes.Buffer
	osCmd.Stdout = &out

	if err := cmdRun(osCmd); err != nil {
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

func CheckFolderExists(pathToFolder string) bool {
	logging.Debug("pathToFolder: %v", pathToFolder)

	fi, err := os.Stat(pathToFolder)
	if fi == nil {
		return false
	}

	logging.Debug("fi: %v", fi)
	logging.Debug("!fi.IsDir(): %v", !fi.IsDir())

	if !fi.IsDir() || os.IsNotExist(err) {
		// Return false if the folder does not exist
		return false
	}

	return true
}

// Removes a folder and all subfolders
func CleanUp(path string) error {

	err = removeAll(path)
	if err != nil {
		return err
	}
	logging.Debug("%s was deleted", path)

	return nil
}
