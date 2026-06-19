package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-virtual-label/ovl/internal/models"
)

// --- sumEntries ---

func TestSumEntries_Empty(t *testing.T) {
	if got := sumEntries(nil); got != 0 {
		t.Errorf("got %.2f, want 0", got)
	}
}

func TestSumEntries_Single(t *testing.T) {
	entries := []models.FinanceEntry{{Amount: 42.50}}
	if got := sumEntries(entries); got != 42.50 {
		t.Errorf("got %.2f, want 42.50", got)
	}
}

func TestSumEntries_Multiple(t *testing.T) {
	entries := []models.FinanceEntry{
		{Amount: 100.00},
		{Amount: 50.25},
		{Amount: 25.75},
	}
	if got := sumEntries(entries); got != 176.00 {
		t.Errorf("got %.2f, want 176.00", got)
	}
}

func TestSumEntries_Negative(t *testing.T) {
	entries := []models.FinanceEntry{
		{Amount: 200.00},
		{Amount: -50.00},
	}
	if got := sumEntries(entries); got != 150.00 {
		t.Errorf("got %.2f, want 150.00", got)
	}
}

// --- readEntries ---

func TestReadEntries_MissingFile(t *testing.T) {
	var entries []models.FinanceEntry
	// Missing file should not panic and should leave entries empty.
	readEntries("/nonexistent/path.json", &entries)
	if len(entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(entries))
	}
}

func TestReadEntries_ValidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "revenue.json")
	raw := `[{"id":"e1","type":"revenue","date":"2024-01-01","amount":100,"currency":"EUR","source":"streaming","schema_version":"1"}]`
	if err := os.WriteFile(path, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}

	var entries []models.FinanceEntry
	readEntries(path, &entries)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Amount != 100 {
		t.Errorf("got amount %.2f, want 100.00", entries[0].Amount)
	}
}

func TestReadEntries_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	var entries []models.FinanceEntry
	// Invalid JSON should not panic; entries remain empty.
	readEntries(path, &entries)
	if len(entries) != 0 {
		t.Errorf("expected empty entries after parse failure, got %d", len(entries))
	}
}

// --- appendFinanceEntry ---

func TestAppendFinanceEntry_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	entry := models.FinanceEntry{
		SchemaVersion: "1",
		ID:            "e1",
		Type:          "revenue",
		Date:          "2024-01-01",
		Amount:        50.00,
		Currency:      "EUR",
		Source:        "streaming",
	}

	if err := appendFinanceEntry(dir, "revenue.json", &entry); err != nil {
		t.Fatalf("appendFinanceEntry: %v", err)
	}

	var entries []models.FinanceEntry
	readEntries(filepath.Join(dir, "finance", "revenue.json"), &entries)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].ID != "e1" {
		t.Errorf("got ID %q, want %q", entries[0].ID, "e1")
	}
}

func TestAppendFinanceEntry_Appends(t *testing.T) {
	dir := t.TempDir()
	for i, id := range []string{"e1", "e2", "e3"} {
		entry := models.FinanceEntry{
			SchemaVersion: "1",
			ID:            id,
			Type:          "revenue",
			Date:          "2024-01-01",
			Amount:        float64(i+1) * 10,
			Currency:      "EUR",
			Source:        "streaming",
		}
		if err := appendFinanceEntry(dir, "revenue.json", &entry); err != nil {
			t.Fatalf("appendFinanceEntry iteration %d: %v", i, err)
		}
	}

	var entries []models.FinanceEntry
	readEntries(filepath.Join(dir, "finance", "revenue.json"), &entries)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}
