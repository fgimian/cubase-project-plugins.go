package parser

import (
	"bytes"
	"errors"

	set "github.com/deckarep/golang-set/v2"
	"github.com/fgimian/cubase-project-plugins.go/models"
)

const (
	PluginUIDSearchTerm  = "Plugin UID\000"
	AppVersionSearchTerm = "PAppVersion\000"
)

var (
	ErrLengthBeyondEOF = errors.New("the length byte requested goes beyond the end of the project")
	ErrTokenBeyondEOF  = errors.New("the token size requested goes beyond the end of the project")
	ErrTokenNulMissing = errors.New("no null terminator was found in the token bytes")
)

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
func (r *Reader) GetProjectDetails() models.Project {
	metadata := models.Metadata{
		Application:  "Cubase",
		Version:      "Unknown",
		ReleaseDate:  "Unknown",
		Architecture: "Unknown",
	}
	plugins := set.NewSet[models.Plugin]()

	index := 0
	for index < len(r.projectBytes) {
		// Check if the current byte matches the letter P which is the first letter of all our
		// search terms.
		if r.projectBytes[index] != 'P' {
			index++
		} else if foundMetadata, updatedIndex, found := r.searchMetadata(index); found {
			// Check whether the next set of bytes are related to the Cubase version.
			metadata = *foundMetadata
			index = updatedIndex
		} else if foundPlugin, updatedIndex, found := r.searchPlugin(index); found {
			// Check whether the next set of bytes relate to a plugin.
			plugins.Add(*foundPlugin)
			index = updatedIndex
		} else {
			index++
		}
	}

	return models.Project{Metadata: metadata, Plugins: plugins}
}

func (r *Reader) searchMetadata(index int) (*models.Metadata, int, bool) {
	readIndex := index

	versionTerm := r.getBytes(readIndex, len(AppVersionSearchTerm))
	if versionTerm == nil || string(versionTerm) != AppVersionSearchTerm {
		return nil, 0, false
	}

	readIndex += len(AppVersionSearchTerm) + 9

	application, length, err := r.getToken(readIndex)
	if err != nil {
		return nil, 0, false
	}

	readIndex += length + 3

	version, length, err := r.getToken(readIndex)
	if err != nil {
		return nil, 0, false
	}

	readIndex += length + 3

	releaseDate, length, err := r.getToken(readIndex)
	if err != nil {
		return nil, 0, false
	}
	readIndex += length + 7

	// Older 32-bit versions of Cubase didn't list the architecture in the project file.
	architecture, length, err := r.getToken(readIndex)
	if err != nil {
		architecture = "Not Specified"
	} else {
		readIndex += length
	}

	metadata := models.Metadata{
		Application:  application,
		Version:      version,
		ReleaseDate:  releaseDate,
		Architecture: architecture,
	}
	return &metadata, readIndex, true
}

func (r *Reader) searchPlugin(index int) (*models.Plugin, int, bool) {
	readIndex := index

	uidTerm := r.getBytes(readIndex, len(PluginUIDSearchTerm))
	if uidTerm == nil || string(uidTerm) != PluginUIDSearchTerm {
		return nil, 0, false
	}

	readIndex += len(PluginUIDSearchTerm) + 22

	guid, length, err := r.getToken(readIndex)
	if err != nil {
		return nil, 0, false
	}

	readIndex += length + 3

	key, length, err := r.getToken(readIndex)
	if err != nil || key != "Plugin Name" {
		return nil, 0, false
	}

	readIndex += length + 5

	name, length, err := r.getToken(readIndex)
	if err != nil {
		return nil, 0, false
	}

	readIndex += length + 3

	key, length, err = r.getToken(readIndex)
	if err != nil {
		return nil, 0, false
	}
	if key == "Original Plugin Name" {
		readIndex += length + 5

		originalName, length, err := r.getToken(readIndex)
		if err == nil {
			readIndex += length
			name = originalName
		}
	}

	plugin := models.Plugin{GUID: guid, Name: name}
	return &plugin, readIndex, true
}

func (r *Reader) getBytes(index int, length int) []byte {
	end := index + length
	if end > len(r.projectBytes) {
		return nil
	}

	return r.projectBytes[index:end]
}

func (r *Reader) getToken(index int) (string, int, error) {
	lenBytes := r.getBytes(index, 1)
	if lenBytes == nil {
		return "", 0, ErrLengthBeyondEOF
	}
	length := int(lenBytes[0])

	tokenBytes := r.getBytes(index+1, length)
	if tokenBytes == nil {
		return "", 0, ErrTokenBeyondEOF
	}

	nullIndex := bytes.Index(tokenBytes, []byte{0})
	if nullIndex == -1 {
		return "", 0, ErrTokenNulMissing
	}

	token := string(tokenBytes[:nullIndex])
	return token, length + 1, nil
}
