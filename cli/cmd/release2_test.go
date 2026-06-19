package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-virtual-label/ovl/internal/models"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
)

func writeMasteringProfile(t *testing.T, wsPath, artistID, profileID string) {
	t.Helper()
	p := models.MasteringProfile{
		SchemaVersion: "1",
		ID:            profileID,
		Name:          "Ambient",
		Targets: models.MasteringTargets{
			IntegratedLUFS: models.LUFSRange{Min: -16, Max: -14},
			TruePeakDBTP:   -1.0,
		},
	}
	if err := ws.WriteJSON(ws.MasteringProfileFile(wsPath, artistID, profileID), p); err != nil {
		t.Fatal(err)
	}
}

// --- runReleaseSetProfile ---

func TestRunReleaseSetProfile_Success(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	writeRelease(t, wsPath, "artist-a", "album-1", minimalRelease("album-1", "artist-a"))
	writeMasteringProfile(t, wsPath, "artist-a", "ambient")

	releaseSetProfileID = "ambient"
	t.Cleanup(func() { releaseSetProfileID = "" })

	got := capStdout(t, func() {
		if err := runReleaseSetProfile(nil, []string{"album-1"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "ambient") {
		t.Errorf("expected profile ID in output, got %q", got)
	}
}

func TestRunReleaseSetProfile_ProfileNotFound(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	writeRelease(t, wsPath, "artist-a", "album-1", minimalRelease("album-1", "artist-a"))

	releaseSetProfileID = "nonexistent-profile"
	t.Cleanup(func() { releaseSetProfileID = "" })

	err := runReleaseSetProfile(nil, []string{"album-1"})
	if err == nil {
		t.Error("expected error for missing profile")
	}
}

// --- runReleaseSetLive ---

func TestRunReleaseSetLive_Success(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	rel := minimalRelease("album-1", "artist-a")
	rel.Status = statusSubmitted
	writeRelease(t, wsPath, "artist-a", "album-1", rel)

	releaseSetLiveDate = "2025-01-01"
	cfg.SkipConfirm = true
	t.Cleanup(func() { releaseSetLiveDate = ""; cfg.SkipConfirm = false })

	got := capStdout(t, func() {
		if err := runReleaseSetLive(nil, []string{"album-1"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "live") {
		t.Errorf("expected 'live' in output, got %q", got)
	}
}

func TestRunReleaseSetLive_WrongStatus(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	writeRelease(t, wsPath, "artist-a", "album-1", minimalRelease("album-1", "artist-a"))

	releaseSetLiveDate = "2025-01-01"
	cfg.SkipConfirm = true
	t.Cleanup(func() { releaseSetLiveDate = ""; cfg.SkipConfirm = false })

	err := runReleaseSetLive(nil, []string{"album-1"})
	if err == nil {
		t.Error("expected error: release must be submitted")
	}
}

// --- runReleaseAddLink ---

func TestRunReleaseAddLink_ValidPlatform(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	rel := minimalRelease("album-1", "artist-a")
	rel.Status = statusLive
	writeRelease(t, wsPath, "artist-a", "album-1", rel)

	releaseAddLinkPlatform = "spotify"
	releaseAddLinkURL = "https://open.spotify.com/album/test"
	t.Cleanup(func() { releaseAddLinkPlatform = ""; releaseAddLinkURL = "" })

	got := capStdout(t, func() {
		if err := runReleaseAddLink(nil, []string{"album-1"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "spotify") {
		t.Errorf("expected platform in output, got %q", got)
	}
}

func TestRunReleaseAddLink_InvalidPlatform(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	writeRelease(t, wsPath, "artist-a", "album-1", minimalRelease("album-1", "artist-a"))

	releaseAddLinkPlatform = "unknown-platform"
	releaseAddLinkURL = "https://example.com"
	t.Cleanup(func() { releaseAddLinkPlatform = ""; releaseAddLinkURL = "" })

	err := runReleaseAddLink(nil, []string{"album-1"})
	if err == nil {
		t.Error("expected error for invalid platform")
	}
}

// --- runFinanceAddRevenue ---

func TestRunFinanceAddRevenue_Success(t *testing.T) {
	setupWorkspace(t)
	finRevSource = "streaming"
	finRevAmount = 150.50
	finRevCurrency = "EUR"
	finRevPeriod = "2024-06"
	finRevArtist = ""
	finRevRelease = ""
	finRevDesc = "Monthly streaming revenue"
	t.Cleanup(func() {
		finRevSource = ""; finRevAmount = 0; finRevCurrency = "EUR"
		finRevPeriod = ""; finRevArtist = ""; finRevRelease = ""; finRevDesc = ""
	})

	got := capStdout(t, func() {
		if err := runFinanceAddRevenue(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "150.50") {
		t.Errorf("expected amount in output, got %q", got)
	}
}

func TestRunFinanceAddRevenue_WithArtistAndRelease(t *testing.T) {
	setupWorkspace(t)
	finRevSource = "streaming"
	finRevAmount = 50.00
	finRevCurrency = "USD"
	finRevPeriod = "2024-07"
	finRevArtist = "artist-a"
	finRevRelease = "album-1"
	finRevDesc = ""
	t.Cleanup(func() {
		finRevSource = ""; finRevAmount = 0; finRevCurrency = "EUR"
		finRevPeriod = ""; finRevArtist = ""; finRevRelease = ""; finRevDesc = ""
	})

	if err := runFinanceAddRevenue(nil, nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// --- runFinanceAddExpense ---

func TestRunFinanceAddExpense_Success(t *testing.T) {
	setupWorkspace(t)
	finExpSource = "studio-time"
	finExpAmount = 200.00
	finExpCurrency = "EUR"
	finExpDate = "2024-06-15"
	finExpArtist = ""
	finExpDesc = "Recording session"
	t.Cleanup(func() {
		finExpSource = ""; finExpAmount = 0; finExpCurrency = "EUR"
		finExpDate = ""; finExpArtist = ""; finExpDesc = ""
	})

	got := capStdout(t, func() {
		if err := runFinanceAddExpense(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "200.00") {
		t.Errorf("expected amount in output, got %q", got)
	}
}

func TestRunFinanceAddExpense_WithArtist(t *testing.T) {
	setupWorkspace(t)
	finExpSource = "mastering"
	finExpAmount = 80.00
	finExpCurrency = "EUR"
	finExpDate = "2024-07-01"
	finExpArtist = "artist-a"
	finExpDesc = ""
	t.Cleanup(func() {
		finExpSource = ""; finExpAmount = 0; finExpCurrency = "EUR"
		finExpDate = ""; finExpArtist = ""; finExpDesc = ""
	})

	if err := runFinanceAddExpense(nil, nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// --- agentStub ---

func TestAgentStub(t *testing.T) {
	got := capStdout(t, func() {
		if err := agentStub("outreach-crm"); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "outreach-crm") {
		t.Errorf("expected agent name in output, got %q", got)
	}
}

// --- runArtistAddAlias ---

func TestRunArtistAddAlias_Success(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "Artist A", DefaultLicense: "CC BY 4.0",
	})

	artistAddAliasName = "DJ Test"
	t.Cleanup(func() { artistAddAliasName = "" })

	got := capStdout(t, func() {
		if err := runArtistAddAlias(nil, []string{"artist-a"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "DJ Test") {
		t.Errorf("expected alias in output, got %q", got)
	}
}

func TestRunArtistAddAlias_Duplicate(t *testing.T) {
	wsPath := setupWorkspace(t)
	a := models.Artist{
		SchemaVersion:  "1",
		DisplayName:    "Artist A",
		DefaultLicense: "CC BY 4.0",
		AlsoKnownAs:    []string{"DJ Test"},
	}
	writeArtist(t, wsPath, "artist-a", a)

	artistAddAliasName = "DJ Test"
	t.Cleanup(func() { artistAddAliasName = "" })

	err := runArtistAddAlias(nil, []string{"artist-a"})
	if err == nil {
		t.Error("expected error for duplicate alias")
	}
}

func TestRunArtistAddAlias_NotFound(t *testing.T) {
	setupWorkspace(t)
	artistAddAliasName = "DJ Test"
	t.Cleanup(func() { artistAddAliasName = "" })

	err := runArtistAddAlias(nil, []string{"nonexistent"})
	if err == nil {
		t.Error("expected error for nonexistent artist")
	}
}

// --- runStatus ---

func TestRunStatus_NoArtists(t *testing.T) {
	setupWorkspace(t)
	got := capStdout(t, func() {
		if err := runStatus(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "No artists") {
		t.Errorf("expected 'No artists' message, got %q", got)
	}
}

func TestRunStatus_WithRelease(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	writeRelease(t, wsPath, "artist-a", "album-1", minimalRelease("album-1", "artist-a"))

	got := capStdout(t, func() {
		if err := runStatus(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Active releases") {
		t.Errorf("expected releases section, got %q", got)
	}
}

func TestRunStatus_WithStateFile(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})

	statePath := ws.StateFile(wsPath)
	if err := os.MkdirAll(filepath.Dir(statePath), 0o755); err != nil {
		t.Fatal(err)
	}
	content := "# State\n\n### Session notes here\n\nSome content.\n"
	if err := os.WriteFile(statePath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	got := capStdout(t, func() {
		if err := runStatus(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Session notes here") {
		t.Errorf("expected state header in output, got %q", got)
	}
}

// --- HasMasteringData ---

func TestHasMasteringData_Nil(t *testing.T) {
	var track models.Track
	if track.HasMasteringData() {
		t.Error("expected false for track with no mastering data")
	}
}

func TestHasMasteringData_WithData(t *testing.T) {
	lufs := -14.5
	peak := -1.0
	track := models.Track{
		Mastering: &models.TrackMastering{
			IntegratedLUFS: &lufs,
			TruePeakDBTP:   &peak,
		},
	}
	if !track.HasMasteringData() {
		t.Error("expected true for track with mastering data")
	}
}
