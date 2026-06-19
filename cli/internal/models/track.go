package models

type Track struct {
	SchemaVersion   string          `json:"schema_version"`
	ID              string          `json:"id"`
	ReleaseID       string          `json:"release_id"`
	Title           string          `json:"title"`
	Position        int             `json:"position"`
	DurationSeconds *int            `json:"duration_seconds"`
	ISRC            *string         `json:"isrc"`
	ISWC            *string         `json:"iswc"`
	MusicalKey      *string         `json:"musical_key"`
	BPM             *float64        `json:"bpm"`
	Collaborators   []Collaborator  `json:"collaborators,omitempty"`
	Files           *TrackFiles     `json:"files,omitempty"`
	Mastering       *TrackMastering `json:"mastering,omitempty"`
	QC              *TrackQC        `json:"qc,omitempty"`
	Notes           string          `json:"notes,omitempty"`
}

type Collaborator struct {
	Name            string  `json:"name"`
	ArtistID        *string `json:"artist_id"`
	Role            string  `json:"role"`
	SplitPercentage float64 `json:"split_percentage"`
	PRO             *string `json:"pro"`
	IPINumber       *string `json:"ipi_number"`
}

type TrackFiles struct {
	MasterWAV          *string `json:"master_wav"`
	StemsZip           *string `json:"stems_zip"`
	ProjectFile        *string `json:"project_file"`
	MP3320             *string `json:"mp3_320"`
	WAVForDistribution *string `json:"wav_for_distribution"`
}

type TrackMastering struct {
	ProfileID      *string  `json:"profile_id"`
	IntegratedLUFS *float64 `json:"integrated_lufs"`
	TruePeakDBTP   *float64 `json:"true_peak_dbtp"`
	LRA            *float64 `json:"lra"`
	SampleRateHz   *int     `json:"sample_rate_hz"`
	BitDepth       *int     `json:"bit_depth"`
	MasteredDate   *string  `json:"mastered_date"`
	MasteredBy     *string  `json:"mastered_by"`
	Notes          string   `json:"notes,omitempty"`
}

type TrackQC struct {
	Passed      bool        `json:"passed"`
	CheckedDate *string     `json:"checked_date"`
	Failures    []string    `json:"failures,omitempty"`
	Override    *QCOverride `json:"override,omitempty"`
}

type QCOverride struct {
	OverriddenBy string `json:"overridden_by"`
	Date         string `json:"date"`
	Reason       string `json:"reason"`
}

// HasMasteringData returns true when the minimum mastering measurements are populated.
func (t *Track) HasMasteringData() bool {
	if t.Mastering == nil {
		return false
	}
	return t.Mastering.IntegratedLUFS != nil && t.Mastering.TruePeakDBTP != nil
}
