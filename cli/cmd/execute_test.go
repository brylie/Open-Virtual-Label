package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/open-virtual-label/ovl/internal/models"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
)

// Ensure ws import is used.
var _ = ws.MasteringProfileFile

// --- Execute error paths ---

func TestExecute_ExitErrorReturnsCode(t *testing.T) {
	// mcp connect always returns ExitError{5}
	rootCmd.SetArgs([]string{"mcp", "connect", "internet-archive"})
	t.Cleanup(func() { rootCmd.SetArgs(nil) })

	code := Execute()
	if code != 5 {
		t.Errorf("expected exit code 5, got %d", code)
	}
}

func TestExecute_GenericErrorReturns1(t *testing.T) {
	// archive push without --release causes cobra to return an error (required flag missing)
	rootCmd.SetArgs([]string{"archive", "push"})
	t.Cleanup(func() { rootCmd.SetArgs(nil) })

	code := Execute()
	if code == 0 {
		t.Error("expected non-zero exit code for missing required flag")
	}
}

// --- outreach subcommands (cover init closure bodies) ---

func TestOutreachSubcommands_AgentStub(t *testing.T) {
	// Call each outreach subcommand through the cobra tree.
	subcmds := []string{"research", "review", "draft", "send",
		"follow-up", "log-response", "score", "intake", "close", "log"}

	for _, sub := range subcmds {
		t.Run(sub, func(t *testing.T) {
			got := capStdout(t, func() {
				rootCmd.SetArgs([]string{"outreach", sub})
				t.Cleanup(func() { rootCmd.SetArgs(nil) })
				if code := Execute(); code != 0 {
					t.Errorf("outreach %s: expected code 0, got %d", sub, code)
				}
			})
			if !strings.Contains(got, "agent") {
				t.Errorf("outreach %s: expected agent message in output, got %q", sub, got)
			}
		})
	}
}

// --- runFinanceSummary with expenses ---

func TestRunFinanceSummary_WithExpenses(t *testing.T) {
	wsPath := setupWorkspace(t)
	finSummaryPeriod = "2024-03"
	finSummaryBrief = false
	t.Cleanup(func() { finSummaryPeriod = ""; finSummaryBrief = false })

	expense := models.FinanceEntry{
		SchemaVersion: "1", ID: "exp-1", Type: "expense",
		Date: "2024-03-15", Amount: 75.00, Currency: "EUR", Source: "mastering",
	}
	if err := appendFinanceEntry(wsPath, "expenses.json", &expense); err != nil {
		t.Fatal(err)
	}

	got := capStdout(t, func() {
		if err := runFinanceSummary(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "mastering") {
		t.Errorf("expected expense source in output, got %q", got)
	}
	if !strings.Contains(got, "Net:") {
		t.Errorf("expected net calculation in output, got %q", got)
	}
}

// --- runReleaseShow no tracks branch ---

func TestRunReleaseShow_NoTracks(t *testing.T) {
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
	if !strings.Contains(got, "none") {
		t.Errorf("expected 'none' for tracks, got %q", got)
	}
}

// --- runAgentsList no agents in skills dir ---

func TestRunAgentsList_EmptySkillsDir(t *testing.T) {
	wsPath := setupWorkspace(t)
	importPath := ws.Root(wsPath) + "/.agents/skills"
	if err := osCreateDir(importPath); err != nil {
		t.Fatal(err)
	}

	got := capStdout(t, func() {
		if err := runAgentsList(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "No agent skills found") {
		t.Errorf("expected no-skills message, got %q", got)
	}
}

// --- commission and content agent-stub commands ---

func TestCommissionAgreement_AgentStub(t *testing.T) {
	rootCmd.SetArgs([]string{"commission", "agreement", "--opportunity", "opp-1"})
	t.Cleanup(func() { rootCmd.SetArgs(nil) })

	got := capStdout(t, func() {
		if code := Execute(); code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
	})
	if !strings.Contains(got, "agent") {
		t.Errorf("expected agent message in output, got %q", got)
	}
}

func TestContentBrief_AgentStub(t *testing.T) {
	rootCmd.SetArgs([]string{"content", "brief", "--release", "album-1"})
	t.Cleanup(func() { rootCmd.SetArgs(nil) })

	got := capStdout(t, func() {
		if code := Execute(); code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
	})
	if !strings.Contains(got, "agent") {
		t.Errorf("expected agent message in output, got %q", got)
	}
}

func TestSocialDraft_AgentStub(t *testing.T) {
	rootCmd.SetArgs([]string{"social", "draft", "--release", "album-1"})
	t.Cleanup(func() { rootCmd.SetArgs(nil) })

	got := capStdout(t, func() {
		if code := Execute(); code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
	})
	if !strings.Contains(got, "agent") {
		t.Errorf("expected agent message in output, got %q", got)
	}
}

// --- metrics snapshot agent stub ---

func TestMetricsSnapshot_AgentStub(t *testing.T) {
	rootCmd.SetArgs([]string{"metrics", "snapshot", "--period", "2024-01"})
	t.Cleanup(func() { rootCmd.SetArgs(nil) })

	got := capStdout(t, func() {
		if code := Execute(); code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
	})
	if !strings.Contains(got, "agent") {
		t.Errorf("expected agent message in output, got %q", got)
	}
}

// --- runStatus with state file containing ### heading ---

func TestRunStatus_StateFileWithHeading(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	stateDir := wsPath + "/state"
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		t.Fatal(err)
	}
	stateContent := "# Label State\n\n### Last session summary\n\nSome notes here.\n"
	if err := os.WriteFile(stateDir+"/label-state.md", []byte(stateContent), 0o644); err != nil {
		t.Fatal(err)
	}

	got := capStdout(t, func() {
		if err := runStatus(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Last session:") {
		t.Errorf("expected 'Last session:' in output, got %q", got)
	}
}

// --- runArchiveStatus with corrupt release JSON ---

func TestRunArchiveStatus_CorruptRelease(t *testing.T) {
	wsPath := setupWorkspace(t)
	writeArtist(t, wsPath, "artist-a", models.Artist{
		SchemaVersion: "1", DisplayName: "A", DefaultLicense: "CC BY 4.0",
	})
	// Create release directory but write invalid JSON so ReadJSON fails.
	relDir := ws.ReleaseDir(wsPath, "artist-a", "broken-release")
	if err := os.MkdirAll(relDir, 0o755); err != nil {
		t.Fatal(err)
	}
	relFile := ws.ReleaseFile(wsPath, "artist-a", "broken-release")
	if err := os.WriteFile(relFile, []byte("{ not valid json"), 0o644); err != nil {
		t.Fatal(err)
	}

	archiveStatusRelease = "broken-release"
	t.Cleanup(func() { archiveStatusRelease = "" })

	err := runArchiveStatus(nil, nil)
	if err == nil {
		t.Error("expected error reading corrupt release JSON")
	}
}

// --- runArtistList error reading artist (continues gracefully) ---

func TestRunArtistList_CorruptArtist(t *testing.T) {
	wsPath := setupWorkspace(t)
	// Create artist dir but with corrupt JSON.
	artistDir := ws.ArtistDir(wsPath, "corrupt-artist")
	if err := osCreateDir(artistDir); err != nil {
		t.Fatal(err)
	}
	if err := ws.WriteJSON(ws.ArtistFile(wsPath, "corrupt-artist"), map[string]string{"broken": "true"}); err != nil {
		t.Fatal(err)
	}

	got := capStdout(t, func() {
		// Should not error — corrupt artist is skipped in the table.
		if err := runArtistList(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	_ = got
}
