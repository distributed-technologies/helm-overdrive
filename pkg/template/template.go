package template

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/distributed-technologies/helm-overdrive/pkg/logging"
	"github.com/spf13/viper"
)

var tmpDir string

const tmpBaseName string = "base"
const tmpEnvName string = "env"

func Template(additionalResourcesFolder,
	applicaitonFolder,
	baseFolder,
	chartName,
	chartVersion,
	envFolder,
	globalFile,
	helmRepo,
	valuesFile string) error {

	tmpDir, err = ioutil.TempDir("", "helm-*")
	if err != nil {
		return err
	}

	defer CleanUp(tmpDir)

	// Get the config needed
	hod := HelmOverDrive{
		additionalResourcesFolder: additionalResourcesFolder,
		applicaitonFolder:         applicaitonFolder,
		baseFolder:                baseFolder,
		chartName:                 chartName,
		chartVersion:              chartVersion,
		envFolder:                 envFolder,
		globalFile:                globalFile,
		helmRepo:                  helmRepo,
		valuesFile:                valuesFile,
	}

	if name := viper.GetString("AppName"); name != "" {
		logging.Debug("name: %v", name)
		hod.AppName = name
	} else if name := os.Getenv("ARGOCD_AppName"); name != "" {
		logging.Debug("name: %v", name)
		hod.AppName = name
	} else {
		hod.AppName = hod.chartName
	}

	var tmpHelms []string = []string{tmpBaseName}
	var outputFiles []string
	var err error

	logging.Debug("hod: %v\n", hod)
	logging.Debug("hod.HasEnvironment(): %v\n", hod.HasEnvironment())

	// Check if required fields exists
	if err = hod.CheckRequired(); err != nil {
		return err
	}

	// Check if an environment is present
	if hod.HasEnvironment() {
		tmpHelms = append(tmpHelms, tmpEnvName)
	}

	/*
		Loop over slice to create helm charts, add value files to the templates folder,
		run `helm template` with `.../<base-folder>/<global-file>` and
		`.../<env-folder>/<env>/<global-file>` as values files.
		save output as files
	*/
	for _, helmName := range tmpHelms {
		if helmName == "" {
			return (errors.New("baseFolder and/or envFolder is missing"))
		}

		hw := TempHelmWorkspace{
			TmpHelmDir:  tmpDir,
			chartName:   helmName,
			ReleaseName: hod.AppName,
		}

		// Create 2 temp helm charts and remove everything in the /templates and /charts folder also cleans the values.yaml file
		if err = hw.CreateHelmChart(); err != nil {
			return (logging.WrapError("Failed creating helm chart %s \n%w", helmName, err))
		}

		// Move `.../<base-folder>/<app>/values.yaml` into chart named <base-folder> and
		// Move `.../<env-folder>/<env>/<app>/values.yaml` into chart named <env-folder>
		if hw.chartName == tmpBaseName {
			err = hw.AddFileToTemplate(hod.GetBaseApplicationValuesFile())
			if err != nil {
				return (logging.WrapError("Failed coping base app values to template folder \n%w", err))
			}
		} else if hw.chartName == tmpEnvName {
			err = hw.AddFileToTemplate(hod.GetEnvApplicationValuesFile())
			if err != nil {
				return (logging.WrapError("Failed coping env app values to template folder \n%w", err))
			}
		}

		// Template both charts with  with `.../<base-folder>/<global-file>` and `.../<env-folder>/<env>/<global-file>`
		var output string
		if hod.HasEnvironment() {
			output, err = hw.TemplateChart(hod.GetBaseGlobalFile(), hod.GetEnvGlobalFile())
			if err != nil {
				return (logging.WrapError("Failed templating %s \n%w", hw.chartName, err))
			}
		} else {
			output, err = hw.TemplateChart(hod.GetBaseGlobalFile())
			if err != nil {
				return (logging.WrapError("Failed templating %s \n%w", hw.chartName, err))
			}
		}

		// Save both outputs as new values files
		appValuesFile := fmt.Sprintf("%s/%s", hw.TmpHelmDir, hw.chartName+".yaml")

		err = WriteOutputToFile(appValuesFile, output)
		if err != nil {
			return (logging.WrapError("Failed wrting %s to a file \n%w", appValuesFile, err))
		}

		outputFiles = append(outputFiles, appValuesFile)
	}

	logging.Debug("outputFiles: %v\n", outputFiles)

	// Pull and unpack the chart to tmpDir
	err = hod.GetHelmChart(tmpDir)
	if err != nil {
		return (logging.WrapError("Failed pulling %s \n%w", hod.chartName, err))
	}

	hw := TempHelmWorkspace{
		chartName:   hod.chartName,
		TmpHelmDir:  tmpDir,
		ReleaseName: hod.AppName,
	}

	// Add additional_resources to the templates folder under the <additional_resources> folder name
	if hod.additionalResourcesFolder != "" && CheckFolderExists(hod.GetBaseApplicationAdditionalResourcesFolder()) {
		err = hw.AddDirToTemplate(hod.GetBaseApplicationAdditionalResourcesFolder())
		if err != nil {
			return (logging.WrapError("Failed adding %s to templatefolder \n%w", hod.GetBaseApplicationAdditionalResourcesFolder(), err))
		}
		if hod.HasEnvironment() && CheckFolderExists(hod.GetEnvApplicationAdditionalResourcesFolder()) {
			err = hw.AddDirToTemplate(hod.GetEnvApplicationAdditionalResourcesFolder())
			if err != nil {
				return (logging.WrapError("Failed adding %s to templatefolder \n%w", hod.GetEnvApplicationAdditionalResourcesFolder(), err))
			}
		}
	} else {
		logging.Debug("additional_resources option is not present, skipping...")
	}

	// Template the chart with the two new values files
	output, err := hw.TemplateChart(outputFiles...)
	if err != nil {
		return err
	}

	// Prints the template output to stdout since this is what the ArgoCD plugin needs
	fmt.Fprintf(os.Stdout, "%v", output)

	return nil
}
