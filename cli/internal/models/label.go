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
	CreatedDate        string        `json:"created_date,omitempty"`
}

type LabelContact struct {
	Email    string `json:"email,omitempty"`
	Website  string `json:"website,omitempty"`
	Location string `json:"location,omitempty"`
}
