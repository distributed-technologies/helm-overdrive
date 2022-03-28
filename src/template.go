package src

import (
	"log"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
)

func GenClient(releaseName string, repoUrl string, settings *cli.EnvSettings) *action.Install {
	actionConfig := new(action.Configuration)
	// You can pass an empty string instead of settings.Namespace() to list
	// all namespaces
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(),
		os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	client := action.NewInstall(actionConfig)
	client.ReleaseName = releaseName
	client.RepoURL = repoUrl
	client.SkipCRDs = false
	client.DryRun = true
	client.ClientOnly = true

	return client
}
