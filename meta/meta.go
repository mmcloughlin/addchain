// Package meta defines properties about this project.
package meta

import "time"

// VersionTagPrefix is the prefix used on Git tags corresponding to semantic
// version releases.
const VersionTagPrefix = "v"

// Properties about this software package.
type Properties struct {
	// BuildVersion is the version that was built. Typically populated at build
	// time and will typically be empty for non-release builds.
	BuildVersion string

	// ReleaseVersion is the version of the most recent release.
	ReleaseVersion string

	// ReleaseDate is the date of the most recent release. (RFC3339 date format.)
	ReleaseDate string

	// ConceptDOI is the DOI for all versions.
	ConceptDOI string

	// DOI for the most recent release.
	DOI string

	// ZenodoID is the Zenodo deposit ID for the most recent release.
	ZenodoID string
}

// Meta defines specific properties for the current version of this software.
var Meta = &Properties{
	BuildVersion:   buildversion,
	ReleaseVersion: releaseversion,
	ReleaseDate:    releasedate,
	ConceptDOI:     conceptdoi,
	DOI:            doi,
	ZenodoID:       zenodoid,
}

// IsRelease reports whether the built version is a release.
func (p *Properties) IsRelease() bool {
	return p.BuildVersion == p.ReleaseVersion
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
func (p *Properties) DOIURL() string { return doiurl(p.DOI) }

// ConceptDOIURL returns the DOI URL corresponding to the most recent release.
func (p *Properties) ConceptDOIURL() string { return doiurl(p.ConceptDOI) }

func doiurl(doi string) string {
	return "https://doi.org/" + doi
}
