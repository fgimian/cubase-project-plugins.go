package main

// Project specific configuration for the tool.
type Projects struct {
	Report32Bit bool `toml:"report_32_bit"`
	Report64Bit bool `toml:"report_64_bit"`
}

// Plugin specific configuration for the tool.
type Plugins struct {
	GuidIgnores []string `toml:"guid_ignores"`
	NameIgnores []string `toml:"name_ignores"`
}

// The main configuration structure for the tool.
type Config struct {
	PathIgnorePatterns []string `toml:"path_ignore_patterns"`
	Projects           Projects `toml:"projects"`
	Plugins            Plugins  `toml:"plugins"`
}
