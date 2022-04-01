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

const tempalteDesc = ``
const tempalteDescDepricated = `
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
const tmpBaseName string = "base"
const tmpEnvName string = "env"

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Merges base and env values and templates a chart using the merged values",
	Long:  tempalteDesc,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the config needed
		hod := src.HelmOverDrive{
			Additional_resources_folder: viper.GetString("ADDITIONAL_RESOURCES"),
			Applicaiton_folder:          viper.GetString("APPLICATION_FOLDER"),
			Base_folder:                 viper.GetString("BASE_FOLDER"),
			Chart_name:                  viper.GetString("CHART_NAME"),
			Chart_version:               viper.GetString("CHART_VERSION"),
			Env_folder:                  viper.GetString("ENV_FOLDER"),
			Global_file:                 viper.GetString("GLOBAL_FILE"),
			Helm_repo:                   viper.GetString("HELM_REPO"),
			Values_file:                 viper.GetString("VALUES_FILE"),
		}

		var tmpHelms []string = []string{tmpBaseName}
		var outputFiles []string
		var err error

		debug("hod: %v\n", hod)
		debug("hod.HasEnvironment(): %v\n", hod.HasEnvironment())

		// Check if required fields exists
		if err = hod.CheckRequired(); err != nil {
			panic(err)
		}

		// Check if an environment is present
		if hod.HasEnvironment() {
			tmpHelms = append(tmpHelms, tmpEnvName)
		}

		// Create a temporary folder to hold files while working
		_, err = os.Stat(tmpDir)
		if os.IsNotExist(err) {
			debug("%s does not exist", tmpDir)

			if err = os.Mkdir(tmpDir, 0777); err != nil {
				panic(wrapError("Failed in creating tmpDir \n%w", err))
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
				panic(errors.New("Base_folder and/or env_folder is missing"))
			}

			hw := src.TempHelmWorkspace{
				TmpHelmDir: tmpDir,
				Chart_name: helm_name,
			}

			// Create 2 temp helm charts and remove everything in the /templates and /charts folder also cleans the values.yaml file
			if err = hw.CreateHelmChart(); err != nil {
				panic(wrapError("Failed creating helm chart %s \n%w", helm_name, err))
			}

			// Move `.../<base-folder>/<app>/values.yaml` into chart named <base-folder> and
			// Move `.../<env-folder>/<env>/<app>/values.yaml` into chart named <env-folder>
			if hw.Chart_name == tmpBaseName {
				err = hw.AddFileToTemplate(hod.GetBaseApplicationValuesFile())
				if err != nil {
					panic(wrapError("Failed coping base app values to template folder \n%w", err))
				}
			} else if hw.Chart_name == tmpEnvName {
				err = hw.AddFileToTemplate(hod.GetEnvApplicationValuesFile())
				if err != nil {
					panic(wrapError("Failed coping env app values to template folder \n%w", err))
				}
			}

			// Template both charts with  with `.../<base-folder>/<global-file>` and `.../<env-folder>/<env>/<global-file>`
			var output string
			if hod.HasEnvironment() {
				output, err = hw.TemplateChart(hod.GetBaseGlobalFile(), hod.GetEnvGlobalFile())
				if err != nil {
					panic(wrapError("Failed templating %s \n%w", hw.Chart_name, err))
				}
			} else {
				output, err = hw.TemplateChart(hod.GetBaseGlobalFile())
				if err != nil {
					panic(wrapError("Failed templating %s \n%w", hw.Chart_name, err))
				}
			}

			// Save both outputs as new values files
			appValuesFile := strings.Join([]string{hw.TmpHelmDir, hw.Chart_name + ".yaml"}, "/")
			err = src.WriteOutputToFile(appValuesFile, output)
			if err != nil {
				panic(wrapError("Failed wrting %s to a file \n%w", appValuesFile, err))
			}

			outputFiles = append(outputFiles, appValuesFile)
		}

		debug("outputFiles: %v\n", outputFiles)

		// Pull and unpack the chart to tmpDir
		err = hod.GetHelmChart(tmpDir)
		if err != nil {
			panic(wrapError("Failed pulling %s \n%w", hod.Chart_name, err))
		}

		hw := src.TempHelmWorkspace{
			Chart_name: hod.Chart_name,
			TmpHelmDir: tmpDir,
		}

		// Add additional_resources to the templates folder under the <additional_resources> folder name
		debug("hod.Additional_resources_folder: %v\n", hod.Additional_resources_folder)
		if hod.Additional_resources_folder != "" {
			err = hw.AddDirToTemplate(hod.GetBaseApplicationAdditionalResourcesFolder())
			if err != nil {
				panic(wrapError("Failed adding %s to templatefolder \n%w", hod.GetBaseApplicationAdditionalResourcesFolder(), err))
			}
			if hod.HasEnvironment() {
				err = hw.AddDirToTemplate(hod.GetEnvApplicationAdditionalResourcesFolder())
				if err != nil {
					panic(wrapError("Failed adding %s to templatefolder \n%w", hod.GetBaseApplicationAdditionalResourcesFolder(), err))
				}
			}
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

	templateCmd.Flags().String("application_folder", "", "Path to the folder that contains the application, It's assumed that the name is the same in base and env folders")
	templateCmd.Flags().String("base_folder", "", "Path the folder containing the base config")
	templateCmd.Flags().StringP("env_folder", "e", "", "Name of the environment folder you with to deploy")
	templateCmd.Flags().StringP("chart_version", "v", "", "Chart version")
	templateCmd.Flags().StringP("chart_name", "n", "", "Chart")
	templateCmd.Flags().String("global_file", "", "Name of the global files, It's assumed that the name is the same in base and env folders")
	templateCmd.Flags().String("helm_repo", "", "Repo url")
	templateCmd.Flags().String("values_file", "", "Name of the value files in the application folder, It's assumed that the name is the same in base and env folders")
	templateCmd.Flags().String("additional_resources", "", "Name of the folder that contains the additional resources, this has to be located within the <application_folder>, It's assumed that the name is the same in base and env folders")

	viper.BindPFlags(templateCmd.Flags())
}
