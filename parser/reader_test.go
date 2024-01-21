package parser_test

import (
	"cmp"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fgimian/cubase-project-plugins/parser"
)

func TestGetProjectDetails(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		Filename               string
		Version                string
		ReleaseDate            string
		Architecture           string
		IncludesChannelPlugins bool
		DitherPluginName       string
	}{
		{
			Name:         "Cubase 4.5 32-bit",
			Filename:     "Example Project (Cubase 4.5 32-bit).cpr",
			Version:      "4.5.2",
			ReleaseDate:  "Sep  2 2008",
			Architecture: "WIN32",
		},
		{
			Name:         "Cubase 4.5 64-bit",
			Filename:     "Example Project (Cubase 4.5 64-bit).cpr",
			Version:      "4.5.2",
			ReleaseDate:  "Sep  2 2008",
			Architecture: "WIN64",
		},
		{
			Name:         "Cubase 5 32-bit",
			Filename:     "Example Project (Cubase 5 32-bit).cpr",
			Version:      "5.5.3",
			ReleaseDate:  "Jan 13 2011",
			Architecture: "WIN32",
		},
		{
			Name:         "Cubase 5 64-bit",
			Filename:     "Example Project (Cubase 5 64-bit).cpr",
			Version:      "5.5.3",
			ReleaseDate:  "Jan 13 2011",
			Architecture: "WIN64",
		},
		{
			Name:         "Cubase 6.5 32-bit",
			Filename:     "Example Project (Cubase 6.5 32-bit).cpr",
			Version:      "6.5.5",
			ReleaseDate:  "Jun 24 2013",
			Architecture: "WIN32",
		},
		{
			Name:         "Cubase 6.5 64-bit",
			Filename:     "Example Project (Cubase 6.5 64-bit).cpr",
			Version:      "6.5.5",
			ReleaseDate:  "Jun 24 2013",
			Architecture: "WIN64",
		},
		{
			Name:                   "Cubase 7 32-bit",
			Filename:               "Example Project (Cubase 7 32-bit).cpr",
			Version:                "7.0.7",
			ReleaseDate:            "Jan 21 2014",
			Architecture:           "WIN32",
			IncludesChannelPlugins: true,
		},
		{
			Name:                   "Cubase 7 64-bit",
			Filename:               "Example Project (Cubase 7 64-bit).cpr",
			Version:                "7.0.7",
			ReleaseDate:            "Jan 21 2014",
			Architecture:           "WIN64",
			IncludesChannelPlugins: true,
		},
		{
			Name:                   "Cubase 8.5 32-bit",
			Filename:               "Example Project (Cubase 8.5 32-bit).cpr",
			Version:                "8.5.30",
			ReleaseDate:            "Feb 22 2017",
			Architecture:           "WIN32",
			IncludesChannelPlugins: true,
		},
		{
			Name:                   "Cubase 8.5 64-bit",
			Filename:               "Example Project (Cubase 8.5 64-bit).cpr",
			Version:                "8.5.30",
			ReleaseDate:            "Feb 22 2017",
			Architecture:           "WIN64",
			IncludesChannelPlugins: true,
		},
		{
			Name:                   "Cubase 9.5",
			Filename:               "Example Project (Cubase 9.5).cpr",
			Version:                "9.5.50",
			ReleaseDate:            "Feb  2 2019",
			Architecture:           "WIN64",
			IncludesChannelPlugins: true,
		},
		{
			Name:                   "Cubase 11",
			Filename:               "Example Project (Cubase 11).cpr",
			Version:                "11.0.41",
			ReleaseDate:            "Sep 27 2021",
			Architecture:           "WIN64",
			IncludesChannelPlugins: true,
		},
		{
			Name:                   "Cubase 13",
			Filename:               "Example Project (Cubase 13).cpr",
			Version:                "13.0.10",
			ReleaseDate:            "Oct 10 2023",
			Architecture:           "WIN64",
			IncludesChannelPlugins: true,
			DitherPluginName:       "Lin Dither",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			projectBytes, err := os.ReadFile(filepath.Join("testdata", tc.Filename))
			require.NoError(t, err)

			reader := parser.NewReader(projectBytes)
			project, err := reader.GetProjectDetails()
			require.NoError(t, err)

			ditherPluginName := "UV22HR"
			if tc.DitherPluginName != "" {
				ditherPluginName = tc.DitherPluginName
			}

			expectedPlugins := []parser.Plugin{
				{GUID: "1C3A662167D347A99F7D797EA4911CDB", Name: "Elephant"},
				{GUID: "44E1149EDB3E4387BDD827FEA3A39EE7", Name: "Standard Panner"},
				{GUID: "565354414152626172747361636F7573", Name: "ArtsAcousticReverb"},
				{GUID: "565354416D62726F6D6E697370686572", Name: "Omnisphere"},
				{GUID: "56535444475443747261636B636F6D70", Name: "TrackComp"},
				{GUID: "56535455564852757632326872000000", Name: ditherPluginName},
				{GUID: "56535473796C3173796C656E74683100", Name: "Sylenth1"},
				{GUID: "77BBA7CA90F14C9BB298BA9010D6DD78", Name: "StereoEnhancer"},
				{GUID: "946051208E29496E804F64A825C8A047", Name: "StudioEQ"},
				{GUID: "D39D5B69D6AF42FA1234567868495645", Name: "Hive"},
			}
			if tc.IncludesChannelPlugins {
				expectedPlugins = append(expectedPlugins, []parser.Plugin{
					{GUID: "297BA567D83144E1AE921DEF07B41156", Name: "EQ"},
					{GUID: "D56B9C6CA4F946018EED73EB83A74B58", Name: "Input Filter"},
				}...)
			}

			slices.SortFunc(project.Plugins, func(a, b parser.Plugin) int {
				return cmp.Compare(a.GUID, b.GUID)
			})

			slices.SortFunc(expectedPlugins, func(a, b parser.Plugin) int {
				return cmp.Compare(a.GUID, b.GUID)
			})

			require.Equal(
				t,
				parser.Project{
					Metadata: parser.Metadata{
						Application:  "Cubase",
						Version:      tc.Version,
						ReleaseDate:  tc.ReleaseDate,
						Architecture: tc.Architecture,
					},
					Plugins: expectedPlugins,
				},
				*project,
			)
		})
	}
}

func TestGetProjectDetailsSX3(t *testing.T) {
	t.Parallel()

	projectBytes, err := os.ReadFile(filepath.Join("testdata", "Example Project (Cubase SX3).cpr"))
	require.NoError(t, err)

	reader := parser.NewReader(projectBytes)
	project, err := reader.GetProjectDetails()
	require.NoError(t, err)

	require.Equal(
		t,
		parser.Project{
			Metadata: parser.Metadata{
				Application:  "Cubase SX",
				Version:      "3.1.1",
				ReleaseDate:  "Oct 13 2005",
				Architecture: "Unspecified",
			},
			Plugins: []parser.Plugin{},
		},
		*project,
	)
}

func TestGetProjectDetailsTruncated(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name     string
		Filename string
		Error    error
	}{
		{
			Name:     "Application",
			Filename: "Truncated Project (Application).cpr",
			Error:    parser.ErrNoApplication,
		},
		{
			Name:     "Version",
			Filename: "Truncated Project (Version).cpr",
			Error:    parser.ErrNoVersion,
		},
		{
			Name:     "Release Date",
			Filename: "Truncated Project (Release Date).cpr",
			Error:    parser.ErrNoReleaseDate,
		},
		{
			Name:     "Plugin GUID",
			Filename: "Truncated Project (Plugin GUID).cpr",
			Error:    parser.ErrNoPluginGUID,
		},
		{
			Name:     "Plugin Name Tag",
			Filename: "Truncated Project (Plugin Name Tag).cpr",
			Error:    parser.ErrNoPluginName,
		},
		{
			Name:     "Plugin Name Value",
			Filename: "Truncated Project (Plugin Name Value).cpr",
			Error:    parser.ErrNoPluginName,
		},
		{
			Name:     "Tag After Plugin Name",
			Filename: "Truncated Project (Tag After Plugin Name).cpr",
			Error:    parser.ErrNoTokenAfterPluginName,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			projectBytes, err := os.ReadFile(filepath.Join("testdata", tc.Filename))
			require.NoError(t, err)

			reader := parser.NewReader(projectBytes)
			project, err := reader.GetProjectDetails()

			require.ErrorIs(t, err, tc.Error)
			require.Nil(t, project)
		})
	}
}

func TestGetProjectDetailsInvalidProject(t *testing.T) {
	t.Parallel()

	projectBytes := []byte{}

	reader := parser.NewReader(projectBytes)
	project, err := reader.GetProjectDetails()

	require.ErrorIs(t, err, parser.ErrCorruptProject)
	require.Nil(t, project)
}
