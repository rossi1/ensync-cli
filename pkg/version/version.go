package version

import "fmt"

var (
	// Version holds the current version number
	version = "dev"
	// Commit holds the git commit hash
	commit = "none"
	// BuildDate holds the build date
	buildDate = "unknown"
)

// Info represents version information
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"buildDate"`
}

// Get returns version information
func Get() Info {
	return Info{
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
	}
}

// String returns version information as a string
func String() string {
	return fmt.Sprintf("Version: %s\nCommit: %s\nBuild Date: %s",
		version, commit, buildDate)
}
