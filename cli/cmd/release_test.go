package cmd

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/open-virtual-label/ovl/internal/models"
	"github.com/open-virtual-label/ovl/internal/output"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
)

func minimalRelease(id, artistID string) models.Release {
	return models.Release{
		SchemaVersion: "1",
		ID:            id,
		Title:         "Test Album",
		ArtistID:      artistID,
		ReleaseType:   "album",
		Status:        statusInProduction,
		License:       "CC BY 4.0",
	}
}

// --- runReleaseList ---

func TestRunReleaseList_Empty(t *testing.T) {
	setupWorkspace(t)
	got := capStdout(t, func() {
		if err := runReleaseList(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "No releases") {
		t.Errorf("expected empty message, got %q", got)
	}
}

func TestRunReleaseList_WithRelease(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "Artist A", DefaultLicense: "CC BY 4.0",
	})
	writeRelease(t, wsPath, "artist-a", "album-1", minimalRelease("album-1", "artist-a"))

	got := capStdout(t, func() {
		if err := runReleaseList(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "album-1") {
		t.Errorf("expected release ID in output, got %q", got)
	}
}

func TestRunReleaseList_FilterByArtist(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	writeArtist(t, wsPath, "artist-b", models.Artist{
		SchemaVersion: "1", DisplayName: "B", DefaultLicense: "CC BY 4.0",
	})
	writeRelease(t, wsPath, "artist-a", "release-a", minimalRelease("release-a", "artist-a"))
	writeRelease(t, wsPath, "artist-b", "release-b", minimalRelease("release-b", "artist-b"))

	cfg.ArtistID = "artist-a"
	got := capStdout(t, func() {
		if err := runReleaseList(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if strings.Contains(got, "release-b") {
		t.Errorf("filtered list should not contain release-b, got %q", got)
	}
}

func TestRunReleaseList_JSON(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	writeRelease(t, wsPath, "artist-a", "album-1", minimalRelease("album-1", "artist-a"))

	output.SetJSON(true)
	t.Cleanup(func() { output.SetJSON(false) })

	got := capStdout(t, func() {
		if err := runReleaseList(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !json.Valid([]byte(got)) {
		t.Errorf("expected valid JSON, got %q", got)
	}
}

func TestRunReleaseList_WithStatusFilter(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	rel := minimalRelease("album-1", "artist-a")
	rel.Status = statusReady
	writeRelease(t, wsPath, "artist-a", "album-1", rel)

	releaseListStatus = statusMastering
	t.Cleanup(func() { releaseListStatus = "" })

	got := capStdout(t, func() {
		if err := runReleaseList(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if strings.Contains(got, "album-1") {
		t.Errorf("status filter should exclude album-1 (ready≠mastering), got %q", got)
	}
}

// --- runReleaseShow ---

func TestRunReleaseShow_Found(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	writeRelease(t, wsPath, "artist-a", "album-1", minimalRelease("album-1", "artist-a"))

	got := capStdout(t, func() {
		if err := runReleaseShow(nil, []string{"album-1"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Test Album") {
		t.Errorf("expected title in output, got %q", got)
	}
}

func TestRunReleaseShow_NotFound(t *testing.T) {
	setupWorkspace(t)
	err := runReleaseShow(nil, []string{"nonexistent"})
	if err == nil {
		t.Error("expected error for nonexistent release")
	}
}

func TestRunReleaseShow_JSON(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	writeRelease(t, wsPath, "artist-a", "album-1", minimalRelease("album-1", "artist-a"))

	output.SetJSON(true)
	t.Cleanup(func() { output.SetJSON(false) })

	got := capStdout(t, func() {
		if err := runReleaseShow(nil, []string{"album-1"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !json.Valid([]byte(got)) {
		t.Errorf("expected valid JSON, got %q", got)
	}
}

func TestRunReleaseShow_WithTargetDate(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	rel := minimalRelease("album-1", "artist-a")
	rel.Dates = &models.ReleaseDates{TargetRelease: "2025-06-01"}
	writeRelease(t, wsPath, "artist-a", "album-1", rel)

	got := capStdout(t, func() {
		if err := runReleaseShow(nil, []string{"album-1"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "2025-06-01") {
		t.Errorf("expected target date in output, got %q", got)
	}
}

// --- runReleaseAdvance ---

func TestRunReleaseAdvance_NoTracks(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	writeRelease(t, wsPath, "artist-a", "album-1", minimalRelease("album-1", "artist-a"))

	releaseAdvanceStatus = statusMastering
	t.Cleanup(func() { releaseAdvanceStatus = "" })

	err := runReleaseAdvance(nil, []string{"album-1"})
	if err == nil {
		t.Error("expected error advancing release with no tracks")
	}
}

// --- resolveArtistID (already in integration_test.go) ---

// Check that release list handles a release with mastering profile ID.
func TestRunReleaseShow_WithMasteringProfile(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	rel := minimalRelease("album-1", "artist-a")
	rel.MasteringProfileID = "ambient"
	writeRelease(t, wsPath, "artist-a", "album-1", rel)

	got := capStdout(t, func() {
		if err := runReleaseShow(nil, []string{"album-1"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "ambient") {
		t.Errorf("expected mastering profile in output, got %q", got)
	}
}

// Verify resolveArtistID helper from release.go also works for release list.
// writeRelease path helper used in setup — verify listRelease with dates nil.
func TestRunReleaseList_WithNilDates(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	rel := minimalRelease("album-1", "artist-a")
	rel.Dates = nil
	writeRelease(t, wsPath, "artist-a", "album-1", rel)

	got := capStdout(t, func() {
		if err := runReleaseList(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "album-1") {
		t.Errorf("expected release in output, got %q", got)
	}
}

// Helper to write a track fixture into a release.
func writeTrack(t *testing.T, wsPath, artistID, releaseID, trackID string, track models.Track) {
	t.Helper()
	track.ID = trackID
	track.ReleaseID = releaseID
	if err := ws.WriteJSON(ws.TrackFile(wsPath, artistID, releaseID, trackID), track); err != nil {
		t.Fatal(err)
	}
}

func TestRunReleaseShow_WithTracks(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	writeRelease(t, wsPath, "artist-a", "album-1", minimalRelease("album-1", "artist-a"))
	writeTrack(t, wsPath, "artist-a", "album-1", "track-01", models.Track{
		SchemaVersion: "1", Title: "My Track", Position: 1,
	})

	got := capStdout(t, func() {
		if err := runReleaseShow(nil, []string{"album-1"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "track-01") {
		t.Errorf("expected track in output, got %q", got)
	}
}
