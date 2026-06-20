package models

type Label struct {
	SchemaVersion      string        `json:"schema_version"`
	ID                 string        `json:"id"`
	Name               string        `json:"name"`
	Description        string        `json:"description,omitempty"`
	Contact            *LabelContact `json:"contact,omitempty"`
	DefaultLicense     string        `json:"default_license"`
	DefaultDistributor string        `json:"default_distributor,omitempty"`
	StyleguidePath     string        `json:"styleguide_path,omitempty"`
	Sites              []LabelSite   `json:"sites,omitempty"`
	CreatedDate        string        `json:"created_date,omitempty"`
}

// LabelSite represents a website target for ovl site sync.
type LabelSite struct {
	ID          string `json:"id"`
	Path        string `json:"path"`
	ArtistID    string `json:"artist_id,omitempty"`
	ArtistsDir  string `json:"artists_dir,omitempty"`
	ReleasesDir string `json:"releases_dir,omitempty"`
	Description string `json:"description,omitempty"`
}

// ResolvedArtistsDir returns the content directory for artist JSON files,
// falling back to the conventional default if not explicitly set.
func (s LabelSite) ResolvedArtistsDir() string {
	if s.ArtistsDir != "" {
		return s.ArtistsDir
	}
	return "src/content/artists"
}

// ResolvedReleasesDir returns the content directory for release JSON files,
// falling back to the conventional default if not explicitly set.
func (s LabelSite) ResolvedReleasesDir() string {
	if s.ReleasesDir != "" {
		return s.ReleasesDir
	}
	return "src/content/releases"
}

type LabelContact struct {
	Email    string `json:"email,omitempty"`
	Website  string `json:"website,omitempty"`
	Location string `json:"location,omitempty"`
}
