package parser

import (
	"bytes"
	"errors"
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
func (r *Reader) GetProjectDetails() Project {
	metadata := Metadata{
		Application:  "Cubase",
		Version:      "Unknown",
		ReleaseDate:  "Unknown",
		Architecture: "Unknown",
	}
	plugins := make(map[Plugin]Nothing)

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
			plugins[*foundPlugin] = Nothing{}
			index = updatedIndex
		} else {
			index++
		}
	}

	return Project{Metadata: metadata, Plugins: plugins}
}

func (r *Reader) searchMetadata(index int) (*Metadata, int, bool) {
	readIndex := index

	versionTerm := r.getBytes(readIndex, len(AppVersionSearchTerm))
	if versionTerm == nil || string(versionTerm) != AppVersionSearchTerm {
		return nil, 0, false
	}

	readIndex += len(AppVersionSearchTerm) + 9

	application, readBytes, err := r.getToken(readIndex)
	if err != nil {
		return nil, 0, false
	}

	readIndex += readBytes + 3

	version, readBytes, err := r.getToken(readIndex)
	if err != nil {
		return nil, 0, false
	}

	readIndex += readBytes + 3

	releaseDate, readBytes, err := r.getToken(readIndex)
	if err != nil {
		return nil, 0, false
	}

	readIndex += readBytes + 7

	// Older 32-bit versions of Cubase didn't list the architecture in the project file.
	architecture, readBytes, err := r.getToken(readIndex)
	if err != nil {
		architecture = "Not Specified"
	} else {
		readIndex += readBytes
	}

	metadata := Metadata{
		Application:  application,
		Version:      version,
		ReleaseDate:  releaseDate,
		Architecture: architecture,
	}

	return &metadata, readIndex, true
}

func (r *Reader) searchPlugin(index int) (*Plugin, int, bool) {
	readIndex := index

	uidTerm := r.getBytes(readIndex, len(PluginUIDSearchTerm))
	if uidTerm == nil || string(uidTerm) != PluginUIDSearchTerm {
		return nil, 0, false
	}

	readIndex += len(PluginUIDSearchTerm) + 22

	guid, readBytes, err := r.getToken(readIndex)
	if err != nil {
		return nil, 0, false
	}

	readIndex += readBytes + 3

	key, readBytes, err := r.getToken(readIndex)
	if err != nil || key != "Plugin Name" {
		return nil, 0, false
	}

	readIndex += readBytes + 5

	name, readBytes, err := r.getToken(readIndex)
	if err != nil {
		return nil, 0, false
	}

	readIndex += readBytes + 3

	key, readBytes, err = r.getToken(readIndex)
	if err != nil {
		return nil, 0, false
	}

	if key == "Original Plugin Name" {
		readIndex += readBytes + 5

		originalName, readBytes, err := r.getToken(readIndex)
		if err == nil {
			readIndex += readBytes
			name = originalName
		}
	}

	plugin := Plugin{GUID: guid, Name: name}

	return &plugin, readIndex, true
}

func (r *Reader) getBytes(index, length int) []byte {
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
