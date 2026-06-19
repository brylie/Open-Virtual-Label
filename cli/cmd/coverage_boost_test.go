package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-virtual-label/ovl/internal/models"
	"github.com/open-virtual-label/ovl/internal/output"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
)

// --- runMasteringProfileList with profiles ---

func TestRunMasteringProfileList_WithProfiles(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	cfg.ArtistID = "artist-a"
	writeMasteringProfile(t, wsPath, "artist-a", "ambient")
	writeMasteringProfile(t, wsPath, "artist-a", "streaming")

	got := capStdout(t, func() {
		if err := runMasteringProfileList(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "ambient") || !strings.Contains(got, "streaming") {
		t.Errorf("expected both profiles in output, got %q", got)
	}
}

// --- runArtistShow branches ---

func TestRunArtistShow_AllFields(t *testing.T) {
	wsPath := setupWorkspace(t)
	a := models.Artist{
		SchemaVersion:  "1",
		DisplayName:    "Full Artist",
		LegalName:      "Legal Name Ltd",
		AlsoKnownAs:    []string{"Alias One", "Alias Two"},
		DefaultLicense: "CC BY 4.0",
		Distribution:   &models.ArtistDistrib{Distributor: "amuse"},
		Rights:         &models.ArtistRights{PRO: "Teosto", IPINumber: "12345"},
		CreatedDate:    "2024-01-01",
	}
	writeArtist(t, wsPath, "full-artist", a)

	got := capStdout(t, func() {
		if err := runArtistShow(nil, []string{"full-artist"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	for _, expected := range []string{"Legal Name", "Alias One", "amuse", "Teosto", "12345", "2024-01-01"} {
		if !strings.Contains(got, expected) {
			t.Errorf("expected %q in output, got %q", expected, got)
		}
	}
}

func TestRunArtistShow_JSON(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})

	output.SetJSON(true)
	t.Cleanup(func() { output.SetJSON(false) })

	got := capStdout(t, func() {
		if err := runArtistShow(nil, []string{"artist-a"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !json.Valid([]byte(got)) {
		t.Errorf("expected valid JSON, got %q", got)
	}
}

// --- runValidate error cases ---

func TestRunValidate_InvalidFile(t *testing.T) {
	wsPath := setupWorkspace(t)
	// Write a malformed artist JSON — valid JSON but missing required fields.
	artistDir := filepath.Join(wsPath, "artists", "bad-artist")
	if err := os.MkdirAll(artistDir, 0o755); err != nil {
		t.Fatal(err)
	}
	badPath := filepath.Join(artistDir, "artist.json")
	if err := os.WriteFile(badPath, []byte(`{"schema_version":"1","id":"x"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	validateAll = false
	t.Cleanup(func() { validateAll = false })

	err := runValidate(nil, []string{badPath})
	if err == nil {
		t.Error("expected ExitError for schema validation failures")
	}
	var e *ExitError
	if errors.As(err, &e) && e.Code != 3 {
		t.Errorf("expected exit code 3, got %d", e.Code)
	}
}

func TestRunValidate_All_WithFailure(t *testing.T) {
	wsPath := setupWorkspace(t)
	// Write an artist.json that fails validation (missing display_name).
	artistDir := filepath.Join(wsPath, "artists", "bad-artist")
	if err := os.MkdirAll(artistDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(artistDir, "artist.json"),
		[]byte(`{"schema_version":"1","id":"x"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	validateAll = true
	t.Cleanup(func() { validateAll = false })

	err := runValidate(nil, nil)
	if err == nil {
		t.Error("expected error for --all with invalid files")
	}
}

func TestRunValidate_All_SkipsUnknownPaths(t *testing.T) {
	wsPath := setupWorkspace(t)
	// Write a JSON file in an unrecognized location — InferSchema will fail.
	unknownDir := filepath.Join(wsPath, "unknown")
	if err := os.MkdirAll(unknownDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(unknownDir, "data.json"),
		[]byte(`{"foo":"bar"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	validateAll = true
	t.Cleanup(func() { validateAll = false })

	// Should succeed: unrecognized files are skipped, not counted as errors.
	if err := runValidate(nil, nil); err != nil {
		t.Errorf("unexpected error for unknown path: %v", err)
	}
}

// --- runTrackShow mastering and QC passed ---

func TestRunTrackShow_WithMastering(t *testing.T) {
	wsPath, artistID, releaseID, trackID := setupReleaseWithTrack(t)
	trackShowRelease = ""

	lufs := -14.5
	peak := -1.0
	lra := 8.0
	var track models.Track
	if err := ws.ReadJSON(ws.TrackFile(wsPath, artistID, releaseID, trackID), &track); err != nil {
		t.Fatal(err)
	}
	track.Mastering = &models.TrackMastering{
		IntegratedLUFS: &lufs,
		TruePeakDBTP:   &peak,
		LRA:            &lra,
	}
	if err := ws.WriteJSON(ws.TrackFile(wsPath, artistID, releaseID, trackID), track); err != nil {
		t.Fatal(err)
	}

	got := capStdout(t, func() {
		if err := runTrackShow(nil, []string{trackID}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Mastering") {
		t.Errorf("expected mastering section in output, got %q", got)
	}
}

func TestRunTrackShow_WithFiles(t *testing.T) {
	wsPath, artistID, releaseID, trackID := setupReleaseWithTrack(t)
	trackShowRelease = ""

	masterPath := "/path/to/master.wav"
	var track models.Track
	if err := ws.ReadJSON(ws.TrackFile(wsPath, artistID, releaseID, trackID), &track); err != nil {
		t.Fatal(err)
	}
	track.Files = &models.TrackFiles{MasterWAV: &masterPath}
	if err := ws.WriteJSON(ws.TrackFile(wsPath, artistID, releaseID, trackID), track); err != nil {
		t.Fatal(err)
	}

	got := capStdout(t, func() {
		if err := runTrackShow(nil, []string{trackID}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "/path/to/master.wav") {
		t.Errorf("expected file path in output, got %q", got)
	}
}

func TestRunTrackShow_QCPassed(t *testing.T) {
	wsPath, artistID, releaseID, trackID := setupReleaseWithTrack(t)
	trackShowRelease = ""

	var track models.Track
	if err := ws.ReadJSON(ws.TrackFile(wsPath, artistID, releaseID, trackID), &track); err != nil {
		t.Fatal(err)
	}
	track.QC = &models.TrackQC{Passed: true}
	if err := ws.WriteJSON(ws.TrackFile(wsPath, artistID, releaseID, trackID), track); err != nil {
		t.Fatal(err)
	}

	got := capStdout(t, func() {
		if err := runTrackShow(nil, []string{trackID}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "passed") {
		t.Errorf("expected 'passed' in output, got %q", got)
	}
}

func TestRunTrackShow_QCFailed(t *testing.T) {
	wsPath, artistID, releaseID, trackID := setupReleaseWithTrack(t)
	trackShowRelease = ""

	var track models.Track
	if err := ws.ReadJSON(ws.TrackFile(wsPath, artistID, releaseID, trackID), &track); err != nil {
		t.Fatal(err)
	}
	track.QC = &models.TrackQC{Passed: false, Failures: []string{"LUFS too high", "peak clip"}}
	if err := ws.WriteJSON(ws.TrackFile(wsPath, artistID, releaseID, trackID), track); err != nil {
		t.Fatal(err)
	}

	got := capStdout(t, func() {
		if err := runTrackShow(nil, []string{trackID}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "failed") {
		t.Errorf("expected 'failed' in QC status, got %q", got)
	}
	if !strings.Contains(got, "2") {
		t.Errorf("expected issue count in QC status, got %q", got)
	}
}

// --- runTrackSetFile all field assignments ---

func TestRunTrackSetFile_AllFields(t *testing.T) {
	fields := []string{"stems_zip", "project_file", "mp3_320", "wav_for_distribution"}
	for _, field := range fields {
		t.Run(field, func(t *testing.T) {
			wsPath, _, _, trackID := setupReleaseWithTrack(t)
			filePath := filepath.Join(wsPath, "test-file")
			if err := os.WriteFile(filePath, []byte("data"), 0o644); err != nil {
				t.Fatal(err)
			}

			trackSetFileRelease = ""
			trackSetFileField = field
			trackSetFilePath = filePath
			t.Cleanup(func() { trackSetFileRelease = ""; trackSetFileField = ""; trackSetFilePath = "" })

			if err := runTrackSetFile(nil, []string{trackID}); err != nil {
				t.Errorf("%s: unexpected error: %v", field, err)
			}
		})
	}
}

// --- runReleaseAddLink all platforms ---

func TestRunReleaseAddLink_AllPlatforms(t *testing.T) {
	platforms := []string{"apple_music", "youtube_music", "bandcamp", "soundcloud", "tidal", "amazon_music"}
	for _, platform := range platforms {
		t.Run(platform, func(t *testing.T) {
			wsPath := setupWorkspace(t)
			writeArtist(t, wsPath, "artist-a", models.Artist{
				SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
			})
			writeRelease(t, wsPath, "artist-a", "album-1", minimalRelease("album-1", "artist-a"))

			releaseAddLinkPlatform = platform
			releaseAddLinkURL = "https://example.com/" + platform
			t.Cleanup(func() { releaseAddLinkPlatform = ""; releaseAddLinkURL = "" })

			if err := runReleaseAddLink(nil, []string{"album-1"}); err != nil {
				t.Errorf("%s: unexpected error: %v", platform, err)
			}
		})
	}
}

// --- runStatus with opportunities ---

func TestRunStatus_WithOpportunity(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	// Write a draft-ready opportunity.
	oppDir := filepath.Join(wsPath, "opportunities")
	if err := os.MkdirAll(oppDir, 0o755); err != nil {
		t.Fatal(err)
	}
	opp := models.Opportunity{
		SchemaVersion: "1",
		ID:            "opp-1",
		Type:          "sync-license",
		Status:        "draft-ready",
		Contact:       models.OpportunityContact{Name: "Alice"},
	}
	if err := ws.WriteJSON(filepath.Join(oppDir, "opp-1.json"), opp); err != nil {
		t.Fatal(err)
	}

	got := capStdout(t, func() {
		if err := runStatus(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Alice") {
		t.Errorf("expected contact name in pending approvals, got %q", got)
	}
}

func TestRunStatus_MasteringReleaseOpenLoop(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	rel := minimalRelease("album-1", "artist-a")
	rel.Status = statusMastering
	writeRelease(t, wsPath, "artist-a", "album-1", rel)

	got := capStdout(t, func() {
		if err := runStatus(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Open loops") {
		t.Errorf("expected open loops section, got %q", got)
	}
}

// --- appendFinanceEntry error path ---

func TestAppendFinanceEntry_MarshalError(t *testing.T) {
	dir := t.TempDir()
	// Create an entry with an invalid path to trigger mkdir failure.
	// This tests the MkdirAll path: use a path where a file blocks dir creation.
	blockFile := filepath.Join(dir, "finance")
	if err := os.WriteFile(blockFile, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	entry := models.FinanceEntry{ID: "e1"}
	err := appendFinanceEntry(dir, "revenue.json", &entry)
	if err == nil {
		t.Error("expected error when finance dir is blocked by a file")
	}
}

// --- runQCCheck with passed QC (stops before confirm prompt) ---

func TestRunQCCheck_AllChecksPassed(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})

	// Create release with all required fields for QC.
	rel := minimalRelease("album-1", "artist-a")
	rel.Artwork = &models.ReleaseArtwork{PrimaryFile: "cover.jpg", DimensionsPx: 3000}
	rel.Dates = &models.ReleaseDates{
		TargetRelease:                 "2025-01-01",
		DistributorSubmissionDeadline: "2024-12-01",
	}
	writeRelease(t, wsPath, "artist-a", "album-1", rel)
	// No tracks — so track section passes with 0 checked.

	qcCheckRelease = "album-1"
	t.Cleanup(func() { qcCheckRelease = "" })

	// The function will display results, then hit prompt.Confirm.
	// Inject "n" so it returns ExitError(4) instead of hanging.
	injectPromptInput(t, "n\n")

	err := runQCCheck(nil, nil)
	// Expected: ExitError with code 4 (user declined)
	var e *ExitError
	if errors.As(err, &e) {
		if e.Code != 4 {
			t.Errorf("expected exit code 4 (canceled), got %d", e.Code)
		}
	} else if err != nil {
		t.Errorf("unexpected error type: %v", err)
	}
}
