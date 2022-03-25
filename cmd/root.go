package cmd

import (
	"fmt"
	"helm-overdrive/src"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "helm-overdrive",
	Short: "Templating multiple environments together",
	Long: `Helm-overdrive is a tool that allows the templating og multiple yaml resources on top of each other.	`,

	Run: func(cmd *cobra.Command, args []string) {
		src.MergeYaml(viper.GetString("BASE_GLOBAL"), viper.GetString("ENV_GLOBAL"))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Persistent Flags will be available to this command and all subcommands to this
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.helm-overdrive.yaml)")

	// rootCmd.PersistentFlags().String("helm_repo", "", "Url to the helm repository containing the chart")
	// rootCmd.PersistentFlags().String("chart_name", "", "Name of the chart to template")
	// rootCmd.PersistentFlags().String("chart_version", "", "Version of the chart to use")
	// rootCmd.PersistentFlags().String("base_values", "", "Path to base global file")
	// rootCmd.PersistentFlags().String("base_global", "", "Path to base values file")

	viper.BindPFlags(rootCmd.Flags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("HO") // Standing for 'helm-overdrive'
	viper.AutomaticEnv()     // read in environment variables that match

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else if viper.GetString("config") != "" {
		// Use config file from environment 'HO_CONFIG'
		viper.SetConfigFile(viper.GetString("config"))
	} else {
		// Look in these paths for a config file
		viper.AddConfigPath("/etc/lib/helm-overdrive") // Checks running dir
		viper.SetConfigType("yaml")
		viper.SetConfigName("helm-overdrive")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
