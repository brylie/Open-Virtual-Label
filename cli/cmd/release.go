package cmd

import (
	"encoding/json"
	"errors"
	"strconv"
	"fmt"
	"strings"
	"time"

	"github.com/open-virtual-label/ovl/internal/models"
	"github.com/open-virtual-label/ovl/internal/output"
	"github.com/open-virtual-label/ovl/internal/prompt"
	"github.com/open-virtual-label/ovl/internal/schema"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
	"github.com/spf13/cobra"
)

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Manage releases",
}

// release create
var releaseCreateType string
var releaseCreateCmd = &cobra.Command{
	Use:   "create \"<title>\"",
	Short: "Create a new release record",
	Args:  cobra.ExactArgs(1),
	RunE:  runReleaseCreate,
}

// release list
var releaseListStatus string
var releaseListCmd = &cobra.Command{
	Use:   cmdUseList,
	Short: "List all releases",
	RunE:  runReleaseList,
}

// release show
var releaseShowCmd = &cobra.Command{
	Use:   "show <release-id>",
	Short: "Display a release and its tracklist",
	Args:  cobra.ExactArgs(1),
	RunE:  runReleaseShow,
}

// release advance
var releaseAdvanceStatus string
var releaseAdvanceCmd = &cobra.Command{
	Use:   "advance <release-id>",
	Short: "Advance a release to the next pipeline status",
	Args:  cobra.ExactArgs(1),
	RunE:  runReleaseAdvance,
}

// release set-profile
var releaseSetProfileID string
var releaseSetProfileCmd = &cobra.Command{
	Use:   "set-profile <release-id>",
	Short: "Assign a mastering profile to a release",
	Args:  cobra.ExactArgs(1),
	RunE:  runReleaseSetProfile,
}

// release set-live
var releaseSetLiveDate string
var releaseSetLiveCmd = &cobra.Command{
	Use:   "set-live <release-id>",
	Short: "Mark a submitted release as live",
	Args:  cobra.ExactArgs(1),
	RunE:  runReleaseSetLive,
}

// release add-link
var releaseAddLinkPlatform, releaseAddLinkURL string
var releaseAddLinkCmd = &cobra.Command{
	Use:   "add-link <release-id>",
	Short: "Add a store link to a live release",
	Args:  cobra.ExactArgs(1),
	RunE:  runReleaseAddLink,
}

// release submit
var releaseSubmitDistributor string
var releaseSubmitCmd = &cobra.Command{
	Use:   "submit <release-id>",
	Short: "Prepare and submit a distribution package",
	Args:  cobra.ExactArgs(1),
	RunE:  runReleaseSubmit,
}

func init() {
	releaseCreateCmd.Flags().StringVar(&releaseCreateType, "type", "album", "release type: album, ep, single, compilation")

	releaseListCmd.Flags().StringVar(&releaseListStatus, "status", "", "filter by status")

	releaseAdvanceCmd.Flags().StringVar(&releaseAdvanceStatus, "status", "", "target status (required)")
	_ = releaseAdvanceCmd.MarkFlagRequired("status")

	releaseSetProfileCmd.Flags().StringVar(&releaseSetProfileID, "profile", "", "mastering profile ID (required)")
	_ = releaseSetProfileCmd.MarkFlagRequired("profile")

	releaseSetLiveCmd.Flags().StringVar(&releaseSetLiveDate, "date", "", "live date YYYY-MM-DD (required)")
	_ = releaseSetLiveCmd.MarkFlagRequired("date")

	releaseAddLinkCmd.Flags().StringVar(&releaseAddLinkPlatform, "platform", "", "platform name (required)")
	releaseAddLinkCmd.Flags().StringVar(&releaseAddLinkURL, "url", "", "store URL (required)")
	_ = releaseAddLinkCmd.MarkFlagRequired("platform")
	_ = releaseAddLinkCmd.MarkFlagRequired("url")

	releaseSubmitCmd.Flags().StringVar(&releaseSubmitDistributor, "distributor", "", "distributor name (required)")
	_ = releaseSubmitCmd.MarkFlagRequired("distributor")

	releaseCmd.AddCommand(releaseCreateCmd, releaseListCmd, releaseShowCmd,
		releaseAdvanceCmd, releaseSetProfileCmd, releaseSetLiveCmd,
		releaseAddLinkCmd, releaseSubmitCmd)
	rootCmd.AddCommand(releaseCmd)
}

func runReleaseCreate(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}

	title := args[0]
	artistID, err := resolveArtistID(wsPath)
	if err != nil {
		return err
	}

	targetDate, err := prompt.Ask("Target release date (YYYY-MM-DD, leave blank to skip)", "")
	if err != nil {
		return err
	}
	deadline, err := prompt.Ask("Distributor submission deadline (YYYY-MM-DD, leave blank to skip)", "")
	if err != nil {
		return err
	}

	// Inherit license from artist
	var artist models.Artist
	_ = ws.ReadJSON(ws.ArtistFile(wsPath, artistID), &artist)
	license := artist.DefaultLicense
	if license == "" {
		license, err = prompt.Ask("License", "CC BY 4.0")
		if err != nil {
			return err
		}
	}

	slug := ws.ToSlug(title)
	today := time.Now().Format("2006-01-02")

	rel := models.Release{
		SchemaVersion: "1",
		ID:            slug,
		Title:         title,
		ArtistID:      artistID,
		ReleaseType:   releaseCreateType,
		Status:        statusInProduction,
		License:       license,
		CreatedDate:   today,
	}
	if targetDate != "" || deadline != "" {
		rel.Dates = &models.ReleaseDates{
			TargetRelease:                 targetDate,
			DistributorSubmissionDeadline: deadline,
		}
	}

	raw, _ := json.Marshal(rel)
	if errs, err := schema.Validate("release", raw); err != nil {
		return err
	} else if len(errs) > 0 {
		return &ExitError{Code: 3, Msg: "SCHEMA_VALIDATION_FAILED:\n" + strings.Join(errs, "\n")}
	}

	relPath := ws.ReleaseFile(wsPath, artistID, slug)
	tracksDir := ws.TracksDir(wsPath, artistID, slug)
	if err := ws.WriteJSON(relPath, rel); err != nil {
		return err
	}
	if err := osCreateDir(tracksDir); err != nil {
		return err
	}

	output.Success("Created %s", relPath)
	fmt.Printf("Release ID: %s\n", slug)
	return nil
}

func runReleaseList(_ *cobra.Command, _ []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	artistIDs, err := ws.ListArtistIDs(wsPath)
	if err != nil {
		return err
	}

	var releases []models.Release
	for _, aID := range artistIDs {
		if cfg.ArtistID != "" && aID != cfg.ArtistID {
			continue
		}
		rIDs, _ := ws.ListReleaseIDs(wsPath, aID)
		for _, rID := range rIDs {
			var rel models.Release
			if err := ws.ReadJSON(ws.ReleaseFile(wsPath, aID, rID), &rel); err == nil {
				if releaseListStatus == "" || rel.Status == releaseListStatus {
					releases = append(releases, rel)
				}
			}
		}
	}

	if output.IsJSON() {
		return output.JSON(releases)
	}

	if len(releases) == 0 {
		fmt.Println("No releases found.")
		return nil
	}

	rows := make([][]string, len(releases))
	for i := range releases {
		target := ""
		if releases[i].Dates != nil {
			target = releases[i].Dates.TargetRelease
		}
		rows[i] = []string{releases[i].ID, releases[i].Title, releases[i].ArtistID, releases[i].ReleaseType, releases[i].Status, target}
	}
	output.Table([]string{"ID", "Title", "Artist", "Type", "Status", "Target Date"}, rows)
	return nil
}

func runReleaseShow(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	releaseID := args[0]
	artistID, err := ws.FindRelease(wsPath, releaseID)
	if err != nil {
		return err
	}

	var rel models.Release
	if err := ws.ReadJSON(ws.ReleaseFile(wsPath, artistID, releaseID), &rel); err != nil {
		return err
	}

	if output.IsJSON() {
		return output.JSON(rel)
	}

	fmt.Printf("Release: %s (%s)\n", rel.Title, rel.ID)
	fmt.Printf("Artist:  %s\n", rel.ArtistID)
	fmt.Printf("Type:    %s\n", rel.ReleaseType)
	fmt.Printf("Status:  %s\n", rel.Status)
	fmt.Printf("License: %s\n", rel.License)
	if rel.Dates != nil && rel.Dates.TargetRelease != "" {
		fmt.Printf("Target date: %s\n", rel.Dates.TargetRelease)
	}
	if rel.MasteringProfileID != "" {
		fmt.Printf("Mastering profile: %s\n", rel.MasteringProfileID)
	}
	fmt.Println()

	// Tracks
	tIDs, _ := ws.ListTrackIDs(wsPath, artistID, releaseID)
	if len(tIDs) == 0 {
		fmt.Println("Tracks: (none — use 'ovl track add' to add tracks)")
		return nil
	}

	rows := make([][]string, 0, len(tIDs))
	for _, tID := range tIDs {
		var t models.Track
		if err := ws.ReadJSON(ws.TrackFile(wsPath, artistID, releaseID, tID), &t); err != nil {
			continue
		}
		isrc := "-"
		if t.ISRC != nil && *t.ISRC != "" {
			isrc = *t.ISRC
		}
		mastering := "✗"
		if t.HasMasteringData() {
			mastering = "✓"
		}
		qc := "✗"
		if t.QC != nil && t.QC.Passed {
			qc = "✓"
		}
		rows = append(rows, []string{
			strconv.Itoa(t.Position), tID, t.Title, isrc, mastering, qc,
		})
	}
	output.Table([]string{"#", "ID", "Title", "ISRC", "Mastering", "QC"}, rows)
	return nil
}

func runReleaseAdvance(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	releaseID := args[0]
	artistID, err := ws.FindRelease(wsPath, releaseID)
	if err != nil {
		return err
	}

	var rel models.Release
	relFile := ws.ReleaseFile(wsPath, artistID, releaseID)
	if err := ws.ReadJSON(relFile, &rel); err != nil {
		return err
	}

	// Validate transition
	allowed, ok := models.ValidTransitions[rel.Status]
	if !ok {
		return fmt.Errorf("cannot advance from status %q via this command", rel.Status)
	}
	if releaseAdvanceStatus != allowed {
		return fmt.Errorf("invalid transition %q→%q; expected %q→%q",
			rel.Status, releaseAdvanceStatus, rel.Status, allowed)
	}

	// Pre-condition checks
	tIDs, _ := ws.ListTrackIDs(wsPath, artistID, releaseID)
	switch rel.Status {
	case statusInProduction:
		if len(tIDs) == 0 {
			return errors.New("cannot advance to mastering: no tracks exist")
		}
	case statusMastering:
		for _, tID := range tIDs {
			var t models.Track
			if err := ws.ReadJSON(ws.TrackFile(wsPath, artistID, releaseID, tID), &t); err != nil {
				continue
			}
			if !t.HasMasteringData() {
				return fmt.Errorf("cannot advance to qc: track %q is missing mastering data", tID)
			}
		}
	case statusQC:
		if rel.QC == nil || !rel.QC.Passed {
			return errors.New("cannot advance to ready: release QC has not passed (run 'ovl qc check')")
		}
	}

	if !cfg.SkipConfirm {
		ok, err := prompt.Confirm(fmt.Sprintf("Advance %s from %q to %q?", releaseID, rel.Status, releaseAdvanceStatus))
		if err != nil {
			return err
		}
		if !ok {
			return &ExitError{Code: 4, Msg: msgCanceled}
		}
	}

	rel.Status = releaseAdvanceStatus
	raw, _ := json.Marshal(rel)
	if errs, err := schema.Validate("release", raw); err != nil {
		return err
	} else if len(errs) > 0 {
		return &ExitError{Code: 3, Msg: strings.Join(errs, "\n")}
	}
	if err := ws.WriteJSON(relFile, rel); err != nil {
		return err
	}
	output.Success("Release %s advanced to %q", releaseID, rel.Status)
	return nil
}

func runReleaseSetProfile(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	releaseID := args[0]
	artistID, err := ws.FindRelease(wsPath, releaseID)
	if err != nil {
		return err
	}

	profilePath := ws.MasteringProfileFile(wsPath, artistID, releaseSetProfileID)
	if _, err := statFile(profilePath); err != nil {
		return fmt.Errorf("mastering profile %q not found for artist %s", releaseSetProfileID, artistID)
	}

	relFile := ws.ReleaseFile(wsPath, artistID, releaseID)
	var rel models.Release
	if err := ws.ReadJSON(relFile, &rel); err != nil {
		return err
	}
	rel.MasteringProfileID = releaseSetProfileID

	if err := ws.WriteJSON(relFile, rel); err != nil {
		return err
	}
	output.Success("Mastering profile set to %q on %s", releaseSetProfileID, releaseID)
	return nil
}

func runReleaseSetLive(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	releaseID := args[0]
	artistID, err := ws.FindRelease(wsPath, releaseID)
	if err != nil {
		return err
	}

	relFile := ws.ReleaseFile(wsPath, artistID, releaseID)
	var rel models.Release
	if err := ws.ReadJSON(relFile, &rel); err != nil {
		return err
	}
	if rel.Status != statusSubmitted {
		return fmt.Errorf("release must be in 'submitted' status; current status: %q", rel.Status)
	}

	if !cfg.SkipConfirm {
		ok, err := prompt.Confirm(fmt.Sprintf("Mark %s as live on %s?", releaseID, releaseSetLiveDate))
		if err != nil {
			return err
		}
		if !ok {
			return &ExitError{Code: 4, Msg: msgCanceled}
		}
	}

	rel.Status = statusLive
	if rel.Dates == nil {
		rel.Dates = &models.ReleaseDates{}
	}
	rel.Dates.Released = releaseSetLiveDate

	if err := ws.WriteJSON(relFile, rel); err != nil {
		return err
	}
	output.Success("Release %s is now live (released: %s)", releaseID, releaseSetLiveDate)
	return nil
}

func runReleaseAddLink(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	releaseID := args[0]
	artistID, err := ws.FindRelease(wsPath, releaseID)
	if err != nil {
		return err
	}

	validPlatforms := map[string]bool{
		platformSpotify: true, platformAppleMusic: true, platformYouTubeMusic: true,
		platformBandcamp: true, platformSoundcloud: true, platformTidal: true, platformAmazonMusic: true,
	}
	if !validPlatforms[releaseAddLinkPlatform] {
		return fmt.Errorf("unknown platform %q; valid: spotify, apple_music, youtube_music, bandcamp, soundcloud, tidal, amazon_music", releaseAddLinkPlatform)
	}

	relFile := ws.ReleaseFile(wsPath, artistID, releaseID)
	var rel models.Release
	if err := ws.ReadJSON(relFile, &rel); err != nil {
		return err
	}
	if rel.StoreLinks == nil {
		rel.StoreLinks = &models.StoreLinks{}
	}

	url := releaseAddLinkURL
	switch releaseAddLinkPlatform {
	case platformSpotify:
		rel.StoreLinks.Spotify = &url
	case platformAppleMusic:
		rel.StoreLinks.AppleMusic = &url
	case platformYouTubeMusic:
		rel.StoreLinks.YouTubeMusic = &url
	case platformBandcamp:
		rel.StoreLinks.Bandcamp = &url
	case platformSoundcloud:
		rel.StoreLinks.Soundcloud = &url
	case platformTidal:
		rel.StoreLinks.Tidal = &url
	case platformAmazonMusic:
		rel.StoreLinks.AmazonMusic = &url
	}

	if err := ws.WriteJSON(relFile, rel); err != nil {
		return err
	}
	output.Success("Added %s link to %s", releaseAddLinkPlatform, releaseID)
	return nil
}

func runReleaseSubmit(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	releaseID := args[0]
	artistID, err := ws.FindRelease(wsPath, releaseID)
	if err != nil {
		return err
	}

	var rel models.Release
	if err := ws.ReadJSON(ws.ReleaseFile(wsPath, artistID, releaseID), &rel); err != nil {
		return err
	}
	if rel.Status != statusReady {
		return fmt.Errorf("QC_NOT_PASSED: release must be in 'ready' status; current: %q", rel.Status)
	}

	// Build summary
	tIDs, _ := ws.ListTrackIDs(wsPath, artistID, releaseID)
	fmt.Printf("\n=== Distribution Summary: %s ===\n\n", rel.Title)
	fmt.Printf("Title:       %s\n", rel.Title)
	fmt.Printf("Type:        %s\n", rel.ReleaseType)
	fmt.Printf("Artist:      %s\n", rel.ArtistID)
	fmt.Printf("License:     %s\n", rel.License)
	fmt.Printf("Distributor: %s\n", releaseSubmitDistributor)
	if rel.Dates != nil && rel.Dates.TargetRelease != "" {
		fmt.Printf("Target date: %s\n", rel.Dates.TargetRelease)
	}
	fmt.Printf("\nTracks (%d):\n", len(tIDs))
	for _, tID := range tIDs {
		var t models.Track
		if err := ws.ReadJSON(ws.TrackFile(wsPath, artistID, releaseID, tID), &t); err != nil {
			continue
		}
		isrc := "(no ISRC)"
		if t.ISRC != nil && *t.ISRC != "" {
			isrc = *t.ISRC
		}
		dur := ""
		if t.DurationSeconds != nil {
			m, s := *t.DurationSeconds/60, *t.DurationSeconds%60
			dur = fmt.Sprintf(" (%d:%02d)", m, s)
		}
		fmt.Printf("  %d. %s — %s%s\n", t.Position, t.Title, isrc, dur)
	}

	// Critical gate — cannot be bypassed with --yes
	fmt.Println()
	ok, err := prompt.Confirm("Submit this release? (this gate cannot be bypassed with --yes)")
	if err != nil {
		return err
	}
	if !ok {
		return &ExitError{Code: 4, Msg: "submission canceled"}
	}

	relFile := ws.ReleaseFile(wsPath, artistID, releaseID)
	rel.Status = statusSubmitted
	today := time.Now().Format("2006-01-02")
	if rel.Dates == nil {
		rel.Dates = &models.ReleaseDates{}
	}
	rel.Dates.Submitted = today
	if rel.Distribution == nil {
		rel.Distribution = &models.ReleaseDistrib{}
	}
	rel.Distribution.Distributor = releaseSubmitDistributor

	if err := ws.WriteJSON(relFile, rel); err != nil {
		return err
	}
	output.Success("Release %s submitted to %s on %s", releaseID, releaseSubmitDistributor, today)
	return nil
}

// resolveArtistID returns the artist ID from --artist flag, env, or prompts when ambiguous.
func resolveArtistID(wsPath string) (string, error) {
	if cfg.ArtistID != "" {
		return cfg.ArtistID, nil
	}
	ids, err := ws.ListArtistIDs(wsPath)
	if err != nil {
		return "", err
	}
	if len(ids) == 0 {
		return "", errors.New("no artists found; run 'ovl artist create' first")
	}
	if len(ids) == 1 {
		return ids[0], nil
	}
	return prompt.Select("Select artist", ids)
}

