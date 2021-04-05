package main

import (
	"errors"
	"regexp"
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
