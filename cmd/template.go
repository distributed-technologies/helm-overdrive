/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"helm-overdrive/src"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const tmpDir string = "tmpHelm"

var err error

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		var tmpHelms []string = []string{hod.Base_folder, hod.Env_folder}
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

		// 1. create 2 temp helm charts and remove everything in the /templates and /charts folder
		_, err = os.Stat(tmpDir)
		if os.IsNotExist(err) {
			debug("folder %s does not exist", tmpDir)

			if err = os.Mkdir(tmpDir, 0777); err != nil {
				panic(err)
			}
		}

		var outputFiles []string
		for _, helm_name := range tmpHelms {
			hw := src.TempHelmWorkspace{
				TmpHelmDir: tmpDir,
				Chart_name: helm_name,
			}

			// 1. create 2 temp helm charts and remove everything in the /templates and /charts folder
			if err = hw.CreateHelmChart(); err != nil {
				debug("%s was not created", helm_name)
				panic(err)
			}

			// 2. move base/<app>/values.yaml into 1. chart and move env/<env>/<app>/values.yaml into 2. chart
			if hw.Chart_name == hod.Base_folder {
				err = hw.AddFileToTemplate(hod.GetBaseApplicationValuesFile())
			} else if hw.Chart_name == hod.Env_folder {
				err = hw.AddFileToTemplate(hod.GetEnvApplicationValuesFile())
			}
			if err != nil {
				panic(err)
			}

			// 3. template both charts with the base/global.yaml and env/<env>/global.yaml
			output, err := hw.TemplateChart(hod.GetBaseGlobalFile(), hod.GetEnvGlobalFile())
			if err != nil {
				panic(err)
			}

			// 4. save both outputs as new values files
			appValuesFile := strings.Join([]string{hw.TmpHelmDir, hw.Chart_name + ".yaml"}, "/")
			err = src.WriteOutputToFile(appValuesFile, output)
			if err != nil {
				panic(err)
			}
			outputFiles = append(outputFiles, appValuesFile)
		}

		debug("outputFiles: %v\n", outputFiles)

		// 5. pull and unpack the chart
		hod.GetHelmChart(tmpDir)

		hw := src.TempHelmWorkspace{
			Chart_name: hod.Chart_name,
			TmpHelmDir: tmpDir,
		}

		// 6. add additional_resources to the templates folder under a uniqe folderName
		if hod.Additional_resources_folder != "" {
			addDirToTemplate(hw, hod.GetBaseApplicationAdditionalResourcesFolder())
			addDirToTemplate(hw, hod.GetEnvApplicationAdditionalResourcesFolder())
		} else {
			debug("additional_resources option is not present, skipping...")
		}

		// 7. template the chart with the two new value files, env value file as the last arg.
		output, err := hw.TemplateChart(outputFiles...)
		if err != nil {
			panic(err)
		}

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

func addDirToTemplate(hw src.TempHelmWorkspace, path string) {
	dir, err := os.Stat(path)
	if err == nil {
		if dir.IsDir() {
			err = hw.AddFileToTemplate(path)
			if err != nil {
				panic(err)
			}
		}
	} else {
		debug("additional_resources folder not found...")
	}
}
