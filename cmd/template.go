/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"fmt"
	"helm-overdrive/src"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
)

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
		// Gets the config needed

		// TODO: Split this into smaller bits, remove tempalte cmd, start working on how to get merge_files working with this.

		namespace := viper.GetString("ARGOCD_APP_NAMESPACE")
		releaseName := viper.GetString("ARGOCD_APP_NAME")
		helm_repo := viper.GetString("HELM_REPO")
		chart_name := viper.GetString("CHART_NAME")
		chart_version := viper.GetString("CHART_VERSION")

		valueOpts := &values.Options{}

		out := os.Stdout

		settings := cli.New()
		settings.Debug = isDebug
		settings.SetNamespace(namespace)

		client := src.GenClient(releaseName, helm_repo, settings)
		client.Version = chart_version

		// origin: https://github.com/helm/helm/blob/ee3f270e1eff0d462312635ad91cecd6f1fce620/cmd/helm/install.go#L190
		cp, err := client.ChartPathOptions.LocateChart(chart_name, settings)
		if err != nil {
			panic(err)
		}

		debug("cp: %v\n", cp)

		p := getter.All(settings)
		// TODO: Get this working with ./src/merge_files.test. This would require building the settings and a provider in the root, so that it can be shared between the merger and the templater.
		vals, err := valueOpts.MergeValues(p)
		if err != nil {
			panic(err)
		}

		debug("%s", vals)

		// Check chart dependencies to make sure all are present in /charts
		chartRequested, err := loader.Load(cp)
		if err != nil {
			panic(err)
		}

		if err := checkIfInstallable(chartRequested); err != nil {
			panic(err)
		}

		if chartRequested.Metadata.Deprecated {
			warning("This chart is deprecated")
		}

		if req := chartRequested.Metadata.Dependencies; req != nil {
			// If CheckDependencies returns an error, we have unfulfilled dependencies.
			// As of Helm 2.4.0, this is treated as a stopping condition:
			// https://github.com/helm/helm/issues/2209
			if err := action.CheckDependencies(chartRequested, req); err != nil {
				err = errors.Wrap(err, "An error occurred while checking for chart dependencies. You may need to run `helm dependency build` to fetch missing dependencies")
				if client.DependencyUpdate {
					man := &downloader.Manager{
						Out:              out,
						ChartPath:        cp,
						Keyring:          client.ChartPathOptions.Keyring,
						SkipUpdate:       false,
						Getters:          p,
						RepositoryConfig: settings.RepositoryConfig,
						RepositoryCache:  settings.RepositoryCache,
						Debug:            settings.Debug,
					}
					if err := man.Update(); err != nil {
						panic(err)
					}
					// Reload the chart with the updated Chart.lock file.
					if chartRequested, err = loader.Load(cp); err != nil {
						panic(errors.Wrap(err, "failed reloading chart after repo update"))
					}
				} else {
					panic(err)
				}
			}
		}

		client.Namespace = settings.Namespace()

		// origin: https://github.com/helm/helm/blob/ee3f270e1eff0d462312635ad91cecd6f1fce620/cmd/helm/template.go#L82
		rel, err := client.Run(chartRequested, vals)
		if err != nil {
			panic(err)
		}

		// origin: https://github.com/helm/helm/blob/ee3f270e1eff0d462312635ad91cecd6f1fce620/cmd/helm/template.go#L84
		if err != nil && !settings.Debug {
			if rel != nil {
				panic(fmt.Errorf("%w\n\nUse --debug flag to render out invalid YAML", err))
			}
			panic(err)
		}

		// origin: https://github.com/helm/helm/blob/ee3f270e1eff0d462312635ad91cecd6f1fce620/cmd/helm/template.go#L93
		if rel != nil {
			var manifests bytes.Buffer
			fmt.Fprintln(&manifests, strings.TrimSpace(rel.Manifest))
			if !client.DisableHooks {
				fileWritten := make(map[string]bool)
				for _, m := range rel.Hooks {
					if isTestHook(m) {
						continue
					}
					if client.OutputDir == "" {
						fmt.Fprintf(&manifests, "---\n# Source: %s\n%s\n", m.Path, m.Manifest)
					} else {
						newDir := client.OutputDir
						if client.UseReleaseName {
							newDir = filepath.Join(client.OutputDir, client.ReleaseName)
						}
						err = writeToFile(newDir, m.Path, m.Manifest, fileWritten[m.Path])
						if err != nil {
							panic(err)
						}
						fileWritten[m.Path] = true
					}

				}
			}

			fmt.Fprintf(out, "%s", manifests.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// templateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// templateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// checkIfInstallable validates if a chart can be installed
//
// Application chart type is only installable
func checkIfInstallable(ch *chart.Chart) error {
	switch ch.Metadata.Type {
	case "", "application":
		return nil
	}
	return errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}

func isTestHook(h *release.Hook) bool {
	for _, e := range h.Events {
		if e == release.HookTest {
			return true
		}
	}
	return false
}

// The following functions (writeToFile, createOrOpenFile, and ensureDirectoryForFile)
// are copied from the actions package. This is part of a change to correct a
// bug introduced by #8156. As part of the todo to refactor renderResources
// this duplicate code should be removed. It is added here so that the API
// surface area is as minimally impacted as possible in fixing the issue.
func writeToFile(outputDir string, name string, data string, append bool) error {
	outfileName := strings.Join([]string{outputDir, name}, string(filepath.Separator))

	err := ensureDirectoryForFile(outfileName)
	if err != nil {
		panic(err)
	}

	f, err := createOrOpenFile(outfileName, append)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("---\n# Source: %s\n%s\n", name, data))

	if err != nil {
		panic(err)
	}

	fmt.Printf("wrote %s\n", outfileName)
	return nil
}

func createOrOpenFile(filename string, append bool) (*os.File, error) {
	if append {
		return os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	}
	return os.Create(filename)
}

func ensureDirectoryForFile(file string) error {
	baseDir := path.Dir(file)
	_, err := os.Stat(baseDir)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	return os.MkdirAll(baseDir, 0755)
}
