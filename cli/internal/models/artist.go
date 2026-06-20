package models

type Artist struct {
	SchemaVersion  string           `json:"schema_version"`
	ID             string           `json:"id"`
	DisplayName    string           `json:"display_name"`
	LegalName      string           `json:"legal_name,omitempty"`
	AlsoKnownAs    []string         `json:"also_known_as,omitempty"`
	Members        []string         `json:"members,omitempty"`
	Bio            *ArtistBio       `json:"bio,omitempty"`
	GenreTags      []string         `json:"genre_tags,omitempty"`
	Contact        *ArtistContact   `json:"contact,omitempty"`
	Location       string           `json:"location,omitempty"`
	Rights         *ArtistRights    `json:"rights,omitempty"`
	Distribution   *ArtistDistrib   `json:"distribution,omitempty"`
	DefaultLicense string           `json:"default_license"`
	Platforms      *ArtistPlatforms `json:"platforms,omitempty"`
	CreatedDate    string           `json:"created_date,omitempty"`
}

type ArtistBio struct {
	Short  string `json:"short,omitempty"`
	Medium string `json:"medium,omitempty"`
	Full   string `json:"full,omitempty"`
}

type ArtistContact struct {
	Email   string `json:"email,omitempty"`
	Website string `json:"website,omitempty"`
}

type ArtistRights struct {
	PRO       string `json:"pro,omitempty"`
	IPINumber string `json:"ipi_number,omitempty"`
	ISNI      string `json:"isni,omitempty"`
}

type ArtistDistrib struct {
	Distributor         string `json:"distributor,omitempty"`
	DistributorArtistID string `json:"distributor_artist_id,omitempty"`
}

type ArtistPlatforms struct {
	SpotifyArtistID      string `json:"spotify_artist_id,omitempty"`
	AppleMusicArtistID   string `json:"apple_music_artist_id,omitempty"`
	YouTubeChannelID     string `json:"youtube_channel_id,omitempty"`
	YouTubeMusicArtistID string `json:"youtube_music_artist_id,omitempty"`
	BandcampURL          string `json:"bandcamp_url,omitempty"`
	SoundcloudURL        string `json:"soundcloud_url,omitempty"`
	InstagramHandle      string `json:"instagram_handle,omitempty"`
	FacebookURL          string `json:"facebook_url,omitempty"`
	TikTokHandle         string `json:"tiktok_handle,omitempty"`
	SubvertFMURL         string `json:"subvert_fm_url,omitempty"`
	FMAURL               string `json:"fma_url,omitempty"`
}
