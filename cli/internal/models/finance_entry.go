package models

type FinanceEntry struct {
	SchemaVersion string         `json:"schema_version"`
	ID            string         `json:"id"`
	Type          string         `json:"type"`
	Date          string         `json:"date"`
	Amount        float64        `json:"amount"`
	Currency      string         `json:"currency"`
	Source        string         `json:"source"`
	ArtistID      *string        `json:"artist_id"`
	ReleaseID     *string        `json:"release_id"`
	OpportunityID *string        `json:"opportunity_id"`
	Period        *FinancePeriod `json:"period,omitempty"`
	Description   string         `json:"description,omitempty"`
	Notes         string         `json:"notes,omitempty"`
}

type FinancePeriod struct {
	Start string `json:"start"`
	End   string `json:"end"`
}
