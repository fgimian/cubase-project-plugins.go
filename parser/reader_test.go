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
		name                   string
		filename               string
		version                string
		releaseDate            string
		architecture           string
		includesChannelPlugins bool
		ditherPluginName       string
	}{
		{
			name:         "Cubase 4.5 32-bit",
			filename:     "Example Project (Cubase 4.5 32-bit).cpr",
			version:      "4.5.2",
			releaseDate:  "Sep  2 2008",
			architecture: "WIN32",
		},
		{
			name:         "Cubase 4.5 64-bit",
			filename:     "Example Project (Cubase 4.5 64-bit).cpr",
			version:      "4.5.2",
			releaseDate:  "Sep  2 2008",
			architecture: "WIN64",
		},
		{
			name:         "Cubase 5 32-bit",
			filename:     "Example Project (Cubase 5 32-bit).cpr",
			version:      "5.5.3",
			releaseDate:  "Jan 13 2011",
			architecture: "WIN32",
		},
		{
			name:         "Cubase 5 64-bit",
			filename:     "Example Project (Cubase 5 64-bit).cpr",
			version:      "5.5.3",
			releaseDate:  "Jan 13 2011",
			architecture: "WIN64",
		},
		{
			name:         "Cubase 6.5 32-bit",
			filename:     "Example Project (Cubase 6.5 32-bit).cpr",
			version:      "6.5.5",
			releaseDate:  "Jun 24 2013",
			architecture: "WIN32",
		},
		{
			name:         "Cubase 6.5 64-bit",
			filename:     "Example Project (Cubase 6.5 64-bit).cpr",
			version:      "6.5.5",
			releaseDate:  "Jun 24 2013",
			architecture: "WIN64",
		},
		{
			name:                   "Cubase 7 32-bit",
			filename:               "Example Project (Cubase 7 32-bit).cpr",
			version:                "7.0.7",
			releaseDate:            "Jan 21 2014",
			architecture:           "WIN32",
			includesChannelPlugins: true,
		},
		{
			name:                   "Cubase 7 64-bit",
			filename:               "Example Project (Cubase 7 64-bit).cpr",
			version:                "7.0.7",
			releaseDate:            "Jan 21 2014",
			architecture:           "WIN64",
			includesChannelPlugins: true,
		},
		{
			name:                   "Cubase 8.5 32-bit",
			filename:               "Example Project (Cubase 8.5 32-bit).cpr",
			version:                "8.5.30",
			releaseDate:            "Feb 22 2017",
			architecture:           "WIN32",
			includesChannelPlugins: true,
		},
		{
			name:                   "Cubase 8.5 64-bit",
			filename:               "Example Project (Cubase 8.5 64-bit).cpr",
			version:                "8.5.30",
			releaseDate:            "Feb 22 2017",
			architecture:           "WIN64",
			includesChannelPlugins: true,
		},
		{
			name:                   "Cubase 9.5",
			filename:               "Example Project (Cubase 9.5).cpr",
			version:                "9.5.50",
			releaseDate:            "Feb  2 2019",
			architecture:           "WIN64",
			includesChannelPlugins: true,
		},
		{
			name:                   "Cubase 11",
			filename:               "Example Project (Cubase 11).cpr",
			version:                "11.0.41",
			releaseDate:            "Sep 27 2021",
			architecture:           "WIN64",
			includesChannelPlugins: true,
		},
		{
			name:                   "Cubase 13",
			filename:               "Example Project (Cubase 13).cpr",
			version:                "13.0.10",
			releaseDate:            "Oct 10 2023",
			architecture:           "WIN64",
			includesChannelPlugins: true,
			ditherPluginName:       "Lin Dither",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			projectBytes, err := os.ReadFile(filepath.Join("testdata", tc.filename))
			require.NoError(t, err)

			reader := parser.NewReader(projectBytes)
			project, err := reader.GetProjectDetails()
			require.NoError(t, err)

			ditherPluginName := "UV22HR"
			if tc.ditherPluginName != "" {
				ditherPluginName = tc.ditherPluginName
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
			if tc.includesChannelPlugins {
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
						Version:      tc.version,
						ReleaseDate:  tc.releaseDate,
						Architecture: tc.architecture,
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
		name          string
		filename      string
		expectedError error
	}{
		{
			name:          "Application",
			filename:      "Truncated Project (Application).cpr",
			expectedError: parser.ErrNoApplication,
		},
		{
			name:          "Version",
			filename:      "Truncated Project (Version).cpr",
			expectedError: parser.ErrNoVersion,
		},
		{
			name:          "Release Date",
			filename:      "Truncated Project (Release Date).cpr",
			expectedError: parser.ErrNoReleaseDate,
		},
		{
			name:          "Plugin GUID",
			filename:      "Truncated Project (Plugin GUID).cpr",
			expectedError: parser.ErrNoPluginGUID,
		},
		{
			name:          "Plugin Name Tag",
			filename:      "Truncated Project (Plugin Name Tag).cpr",
			expectedError: parser.ErrNoPluginName,
		},
		{
			name:          "Plugin Name Value",
			filename:      "Truncated Project (Plugin Name Value).cpr",
			expectedError: parser.ErrNoPluginName,
		},
		{
			name:          "Tag After Plugin Name",
			filename:      "Truncated Project (Tag After Plugin Name).cpr",
			expectedError: parser.ErrNoTokenAfterPluginName,
		},
		{
			name:          "Original Plugin Name",
			filename:      "Truncated Project (Original Plugin Name).cpr",
			expectedError: parser.ErrNoOriginalPluginName,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			projectBytes, err := os.ReadFile(filepath.Join("testdata", tc.filename))
			require.NoError(t, err)

			reader := parser.NewReader(projectBytes)
			project, err := reader.GetProjectDetails()

			require.ErrorIs(t, err, tc.expectedError)
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
