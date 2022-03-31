/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"helm-overdrive/src"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const tempalteDesc = `
This command makes two temporary Helm charts, adds the '.../<base_folder>/<application>/<values_file>' to one
and '.../<env_folder>/<env>/<application>/<values_file>' to the other

then using Helm Cli tool, it templates both charts using '.../<base_folder>/<global_file>'
and '.../<env_folder>/<env>/<global_file>', saves the output to two files

then it pulls the <chart_name> chart using the Helm CLi and
copies the '.../<base_folder>/<application>/<additional_resources>' folder into the templates folder of the chart
this is also done with the '.../<env_folder>/<env>/<application>/<additional_resources>' folder.

The chart is then templated with the two values files that was generated earlier.
`

const tmpDir string = "tmpHelm"

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Merges base and env values and templates a chart using the merged values",
	Long:  tempalteDesc,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the config needed
		hod := src.HelmOverDrive{
			Additional_resources_folder: viper.GetString("Additional_resources"),
			Applicaiton_folder:          viper.GetString("APPLICATION_FOLDER"),
			Application_name:            viper.GetString("ARGOCD_APP_NAME"),
			Base_folder:                 viper.GetString("BASE_FOLDER"),
			Chart_name:                  viper.GetString("CHART_NAME"),
			Chart_version:               viper.GetString("CHART_VERSION"),
			Env_folder:                  viper.GetString("ENV_FOLDER"),
			Environment:                 viper.GetString("ENVIRONMENT"),
			Global_file:                 viper.GetString("GLOBAL_FILE"),
			Helm_repo:                   viper.GetString("HELM_REPO"),
			Values_file:                 viper.GetString("VALUES_FILE"),
			Root_path:                   viper.GetString("ROOT_PATH"),
		}

		// Makes a slice that contains the names of the base folder and env folder
		var tmpHelms []string = []string{hod.Base_folder, hod.Env_folder}
		var outputFiles []string
		var err error

		/*

			1. create 2 temp helm charts and remove everything in the /templates and /charts folder
			2. move base/<app>/values.yaml into 1. chart and move env/<env>/<app>/values.yaml into 2. chart
			3. template both charts with the base/global.yaml and env/<env>/global.yaml
			4. save both outputs as new values files
			5. pull and unpack the chart
			6. add additional_resources to the templates folder under a uniqe folderName
			7. template the chart with the two new value files, env value file as the last arg.

		*/

		// Create a temporary folder to hold files while working
		_, err = os.Stat(tmpDir)
		if os.IsNotExist(err) {
			debug("folder %s does not exist", tmpDir)

			if err = os.Mkdir(tmpDir, 0777); err != nil {
				panic(err)
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
				panic(errors.New("base_folder and/or env_folder is missing "))
			}

			hw := src.TempHelmWorkspace{
				TmpHelmDir: tmpDir,
				Chart_name: helm_name,
			}

			// Create 2 temp helm charts and remove everything in the /templates and /charts folder also cleans the values.yaml file
			if err = hw.CreateHelmChart(); err != nil {
				debug("%s was not created", helm_name)
				panic(err)
			}

			// Move `.../<base-folder>/<app>/values.yaml` into chart named <base-folder> and
			// Move `.../<env-folder>/<env>/<app>/values.yaml` into chart named <env-folder>
			if hw.Chart_name == hod.Base_folder {
				err = hw.AddFileToTemplate(hod.GetBaseApplicationValuesFile())
			} else if hw.Chart_name == hod.Env_folder {
				err = hw.AddFileToTemplate(hod.GetEnvApplicationValuesFile())
			}
			if err != nil {
				panic(err)
			}

			// Template both charts with  with `.../<base-folder>/<global-file>` and `.../<env-folder>/<env>/<global-file>`
			output, err := hw.TemplateChart(hod.GetBaseGlobalFile(), hod.GetEnvGlobalFile())
			if err != nil {
				panic(err)
			}

			// Save both outputs as new values files
			appValuesFile := strings.Join([]string{hw.TmpHelmDir, hw.Chart_name + ".yaml"}, "/")
			err = src.WriteOutputToFile(appValuesFile, output)
			if err != nil {
				panic(err)
			}

			outputFiles = append(outputFiles, appValuesFile)
		}

		debug("outputFiles: %v\n", outputFiles)

		// Pull and unpack the chart to tmpDir
		hod.GetHelmChart(tmpDir)

		hw := src.TempHelmWorkspace{
			Chart_name: hod.Chart_name,
			TmpHelmDir: tmpDir,
		}

		// Add additional_resources to the templates folder under the <additional_resources> folder name
		if hod.Additional_resources_folder != "" {
			hw.AddDirToTemplate(hod.GetBaseApplicationAdditionalResourcesFolder())
			hw.AddDirToTemplate(hod.GetEnvApplicationAdditionalResourcesFolder())
		} else {
			debug("additional_resources option is not present, skipping...")
		}

		// Template the chart with the two new values files
		output, err := hw.TemplateChart(outputFiles...)
		if err != nil {
			panic(err)
		}

		// Prints the template output to stdout since this is what the ArgoCD plugin needs
		fmt.Fprintf(os.Stdout, "%v", output)

		// Clean up the tmpDir folder
		err = os.RemoveAll(tmpDir)
		if err != nil {
			panic(err)
		} else {
			debug("%s was deleted", tmpDir)
		}
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)
}
