package schema_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-virtual-label/ovl/internal/schema"
)

func validArtistJSON() []byte {
	return []byte(`{"schema_version":"1","id":"test-artist","display_name":"Test Artist","default_license":"CC BY 4.0"}`)
}

// --- Validate ---

func TestValidate_UnknownSchema(t *testing.T) {
	_, err := schema.Validate("does-not-exist", []byte(`{}`))
	if err == nil {
		t.Error("expected error for unknown schema name")
	}
}

func TestValidate_InvalidJSON(t *testing.T) {
	errs, err := schema.Validate("artist", []byte(`not json`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) == 0 {
		t.Error("expected validation errors for non-JSON input")
	}
}

func TestValidate_Artist_Valid(t *testing.T) {
	errs, err := schema.Validate("artist", validArtistJSON())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) != 0 {
		t.Errorf("unexpected validation errors: %v", errs)
	}
}

func TestValidate_Artist_MissingRequired(t *testing.T) {
	raw := []byte(`{"schema_version":"1","id":"x"}`)
	errs, err := schema.Validate("artist", raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) == 0 {
		t.Error("expected validation errors for missing display_name and default_license")
	}
}

func TestValidate_Label_Valid(t *testing.T) {
	raw := []byte(`{"schema_version":"1","id":"my-label","name":"My Label","default_license":"CC BY 4.0"}`)
	errs, err := schema.Validate("label", raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) != 0 {
		t.Errorf("unexpected validation errors: %v", errs)
	}
}

func TestValidate_Label_MissingRequired(t *testing.T) {
	raw := []byte(`{"schema_version":"1","id":"x"}`)
	errs, err := schema.Validate("label", raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) == 0 {
		t.Error("expected validation errors")
	}
}

func TestValidate_Release_Valid(t *testing.T) {
	raw := []byte(`{"schema_version":"1","id":"release-1","title":"My Album","artist_id":"artist-a","release_type":"album","status":"in-production","license":"CC BY 4.0"}`)
	errs, err := schema.Validate("release", raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) != 0 {
		t.Errorf("unexpected validation errors: %v", errs)
	}
}

func TestValidate_Track_Valid(t *testing.T) {
	raw := []byte(`{"schema_version":"1","id":"track-01","release_id":"release-1","title":"My Track","position":1}`)
	errs, err := schema.Validate("track", raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) != 0 {
		t.Errorf("unexpected validation errors: %v", errs)
	}
}

func TestValidate_Track_MissingRequired(t *testing.T) {
	raw := []byte(`{"schema_version":"1","id":"track-01"}`)
	errs, err := schema.Validate("track", raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) == 0 {
		t.Error("expected validation errors for missing title, release_id, position")
	}
}

func TestValidate_MasteringProfile_Valid(t *testing.T) {
	raw := []byte(`{"schema_version":"1","id":"ambient","name":"Ambient","targets":{"integrated_lufs":{"min":-16,"max":-14},"true_peak_dbtp":-1.0}}`)
	errs, err := schema.Validate("mastering-profile", raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) != 0 {
		t.Errorf("unexpected validation errors: %v", errs)
	}
}

func TestValidate_Opportunity_Valid(t *testing.T) {
	raw := []byte(`{"schema_version":"1","id":"opp-1","type":"sync-license","status":"identified","contact":{"name":"Alice"}}`)
	errs, err := schema.Validate("opportunity", raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) != 0 {
		t.Errorf("unexpected validation errors: %v", errs)
	}
}

func TestValidate_FinanceEntry_Valid(t *testing.T) {
	raw := []byte(`{"schema_version":"1","id":"rev-2024-01","type":"revenue","date":"2024-01-01","amount":100.00,"currency":"EUR","source":"streaming"}`)
	errs, err := schema.Validate("finance-entry", raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) != 0 {
		t.Errorf("unexpected validation errors: %v", errs)
	}
}

func TestValidate_FinanceEntry_MissingRequired(t *testing.T) {
	raw := []byte(`{"schema_version":"1","id":"x","type":"revenue"}`)
	errs, err := schema.Validate("finance-entry", raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) == 0 {
		t.Error("expected validation errors for missing date, amount, currency, source")
	}
}

// Calling Validate twice exercises the cache path.
func TestValidate_UsesCache(t *testing.T) {
	for range 2 {
		errs, err := schema.Validate("artist", validArtistJSON())
		if err != nil {
			t.Fatalf("unexpected error on cached call: %v", err)
		}
		if len(errs) != 0 {
			t.Errorf("unexpected errors: %v", errs)
		}
	}
}

// --- ValidateFile ---

func TestValidateFile_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "artist.json")
	if err := os.WriteFile(path, validArtistJSON(), 0o644); err != nil {
		t.Fatal(err)
	}
	errs, err := schema.ValidateFile("artist", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) != 0 {
		t.Errorf("unexpected errors: %v", errs)
	}
}

func TestValidateFile_NotFound(t *testing.T) {
	_, err := schema.ValidateFile("artist", "/nonexistent/file.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

// --- InferSchema ---

func TestInferSchema_Label(t *testing.T) {
	wsPath := "/ws"
	got, err := schema.InferSchema(wsPath, "/ws/label/profile.json")
	if err != nil {
		t.Fatal(err)
	}
	if got != "label" {
		t.Errorf("got %q, want %q", got, "label")
	}
}

func TestInferSchema_Artist(t *testing.T) {
	got, err := schema.InferSchema("/ws", "/ws/artists/my-artist/artist.json")
	if err != nil {
		t.Fatal(err)
	}
	if got != "artist" {
		t.Errorf("got %q, want %q", got, "artist")
	}
}

func TestInferSchema_Release(t *testing.T) {
	got, err := schema.InferSchema("/ws", "/ws/artists/a/releases/r/release.json")
	if err != nil {
		t.Fatal(err)
	}
	if got != "release" {
		t.Errorf("got %q, want %q", got, "release")
	}
}

func TestInferSchema_Track(t *testing.T) {
	got, err := schema.InferSchema("/ws", "/ws/artists/a/releases/r/tracks/track-01.json")
	if err != nil {
		t.Fatal(err)
	}
	if got != "track" {
		t.Errorf("got %q, want %q", got, "track")
	}
}

func TestInferSchema_MasteringProfile(t *testing.T) {
	got, err := schema.InferSchema("/ws", "/ws/artists/a/mastering-profiles/ambient.json")
	if err != nil {
		t.Fatal(err)
	}
	if got != "mastering-profile" {
		t.Errorf("got %q, want %q", got, "mastering-profile")
	}
}

func TestInferSchema_Opportunity(t *testing.T) {
	got, err := schema.InferSchema("/ws", "/ws/opportunities/opp-1.json")
	if err != nil {
		t.Fatal(err)
	}
	if got != "opportunity" {
		t.Errorf("got %q, want %q", got, "opportunity")
	}
}

func TestInferSchema_Finance(t *testing.T) {
	got, err := schema.InferSchema("/ws", "/ws/finance/revenue.json")
	if err != nil {
		t.Fatal(err)
	}
	if got != "finance-entry" {
		t.Errorf("got %q, want %q", got, "finance-entry")
	}
}

func TestInferSchema_Unknown(t *testing.T) {
	_, err := schema.InferSchema("/ws", "/ws/unknown/file.json")
	if err == nil {
		t.Error("expected error for unrecognized path")
	}
}
