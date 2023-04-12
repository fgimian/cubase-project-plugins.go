package models

// Project specific configuration for the tool.
type Projects struct {
	Report32Bit bool `toml:"report_32_bit"` // whether or not to include 32-bit projects in output
	Report64Bit bool `toml:"report_64_bit"` // whether or not to include 64-bit projects in output
}

// Plugin specific configuration for the tool.
type Plugins struct {
	GUIDIgnores []string `toml:"guid_ignores"` // plugin GUIDs which should be ignored
	NameIgnores []string `toml:"name_ignores"` // plugin names which should be ignored
}

// The main configuration structure for the tool.
type Config struct {
	PathIgnorePatterns []string `toml:"path_ignore_patterns"` // path patterns to ignore
	Projects           Projects `toml:"projects"`             // project specific configuration
	Plugins            Plugins  `toml:"plugins"`              // plugin specific configuration
}
