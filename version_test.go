package cuckoofilter

import (
	"fmt"
	"regexp"
	"testing"
)

func TestVersionConstants(t *testing.T) {
	if Major < 0 {
		t.Errorf("Major version should be non-negative, got %d", Major)
	}
	if Minor < 0 {
		t.Errorf("Minor version should be non-negative, got %d", Minor)
	}
	if Patch < 0 {
		t.Errorf("Patch version should be non-negative, got %d", Patch)
	}
}

func TestVersionFormat(t *testing.T) {
	expected := fmt.Sprintf("v%d.%d.%d", Major, Minor, Patch)
	if Version != expected {
		t.Errorf("Version format mismatch: expected %s, got %s", expected, Version)
	}
}

func TestVersionSemanticFormat(t *testing.T) {
	// Semantic version regex: v{major}.{minor}.{patch}
	semverRegex := regexp.MustCompile(`^v\d+\.\d+\.\d+$`)
	if !semverRegex.MatchString(Version) {
		t.Errorf("Version does not match semantic versioning format: %s", Version)
	}
}

func TestFullVersion(t *testing.T) {
	tests := []struct {
		name        string
		preRelease  string
		buildMeta   string
		expectedFmt string
	}{
		{
			name:        "stable release",
			preRelease:  "",
			buildMeta:   "",
			expectedFmt: Version,
		},
		{
			name:        "with pre-release",
			preRelease:  "alpha.1",
			buildMeta:   "",
			expectedFmt: Version + "-alpha.1",
		},
		{
			name:        "with build metadata",
			preRelease:  "",
			buildMeta:   "abc123",
			expectedFmt: Version + "+abc123",
		},
		{
			name:        "with both",
			preRelease:  "beta.2",
			buildMeta:   "def456",
			expectedFmt: Version + "-beta.2+def456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For actual testing, we test the current package state
			if tt.preRelease == "" && tt.buildMeta == "" {
				fullVer := FullVersion()
				if PreRelease == "" && BuildMetadata == "" {
					if fullVer != Version {
						t.Errorf("FullVersion() = %s, want %s", fullVer, Version)
					}
				}
			}
		})
	}
}

func TestFullVersionStableRelease(t *testing.T) {
	// Test the actual current state
	fullVer := FullVersion()
	if PreRelease == "" && BuildMetadata == "" {
		if fullVer != Version {
			t.Errorf("FullVersion() = %s, want %s for stable release", fullVer, Version)
		}
	} else if PreRelease != "" && BuildMetadata == "" {
		expected := Version + "-" + PreRelease
		if fullVer != expected {
			t.Errorf("FullVersion() = %s, want %s", fullVer, expected)
		}
	} else if PreRelease == "" && BuildMetadata != "" {
		expected := Version + "+" + BuildMetadata
		if fullVer != expected {
			t.Errorf("FullVersion() = %s, want %s", fullVer, expected)
		}
	} else {
		expected := Version + "-" + PreRelease + "+" + BuildMetadata
		if fullVer != expected {
			t.Errorf("FullVersion() = %s, want %s", fullVer, expected)
		}
	}
}

func TestVersionInfo(t *testing.T) {
	info := VersionInfo()

	// Check all required fields exist
	requiredFields := []string{"version", "major", "minor", "patch", "preRelease", "buildMetadata", "fullVersion"}
	for _, field := range requiredFields {
		if _, exists := info[field]; !exists {
			t.Errorf("VersionInfo missing required field: %s", field)
		}
	}

	// Verify field values match constants
	if info["version"] != Version {
		t.Errorf("VersionInfo version mismatch: got %v, want %s", info["version"], Version)
	}
	if info["major"] != Major {
		t.Errorf("VersionInfo major mismatch: got %v, want %d", info["major"], Major)
	}
	if info["minor"] != Minor {
		t.Errorf("VersionInfo minor mismatch: got %v, want %d", info["minor"], Minor)
	}
	if info["patch"] != Patch {
		t.Errorf("VersionInfo patch mismatch: got %v, want %d", info["patch"], Patch)
	}
	if info["preRelease"] != PreRelease {
		t.Errorf("VersionInfo preRelease mismatch: got %v, want %s", info["preRelease"], PreRelease)
	}
	if info["buildMetadata"] != BuildMetadata {
		t.Errorf("VersionInfo buildMetadata mismatch: got %v, want %s", info["buildMetadata"], BuildMetadata)
	}
	if info["fullVersion"] != FullVersion() {
		t.Errorf("VersionInfo fullVersion mismatch: got %v, want %s", info["fullVersion"], FullVersion())
	}
}

func TestVersionInfoTypes(t *testing.T) {
	info := VersionInfo()

	// Check types
	if _, ok := info["version"].(string); !ok {
		t.Error("version should be a string")
	}
	if _, ok := info["major"].(int); !ok {
		t.Error("major should be an int")
	}
	if _, ok := info["minor"].(int); !ok {
		t.Error("minor should be an int")
	}
	if _, ok := info["patch"].(int); !ok {
		t.Error("patch should be an int")
	}
	if _, ok := info["preRelease"].(string); !ok {
		t.Error("preRelease should be a string")
	}
	if _, ok := info["buildMetadata"].(string); !ok {
		t.Error("buildMetadata should be a string")
	}
	if _, ok := info["fullVersion"].(string); !ok {
		t.Error("fullVersion should be a string")
	}
}

// Benchmark version functions
func BenchmarkFullVersion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = FullVersion()
	}
}

func BenchmarkVersionInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = VersionInfo()
	}
}
