package config

// Project specific configuration for the tool.
type Projects struct {
	Report32Bit bool `toml:"report_32_bit"` // whether 32-bit projects should be reported.
	Report64Bit bool `toml:"report_64_bit"` // whether 64-bit projects should be reported.
}

// Plugin specific configuration for the tool.
type Plugins struct {
	GUIDIgnores []string `toml:"guid_ignores"` // plugin GUIDs which should be ignored
	NameIgnores []string `toml:"name_ignores"` // plugin names which should be ignored
}

// The main configuration structure for the tool.
type Config struct {
	PathIgnorePatterns []string `toml:"path_ignore_patterns"` // project path patterns to skip
	Projects           Projects `toml:"projects"`             // configuration related to projects
	Plugins            Plugins  `toml:"plugins"`              // configuration related to plugins
}
