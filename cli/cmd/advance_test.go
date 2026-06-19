package cmd

import (
	"strings"
	"testing"
	"time"

	"github.com/open-virtual-label/ovl/internal/models"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
)

// --- runReleaseAdvance transitions ---

func TestRunReleaseAdvance_InProductionToMastering(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	rel := minimalRelease("album-1", "artist-a")
	rel.Status = statusInProduction
	writeRelease(t, wsPath, "artist-a", "album-1", rel)
	// Add a track so the pre-condition passes.
	track := models.Track{
		SchemaVersion: "1", ID: "track-01", ReleaseID: "album-1",
		Title: "T1", Position: 1,
	}
	if err := ws.WriteJSON(ws.TrackFile(wsPath, "artist-a", "album-1", "track-01"), track); err != nil {
		t.Fatal(err)
	}

	releaseAdvanceStatus = statusMastering
	cfg.SkipConfirm = true
	t.Cleanup(func() { releaseAdvanceStatus = ""; cfg.SkipConfirm = false })

	got := capStdout(t, func() {
		if err := runReleaseAdvance(nil, []string{"album-1"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "mastering") {
		t.Errorf("expected 'mastering' in output, got %q", got)
	}
}

func TestRunReleaseAdvance_MasteringToQC(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	rel := minimalRelease("album-1", "artist-a")
	rel.Status = statusMastering
	writeRelease(t, wsPath, "artist-a", "album-1", rel)

	// Add a track with mastering data.
	lufs := -14.5
	peak := -1.0
	track := models.Track{
		SchemaVersion: "1", ID: "track-01", ReleaseID: "album-1",
		Title: "T1", Position: 1,
		Mastering: &models.TrackMastering{
			IntegratedLUFS: &lufs,
			TruePeakDBTP:   &peak,
		},
	}
	if err := ws.WriteJSON(ws.TrackFile(wsPath, "artist-a", "album-1", "track-01"), track); err != nil {
		t.Fatal(err)
	}

	releaseAdvanceStatus = statusQC
	cfg.SkipConfirm = true
	t.Cleanup(func() { releaseAdvanceStatus = ""; cfg.SkipConfirm = false })

	got := capStdout(t, func() {
		if err := runReleaseAdvance(nil, []string{"album-1"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "qc") {
		t.Errorf("expected 'qc' in output, got %q", got)
	}
}

func TestRunReleaseAdvance_MasteringToQC_MissingMasteringData(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	rel := minimalRelease("album-1", "artist-a")
	rel.Status = statusMastering
	writeRelease(t, wsPath, "artist-a", "album-1", rel)

	// Track without mastering data.
	track := models.Track{
		SchemaVersion: "1", ID: "track-01", ReleaseID: "album-1",
		Title: "T1", Position: 1,
	}
	if err := ws.WriteJSON(ws.TrackFile(wsPath, "artist-a", "album-1", "track-01"), track); err != nil {
		t.Fatal(err)
	}

	releaseAdvanceStatus = statusQC
	cfg.SkipConfirm = true
	t.Cleanup(func() { releaseAdvanceStatus = ""; cfg.SkipConfirm = false })

	err := runReleaseAdvance(nil, []string{"album-1"})
	if err == nil {
		t.Error("expected error: mastering data incomplete")
	}
}

func TestRunReleaseAdvance_QCToReady(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	rel := minimalRelease("album-1", "artist-a")
	rel.Status = statusQC
	today := time.Now().Format("2006-01-02")
	rel.QC = &models.ReleaseQC{Passed: true, CheckedDate: &today}
	writeRelease(t, wsPath, "artist-a", "album-1", rel)

	releaseAdvanceStatus = statusReady
	cfg.SkipConfirm = true
	t.Cleanup(func() { releaseAdvanceStatus = ""; cfg.SkipConfirm = false })

	got := capStdout(t, func() {
		if err := runReleaseAdvance(nil, []string{"album-1"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "ready") {
		t.Errorf("expected 'ready' in output, got %q", got)
	}
}

func TestRunReleaseAdvance_QCToReady_QCNotPassed(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	rel := minimalRelease("album-1", "artist-a")
	rel.Status = statusQC
	rel.QC = nil // QC not run
	writeRelease(t, wsPath, "artist-a", "album-1", rel)

	releaseAdvanceStatus = statusReady
	cfg.SkipConfirm = true
	t.Cleanup(func() { releaseAdvanceStatus = ""; cfg.SkipConfirm = false })

	err := runReleaseAdvance(nil, []string{"album-1"})
	if err == nil {
		t.Error("expected error: QC not passed")
	}
}

func TestRunReleaseAdvance_InvalidTransitionTarget(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	rel := minimalRelease("album-1", "artist-a")
	rel.Status = statusInProduction
	writeRelease(t, wsPath, "artist-a", "album-1", rel)

	releaseAdvanceStatus = statusReady // wrong target for in-production
	cfg.SkipConfirm = true
	t.Cleanup(func() { releaseAdvanceStatus = ""; cfg.SkipConfirm = false })

	track := models.Track{
		SchemaVersion: "1", ID: "track-01", ReleaseID: "album-1",
		Title: "T1", Position: 1,
	}
	if err := ws.WriteJSON(ws.TrackFile(wsPath, "artist-a", "album-1", "track-01"), track); err != nil {
		t.Fatal(err)
	}

	err := runReleaseAdvance(nil, []string{"album-1"})
	if err == nil {
		t.Error("expected error for invalid transition target")
	}
}

func TestRunReleaseAdvance_InvalidCurrentStatus(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	rel := minimalRelease("album-1", "artist-a")
	rel.Status = statusLive // live cannot be advanced
	writeRelease(t, wsPath, "artist-a", "album-1", rel)

	releaseAdvanceStatus = statusInProduction
	t.Cleanup(func() { releaseAdvanceStatus = "" })

	err := runReleaseAdvance(nil, []string{"album-1"})
	if err == nil {
		t.Error("expected error: live status has no advance transition")
	}
}

// --- runReleaseSubmit error cases (requires 'ready' status + QC passed) ---

func TestRunReleaseSubmit_WrongStatus(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	writeRelease(t, wsPath, "artist-a", "album-1", minimalRelease("album-1", "artist-a"))

	releaseSubmitDistributor = "amuse"
	t.Cleanup(func() { releaseSubmitDistributor = "" })

	err := runReleaseSubmit(nil, []string{"album-1"})
	if err == nil {
		t.Error("expected error: must be in ready status")
	}
}
