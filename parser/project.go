package parser

// Contains information about the Cubase version used to create the project.
type Metadata struct {
	Application  string // application name (this is always "Cubase")
	Version      string // version of Cubase used to create the project
	ReleaseDate  string // release date of the Cubase version used
	Architecture string // system architecture used to create the project
}

// Represents a plugin within a Cubase project.
type Plugin struct {
	GUID string // globally unique identifier for the plugin
	Name string // name of the plugin
}

// Captures the Cubase version and all plugins used for a Cubase project.
type Project struct {
	Metadata Metadata // metadata describing the Cubase version used to create the project
	Plugins  []Plugin // plugins used in the project
}
