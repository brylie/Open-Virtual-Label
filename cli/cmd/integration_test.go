package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-virtual-label/ovl/internal/models"
	"github.com/open-virtual-label/ovl/internal/output"
	"github.com/open-virtual-label/ovl/internal/prompt"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
)

// injectPromptInput replaces the prompt scanner with one reading from s.
func injectPromptInput(t *testing.T, s string) {
	t.Helper()
	prompt.SetReader(strings.NewReader(s))
	t.Cleanup(func() { prompt.SetReader(os.Stdin) })
}

// --- Test helpers ---

// setupWorkspace creates a temp directory and configures cfg to use it.
// Returns the workspace path and a cleanup function.
func setupWorkspace(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "artists"), 0o755); err != nil {
		t.Fatal(err)
	}
	cfg.WorkspacePath = dir
	cfg.ArtistID = ""
	t.Cleanup(func() {
		cfg.WorkspacePath = ""
		cfg.ArtistID = ""
	})
	return dir
}

// writeArtist writes an artist JSON fixture into wsPath.
func writeArtist(t *testing.T, wsPath, artistID string, a models.Artist) {
	t.Helper()
	a.ID = artistID
	if err := ws.WriteJSON(ws.ArtistFile(wsPath, artistID), a); err != nil {
		t.Fatal(err)
	}
}

// writeRelease writes a release JSON fixture.
func writeRelease(t *testing.T, wsPath, artistID, releaseID string, r models.Release) {
	t.Helper()
	r.ID = releaseID
	r.ArtistID = artistID
	if err := ws.WriteJSON(ws.ReleaseFile(wsPath, artistID, releaseID), r); err != nil {
		t.Fatal(err)
	}
}

// capStdout captures os.Stdout for the duration of f.
func capStdout(t *testing.T, f func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	f()

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout = old

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

// --- applyEnvDefaults ---

func TestApplyEnvDefaults_WorkspacePath(t *testing.T) {
	cfg.WorkspacePath = ""
	t.Setenv("OVL_WORKSPACE", "/tmp/test-ws")
	applyEnvDefaults()
	if cfg.WorkspacePath != "/tmp/test-ws" {
		t.Errorf("got %q, want %q", cfg.WorkspacePath, "/tmp/test-ws")
	}
	cfg.WorkspacePath = ""
}

func TestApplyEnvDefaults_ArtistID(t *testing.T) {
	cfg.ArtistID = ""
	t.Setenv("OVL_ARTIST", "my-artist")
	applyEnvDefaults()
	if cfg.ArtistID != "my-artist" {
		t.Errorf("got %q, want %q", cfg.ArtistID, "my-artist")
	}
	cfg.ArtistID = ""
}

func TestApplyEnvDefaults_JSONFlag(t *testing.T) {
	cfg.OutputJSON = false
	t.Setenv("OVL_JSON", "1")
	applyEnvDefaults()
	if !cfg.OutputJSON {
		t.Error("expected OutputJSON=true")
	}
	cfg.OutputJSON = false
}

func TestApplyEnvDefaults_YesFlag(t *testing.T) {
	cfg.SkipConfirm = false
	t.Setenv("OVL_YES", "1")
	applyEnvDefaults()
	if !cfg.SkipConfirm {
		t.Error("expected SkipConfirm=true")
	}
	cfg.SkipConfirm = false
}

func TestApplyEnvDefaults_DoesNotOverrideExisting(t *testing.T) {
	cfg.WorkspacePath = "/already/set"
	t.Setenv("OVL_WORKSPACE", "/should-not-apply")
	applyEnvDefaults()
	if cfg.WorkspacePath != "/already/set" {
		t.Errorf("env should not override explicit flag; got %q", cfg.WorkspacePath)
	}
	cfg.WorkspacePath = ""
}

// --- resolveWorkspace ---

func TestResolveWorkspace_Valid(t *testing.T) {
	dir := t.TempDir()
	cfg.WorkspacePath = dir
	t.Cleanup(func() { cfg.WorkspacePath = "" })

	got, err := resolveWorkspace()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == "" {
		t.Error("expected non-empty workspace path")
	}
}

func TestResolveWorkspace_NotFound(t *testing.T) {
	dir := t.TempDir()
	cfg.WorkspacePath = ""
	t.Cleanup(func() { cfg.WorkspacePath = "" })

	orig, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(orig) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	_, err := resolveWorkspace()
	if err == nil {
		t.Error("expected ExitError for missing workspace")
	}
	var exitErr *ExitError
	if errors.As(err, &exitErr) && exitErr.Code != 2 {
		t.Errorf("expected exit code 2, got %d", exitErr.Code)
	}
}

// --- resolveArtistID ---

func TestResolveArtistID_FromFlag(t *testing.T) {
	wsPath := setupWorkspace(t)
	cfg.ArtistID = "explicit-artist"

	got, err := resolveArtistID(wsPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "explicit-artist" {
		t.Errorf("got %q, want %q", got, "explicit-artist")
	}
}

func TestResolveArtistID_SingleArtist(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "only-artist", models.Artist{
		SchemaVersion: "1", DisplayName: "Only", DefaultLicense: "CC BY 4.0",
	})

	got, err := resolveArtistID(wsPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "only-artist" {
		t.Errorf("got %q, want %q", got, "only-artist")
	}
}

func TestResolveArtistID_NoArtists(t *testing.T) {
	wsPath := setupWorkspace(t)
	_, err := resolveArtistID(wsPath)
	if err == nil {
		t.Error("expected error when no artists exist")
	}
}

// --- runArtistList ---

func TestRunArtistList_Empty(t *testing.T) {
	setupWorkspace(t)
	got := capStdout(t, func() {
		err := runArtistList(nil, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "No artists") {
		t.Errorf("expected empty message, got %q", got)
	}
}

func TestRunArtistList_WithArtists(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "Artist A", DefaultLicense: "CC BY 4.0",
	})

	got := capStdout(t, func() {
		err := runArtistList(nil, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Artist A") {
		t.Errorf("expected Artist A in output, got %q", got)
	}
}

func TestRunArtistList_JSON(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "Artist A", DefaultLicense: "CC BY 4.0",
	})

	output.SetJSON(true)
	t.Cleanup(func() { output.SetJSON(false) })

	got := capStdout(t, func() {
		_ = runArtistList(nil, nil)
	})
	if !json.Valid([]byte(got)) {
		t.Errorf("expected valid JSON output, got %q", got)
	}
}

// --- runArtistShow ---

func TestRunArtistShow_Found(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "Test Artist", DefaultLicense: "CC BY 4.0",
	})

	got := capStdout(t, func() {
		err := runArtistShow(nil, []string{"artist-a"})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Test Artist") {
		t.Errorf("expected artist name in output, got %q", got)
	}
}

func TestRunArtistShow_NotFound(t *testing.T) {
	setupWorkspace(t)
	err := runArtistShow(nil, []string{"nonexistent"})
	if err == nil {
		t.Error("expected error for nonexistent artist")
	}
}

// --- runFinanceSummary ---

func TestRunFinanceSummary_NoData(t *testing.T) {
	setupWorkspace(t)
	finSummaryPeriod = "2024-01"
	finSummaryBrief = false
	t.Cleanup(func() { finSummaryPeriod = ""; finSummaryBrief = false })

	got := capStdout(t, func() {
		err := runFinanceSummary(nil, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "2024-01") {
		t.Errorf("expected period in output, got %q", got)
	}
}

func TestRunFinanceSummary_WithRevenue(t *testing.T) {
	wsPath := setupWorkspace(t)
	finSummaryPeriod = "2024-01"
	finSummaryBrief = false
	t.Cleanup(func() { finSummaryPeriod = ""; finSummaryBrief = false })

	entry := models.FinanceEntry{
		SchemaVersion: "1", ID: "e1", Type: "revenue",
		Date: "2024-01-01", Amount: 150.00, Currency: "EUR", Source: "streaming",
	}
	if err := appendFinanceEntry(wsPath, "revenue.json", &entry); err != nil {
		t.Fatal(err)
	}

	got := capStdout(t, func() {
		err := runFinanceSummary(nil, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "streaming") {
		t.Errorf("expected source in output, got %q", got)
	}
}

func TestRunFinanceSummary_Brief(t *testing.T) {
	setupWorkspace(t)
	finSummaryPeriod = "2024-02"
	finSummaryBrief = true
	t.Cleanup(func() { finSummaryPeriod = ""; finSummaryBrief = false })

	got := capStdout(t, func() {
		err := runFinanceSummary(nil, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "2024-02") {
		t.Errorf("expected period in brief output, got %q", got)
	}
}

// --- runMasteringProfileList ---

func TestRunMasteringProfileList_Empty(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	cfg.ArtistID = "artist-a"

	got := capStdout(t, func() {
		err := runMasteringProfileList(nil, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "No mastering profiles") {
		t.Errorf("expected empty message, got %q", got)
	}
}

// --- runArtistList with bad workspace ---

func TestRunArtistList_BadWorkspace(t *testing.T) {
	cfg.WorkspacePath = ""
	dir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() {
		_ = os.Chdir(orig)
		cfg.WorkspacePath = ""
	})
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	err := runArtistList(nil, nil)
	if err == nil {
		t.Error("expected error for missing workspace")
	}
}

// --- sumEntries is already tested in finance_test.go ---

// --- Slug helpers used in init flow ---

func TestToSlug_ViaWorkspacePackage(t *testing.T) {
	// Exercises ws.ToSlug which is used in most create commands.
	got := ws.ToSlug("My Label Name!")
	if got != "my-label-name" {
		t.Errorf("got %q, want %q", got, "my-label-name")
	}
}
