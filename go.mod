module github.com/fgimian/cubase-project-plugins

go 1.20

require (
	github.com/BurntSushi/toml v1.2.1
	github.com/bmatcuk/doublestar/v4 v4.6.0
	github.com/fatih/color v1.15.0
	github.com/spf13/cobra v1.6.1
	golang.org/x/exp v0.0.0-20230321023759-10a507213a29
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.6.0 // indirect
)

retract (
	v1.0.0 // version not correctly set for program
	v1.0.1 // debug build to see build info
	v1.0.2 // debug build to confirm build info version works
)
