package cmd

// Pipeline status strings used in release transitions and QC checks.
const (
	statusInProduction = "in-production"
	statusMastering    = "mastering"
	statusQC           = "qc"
	statusReady        = "ready"
	statusSubmitted    = "submitted"
	statusLive         = "live"
)

// msgCanceled is the message used when a user declines a confirmation prompt.
const msgCanceled = "canceled"

// cmdUseList is the Use field for all "list" subcommands.
const cmdUseList = "list"

// cmdUseDraft is the Use field for "draft" subcommands.
const cmdUseDraft = "draft"

// Track file field names used in set-file validation and assignment.
const (
	fieldMasterWAV          = "master_wav"
	fieldStemsZip           = "stems_zip"
	fieldProjectFile        = "project_file"
	fieldMP3320             = "mp3_320"
	fieldWAVForDistribution = "wav_for_distribution"
)

// Release link platform identifiers.
const (
	platformSpotify     = "spotify"
	platformAppleMusic  = "apple_music"
	platformYouTubeMusic = "youtube_music"
	platformBandcamp    = "bandcamp"
	platformSoundcloud  = "soundcloud"
	platformTidal       = "tidal"
	platformAmazonMusic = "amazon_music"
)

// MCP integration identifiers.
const (
	mcpInternetArchive = "internet-archive"
	mcpAmuse           = "amuse"
)

// Finance entry types.
const (
	financeTypeRevenue = "revenue"
	financeTypeExpense = "expense"
)
