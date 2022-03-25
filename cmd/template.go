/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/distributed-technologies/helm-overdrive/pkg/template"
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

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Merges base and env values and templates a chart using the merged values",
	Long:  tempalteDesc,
	Run: func(cmd *cobra.Command, args []string) {

		additional_resources_folder := viper.GetString("ADDITIONAL_RESOURCES")
		applicaiton_folder := viper.GetString("APPLICATION_FOLDER")
		base_folder := viper.GetString("BASE_FOLDER")
		chart_name := viper.GetString("CHART_NAME")
		chart_version := viper.GetString("CHART_VERSION")
		env_folder := viper.GetString("ENV_FOLDER")
		global_file := viper.GetString("GLOBAL_FILE")
		helm_repo := viper.GetString("HELM_REPO")
		values_file := viper.GetString("VALUES_FILE")

		err := template.Template(additional_resources_folder,
			applicaiton_folder,
			base_folder, chart_name,
			chart_version, env_folder,
			global_file, helm_repo,
			values_file)

		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)

	templateCmd.Flags().String("application_folder", "", "Path to the folder that contains the application, It's assumed that the name is the same in base and env folders")
	templateCmd.Flags().String("app_name", "", "Name of the release")
	templateCmd.Flags().String("base_folder", "", "Path the folder containing the base config")
	templateCmd.Flags().StringP("env_folder", "e", "", "Name of the environment folder you with to deploy")
	templateCmd.Flags().StringP("chart_version", "v", "", "Chart version")
	templateCmd.Flags().StringP("chart_name", "n", "", "Chart")
	templateCmd.Flags().String("global_file", "", "Name of the global files, It's assumed that the name is the same in base and env folders")
	templateCmd.Flags().String("helm_repo", "", "Repo url")
	templateCmd.Flags().String("values_file", "", "Name of the value files in the application folder, It's assumed that the name is the same in base and env folders")
	templateCmd.Flags().String("additional_resources", "", "Path to the folder that contains the additional resources, this has to be located within the <application_folder>, It's assumed that the name is the same in base and env folders")

	viper.BindPFlags(templateCmd.Flags())
}
