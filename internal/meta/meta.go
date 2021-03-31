package meta

import "time"

const VersionTagPrefix = "v"

type Properties struct {
	// ReleaseVersion is the version of the most recent release.
	ReleaseVersion string

	// ReleaseDate is the date of the most recent release. (RFC3339 date format.)
	ReleaseDate string

	// DOI for the most recent release.
	DOI string

	// Zenodo record ID for the most recent release.
	ZenodoRecordID string
}

var Meta = &Properties{
	ReleaseVersion: releaseversion,
	ReleaseDate:    releasedate,
	DOI:            doi,
	ZenodoRecordID: zenodorecordid,
}

func (p *Properties) ReleaseTag() string {
	return VersionTagPrefix + p.ReleaseVersion
}

// ReleaseTime returns the release date as a time object.
func (p *Properties) ReleaseTime() (time.Time, error) {
	return time.Parse("2006-01-02", p.ReleaseDate)
}

func (p *Properties) DOIURL() string {
	return "https://doi.org/" + p.DOI
}
