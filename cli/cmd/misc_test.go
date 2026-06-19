package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-virtual-label/ovl/internal/models"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
)

// --- helpers.go ---

func TestStatFile_Exists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	info, err := statFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Name() != "file.txt" {
		t.Errorf("got %q, want %q", info.Name(), "file.txt")
	}
}

func TestStatFile_Missing(t *testing.T) {
	_, err := statFile("/nonexistent/path")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestOsCreateDir(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "a", "b", "c")
	if err := osCreateDir(target); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(target); err != nil {
		t.Errorf("directory not created: %v", err)
	}
}

// --- archive.go ---

func TestBoolStr_True(t *testing.T) {
	if boolStr(true) != "✓" {
		t.Errorf("expected ✓ for true, got %q", boolStr(true))
	}
}

func TestBoolStr_False(t *testing.T) {
	if boolStr(false) != "✗" {
		t.Errorf("expected ✗ for false, got %q", boolStr(false))
	}
}

func TestRunArchivePush_ReturnsExitError(t *testing.T) {
	err := runArchivePush(nil, nil)
	if err == nil {
		t.Fatal("expected ExitError from runArchivePush")
	}
	var e *ExitError
	if !errors.As(err, &e) {
		t.Fatalf("expected *ExitError, got %T", err)
	}
	if e.Code != 5 {
		t.Errorf("expected code 5, got %d", e.Code)
	}
}

func TestRunArchiveStatus_NoArchive(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	writeRelease(t, wsPath, "artist-a", "album-1", minimalRelease("album-1", "artist-a"))

	archiveStatusRelease = "album-1"
	t.Cleanup(func() { archiveStatusRelease = "" })

	got := capStdout(t, func() {
		if err := runArchiveStatus(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "not been archived") {
		t.Errorf("expected 'not been archived' message, got %q", got)
	}
}

func TestRunArchiveStatus_WithArchive(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	rel := minimalRelease("album-1", "artist-a")
	iaID := "test-ia-id"
	iaURL := "https://archive.org/test"
	archDate := "2024-01-15"
	rel.Archive = &models.ReleaseArchive{
		MastersArchived:      true,
		StemsArchived:        false,
		ChecksumsVerified:    true,
		InternetArchiveID:    &iaID,
		InternetArchiveURL:   &iaURL,
		ArchiveDate:          &archDate,
	}
	writeRelease(t, wsPath, "artist-a", "album-1", rel)

	archiveStatusRelease = "album-1"
	t.Cleanup(func() { archiveStatusRelease = "" })

	got := capStdout(t, func() {
		if err := runArchiveStatus(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "test-ia-id") {
		t.Errorf("expected IA ID in output, got %q", got)
	}
}

func TestRunArchiveStatus_NotFound(t *testing.T) {
	setupWorkspace(t)
	archiveStatusRelease = "nonexistent"
	t.Cleanup(func() { archiveStatusRelease = "" })

	err := runArchiveStatus(nil, nil)
	if err == nil {
		t.Error("expected error for missing release")
	}
}

// --- agents.go ---

func TestParseSkillDescription_ValidFrontmatter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "SKILL.md")
	content := "---\nname: my-agent\ndescription: Handles music mastering tasks\n---\n\n# My Agent\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	got := parseSkillDescription(path)
	if got != "Handles music mastering tasks" {
		t.Errorf("got %q, want %q", got, "Handles music mastering tasks")
	}
}

func TestParseSkillDescription_MissingFile(t *testing.T) {
	got := parseSkillDescription("/nonexistent/SKILL.md")
	if got != "(no SKILL.md)" {
		t.Errorf("got %q, want %q", got, "(no SKILL.md)")
	}
}

func TestParseSkillDescription_NoFrontmatter_FallsBackToFirstLine(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "SKILL.md")
	content := "# Heading\nThis is the description.\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	got := parseSkillDescription(path)
	if got != "This is the description." {
		t.Errorf("got %q, want %q", got, "This is the description.")
	}
}

func TestParseSkillDescription_LongLineTruncated(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "SKILL.md")
	long := strings.Repeat("a", 100)
	if err := os.WriteFile(path, []byte(long+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got := parseSkillDescription(path)
	if len(got) > 80 {
		t.Errorf("expected truncation at 80 chars, got len=%d: %q", len(got), got)
	}
}

func TestParseSkillDescription_AllHeadings_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "SKILL.md")
	// Only headings and blank lines — no description in frontmatter, no fallback text.
	content := "# Title\n\n## Section\n\n---\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	got := parseSkillDescription(path)
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestRunAgentsList_NoSkillsDir(t *testing.T) {
	wsPath := setupWorkspace(t)
	// Root has no .agents/skills dir — should print guidance.
	_ = wsPath
	got := capStdout(t, func() {
		if err := runAgentsList(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "No agent skills") {
		t.Errorf("expected guidance message, got %q", got)
	}
}

func TestRunAgentsList_WithSkills(t *testing.T) {
	wsPath := setupWorkspace(t)
	// Create .agents/skills/<agent>/SKILL.md relative to workspace parent
	root := ws.Root(wsPath)
	skillDir := filepath.Join(root, ".agents", "skills", "mastering-companion")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	skillFile := filepath.Join(skillDir, "SKILL.md")
	content := "---\ndescription: Mastering companion agent\n---\n"
	if err := os.WriteFile(skillFile, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	got := capStdout(t, func() {
		if err := runAgentsList(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "mastering-companion") {
		t.Errorf("expected agent name in output, got %q", got)
	}
}

// --- mcp.go ---

func TestRunMCPList(t *testing.T) {
	got := capStdout(t, func() {
		if err := runMCPList(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "MCP") {
		t.Errorf("expected MCP table in output, got %q", got)
	}
}

func TestRunMCPConnect_ReturnsExitError(t *testing.T) {
	err := runMCPConnect(nil, []string{"internet-archive"})
	if err == nil {
		t.Fatal("expected ExitError from runMCPConnect")
	}
	var e *ExitError
	if !errors.As(err, &e) {
		t.Fatalf("expected *ExitError, got %T", err)
	}
	if e.Code != 5 {
		t.Errorf("expected code 5, got %d", e.Code)
	}
}

func TestRunMCPDisconnect_NoPanic(t *testing.T) {
	got := capStdout(t, func() {
		if err := runMCPDisconnect(nil, []string{"internet-archive"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "internet-archive") {
		t.Errorf("expected MCP name in output, got %q", got)
	}
}

// --- state.go ---

func TestRunStateShow_Found(t *testing.T) {
	wsPath := setupWorkspace(t)
	statePath := ws.StateFile(wsPath)
	if err := os.MkdirAll(filepath.Dir(statePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(statePath, []byte("# Label State\n\nSome content here.\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := capStdout(t, func() {
		if err := runStateShow(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Some content here") {
		t.Errorf("expected state content in output, got %q", got)
	}
}

func TestRunStateShow_Missing(t *testing.T) {
	setupWorkspace(t)
	err := runStateShow(nil, nil)
	if err == nil {
		t.Error("expected error for missing state file")
	}
}

// --- validate.go ---

func TestRunValidate_NoArgs(t *testing.T) {
	setupWorkspace(t)
	validateAll = false
	t.Cleanup(func() { validateAll = false })

	err := runValidate(nil, nil)
	if err == nil {
		t.Error("expected error when no args and no --all flag")
	}
}

func TestRunValidate_SingleFile_Valid(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	filePath := ws.ArtistFile(wsPath, "artist-a")

	validateAll = false
	t.Cleanup(func() { validateAll = false })

	got := capStdout(t, func() {
		if err := runValidate(nil, []string{filePath}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "valid") && !strings.Contains(got, "✓") {
		t.Errorf("expected valid indicator in output, got %q", got)
	}
}

func TestRunValidate_All_Valid(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})

	validateAll = true
	t.Cleanup(func() { validateAll = false })

	got := capStdout(t, func() {
		if err := runValidate(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	_ = got // Success message or nothing
}

func TestRelPath_SubPath(t *testing.T) {
	got := relPath("/workspace", "/workspace/artists/a/artist.json")
	want := "artists/a/artist.json"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRelPath_Unrelated(t *testing.T) {
	// When base and target are unrelated, returns the target as-is.
	got := relPath("/workspace", "/other/path")
	if got == "" {
		t.Error("expected non-empty path")
	}
}

// --- mastering.go ---

func TestRunMasteringStart_PrintsMessage(t *testing.T) {
	got := capStdout(t, func() {
		if err := runMasteringStart(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Mastering sessions require agent integration") {
		t.Errorf("expected stub message in output, got %q", got)
	}
}

// --- qc.go ---

func TestRunQCCheck_WithFailures(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	writeRelease(t, wsPath, "artist-a", "album-1", minimalRelease("album-1", "artist-a"))
	// Write a track without mastering data or ISRC.
	track := models.Track{
		SchemaVersion: "1", ID: "track-01", ReleaseID: "album-1",
		Title: "Track One", Position: 1,
	}
	if err := ws.WriteJSON(ws.TrackFile(wsPath, "artist-a", "album-1", "track-01"), track); err != nil {
		t.Fatal(err)
	}

	qcCheckRelease = "album-1"
	t.Cleanup(func() { qcCheckRelease = "" })

	got := capStdout(t, func() {
		if err := runQCCheck(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "failure") {
		t.Errorf("expected failure message in output, got %q", got)
	}
}

func TestRunQCCheck_ReleaseNotFound(t *testing.T) {
	setupWorkspace(t)
	qcCheckRelease = "nonexistent"
	t.Cleanup(func() { qcCheckRelease = "" })

	err := runQCCheck(nil, nil)
	if err == nil {
		t.Error("expected error for nonexistent release")
	}
}
