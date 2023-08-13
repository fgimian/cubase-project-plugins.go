module github.com/fgimian/cubase-project-plugins

go 1.21

require (
	github.com/BurntSushi/toml v1.3.2
	github.com/bmatcuk/doublestar/v4 v4.6.0
	github.com/fatih/color v1.15.0
	github.com/spf13/cobra v1.7.0
	golang.org/x/exp v0.0.0-20230811145659-89c5cff77bcb
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.11.0 // indirect
)

retract (
	v1.0.2 // debug build to confirm build info version works
	v1.0.1 // debug build to see build info
	v1.0.0 // version not correctly set for program
)
