package cmd

import (
	"encoding/json"
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

var qcCmd = &cobra.Command{
	Use:   "qc",
	Short: "Quality control checks",
}

var qcCheckRelease string
var qcCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Run the pre-release quality check on a release",
	RunE:  runQCCheck,
}

func init() {
	qcCheckCmd.Flags().StringVar(&qcCheckRelease, "release", "", "release ID (required)")
	_ = qcCheckCmd.MarkFlagRequired("release")

	qcCmd.AddCommand(qcCheckCmd)
	rootCmd.AddCommand(qcCmd)
}

func runQCCheck(_ *cobra.Command, _ []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	releaseID := qcCheckRelease

	artistID, err := ws.FindRelease(wsPath, releaseID)
	if err != nil {
		return err
	}

	var rel models.Release
	relFile := ws.ReleaseFile(wsPath, artistID, releaseID)
	if err := ws.ReadJSON(relFile, &rel); err != nil {
		return err
	}

	fmt.Printf("QC Report — %s\n\n", rel.Title)

	var failures []string

	// Track-level checks
	tIDs, _ := ws.ListTrackIDs(wsPath, artistID, releaseID)
	fmt.Printf("Tracks (%d checked):\n", len(tIDs))
	for _, tID := range tIDs {
		var t models.Track
		if err := ws.ReadJSON(ws.TrackFile(wsPath, artistID, releaseID, tID), &t); err != nil {
			continue
		}
		isrcOK := t.ISRC != nil && *t.ISRC != ""
		masteringOK := t.HasMasteringData()
		masterWAVOK := t.Files != nil && t.Files.MasterWAV != nil && *t.Files.MasterWAV != ""

		trackOK := isrcOK && masteringOK && masterWAVOK
		mark := "✓"
		if !trackOK {
			mark = "✗"
		}

		isrcMark, masteringMark, fileMark := "✓", "✓", "✓"
		if !isrcOK {
			isrcMark = "✗"
			failures = append(failures, fmt.Sprintf("track %s: ISRC not assigned", tID))
		}
		if !masteringOK {
			masteringMark = "✗"
			failures = append(failures, fmt.Sprintf("track %s: mastering data incomplete", tID))
		}
		if !masterWAVOK {
			fileMark = "✗"
			failures = append(failures, fmt.Sprintf("track %s: master_wav not set", tID))
		}

		fmt.Printf("  %s %-24s  ISRC %s  Mastering %s  File %s\n",
			mark, tID, isrcMark, masteringMark, fileMark)
	}
	fmt.Println()

	// Release-level checks
	fmt.Println("Release:")
	artworkOK := rel.Artwork != nil && rel.Artwork.PrimaryFile != "" && rel.Artwork.DimensionsPx >= 3000
	artworkMark := "✓"
	if !artworkOK {
		artworkMark = "✗"
		failures = append(failures, "artwork: missing or dimensions < 3000px")
	}
	fmt.Printf("  %s Artwork\n", artworkMark)

	licenseMark := "✓"
	if rel.License == "" {
		licenseMark = "✗"
		failures = append(failures, "license: not set")
	}
	fmt.Printf("  %s License: %s\n", licenseMark, rel.License)

	dateMark := "✓"
	targetDate := ""
	if rel.Dates != nil {
		targetDate = rel.Dates.TargetRelease
	}
	if targetDate == "" {
		dateMark = "✗"
		failures = append(failures, "dates.target_release: not set")
	}
	fmt.Printf("  %s Target date: %s\n", dateMark, targetDate)

	deadlineMark := "✓"
	deadline := ""
	if rel.Dates != nil {
		deadline = rel.Dates.DistributorSubmissionDeadline
	}
	if deadline == "" {
		deadlineMark = "✗"
		failures = append(failures, "dates.distributor_submission_deadline: not set")
	}
	fmt.Printf("  %s Submission deadline: %s\n", deadlineMark, deadline)
	fmt.Println()

	if len(failures) > 0 {
		fmt.Printf("%d failure(s). Resolve before advancing to ready.\n", len(failures))
		for _, f := range failures {
			fmt.Printf("  - %s\n", f)
		}
		return nil
	}

	fmt.Println("All checks passed.")

	// Confirmation gate to set qc.passed
	ok, err := prompt.Confirm("Set release QC as passed?")
	if err != nil {
		return err
	}
	if !ok {
		return &ExitError{Code: 4, Msg: "QC not confirmed"}
	}

	today := time.Now().Format("2006-01-02")
	rel.QC = &models.ReleaseQC{
		Passed:      true,
		CheckedDate: &today,
	}

	raw, _ := json.Marshal(rel)
	if errs, err := schema.Validate("release", raw); err != nil {
		return err
	} else if len(errs) > 0 {
		return &ExitError{Code: 3, Msg: strings.Join(errs, "\n")}
	}

	if err := ws.WriteJSON(relFile, rel); err != nil {
		return err
	}
	output.Success("QC passed for %s", releaseID)
	return nil
}
