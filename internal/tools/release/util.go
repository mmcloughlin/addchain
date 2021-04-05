package main

import (
	"errors"
	"path/filepath"
	"regexp"
	"runtime"
)

// ValidateVersion checks if version is of the form "MAJOR.MINOR.PATCH". That
// is, a simple semver format without "v" prefix, pre-release or build suffixes.
func ValidateVersion(version string) error {
	if !versionrx.MatchString(version) {
		return errors.New("version must be of the form MAJOR.MINOR.PATCH")
	}
	return nil
}

var versionrx = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

// RepoPath returns a full file path given a path relative to the repository
// root. Returns empty string if it cannot be determined.
func RepoPath(rel string) string {
	_, self, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	path := filepath.Join(filepath.Dir(self), "../../..", rel)
	return filepath.Clean(path)
}
