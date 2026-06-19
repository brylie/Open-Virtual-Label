package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/open-virtual-label/ovl/internal/models"
	"github.com/open-virtual-label/ovl/internal/output"
	"github.com/open-virtual-label/ovl/internal/schema"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
	"github.com/spf13/cobra"
)

var trackCmd = &cobra.Command{
	Use:   "track",
	Short: "Manage tracks within a release",
}

// track add
var trackAddRelease string
var trackAddPosition int
var trackAddCmd = &cobra.Command{
	Use:   "add \"<title>\"",
	Short: "Add a track to a release",
	Args:  cobra.ExactArgs(1),
	RunE:  runTrackAdd,
}

// track show
var trackShowRelease string
var trackShowCmd = &cobra.Command{
	Use:   "show <track-id>",
	Short: "Display a track record",
	Args:  cobra.ExactArgs(1),
	RunE:  runTrackShow,
}

// track set-file
var trackSetFileRelease, trackSetFileField, trackSetFilePath string
var trackSetFileCmd = &cobra.Command{
	Use:   "set-file <track-id>",
	Short: "Set a file path on a track record",
	Args:  cobra.ExactArgs(1),
	RunE:  runTrackSetFile,
}

func init() {
	trackAddCmd.Flags().StringVar(&trackAddRelease, "release", "", "release ID (required)")
	trackAddCmd.Flags().IntVar(&trackAddPosition, "position", 0, "track number (1-based); defaults to next position")
	_ = trackAddCmd.MarkFlagRequired("release")

	trackShowCmd.Flags().StringVar(&trackShowRelease, "release", "", "release ID (disambiguate if track slug exists in multiple releases)")

	trackSetFileCmd.Flags().StringVar(&trackSetFileRelease, "release", "", "release ID")
	trackSetFileCmd.Flags().StringVar(&trackSetFileField, "field", "", "file field: master_wav, stems_zip, project_file, mp3_320, wav_for_distribution (required)")
	trackSetFileCmd.Flags().StringVar(&trackSetFilePath, "path", "", "path to the file (required)")
	_ = trackSetFileCmd.MarkFlagRequired("field")
	_ = trackSetFileCmd.MarkFlagRequired("path")

	trackCmd.AddCommand(trackAddCmd, trackShowCmd, trackSetFileCmd)
	rootCmd.AddCommand(trackCmd)
}

func runTrackAdd(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	title := args[0]
	releaseID := trackAddRelease

	artistID, err := ws.FindRelease(wsPath, releaseID)
	if err != nil {
		return err
	}

	// Determine position
	position := trackAddPosition
	if position == 0 {
		tIDs, _ := ws.ListTrackIDs(wsPath, artistID, releaseID)
		position = len(tIDs) + 1
	}

	slug := ws.ToSlug(title)

	track := models.Track{
		SchemaVersion: "1",
		ID:            slug,
		ReleaseID:     releaseID,
		Title:         title,
		Position:      position,
	}

	raw, _ := json.Marshal(track)
	if errs, err := schema.Validate("track", raw); err != nil {
		return err
	} else if len(errs) > 0 {
		return &ExitError{Code: 3, Msg: "SCHEMA_VALIDATION_FAILED:\n" + strings.Join(errs, "\n")}
	}

	trackPath := ws.TrackFile(wsPath, artistID, releaseID, slug)
	if err := ws.WriteJSON(trackPath, track); err != nil {
		return err
	}

	// Append track ID to release.tracks[]
	relFile := ws.ReleaseFile(wsPath, artistID, releaseID)
	var rel models.Release
	if err := ws.ReadJSON(relFile, &rel); err == nil {
		rel.Tracks = append(rel.Tracks, slug)
		_ = ws.WriteJSON(relFile, rel)
	}

	output.Success("Created %s", trackPath)
	fmt.Printf("Track ID: %s  (position %d)\n", slug, position)
	return nil
}

func runTrackShow(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	trackID := args[0]
	artistID, releaseID, err := ws.FindTrack(wsPath, trackID, trackShowRelease)
	if err != nil {
		return err
	}

	var t models.Track
	if err := ws.ReadJSON(ws.TrackFile(wsPath, artistID, releaseID, trackID), &t); err != nil {
		return err
	}

	if output.IsJSON() {
		return output.JSON(t)
	}

	fmt.Printf("Track: %s (%s)\n", t.Title, t.ID)
	fmt.Printf("Release: %s / %s\n", artistID, releaseID)
	fmt.Printf("Position: %d\n", t.Position)

	isrc := "(not assigned)"
	if t.ISRC != nil && *t.ISRC != "" {
		isrc = *t.ISRC
	}
	fmt.Printf("ISRC: %s\n", isrc)

	if t.Files != nil {
		fmt.Println("\nFiles:")
		printFileField("  master_wav", t.Files.MasterWAV)
		printFileField("  stems_zip", t.Files.StemsZip)
		printFileField("  project_file", t.Files.ProjectFile)
		printFileField("  mp3_320", t.Files.MP3320)
		printFileField("  wav_for_distribution", t.Files.WAVForDistribution)
	}

	if t.Mastering != nil {
		fmt.Println("\nMastering:")
		if t.Mastering.IntegratedLUFS != nil {
			fmt.Printf("  Integrated LUFS: %.1f\n", *t.Mastering.IntegratedLUFS)
		}
		if t.Mastering.TruePeakDBTP != nil {
			fmt.Printf("  True peak: %.1f dBTP\n", *t.Mastering.TruePeakDBTP)
		}
		if t.Mastering.LRA != nil {
			fmt.Printf("  LRA: %.1f LU\n", *t.Mastering.LRA)
		}
	}

	qcStatus := "not run"
	if t.QC != nil {
		if t.QC.Passed {
			qcStatus = "passed"
		} else {
			qcStatus = fmt.Sprintf("failed (%d issue(s))", len(t.QC.Failures))
		}
	}
	fmt.Printf("\nQC: %s\n", qcStatus)
	return nil
}

func runTrackSetFile(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	trackID := args[0]
	artistID, releaseID, err := ws.FindTrack(wsPath, trackID, trackSetFileRelease)
	if err != nil {
		return err
	}

	validFields := map[string]bool{
		fieldMasterWAV: true, fieldStemsZip: true, fieldProjectFile: true,
		fieldMP3320: true, fieldWAVForDistribution: true,
	}
	if !validFields[trackSetFileField] {
		return fmt.Errorf("invalid field %q; valid: master_wav, stems_zip, project_file, mp3_320, wav_for_distribution", trackSetFileField)
	}

	// Validate path exists
	if _, err := os.Stat(trackSetFilePath); err != nil {
		return fmt.Errorf("file not found: %s", trackSetFilePath)
	}

	trackPath := ws.TrackFile(wsPath, artistID, releaseID, trackID)
	var t models.Track
	if err := ws.ReadJSON(trackPath, &t); err != nil {
		return err
	}

	if t.Files == nil {
		t.Files = &models.TrackFiles{}
	}

	p := trackSetFilePath
	switch trackSetFileField {
	case fieldMasterWAV:
		t.Files.MasterWAV = &p
	case fieldStemsZip:
		t.Files.StemsZip = &p
	case fieldProjectFile:
		t.Files.ProjectFile = &p
	case fieldMP3320:
		t.Files.MP3320 = &p
	case fieldWAVForDistribution:
		t.Files.WAVForDistribution = &p
	}

	raw, _ := json.Marshal(t)
	if errs, err := schema.Validate("track", raw); err != nil {
		return err
	} else if len(errs) > 0 {
		return &ExitError{Code: 3, Msg: strings.Join(errs, "\n")}
	}

	if err := ws.WriteJSON(trackPath, t); err != nil {
		return err
	}
	output.Success("Set %s on track %s", trackSetFileField, trackID)
	return nil
}

func printFileField(label string, v *string) {
	if v != nil && *v != "" {
		fmt.Printf("%s: %s\n", label, *v)
	} else {
		fmt.Printf("%s: (not set)\n", label)
	}
}
