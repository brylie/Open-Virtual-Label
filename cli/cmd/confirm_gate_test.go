package cmd

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/open-virtual-label/ovl/internal/models"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
)

// --- runReleaseSubmit happy path ---

func TestRunReleaseSubmit_Success(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	rel := minimalRelease("album-1", "artist-a")
	rel.Status = statusReady
	writeRelease(t, wsPath, "artist-a", "album-1", rel)

	releaseSubmitDistributor = "amuse"
	t.Cleanup(func() { releaseSubmitDistributor = "" })

	injectPromptInput(t, "y\n")

	got := capStdout(t, func() {
		if err := runReleaseSubmit(nil, []string{"album-1"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "amuse") {
		t.Errorf("expected distributor in output, got %q", got)
	}
}

func TestRunReleaseSubmit_Declined(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	rel := minimalRelease("album-1", "artist-a")
	rel.Status = statusReady
	writeRelease(t, wsPath, "artist-a", "album-1", rel)

	releaseSubmitDistributor = "amuse"
	t.Cleanup(func() { releaseSubmitDistributor = "" })

	injectPromptInput(t, "n\n")

	err := runReleaseSubmit(nil, []string{"album-1"})
	var e *ExitError
	if !errors.As(err, &e) {
		t.Errorf("expected ExitError, got %v", err)
	} else if e.Code != 4 {
		t.Errorf("expected exit code 4 (canceled), got %d", e.Code)
	}
}

func TestRunReleaseSubmit_WithTracks(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	rel := minimalRelease("album-1", "artist-a")
	rel.Status = statusReady
	rel.Dates = &models.ReleaseDates{TargetRelease: "2025-06-01"}
	writeRelease(t, wsPath, "artist-a", "album-1", rel)

	isrc := "GB-EMI-21-00001"
	dur := 240
	track := models.Track{
		SchemaVersion: "1", ID: "track-01", ReleaseID: "album-1",
		Title: "My Track", Position: 1,
		ISRC: &isrc, DurationSeconds: &dur,
	}
	if err := ws.WriteJSON(ws.TrackFile(wsPath, "artist-a", "album-1", "track-01"), track); err != nil {
		t.Fatal(err)
	}

	releaseSubmitDistributor = "distrokid"
	t.Cleanup(func() { releaseSubmitDistributor = "" })

	injectPromptInput(t, "n\n") // decline to avoid writing state

	capStdout(t, func() {
		_ = runReleaseSubmit(nil, []string{"album-1"})
	})
}

// --- runQCCheck all-checks-passed path ---

func TestRunQCCheck_AllChecksPass_Confirmed(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	rel := minimalRelease("album-1", "artist-a")
	rel.Artwork = &models.ReleaseArtwork{PrimaryFile: "cover.jpg", DimensionsPx: 3000}
	rel.Dates = &models.ReleaseDates{
		TargetRelease:                 "2025-01-01",
		DistributorSubmissionDeadline: "2024-12-01",
	}
	writeRelease(t, wsPath, "artist-a", "album-1", rel)

	qcCheckRelease = "album-1"
	t.Cleanup(func() { qcCheckRelease = "" })

	injectPromptInput(t, "y\n")

	got := capStdout(t, func() {
		if err := runQCCheck(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "All checks passed") {
		t.Errorf("expected 'All checks passed', got %q", got)
	}
}

// --- runReleaseSetLive cancel path ---

func TestRunReleaseSetLive_Canceled(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	rel := minimalRelease("album-1", "artist-a")
	rel.Status = statusSubmitted
	writeRelease(t, wsPath, "artist-a", "album-1", rel)

	releaseSetLiveDate = "2025-01-01"
	cfg.SkipConfirm = false
	t.Cleanup(func() { releaseSetLiveDate = ""; cfg.SkipConfirm = false })

	injectPromptInput(t, "n\n")

	err := runReleaseSetLive(nil, []string{"album-1"})
	var e *ExitError
	if !errors.As(err, &e) {
		t.Errorf("expected ExitError{4}, got %v", err)
	} else if e.Code != 4 {
		t.Errorf("expected exit code 4 (canceled), got %d", e.Code)
	}
}

// --- runReleaseAdvance cancel path ---

func TestRunReleaseAdvance_Canceled(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	rel := minimalRelease("album-1", "artist-a")
	rel.Status = statusInProduction
	writeRelease(t, wsPath, "artist-a", "album-1", rel)
	track := models.Track{
		SchemaVersion: "1", ID: "track-01", ReleaseID: "album-1",
		Title: "T1", Position: 1,
	}
	if err := ws.WriteJSON(ws.TrackFile(wsPath, "artist-a", "album-1", "track-01"), track); err != nil {
		t.Fatal(err)
	}

	releaseAdvanceStatus = statusMastering
	cfg.SkipConfirm = false
	t.Cleanup(func() { releaseAdvanceStatus = ""; cfg.SkipConfirm = false })

	injectPromptInput(t, "n\n")

	err := runReleaseAdvance(nil, []string{"album-1"})
	var e *ExitError
	if !errors.As(err, &e) {
		t.Errorf("expected ExitError{4}, got %v", err)
	} else if e.Code != 4 {
		t.Errorf("expected exit code 4, got %d", e.Code)
	}
}

// --- runStatus with live release (should be skipped) ---

func TestRunStatus_LiveReleaseSkipped(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	rel := minimalRelease("album-1", "artist-a")
	rel.Status = statusLive
	writeRelease(t, wsPath, "artist-a", "album-1", rel)

	got := capStdout(t, func() {
		if err := runStatus(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if strings.Contains(got, "album-1") {
		t.Errorf("live release should be excluded from active list, got %q", got)
	}
	if !strings.Contains(got, "Active releases (0)") {
		t.Errorf("expected 0 active releases, got %q", got)
	}
}

// --- runMasteringProfileList error reading profile ---

func TestRunMasteringProfileList_SkipsUnreadable(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	cfg.ArtistID = "artist-a"
	// Write an unreadable profile (bad JSON).
	profDir := ws.ArtistDir(wsPath, "artist-a")
	path := ws.MasteringProfileFile(wsPath, "artist-a", "bad-profile")
	if err := ws.WriteJSON(path, map[string]any{"id": "bad-profile", "schema_version": "1"}); err != nil {
		t.Fatal(err)
	}
	_ = profDir

	// Also add a valid profile so the table renders.
	writeMasteringProfile(t, wsPath, "artist-a", "ambient")

	got := capStdout(t, func() {
		if err := runMasteringProfileList(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "ambient") {
		t.Errorf("expected valid profile in output, got %q", got)
	}
}

// --- resolveArtistID with filter ---

func TestRunReleaseList_FilterByArtistViaConfig(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	writeRelease(t, wsPath, "artist-a", "album-1", minimalRelease("album-1", "artist-a"))

	cfg.ArtistID = "artist-b" // non-existent: filter excludes artist-a
	got := capStdout(t, func() {
		if err := runReleaseList(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if strings.Contains(got, "album-1") {
		t.Errorf("expected album-1 filtered out, got %q", got)
	}
}

// --- Execute basic smoke test ---

func TestExecute_HelpReturnsZero(t *testing.T) {
	rootCmd.SetArgs([]string{"--help"})
	t.Cleanup(func() { rootCmd.SetArgs(nil) })

	code := Execute()
	if code != 0 {
		t.Errorf("expected exit code 0 for --help, got %d", code)
	}
}

// --- runStatus with follow-up opportunity ---

func TestRunStatus_FollowUpOpportunity(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})

	followUp := "2025-03-01"
	opp := models.Opportunity{
		SchemaVersion: "1",
		ID:            "opp-2",
		Type:          "sync-license",
		Status:        "sent",
		Contact:       models.OpportunityContact{Name: "Bob"},
		FollowUpDue:   &followUp,
	}
	importPath := wsPath + "/opportunities/opp-2.json"
	if err := ws.WriteJSON(importPath, opp); err != nil {
		t.Fatal(err)
	}

	got := capStdout(t, func() {
		if err := runStatus(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Bob") {
		t.Errorf("expected follow-up contact in output, got %q", got)
	}
}

// use time import to avoid unused
var _ = time.Now
