# helm-overdrive

## File structure
The file structure currently needed to use this app
``` tree
📦config-example
 ┣ 📂base
 ┃ ┣ 📂applications
 ┃ ┃ ┗ 📂hello-world
 ┃ ┃ ┃ ┣ 📂additional_resources (optional)
 ┃ ┃ ┃ ┃ ┗ 📜cm.yaml
 ┃ ┃ ┃ ┗ 📜values.yaml
 ┃ ┃ ┗ 📂 nidhogg.yaml
       ┗ 📜nidhogg.yaml
 ┃ ┣ 📜global.yaml
 ┣ 📂env
 ┃ ┗ 📂test
 ┃ ┃ ┣ 📂applications
 ┃ ┃ ┃ ┗ 📂hello-world
 ┃ ┃ ┃ ┃ ┣ 📂additional_resources (optional)
 ┃ ┃ ┃ ┃ ┃ ┗ 📜sec.yaml
 ┃ ┃ ┃ ┃ ┗ 📜values.yaml
 ┃ ┃ ┗ 📜global.yaml
 ┗ 📜helm-overdrive.yaml
```

go run main.go template -c ./config-example/helm-overdrive.yaml -n hello-world -v 0.1.0 --helm_repo https://helm.github.io/examples

helm-overdrive template -v 1.0 --chart nidhogg --repo deptech --environment test --applicationsfolder applications/tes/42/420/69 --application hello-world --valuefile values.yaml

helm-overdrive template -v 1.0 --chart nidhogg --repo deptech --environment test --application applications/hello-world/aerg/ghear/hrage --valuefile values.yaml


./helm-overdrive template \
-n hello-world \
-v 0.1.0 \
--helm_repo https://helm.github.io/examples \
--base_folder config-example/base \
--env_folder config-example/env/test \
--global_file global.yaml \
--values_file values.yaml \
--application_folder applications/hello-world \
--additional_resources additional_resources \
--debug

./helm-overdrive template \
-c ./config-example/helm-overdrive.yaml \
-n hello-world \
-v 0.1.0 \
--helm_repo https://helm.github.io/examples
