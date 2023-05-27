package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/fgimian/cubase-project-plugins/config"
	"github.com/fgimian/cubase-project-plugins/parser"
)

var (
	ErrOpenConfigFile  = errors.New("unable to open the config file requested")
	ErrParseConfigFile = errors.New("unable to parse the config file requested")
	ErrWalkDir         = errors.New("unable to walk one or more directories requested")
)

var configPath string

var rootCmd = &cobra.Command{
	Use: "cubase-project-plugins [flags] [project path]...",
	Short: "Displays all plugins used in your Cubase projects along with the Cubase version " +
		"the project was created with.",
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		heading := color.New(color.BgRed, color.FgHiWhite)
		subHeading := color.New(color.FgHiBlue)

		config := config.Config{
			Projects: config.Projects{
				Report32Bit: true,
				Report64Bit: true,
			},
		}

		if configPath == "" {
			defaultConfigPath := getDefaultConfigPath()
			if defaultConfigPath != "" {
				if _, err := os.Stat(defaultConfigPath); err == nil {
					configPath = defaultConfigPath
				}
			}
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

		pluginCounts := make(map[parser.Plugin]int)
		pluginCounts32 := make(map[parser.Plugin]int)
		pluginCounts64 := make(map[parser.Plugin]int)

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

					reader := parser.NewReader(projectBytes)
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

					var displayPlugins []parser.Plugin

					for _, plugin := range maps.Keys(project.Plugins) {
						if slices.Contains(config.Plugins.GUIDIgnores, plugin.GUID) ||
							slices.Contains(config.Plugins.NameIgnores, plugin.Name) {
							continue
						}

						displayPlugins = append(displayPlugins, plugin)
					}

					if len(displayPlugins) == 0 {
						return nil
					}

					slices.SortFunc(displayPlugins, func(a, b parser.Plugin) bool {
						return strings.ToLower(a.Name) < strings.ToLower(b.Name)
					})

					fmt.Println()
					for _, plugin := range displayPlugins {
						pluginCounts[plugin]++
						if is64Bit {
							pluginCounts64[plugin]++
						} else {
							pluginCounts32[plugin]++
						}

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
	if info, ok := debug.ReadBuildInfo(); ok {
		rootCmd.Version = info.Main.Version
	}

	_ = rootCmd.MarkFlagRequired("project-path")
	rootCmd.Flags().
		StringVarP(&configPath, "config", "c", "", "config file `path`")
}

func getDefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(home, ".config", "cubase-project-plugins.toml")
}

func printSummary(pluginCounts map[parser.Plugin]int, description string, heading *color.Color) {
	if len(pluginCounts) == 0 {
		return
	}

	fmt.Println()
	heading.Printf("Summary: Plugins Used In %s Projects", description)
	fmt.Println()
	fmt.Println()

	plugins := make([]parser.Plugin, 0, len(pluginCounts))
	for plugin := range pluginCounts {
		plugins = append(plugins, plugin)
	}

	slices.SortFunc(plugins, func(a, b parser.Plugin) bool {
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})

	for _, plugin := range plugins {
		count := pluginCounts[plugin]
		fmt.Printf("    > %s : %s (%d)\n", plugin.GUID, plugin.Name, count)
	}
}
