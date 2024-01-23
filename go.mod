module github.com/fgimian/cubase-project-plugins

go 1.21

require (
	github.com/BurntSushi/toml v1.3.2
	github.com/bmatcuk/doublestar/v4 v4.6.1
	github.com/fatih/color v1.16.0
	github.com/spf13/cobra v1.8.0
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.16.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	v1.0.2 // debug build to confirm build info version works
	v1.0.1 // debug build to see build info
	v1.0.0 // version not correctly set for program
)
