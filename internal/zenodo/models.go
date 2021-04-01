package zenodo

import "time"

type Deposition struct {
	Conceptdoi   string              `json:"conceptdoi,omitempty"`
	Conceptrecid string              `json:"conceptrecid,omitempty"`
	Created      *time.Time          `json:"created,omitempty"`
	Doi          string              `json:"doi,omitempty"`
	DoiURL       string              `json:"doi_url,omitempty"`
	Files        []*DepositionFile   `json:"files,omitempty"`
	ID           int                 `json:"id,omitempty"`
	Links        map[string]string   `json:"links,omitempty"`
	Metadata     *DepositionMetadata `json:"metadata,omitempty"`
	Modified     *time.Time          `json:"modified,omitempty"`
	Owner        int                 `json:"owner,omitempty"`
	RecordID     int                 `json:"record_id,omitempty"`
	State        string              `json:"state,omitempty"`
	Submitted    bool                `json:"submitted,omitempty"`
	Title        string              `json:"title,omitempty"`
}

type DepositionFile struct {
	Checksum string            `json:"checksum"`
	Filename string            `json:"filename"`
	Filesize int               `json:"filesize"`
	ID       string            `json:"id"`
	Links    map[string]string `json:"links"`
}

type DepositionMetadata struct {
	AccessRight     string         `json:"access_right,omitempty"`
	Communities     []*Community   `json:"communities,omitempty"`
	Creators        []*Creator     `json:"creators,omitempty"`
	Description     string         `json:"description,omitempty"`
	Doi             string         `json:"doi,omitempty"`
	License         string         `json:"license,omitempty"`
	PrereserveDoi   *PrereserveDoi `json:"prereserve_doi,omitempty"`
	PublicationDate string         `json:"publication_date,omitempty"`
	Title           string         `json:"title,omitempty"`
	UploadType      string         `json:"upload_type,omitempty"`
}

type Community struct {
	Identifier string `json:"identifier"`
}

type Creator struct {
	Name string `json:"name"`
}

type PrereserveDoi struct {
	Doi   string `json:"doi"`
	Recid int    `json:"recid"`
}
