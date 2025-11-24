package cuckoofilter

// Version information for the SIMDCuckooFilter package
const (
	// Major version number - incremented for breaking changes
	Major = 0

	// Minor version number - incremented for new features (backward compatible)
	Minor = 1

	// Patch version number - incremented for bug fixes
	Patch = 0

	// Version is the full semantic version string
	Version = "v0.1.0"

	// PreRelease indicates pre-release version info (e.g., "alpha.1", "beta.2")
	// Empty string for stable releases
	PreRelease = ""

	// BuildMetadata contains build metadata (e.g., commit hash, build date)
	// Empty string if not provided
	BuildMetadata = ""
)

// FullVersion returns the complete version string including pre-release and build metadata
func FullVersion() string {
	v := Version
	if PreRelease != "" {
		v += "-" + PreRelease
	}
	if BuildMetadata != "" {
		v += "+" + BuildMetadata
	}
	return v
}

// VersionInfo returns a structured version information
func VersionInfo() map[string]interface{} {
	return map[string]interface{}{
		"version":       Version,
		"major":         Major,
		"minor":         Minor,
		"patch":         Patch,
		"preRelease":    PreRelease,
		"buildMetadata": BuildMetadata,
		"fullVersion":   FullVersion(),
	}
}
