package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var (
	ErrOpenConfigFile  = errors.New("unable to open the config file requested")
	ErrParseConfigFile = errors.New("unable to parse the config file requested")
	ErrWalkDir         = errors.New("unable to walk one or more directories requested")
)

// See https://github.com/golang/go/issues/50603 for a better way of determining the version
// in future.
var (
	version    string
	configPath string
)

var rootCmd = &cobra.Command{
	Use:     "cubase-project-plugins [flags] [project path]...",
	Version: version,
	Short: "Displays all plugins used in your Cubase projects along with the Cubase version " +
		"the project was created with.",
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		heading := color.New(color.BgRed, color.FgHiWhite)
		subHeading := color.New(color.FgHiBlue)

		config := Config{
			Projects: Projects{
				Report32Bit: true,
				Report64Bit: true,
			},
		}
		if configPath != "" {
			f, err := os.Open(configPath)
			if err != nil {
				return ErrOpenConfigFile
			}
			defer f.Close()

			_, err = toml.NewDecoder(f).Decode(&config)
			if err != nil {
				return ErrParseConfigFile
			}
		}

		pluginCounts := make(map[Plugin]int)
		pluginCounts32 := make(map[Plugin]int)
		pluginCounts64 := make(map[Plugin]int)

		for _, projectPath := range args {
			err := filepath.Walk(
				projectPath,
				func(path string, info fs.FileInfo, err error) error {
					if err != nil || filepath.Ext(path) != ".cpr" {
						return nil
					}

					for _, pathIgnorePattern := range config.PathIgnorePatterns {
						match, err := doublestar.Match(
							filepath.ToSlash(pathIgnorePattern),
							filepath.ToSlash(path),
						)
						if err == nil && match {
							return nil
						}
					}

					projectBytes, err := os.ReadFile(path)
					if err != nil {
						return nil
					}

					reader := NewReader(projectBytes)
					project := reader.GetProjectDetails()

					is64Bit := project.Metadata.Architecture == "WIN64" ||
						project.Metadata.Architecture == "MAC64 LE"

					if is64Bit && !config.Projects.Report64Bit ||
						!is64Bit && !config.Projects.Report32Bit {
						return nil
					}

					fmt.Println()
					heading.Printf("Path: %s", path)
					fmt.Println()

					fmt.Println()
					subHeading.Printf(
						"%s %s (%s)",
						project.Metadata.Application,
						project.Metadata.Version,
						project.Metadata.Architecture,
					)
					fmt.Println()

					var displayPlugins []Plugin

					for plugin := range project.Plugins.Iterator().C {
						if slices.Contains(config.Plugins.GUIDIgnores, plugin.GUID) {
							continue
						}

						if slices.Contains(config.Plugins.NameIgnores, plugin.Name) {
							continue
						}

						displayPlugins = append(displayPlugins, plugin)
					}

					if len(displayPlugins) == 0 {
						return nil
					}

					slices.SortFunc(displayPlugins, func(a, b Plugin) bool {
						return strings.ToLower(a.Name) < strings.ToLower(b.Name)
					})

					fmt.Println()
					for _, plugin := range displayPlugins {
						if is64Bit {
							pluginCounts64[plugin]++
						} else {
							pluginCounts32[plugin]++
						}

						pluginCounts[plugin]++

						fmt.Printf("    > %s : %s\n", plugin.GUID, plugin.Name)
					}

					return nil
				},
			)
			if err != nil {
				return ErrWalkDir
			}
		}

		printSummary(pluginCounts32, "32-bit", heading)
		printSummary(pluginCounts64, "64-bit", heading)
		printSummary(pluginCounts, "All", heading)

		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	_ = rootCmd.MarkFlagRequired("project-path")
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "config file `path`")
}

func printSummary(pluginCounts map[Plugin]int, description string, heading *color.Color) {
	if len(pluginCounts) == 0 {
		return
	}
	fmt.Println()
	heading.Printf("Summary: Plugins Used In %s Projects", description)
	fmt.Println()
	fmt.Println()

	plugins := make([]Plugin, 0, len(pluginCounts))
	for plugin := range pluginCounts {
		plugins = append(plugins, plugin)
	}

	slices.SortFunc(plugins, func(a, b Plugin) bool {
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})

	for _, plugin := range plugins {
		count := pluginCounts[plugin]
		fmt.Printf("    > %s : %s (%d)\n", plugin.GUID, plugin.Name, count)
	}
}
