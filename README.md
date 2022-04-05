# helm-overdrive

## File structure
The file structure currently needed to use this app
``` tree
📦config-example
 ┣ 📂base
 ┃ ┣ 📂applications
 ┃ ┃ ┗ 📂hello-world
 ┃ ┃   ┣ 📂additional_resources (optional)
 ┃ ┃   ┃ ┗ 📜cm.yaml
 ┃ ┃   ┗ 📜values.yaml
 ┃ ┣ 📜global.yaml
 ┣ 📂env (optional)
 ┃ ┗ 📂test
 ┃   ┣ 📂applications
 ┃   ┃ ┗ 📂hello-world
 ┃   ┃   ┣ 📂additional_resources (optional)
 ┃   ┃   ┃ ┗ 📜sec.yaml
 ┃   ┃   ┗ 📜values.yaml
 ┃   ┗ 📜global.yaml
 ┗ 📜helm-overdrive.yaml
```

## code flow-chart
The diagram for the code flow can be seen here:<br>
[code-flow](docs/code-diagram.drawio.svg)

## Running the code
After having downloaded the binary or build it yourself, you can run the code using flags, environment variables, the config file or a mix of them.

The app is using the [viper](https://github.com/spf13/viper/) config package
so the priority for the input methods is:
  - flag
  - env
  - config

### Using flags
This templates the base on the chart specified.
```
./helm-overdrive template \
--chart_name hello-world \
--chart_version 0.1.0 \
--helm_repo https://helm.github.io/examples \
--base_folder config-example/base \
--global_file global.yaml \
--values_file values.yaml \
--application_folder applications/hello-world
```

Look at all the configs flag names [here](#config)

### Using environment variables
OBS! All envs for this app is prefixed with `HO`<br>
Available environment vars can be found [here](#config)

## config
| flag | env | config | description |
|------|-----|--------|-------------|
| --additional_resources | HO_ADDITIONAL_RESOURCES | additional_resources | Path to the folder that contains the additional resources, this has to be located within the <application_folder>, Same in base and env folders |
| --application_folder | HO_APPLICATION_FOLDER | application_folder | Path to the folder that contains the application, Same in base and env folders |
| --app_name | HO_APP_NAME | app_name | Name of the release |
| --base_folder | HO_BASE_FOLDER | base_folder | Path the folder containing the base config |
| --env_folder / -e | HO_ENV_FOLDER | env_folder | Name of the environment folder you with to deploy |
| --chart_version / -v | HO_CHART_VERSION | chart_version | Chart version |
| --chart_name / -n | HO_CHART_NAME | chart_name | Chart |
| --global_file | HO_GLOBAL_FILE | global_file | Name of the global files, same in base and env folders |
| --helm_repo | HO_HELM_REPO | helm_repo | Repo url |
| --values_file | HO_VALUE_FILES | values_file | Name of the value files in the application folder, Same in base and env folders |
