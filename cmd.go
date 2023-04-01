package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var configPath string

var rootCmd = &cobra.Command{
	Use:     "cubase-project-plugins [flags] [project path]...",
	Version: "1.0.0",
	Short:   "Displays all plugins used in your Cubase projects along with the Cubase version the project was created with.",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		failure := color.New(color.FgRed)
		heading := color.New(color.BgRed, color.FgHiWhite)
		subHeading := color.New(color.FgHiBlue)

		config := Config{
			Projects: Projects{
				Report32Bit: true,
				Report64Bit: true,
			},
		}
		if configPath != "" {
			_, err := toml.DecodeFile(configPath, &config)
			if err != nil {
				failure.Fprintf(
					os.Stderr,
					"Error: Unable to open the config file at %s\n",
					configPath,
				)
				os.Exit(1)
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
						if slices.Contains(config.Plugins.GuidIgnores, plugin.Guid) {
							continue
						}

						if slices.Contains(config.Plugins.NameIgnores, plugin.Name) {
							continue
						}

						displayPlugins = append(displayPlugins, plugin)
					}

					slices.SortFunc(displayPlugins, func(a, b Plugin) bool {
						return strings.ToLower(a.Name) < strings.ToLower(b.Name)
					})

					if len(displayPlugins) > 0 {
						fmt.Println()
						for _, plugin := range displayPlugins {
							pluginCounts[plugin]++

							if is64Bit {
								pluginCounts64[plugin]++
							} else {
								pluginCounts32[plugin]++
							}

							fmt.Printf("    > %s : %s\n", plugin.Guid, plugin.Name)
						}
					}

					return nil
				},
			)

			if err != nil {
				failure.Fprintf(
					os.Stderr,
					"Error: Unable to walk the directory %s\n",
					projectPath,
				)
				os.Exit(1)
			}
		}

		if len(pluginCounts32) != 0 {
			fmt.Println()
			heading.Printf("Summary: Plugins Used In 32-bit Projects")
			fmt.Println()
			fmt.Println()

			plugins := make([]Plugin, 0, len(pluginCounts32))
			for plugin := range pluginCounts32 {
				plugins = append(plugins, plugin)
			}

			slices.SortFunc(plugins, func(a, b Plugin) bool {
				return strings.ToLower(a.Name) < strings.ToLower(b.Name)
			})

			for _, plugin := range plugins {
				count := pluginCounts32[plugin]
				fmt.Printf("    > %s : %s (%d)\n", plugin.Guid, plugin.Name, count)
			}
		}

		if len(pluginCounts64) != 0 {
			fmt.Println()
			heading.Printf("Summary: Plugins Used In 64-bit Projects")
			fmt.Println()
			fmt.Println()

			plugins := make([]Plugin, 0, len(pluginCounts64))
			for plugin := range pluginCounts64 {
				plugins = append(plugins, plugin)
			}

			slices.SortFunc(plugins, func(a, b Plugin) bool {
				return strings.ToLower(a.Name) < strings.ToLower(b.Name)
			})

			for _, plugin := range plugins {
				count := pluginCounts64[plugin]
				fmt.Printf("    > %s : %s (%d)\n", plugin.Guid, plugin.Name, count)
			}
		}

		if len(pluginCounts) != 0 {
			fmt.Println()
			heading.Printf("Summary: Plugins Used In All Projects")
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
				fmt.Printf("    > %s : %s (%d)\n", plugin.Guid, plugin.Name, count)
			}
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.MarkFlagRequired("project-path")
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "config file `path`")
}
