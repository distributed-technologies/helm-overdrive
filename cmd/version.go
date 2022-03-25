package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Helm-overdrive",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Helm-overdrive v0.0.1")
	},
}

// Adds suggestion to this command if release is given as cmd
func init() {
	versionCmd.SuggestFor = append(versionCmd.SuggestFor, "release")
}
