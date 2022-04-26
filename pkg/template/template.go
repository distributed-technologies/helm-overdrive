package template

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/distributed-technologies/helm-overdrive/pkg/logging"
	"github.com/spf13/viper"
)

const tmpDir string = "tmpHelm"
const tmpBaseName string = "base"
const tmpEnvName string = "env"

func Template(additional_resources_folder,
	applicaiton_folder,
	base_folder,
	chart_name,
	chart_version,
	env_folder,
	global_file,
	helm_repo,
	values_file string) error {

	// Get the config needed
	hod := HelmOverDrive{
		Additional_resources_folder: additional_resources_folder,
		Applicaiton_folder:          applicaiton_folder,
		Base_folder:                 base_folder,
		Chart_name:                  chart_name,
		Chart_version:               chart_version,
		Env_folder:                  env_folder,
		Global_file:                 global_file,
		Helm_repo:                   helm_repo,
		Values_file:                 values_file,
	}

	if name := viper.GetString("APP_NAME"); name != "" {
		logging.Debug("name: %v", name)
		hod.App_name = name
	} else if name := os.Getenv("ARGOCD_APP_NAME"); name != "" {
		logging.Debug("name: %v", name)
		hod.App_name = name
	} else {
		hod.App_name = hod.Chart_name
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

	// Create a temporary folder to hold files while working
	_, err = os.Stat(tmpDir)
	if os.IsNotExist(err) {
		logging.Debug("%s does not exist", tmpDir)

		if err = os.Mkdir(tmpDir, 0777); err != nil {
			return (logging.WrapError("Failed in creating tmpDir \n%w", err))
		}
	}

	/*
		Loop over slice to create helm charts, add value files to the templates folder,
		run `helm template` with `.../<base-folder>/<global-file>` and
		`.../<env-folder>/<env>/<global-file>` as values files.
		save output as files
	*/
	for _, helm_name := range tmpHelms {
		if helm_name == "" {
			return (errors.New("Base_folder and/or env_folder is missing"))
		}

		hw := TempHelmWorkspace{
			Tmp_helm_dir: tmpDir,
			Chart_name:   helm_name,
			Release_name: hod.App_name,
		}

		// Create 2 temp helm charts and remove everything in the /templates and /charts folder also cleans the values.yaml file
		if err = hw.CreateHelmChart(); err != nil {
			return (logging.WrapError("Failed creating helm chart %s \n%w", helm_name, err))
		}

		// Move `.../<base-folder>/<app>/values.yaml` into chart named <base-folder> and
		// Move `.../<env-folder>/<env>/<app>/values.yaml` into chart named <env-folder>
		if hw.Chart_name == tmpBaseName {
			err = hw.AddFileToTemplate(hod.GetBaseApplicationValuesFile())
			if err != nil {
				return (logging.WrapError("Failed coping base app values to template folder \n%w", err))
			}
		} else if hw.Chart_name == tmpEnvName {
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
				return (logging.WrapError("Failed templating %s \n%w", hw.Chart_name, err))
			}
		} else {
			output, err = hw.TemplateChart(hod.GetBaseGlobalFile())
			if err != nil {
				return (logging.WrapError("Failed templating %s \n%w", hw.Chart_name, err))
			}
		}

		// Save both outputs as new values files
		appValuesFile := strings.Join([]string{hw.Tmp_helm_dir, hw.Chart_name + ".yaml"}, "/")
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
		return (logging.WrapError("Failed pulling %s \n%w", hod.Chart_name, err))
	}

	hw := TempHelmWorkspace{
		Chart_name:   hod.Chart_name,
		Tmp_helm_dir: tmpDir,
		Release_name: hod.App_name,
	}

	// Add additional_resources to the templates folder under the <additional_resources> folder name
	if hod.Additional_resources_folder != "" && hod.CheckFolderExists(hod.GetBaseApplicationAdditionalResourcesFolder()) {
		err = hw.AddDirToTemplate(hod.GetBaseApplicationAdditionalResourcesFolder())
		if err != nil {
			return (logging.WrapError("Failed adding %s to templatefolder \n%w", hod.GetBaseApplicationAdditionalResourcesFolder(), err))
		}
		if hod.HasEnvironment() && hod.CheckFolderExists(hod.GetEnvApplicationAdditionalResourcesFolder()) {
			err = hw.AddDirToTemplate(hod.GetEnvApplicationAdditionalResourcesFolder())
			if err != nil {
				return (logging.WrapError("Failed adding %s to templatefolder \n%w", hod.GetBaseApplicationAdditionalResourcesFolder(), err))
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

	// Clean up the tmpDir folder
	err = os.RemoveAll(tmpDir)
	if err != nil {
		return err
	} else {
		logging.Debug("%s was deleted", tmpDir)
	}

	return nil
}
