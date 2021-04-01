package main

import (
	"errors"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/mmcloughlin/addchain/internal/metavars"
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

// DefaultMetaVarsPath returns the path to the default variables file in the
// meta package.
func DefaultMetaVarsPath() string {
	_, self, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(self), "../../meta/vars.go")
	return filepath.Clean(path)
}

// SetMetaVar sets the given variable in a meta variables file.
func SetMetaVar(filename, name, value string) error {
	f, err := metavars.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := f.Set(name, value); err != nil {
		return err
	}

	if err := metavars.WriteFile(filename, f); err != nil {
		return err
	}

	return nil
}
