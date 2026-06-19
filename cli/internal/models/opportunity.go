package models

type Opportunity struct {
	SchemaVersion   string             `json:"schema_version"`
	ID              string             `json:"id"`
	Type            string             `json:"type"`
	Status          string             `json:"status"`
	ArtistID        string             `json:"artist_id,omitempty"`
	Contact         OpportunityContact `json:"contact"`
	Match           *OpportunityMatch  `json:"match,omitempty"`
	TracksSuggested []string           `json:"tracks_suggested,omitempty"`
	ValueEstimate   *ValueEstimate     `json:"value_estimate,omitempty"`
	OutreachHistory []OutreachAction   `json:"outreach_history,omitempty"`
	FollowUpDue     *string            `json:"follow_up_due"`
	CreatedDate     string             `json:"created_date,omitempty"`
	Notes           string             `json:"notes,omitempty"`
}

type OpportunityContact struct {
	Name          string            `json:"name"`
	Role          string            `json:"role,omitempty"`
	Email         *string           `json:"email"`
	URL           *string           `json:"url"`
	SocialHandles map[string]string `json:"social_handles,omitempty"`
	Notes         string            `json:"notes,omitempty"`
}

type OpportunityMatch struct {
	Score     int    `json:"score"`
	Rationale string `json:"rationale,omitempty"`
}

type ValueEstimate struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Basis    string  `json:"basis,omitempty"`
}

type OutreachAction struct {
	Date        string  `json:"date"`
	Action      string  `json:"action"`
	Note        string  `json:"note,omitempty"`
	ApprovedBy  *string `json:"approved_by"`
	FollowUpDue *string `json:"follow_up_due"`
}
