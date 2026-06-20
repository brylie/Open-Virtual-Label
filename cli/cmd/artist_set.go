package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/open-virtual-label/ovl/internal/models"
	"github.com/open-virtual-label/ovl/internal/output"
	"github.com/open-virtual-label/ovl/internal/prompt"
	"github.com/open-virtual-label/ovl/internal/schema"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
	"github.com/spf13/cobra"
)

const (
	platformAppleMusicArtist = "apple-music"
	platformYouTubeArtist    = "youtube"
	platformYouTubeMusicSlug = "youtube-music"
	platformInstagram        = "instagram"
	platformFacebookArtist   = "facebook"
	platformTikTok           = "tiktok"
	platformSubvertFMArtist  = "subvert-fm"
	platformFMAArtist        = "fma"
)

var validArtistPlatforms = map[string]bool{
	platformSpotify:          true,
	platformAppleMusicArtist: true,
	platformYouTubeArtist:    true,
	platformYouTubeMusicSlug: true,
	platformBandcamp:         true,
	platformSoundcloud:       true,
	platformInstagram:        true,
	platformFacebookArtist:   true,
	platformTikTok:           true,
	platformSubvertFMArtist:  true,
	platformFMAArtist:        true,
}

var (
	artistSetBioShort  string
	artistSetBioMedium string
	artistSetBioFull   string
)

var artistSetBioCmd = &cobra.Command{
	Use:   "set-bio <artist-id>",
	Short: "Set or update the bio for an artist",
	Args:  cobra.ExactArgs(1),
	RunE:  runArtistSetBio,
}

var (
	artistSetContactEmail   string
	artistSetContactWebsite string
)

var artistSetContactCmd = &cobra.Command{
	Use:   "set-contact <artist-id>",
	Short: "Update contact details for an artist",
	Args:  cobra.ExactArgs(1),
	RunE:  runArtistSetContact,
}

var (
	artistSetPlatformName  string
	artistSetPlatformValue string
)

var artistSetPlatformCmd = &cobra.Command{
	Use:   "set-platform <artist-id>",
	Short: "Set a platform identifier or URL on an artist profile",
	Args:  cobra.ExactArgs(1),
	RunE:  runArtistSetPlatform,
}

var (
	artistSetRightsPRO  string
	artistSetRightsIPI  string
	artistSetRightsISNI string
)

var artistSetRightsCmd = &cobra.Command{
	Use:   "set-rights <artist-id>",
	Short: "Update performing rights registration details for an artist",
	Args:  cobra.ExactArgs(1),
	RunE:  runArtistSetRights,
}

var (
	artistSetDistributorName string
	artistSetDistributorID   string
)

var artistSetDistributorCmd = &cobra.Command{
	Use:   "set-distributor <artist-id>",
	Short: "Set distribution details for an artist",
	Args:  cobra.ExactArgs(1),
	RunE:  runArtistSetDistributor,
}

var artistSetLicenseCmd = &cobra.Command{
	Use:   "set-license <artist-id> <license>",
	Short: "Set the default license for an artist's releases",
	Args:  cobra.ExactArgs(2),
	RunE:  runArtistSetLicense,
}

var artistSetLocationCmd = &cobra.Command{
	Use:   "set-location <artist-id> <location>",
	Short: "Set the artist's location",
	Args:  cobra.ExactArgs(2),
	RunE:  runArtistSetLocation,
}

func init() {
	artistSetBioCmd.Flags().StringVar(&artistSetBioShort, "short", "", "one or two sentences (max 280 characters)")
	artistSetBioCmd.Flags().StringVar(&artistSetBioMedium, "medium", "", "one paragraph (max 1024 characters)")
	artistSetBioCmd.Flags().StringVar(&artistSetBioFull, "full", "", "full bio")

	artistSetContactCmd.Flags().StringVar(&artistSetContactEmail, "email", "", "public or licensing contact email")
	artistSetContactCmd.Flags().StringVar(&artistSetContactWebsite, "website", "", "artist website URL")

	artistSetPlatformCmd.Flags().StringVar(&artistSetPlatformName, "platform", "", "platform name (required)")
	artistSetPlatformCmd.Flags().StringVar(&artistSetPlatformValue, "value", "", "platform value (required)")
	_ = artistSetPlatformCmd.MarkFlagRequired("platform")
	_ = artistSetPlatformCmd.MarkFlagRequired("value")

	artistSetRightsCmd.Flags().StringVar(&artistSetRightsPRO, "pro", "", "performing rights organization")
	artistSetRightsCmd.Flags().StringVar(&artistSetRightsIPI, "ipi", "", "IPI number assigned by the PRO")
	artistSetRightsCmd.Flags().StringVar(&artistSetRightsISNI, "isni", "", "International Standard Name Identifier")

	artistSetDistributorCmd.Flags().StringVar(&artistSetDistributorName, "distributor", "", "distribution platform (required)")
	artistSetDistributorCmd.Flags().StringVar(&artistSetDistributorID, "id", "", "artist identifier within the distributor's system")
	_ = artistSetDistributorCmd.MarkFlagRequired("distributor")

	artistCmd.AddCommand(artistSetBioCmd)
	artistCmd.AddCommand(artistSetContactCmd)
	artistCmd.AddCommand(artistSetPlatformCmd)
	artistCmd.AddCommand(artistSetRightsCmd)
	artistCmd.AddCommand(artistSetDistributorCmd)
	artistCmd.AddCommand(artistSetLicenseCmd)
	artistCmd.AddCommand(artistSetLocationCmd)
}

func readArtist(wsPath, artistID string) (models.Artist, error) {
	var a models.Artist
	if err := ws.ReadJSON(ws.ArtistFile(wsPath, artistID), &a); err != nil {
		return a, fmt.Errorf("ARTIST_NOT_FOUND: artist %q not found", artistID)
	}
	return a, nil
}

func saveArtist(wsPath string, a *models.Artist) error {
	raw, err := json.Marshal(a)
	if err != nil {
		return err
	}
	errs, err := schema.Validate("artist", raw)
	if err != nil {
		return err
	}
	if len(errs) > 0 {
		return &ExitError{Code: 3, Msg: "VALIDATION_FAILED: " + strings.Join(errs, "; ")}
	}
	return ws.WriteJSON(ws.ArtistFile(wsPath, a.ID), a)
}

func promptArtistBio(bio *models.ArtistBio) error {
	short, err := prompt.Ask("Short bio (max 280 chars)", bio.Short)
	if err != nil {
		return err
	}
	medium, err := prompt.Ask("Medium bio (max 1024 chars)", bio.Medium)
	if err != nil {
		return err
	}
	full, err := prompt.Ask("Full bio", bio.Full)
	if err != nil {
		return err
	}
	bio.Short, bio.Medium, bio.Full = short, medium, full
	return nil
}

func applyArtistBioFlags(cmd *cobra.Command, bio *models.ArtistBio) []string {
	var updated []string
	if cmd.Flags().Changed("short") {
		bio.Short = artistSetBioShort
		updated = append(updated, "short")
	}
	if cmd.Flags().Changed("medium") {
		bio.Medium = artistSetBioMedium
		updated = append(updated, "medium")
	}
	if cmd.Flags().Changed("full") {
		bio.Full = artistSetBioFull
		updated = append(updated, "full")
	}
	return updated
}

func runArtistSetBio(cmd *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	a, err := readArtist(wsPath, args[0])
	if err != nil {
		return err
	}
	if a.Bio == nil {
		a.Bio = &models.ArtistBio{}
	}

	var updated []string
	if !cmd.Flags().Changed("short") && !cmd.Flags().Changed("medium") && !cmd.Flags().Changed("full") {
		if err := promptArtistBio(a.Bio); err != nil {
			return err
		}
		updated = []string{"short", "medium", "full"}
	} else {
		updated = applyArtistBioFlags(cmd, a.Bio)
	}

	if err := saveArtist(wsPath, &a); err != nil {
		return err
	}
	output.Success("Updated bio (%s) for %s", strings.Join(updated, ", "), a.ID)
	return nil
}

func runArtistSetContact(cmd *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	a, err := readArtist(wsPath, args[0])
	if err != nil {
		return err
	}
	if a.Contact == nil {
		a.Contact = &models.ArtistContact{}
	}
	if cmd.Flags().Changed("email") {
		a.Contact.Email = artistSetContactEmail
	}
	if cmd.Flags().Changed("website") {
		a.Contact.Website = artistSetContactWebsite
	}

	if err := saveArtist(wsPath, &a); err != nil {
		return err
	}
	output.Success("Updated contact details for %s", a.ID)
	return nil
}

func runArtistSetPlatform(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	a, err := readArtist(wsPath, args[0])
	if err != nil {
		return err
	}

	if !validArtistPlatforms[artistSetPlatformName] {
		return &ExitError{Code: 1, Msg: fmt.Sprintf("UNKNOWN_PLATFORM: %q is not a valid platform", artistSetPlatformName)}
	}
	if a.Platforms == nil {
		a.Platforms = &models.ArtistPlatforms{}
	}
	setArtistPlatformField(a.Platforms, artistSetPlatformName, artistSetPlatformValue)

	if err := saveArtist(wsPath, &a); err != nil {
		return err
	}
	output.Success("Set %s = %q for %s", artistSetPlatformName, artistSetPlatformValue, a.ID)
	return nil
}

func setArtistPlatformField(p *models.ArtistPlatforms, platform, value string) {
	switch platform {
	case platformSpotify:
		p.SpotifyArtistID = value
	case platformAppleMusicArtist:
		p.AppleMusicArtistID = value
	case platformYouTubeArtist:
		p.YouTubeChannelID = value
	case platformYouTubeMusicSlug:
		p.YouTubeMusicArtistID = value
	case platformBandcamp:
		p.BandcampURL = value
	case platformSoundcloud:
		p.SoundcloudURL = value
	case platformInstagram:
		p.InstagramHandle = value
	case platformFacebookArtist:
		p.FacebookURL = value
	case platformTikTok:
		p.TikTokHandle = value
	case platformSubvertFMArtist:
		p.SubvertFMURL = value
	case platformFMAArtist:
		p.FMAURL = value
	}
}

func runArtistSetRights(cmd *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	a, err := readArtist(wsPath, args[0])
	if err != nil {
		return err
	}
	if a.Rights == nil {
		a.Rights = &models.ArtistRights{}
	}
	if cmd.Flags().Changed("pro") {
		a.Rights.PRO = artistSetRightsPRO
	}
	if cmd.Flags().Changed("ipi") {
		a.Rights.IPINumber = artistSetRightsIPI
	}
	if cmd.Flags().Changed("isni") {
		a.Rights.ISNI = artistSetRightsISNI
	}

	if err := saveArtist(wsPath, &a); err != nil {
		return err
	}
	output.Success("Updated rights details for %s", a.ID)
	return nil
}

func runArtistSetDistributor(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	a, err := readArtist(wsPath, args[0])
	if err != nil {
		return err
	}
	a.Distribution = &models.ArtistDistrib{
		Distributor:         artistSetDistributorName,
		DistributorArtistID: artistSetDistributorID,
	}

	if err := saveArtist(wsPath, &a); err != nil {
		return err
	}
	output.Success("Set distributor to %q for %s", artistSetDistributorName, a.ID)
	return nil
}

func runArtistSetLicense(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	a, err := readArtist(wsPath, args[0])
	if err != nil {
		return err
	}
	license := args[1]
	if license == "" {
		return errors.New("VALIDATION_FAILED: license is required")
	}
	a.DefaultLicense = license

	if err := saveArtist(wsPath, &a); err != nil {
		return err
	}
	output.Success("Set default license to %q for %s", license, a.ID)
	return nil
}

func runArtistSetLocation(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	a, err := readArtist(wsPath, args[0])
	if err != nil {
		return err
	}
	a.Location = args[1]

	if err := saveArtist(wsPath, &a); err != nil {
		return err
	}
	output.Success("Set location to %q for %s", a.Location, a.ID)
	return nil
}
