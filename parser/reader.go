package parser

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

const (
	PluginUIDSearchTerm  = "Plugin UID\000"
	AppVersionSearchTerm = "PAppVersion\000"
)

var (
	ErrLengthBeyondEOF        = errors.New("the length byte goes beyond the end of the project")
	ErrTokenBeyondEOF         = errors.New("the token size goes beyond the end of the project")
	ErrCorruptProject         = errors.New("the project has no metadata and appears to be corrupt")
	ErrNoApplication          = errors.New("unable to obtain the application name")
	ErrNoVersion              = errors.New("unable to obtain the application version")
	ErrNoReleaseDate          = errors.New("unable to obtain the application release date")
	ErrNoPluginGUID           = errors.New("unable to obtain a plugin GUID")
	ErrNoPluginName           = errors.New("unable to obtain a plugin name")
	ErrNoTokenAfterPluginName = errors.New("unable to obtain the token after a plugin name")
	ErrNoOriginalPluginName   = errors.New("unable to obtain an original plugin name")
)

type Nothing struct{}

// Determines the used plugins in a Cubase project along with related version of Cubase which the
// project was created on by parsing the binary in a *.cpr file.
type Reader struct {
	projectBytes []byte
}

// NewReader returns a new reader that parses the given project bytes.
func NewReader(projectBytes []byte) Reader {
	return Reader{projectBytes: projectBytes}
}

// GetProjectDetails obtains all project details including Cubase version and plugins used and
// returns an instance of Project containing project details.
func (r *Reader) GetProjectDetails() (*Project, error) {
	var metadata *Metadata

	uniquePlugins := make(map[Plugin]Nothing)

	index := 0
	for index < len(r.projectBytes) {
		// Check if the current byte matches the letter P which is the first letter of all our
		// search terms.
		if r.projectBytes[index] != 'P' {
			index++
			continue
		}

		// Check whether the next set of bytes are related to the Cubase version.
		if metadata == nil {
			foundMetadata, updatedIndex, err := r.searchMetadata(index)
			if err != nil {
				return nil, fmt.Errorf("the project is corrupted: %w", err)
			}

			if foundMetadata != nil {
				metadata = foundMetadata
				index = updatedIndex

				continue
			}
		}

		// Check whether the next set of bytes relate to a plugin.
		foundPlugin, updatedIndex, err := r.searchPlugin(index)
		if err != nil {
			return nil, fmt.Errorf("the project is corrupted: %w", err)
		}

		if foundPlugin != nil {
			uniquePlugins[*foundPlugin] = Nothing{}
			index = updatedIndex

			continue
		}

		index++
	}

	if metadata == nil {
		return nil, ErrCorruptProject
	}

	plugins := make([]Plugin, 0, len(uniquePlugins))
	for plugin := range uniquePlugins {
		plugins = append(plugins, plugin)
	}

	return &Project{Metadata: *metadata, Plugins: plugins}, nil
}

func (r *Reader) searchMetadata(index int) (*Metadata, int, error) {
	versionTerm := r.getBytes(index, len(AppVersionSearchTerm))
	if versionTerm == nil || string(versionTerm) != AppVersionSearchTerm {
		return nil, 0, nil
	}

	index += len(AppVersionSearchTerm) + 9

	application, readBytes, err := r.getToken(index)
	if err != nil {
		return nil, 0, ErrNoApplication
	}

	index += readBytes + 3

	version, readBytes, err := r.getToken(index)
	if err != nil {
		return nil, 0, ErrNoVersion
	}

	version = strings.TrimPrefix(version, "Version ")

	index += readBytes + 3

	releaseDate, readBytes, err := r.getToken(index)
	if err != nil {
		return nil, 0, ErrNoReleaseDate
	}

	index += readBytes + 7

	// Older 32-bit versions of Cubase didn't list the architecture in the project file.
	architecture, readBytes, err := r.getToken(index)
	if err != nil {
		architecture = "Unspecified"
	} else {
		index += readBytes
	}

	metadata := Metadata{
		Application:  application,
		Version:      version,
		ReleaseDate:  releaseDate,
		Architecture: architecture,
	}

	return &metadata, index, nil
}

func (r *Reader) searchPlugin(index int) (*Plugin, int, error) {
	uidTerm := r.getBytes(index, len(PluginUIDSearchTerm))
	if uidTerm == nil || string(uidTerm) != PluginUIDSearchTerm {
		return nil, 0, nil
	}

	index += len(PluginUIDSearchTerm) + 22

	guid, readBytes, err := r.getToken(index)
	if err != nil {
		return nil, 0, ErrNoPluginGUID
	}

	index += readBytes + 3

	key, readBytes, err := r.getToken(index)
	if err != nil || key != "Plugin Name" {
		return nil, 0, ErrNoPluginName
	}

	index += readBytes + 5

	name, readBytes, err := r.getToken(index)
	if err != nil {
		return nil, 0, ErrNoPluginName
	}

	index += readBytes + 3

	key, readBytes, err = r.getToken(index)
	if err != nil {
		return nil, 0, ErrNoTokenAfterPluginName
	}

	// In Cubase 8.x and above, in cases where an instrument track has been renamed using
	// Shift+Enter, the name retrieved above will be the track title and the name of the plugin
	// will follow under the key "Original Plugin Name".
	if key == "Original Plugin Name" {
		index += readBytes + 5

		name, readBytes, err = r.getToken(index)
		if err != nil {
			return nil, 0, ErrNoOriginalPluginName
		}

		index += readBytes
	}

	plugin := Plugin{GUID: guid, Name: name}

	return &plugin, index, nil
}

func (r *Reader) getBytes(index, length int) []byte {
	end := index + length
	if end > len(r.projectBytes) {
		return nil
	}

	return r.projectBytes[index:end]
}

func (r *Reader) getToken(index int) (token string, readBytes int, err error) {
	lenBytes := r.getBytes(index, 1)
	if lenBytes == nil {
		return "", 0, ErrLengthBeyondEOF
	}

	length := int(lenBytes[0])

	tokenBytes := r.getBytes(index+1, length)
	if tokenBytes == nil {
		return "", 0, ErrTokenBeyondEOF
	}

	// Older versions of before Cubase 5 didn't always provide nul terminators in token strings.
	nulIndex := bytes.Index(tokenBytes, []byte{0})
	if nulIndex == -1 {
		token = string(tokenBytes)
	} else {
		token = string(tokenBytes[:nulIndex])
	}

	readBytes = length + 1

	return token, readBytes, nil
}
