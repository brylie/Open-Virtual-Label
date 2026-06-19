package models

type MasteringProfile struct {
	SchemaVersion string            `json:"schema_version"`
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Description   string            `json:"description,omitempty"`
	Targets       MasteringTargets  `json:"targets"`
	PlatformNotes map[string]string `json:"platform_notes,omitempty"`
	Checklist     []string          `json:"checklist,omitempty"`
	SessionNotes  []SessionNote     `json:"session_notes,omitempty"`
	CreatedDate   string            `json:"created_date,omitempty"`
	LastUpdated   string            `json:"last_updated,omitempty"`
}

type MasteringTargets struct {
	IntegratedLUFS LUFSRange `json:"integrated_lufs"`
	TruePeakDBTP   float64   `json:"true_peak_dbtp"`
	LRAMin         *float64  `json:"lra_min"`
	LRAMax         *float64  `json:"lra_max"`
	SampleRateHz   int       `json:"sample_rate_hz"`
	BitDepth       int       `json:"bit_depth"`
}

type LUFSRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type SessionNote struct {
	Date    string `json:"date"`
	TrackID string `json:"track_id"`
	Note    string `json:"note"`
}
