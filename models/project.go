package models

import (
	set "github.com/deckarep/golang-set/v2"
)

// Contains information about the Cubase version used to create the project.
type Metadata struct {
	Application  string
	Version      string
	ReleaseDate  string
	Architecture string
}

// Represents a plugin within a Cubase project.
type Plugin struct {
	GUID string
	Name string
}

// Captures the Cubase version and all plugins used for a Cubase project.
type Project struct {
	Metadata Metadata
	Plugins  set.Set[Plugin]
}
