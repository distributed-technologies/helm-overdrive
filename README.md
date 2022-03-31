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
 ┃ ┣ 📜global.yaml
 ┃ ┣ 📜nidhogg.yaml
 ┃ ┗ 📜yggdrasil.yaml
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
