package models

type Release struct {
	SchemaVersion      string              `json:"schema_version"`
	ID                 string              `json:"id"`
	Title              string              `json:"title"`
	ArtistID           string              `json:"artist_id"`
	ReleaseType        string              `json:"release_type"`
	Status             string              `json:"status"`
	License            string              `json:"license"`
	Description        *ReleaseDescription `json:"description,omitempty"`
	GenreTags          []string            `json:"genre_tags,omitempty"`
	Tracks             []string            `json:"tracks,omitempty"`
	Dates              *ReleaseDates       `json:"dates,omitempty"`
	Distribution       *ReleaseDistrib     `json:"distribution,omitempty"`
	MasteringProfileID string              `json:"mastering_profile_id,omitempty"`
	Artwork            *ReleaseArtwork     `json:"artwork,omitempty"`
	Archive            *ReleaseArchive     `json:"archive,omitempty"`
	StoreLinks         *StoreLinks         `json:"store_links,omitempty"`
	QC                 *ReleaseQC          `json:"qc,omitempty"`
	CreatedDate        string              `json:"created_date,omitempty"`
	Notes              string              `json:"notes,omitempty"`
}

type ReleaseDescription struct {
	Short string `json:"short,omitempty"`
	Full  string `json:"full,omitempty"`
}

type ReleaseDates struct {
	TargetRelease                 string `json:"target_release,omitempty"`
	DistributorSubmissionDeadline string `json:"distributor_submission_deadline,omitempty"`
	Submitted                     string `json:"submitted,omitempty"`
	Released                      string `json:"released,omitempty"`
}

type ReleaseDistrib struct {
	Distributor   string `json:"distributor,omitempty"`
	UPC           string `json:"upc,omitempty"`
	CatalogNumber string `json:"catalog_number,omitempty"`
}

type ReleaseArtwork struct {
	PrimaryFile  string `json:"primary_file,omitempty"`
	DimensionsPx int    `json:"dimensions_px,omitempty"`
	Format       string `json:"format,omitempty"`
	QCPassed     bool   `json:"qc_passed,omitempty"`
}

type ReleaseArchive struct {
	InternetArchiveID    *string `json:"internet_archive_id"`
	InternetArchiveURL   *string `json:"internet_archive_url"`
	ObjectStoragePath    *string `json:"object_storage_path"`
	MastersArchived      bool    `json:"masters_archived"`
	StemsArchived        bool    `json:"stems_archived"`
	ProjectFilesArchived bool    `json:"project_files_archived"`
	ArchiveDate          *string `json:"archive_date"`
	ChecksumsVerified    bool    `json:"checksums_verified"`
}

type StoreLinks struct {
	Spotify      *string `json:"spotify"`
	AppleMusic   *string `json:"apple_music"`
	YouTubeMusic *string `json:"youtube_music"`
	Bandcamp     *string `json:"bandcamp"`
	Soundcloud   *string `json:"soundcloud"`
	Tidal        *string `json:"tidal"`
	AmazonMusic  *string `json:"amazon_music"`
}

type ReleaseQC struct {
	Passed      bool    `json:"passed"`
	CheckedDate *string `json:"checked_date"`
	Notes       string  `json:"notes,omitempty"`
}

// ValidTransitions maps current status to the set of valid next statuses via `release advance`.
// Transitions to submitted and live are handled by dedicated commands only.
var ValidTransitions = map[string]string{
	"in-production": "mastering",
	"mastering":     "qc",
	"qc":            "ready",
}
