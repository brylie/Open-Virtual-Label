package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-virtual-label/ovl/internal/models"
	"github.com/open-virtual-label/ovl/internal/output"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
)

func setupReleaseWithTrack(t *testing.T) (wsPath, artistID, releaseID, trackID string) {
	t.Helper()
	wsPath = setupWorkspace(t)
	artistID = "artist-a"
	releaseID = "album-1"
	trackID = "my-track"

	writeArtist(t, wsPath, artistID, models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	writeRelease(t, wsPath, artistID, releaseID, minimalRelease(releaseID, artistID))

	track := models.Track{
		SchemaVersion: "1",
		ID:            trackID,
		ReleaseID:     releaseID,
		Title:         "My Track",
		Position:      1,
	}
	if err := ws.WriteJSON(ws.TrackFile(wsPath, artistID, releaseID, trackID), track); err != nil {
		t.Fatal(err)
	}
	return wsPath, artistID, releaseID, trackID
}

// --- printFileField ---

func TestPrintFileField_Set(t *testing.T) {
	v := "/path/to/file.wav"
	got := capStdout(t, func() {
		printFileField("  master_wav", &v)
	})
	if !strings.Contains(got, "/path/to/file.wav") {
		t.Errorf("expected path in output, got %q", got)
	}
}

func TestPrintFileField_Nil(t *testing.T) {
	got := capStdout(t, func() {
		printFileField("  master_wav", nil)
	})
	if !strings.Contains(got, "(not set)") {
		t.Errorf("expected '(not set)' for nil field, got %q", got)
	}
}

func TestPrintFileField_EmptyString(t *testing.T) {
	v := ""
	got := capStdout(t, func() {
		printFileField("  master_wav", &v)
	})
	if !strings.Contains(got, "(not set)") {
		t.Errorf("expected '(not set)' for empty string, got %q", got)
	}
}

// --- runTrackAdd ---

func TestRunTrackAdd_CreatesTrack(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	writeRelease(t, wsPath, "artist-a", "album-1", minimalRelease("album-1", "artist-a"))

	trackAddRelease = "album-1"
	trackAddPosition = 1
	t.Cleanup(func() { trackAddRelease = ""; trackAddPosition = 0 })

	got := capStdout(t, func() {
		if err := runTrackAdd(nil, []string{"My New Track"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "my-new-track") {
		t.Errorf("expected track ID in output, got %q", got)
	}
}

func TestRunTrackAdd_AutoPosition(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	writeRelease(t, wsPath, "artist-a", "album-1", minimalRelease("album-1", "artist-a"))

	// Write an existing track at position 1
	existing := models.Track{
		SchemaVersion: "1", ID: "track-01", ReleaseID: "album-1", Title: "First", Position: 1,
	}
	if err := ws.WriteJSON(ws.TrackFile(wsPath, "artist-a", "album-1", "track-01"), existing); err != nil {
		t.Fatal(err)
	}

	trackAddRelease = "album-1"
	trackAddPosition = 0 // auto
	t.Cleanup(func() { trackAddRelease = ""; trackAddPosition = 0 })

	got := capStdout(t, func() {
		if err := runTrackAdd(nil, []string{"Second Track"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "position 2") {
		t.Errorf("expected position 2 in output, got %q", got)
	}
}

func TestRunTrackAdd_InvalidRelease(t *testing.T) {
	setupWorkspace(t)
	trackAddRelease = "nonexistent-release"
	t.Cleanup(func() { trackAddRelease = "" })

	err := runTrackAdd(nil, []string{"Some Track"})
	if err == nil {
		t.Error("expected error for nonexistent release")
	}
}

// --- runTrackShow ---

func TestRunTrackShow_Found(t *testing.T) {
	_, _, _, trackID := setupReleaseWithTrack(t)
	trackShowRelease = ""

	got := capStdout(t, func() {
		if err := runTrackShow(nil, []string{trackID}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "My Track") {
		t.Errorf("expected track title in output, got %q", got)
	}
}

func TestRunTrackShow_JSON(t *testing.T) {
	_, _, _, trackID := setupReleaseWithTrack(t)
	trackShowRelease = ""

	output.SetJSON(true)
	t.Cleanup(func() { output.SetJSON(false) })

	got := capStdout(t, func() {
		if err := runTrackShow(nil, []string{trackID}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, `"title"`) {
		t.Errorf("expected JSON with title field, got %q", got)
	}
}

func TestRunTrackShow_NotFound(t *testing.T) {
	setupWorkspace(t)
	err := runTrackShow(nil, []string{"nonexistent"})
	if err == nil {
		t.Error("expected error for nonexistent track")
	}
}

func TestRunTrackShow_WithISRC(t *testing.T) {
	wsPath, artistID, releaseID, trackID := setupReleaseWithTrack(t)
	trackShowRelease = ""

	// Add ISRC to the track
	var t2 models.Track
	if err := ws.ReadJSON(ws.TrackFile(wsPath, artistID, releaseID, trackID), &t2); err != nil {
		t.Fatal(err)
	}
	isrc := "GB-EMI-21-00001"
	t2.ISRC = &isrc
	if err := ws.WriteJSON(ws.TrackFile(wsPath, artistID, releaseID, trackID), t2); err != nil {
		t.Fatal(err)
	}

	got := capStdout(t, func() {
		if err := runTrackShow(nil, []string{trackID}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "GB-EMI-21-00001") {
		t.Errorf("expected ISRC in output, got %q", got)
	}
}

func TestRunTrackShow_WithQCFailed(t *testing.T) {
	wsPath, artistID, releaseID, trackID := setupReleaseWithTrack(t)
	trackShowRelease = ""

	var track models.Track
	if err := ws.ReadJSON(ws.TrackFile(wsPath, artistID, releaseID, trackID), &track); err != nil {
		t.Fatal(err)
	}
	track.QC = &models.TrackQC{Passed: false, Failures: []string{"low lufs"}}
	if err := ws.WriteJSON(ws.TrackFile(wsPath, artistID, releaseID, trackID), track); err != nil {
		t.Fatal(err)
	}

	got := capStdout(t, func() {
		if err := runTrackShow(nil, []string{trackID}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "failed") {
		t.Errorf("expected QC failed in output, got %q", got)
	}
}

// --- runTrackSetFile ---

func TestRunTrackSetFile_ValidField(t *testing.T) {
	wsPath, _, _, trackID := setupReleaseWithTrack(t)
	trackShowRelease = ""

	// Create a real file to reference
	filePath := filepath.Join(wsPath, "test.wav")
	if err := os.WriteFile(filePath, []byte("fake"), 0o644); err != nil {
		t.Fatal(err)
	}

	trackSetFileRelease = ""
	trackSetFileField = "master_wav"
	trackSetFilePath = filePath
	t.Cleanup(func() { trackSetFileRelease = ""; trackSetFileField = ""; trackSetFilePath = "" })

	got := capStdout(t, func() {
		if err := runTrackSetFile(nil, []string{trackID}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "master_wav") {
		t.Errorf("expected field name in output, got %q", got)
	}
}

func TestRunTrackSetFile_InvalidField(t *testing.T) {
	_, _, _, trackID := setupReleaseWithTrack(t)

	trackSetFileRelease = ""
	trackSetFileField = "invalid_field"
	trackSetFilePath = "/some/path"
	t.Cleanup(func() { trackSetFileRelease = ""; trackSetFileField = ""; trackSetFilePath = "" })

	err := runTrackSetFile(nil, []string{trackID})
	if err == nil {
		t.Error("expected error for invalid field name")
	}
}

func TestRunTrackSetFile_MissingFile(t *testing.T) {
	_, _, _, trackID := setupReleaseWithTrack(t)

	trackSetFileRelease = ""
	trackSetFileField = "master_wav"
	trackSetFilePath = "/nonexistent/file.wav"
	t.Cleanup(func() { trackSetFileRelease = ""; trackSetFileField = ""; trackSetFilePath = "" })

	err := runTrackSetFile(nil, []string{trackID})
	if err == nil {
		t.Error("expected error for missing file")
	}
}
