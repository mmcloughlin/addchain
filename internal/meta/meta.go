// Package meta defines properties about this project.
package meta

import "time"

// VersionTagPrefix is the prefix used on Git tags corresponding to semantic
// version releases.
const VersionTagPrefix = "v"

// Properties about this software package.
type Properties struct {
	// ReleaseVersion is the version of the most recent release.
	ReleaseVersion string

	// ReleaseDate is the date of the most recent release. (RFC3339 date format.)
	ReleaseDate string

	// DOI for the most recent release.
	DOI string
}

// Meta defines specific properties for the current version of this software.
var Meta = &Properties{
	ReleaseVersion: releaseversion,
	ReleaseDate:    releasedate,
	DOI:            doi,
}

// ReleaseTag returns the release tag corresponding to the most recent release.
func (p *Properties) ReleaseTag() string {
	return VersionTagPrefix + p.ReleaseVersion
}

// GithubReleaseURL returns the URL to the release page on Github.
func (p *Properties) GithubReleaseURL() string {
	return "https://github.com/mmcloughlin/addchain/releases/tag/" + p.ReleaseTag()
}

// ReleaseTime returns the release date as a time object.
func (p *Properties) ReleaseTime() (time.Time, error) {
	return time.Parse("2006-01-02", p.ReleaseDate)
}

// DOIURL returns the DOI URL corresponding to the most recent release.
func (p *Properties) DOIURL() string {
	return "https://doi.org/" + p.DOI
}
