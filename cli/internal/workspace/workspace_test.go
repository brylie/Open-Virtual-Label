package workspace_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	ws "github.com/open-virtual-label/ovl/internal/workspace"
)

// --- ToSlug ---

func TestToSlug_BasicLowercasing(t *testing.T) {
	if got := ws.ToSlug("Hello World"); got != "hello-world" {
		t.Errorf("got %q, want %q", got, "hello-world")
	}
}

func TestToSlug_SpecialCharsCollapsed(t *testing.T) {
	if got := ws.ToSlug("A/B  C___D"); got != "a-b-c-d" {
		t.Errorf("got %q, want %q", got, "a-b-c-d")
	}
}

func TestToSlug_LeadingTrailingSpecial(t *testing.T) {
	if got := ws.ToSlug("  !! test !!  "); got != "test" {
		t.Errorf("got %q, want %q", got, "test")
	}
}

func TestToSlug_AlreadySlug(t *testing.T) {
	if got := ws.ToSlug("my-artist"); got != "my-artist" {
		t.Errorf("got %q, want %q", got, "my-artist")
	}
}

func TestToSlug_Truncation(t *testing.T) {
	long := "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqr" // 70 chars
	got := ws.ToSlug(long)
	if len(got) > 64 {
		t.Errorf("slug length %d exceeds 64", len(got))
	}
}

func TestToSlug_Empty(t *testing.T) {
	if got := ws.ToSlug(""); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestToSlug_Numbers(t *testing.T) {
	if got := ws.ToSlug("Track 01"); got != "track-01" {
		t.Errorf("got %q, want %q", got, "track-01")
	}
}

// --- Root ---

func TestRoot(t *testing.T) {
	got := ws.Root("/home/user/project/workspace")
	if got != "/home/user/project" {
		t.Errorf("got %q, want %q", got, "/home/user/project")
	}
}

// --- Find ---

func TestFind_WithHint(t *testing.T) {
	dir := t.TempDir()
	got, err := ws.Find(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != dir {
		t.Errorf("got %q, want %q", got, dir)
	}
}

func TestFind_WalksUpToWorkspaceDir(t *testing.T) {
	root := t.TempDir()
	wsDir := filepath.Join(root, "workspace")
	if err := os.Mkdir(wsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Set cwd to a subdir of root (not workspace itself) — Find should walk up.
	subDir := filepath.Join(root, "subdir")
	if err := os.Mkdir(subDir, 0o755); err != nil {
		t.Fatal(err)
	}
	orig, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(orig) })
	if err := os.Chdir(subDir); err != nil {
		t.Fatal(err)
	}

	got, err := ws.Find("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Resolve symlinks on both sides (macOS /var → /private/var).
	gotReal, _ := filepath.EvalSymlinks(got)
	wantReal, _ := filepath.EvalSymlinks(wsDir)
	if gotReal != wantReal {
		t.Errorf("got %q, want %q", got, wsDir)
	}
}

func TestFind_NotFound(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(orig) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	_, err := ws.Find("")
	if err == nil {
		t.Error("expected error for missing workspace, got nil")
	}
}

// --- Path helpers ---

func TestPathHelpers(t *testing.T) {
	base := "/ws"

	tests := []struct {
		name string
		got  string
		want string
	}{
		{"LabelFile", ws.LabelFile(base), "/ws/label/profile.json"},
		{"StateFile", ws.StateFile(base), "/ws/state/label-state.md"},
		{"ArtistDir", ws.ArtistDir(base, "artist-1"), "/ws/artists/artist-1"},
		{"ArtistFile", ws.ArtistFile(base, "artist-1"), "/ws/artists/artist-1/artist.json"},
		{"ReleaseDir", ws.ReleaseDir(base, "artist-1", "release-a"), "/ws/artists/artist-1/releases/release-a"},
		{"ReleaseFile", ws.ReleaseFile(base, "artist-1", "release-a"), "/ws/artists/artist-1/releases/release-a/release.json"},
		{"TracksDir", ws.TracksDir(base, "artist-1", "release-a"), "/ws/artists/artist-1/releases/release-a/tracks"},
		{"TrackFile", ws.TrackFile(base, "artist-1", "release-a", "track-01"), "/ws/artists/artist-1/releases/release-a/tracks/track-01.json"},
		{"MasteringProfileFile", ws.MasteringProfileFile(base, "artist-1", "ambient"), "/ws/artists/artist-1/mastering-profiles/ambient.json"},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s: got %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

// --- WriteJSON / ReadJSON ---

type testStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestWriteReadJSON_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "data.json")

	in := testStruct{Name: "test", Value: 42}
	if err := ws.WriteJSON(path, in); err != nil {
		t.Fatalf("WriteJSON: %v", err)
	}

	var out testStruct
	if err := ws.ReadJSON(path, &out); err != nil {
		t.Fatalf("ReadJSON: %v", err)
	}
	if out.Name != in.Name || out.Value != in.Value {
		t.Errorf("round-trip mismatch: got %+v, want %+v", out, in)
	}
}

func TestWriteJSON_CreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a", "b", "c", "file.json")
	if err := ws.WriteJSON(path, testStruct{Name: "x"}); err != nil {
		t.Fatalf("WriteJSON should create parent dirs: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}

func TestReadJSON_FileNotFound(t *testing.T) {
	err := ws.ReadJSON("/nonexistent/path.json", &testStruct{})
	if err == nil {
		t.Error("expected error reading nonexistent file")
	}
}

func TestReadJSON_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := ws.ReadJSON(path, &testStruct{})
	if err == nil {
		t.Error("expected JSON parse error")
	}
}

// --- List helpers ---

func makeArtistDir(t *testing.T, wsPath, artistID string) {
	t.Helper()
	dir := ws.ArtistDir(wsPath, artistID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
}

func makeTrackFile(t *testing.T, wsPath, artistID, releaseID, trackID string) {
	t.Helper()
	path := ws.TrackFile(wsPath, artistID, releaseID, trackID)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestListArtistIDs_Empty(t *testing.T) {
	dir := t.TempDir()
	ids, err := ws.ListArtistIDs(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("expected empty, got %v", ids)
	}
}

func TestListArtistIDs_NonEmpty(t *testing.T) {
	dir := t.TempDir()
	makeArtistDir(t, dir, "artist-a")
	makeArtistDir(t, dir, "artist-b")

	ids, err := ws.ListArtistIDs(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 2 {
		t.Errorf("expected 2 artists, got %d: %v", len(ids), ids)
	}
}

func TestListTrackIDs_Empty(t *testing.T) {
	dir := t.TempDir()
	ids, err := ws.ListTrackIDs(dir, "artist-a", "release-x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("expected empty, got %v", ids)
	}
}

func TestListTrackIDs_NonEmpty(t *testing.T) {
	dir := t.TempDir()
	makeTrackFile(t, dir, "artist-a", "release-x", "track-01")
	makeTrackFile(t, dir, "artist-a", "release-x", "track-02")

	ids, err := ws.ListTrackIDs(dir, "artist-a", "release-x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 2 {
		t.Errorf("expected 2 tracks, got %d: %v", len(ids), ids)
	}
}

func TestListMasteringProfileIDs_Empty(t *testing.T) {
	dir := t.TempDir()
	ids, err := ws.ListMasteringProfileIDs(dir, "artist-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("expected empty, got %v", ids)
	}
}

func TestListMasteringProfileIDs_NonEmpty(t *testing.T) {
	dir := t.TempDir()
	profDir := filepath.Join(ws.ArtistDir(dir, "artist-a"), "mastering-profiles")
	if err := os.MkdirAll(profDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"ambient.json", "streaming.json"} {
		if err := os.WriteFile(filepath.Join(profDir, name), []byte("{}"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	ids, err := ws.ListMasteringProfileIDs(dir, "artist-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 2 {
		t.Errorf("expected 2 profiles, got %d: %v", len(ids), ids)
	}
}

// --- FindRelease ---

func makeReleaseDir(t *testing.T, wsPath, artistID, releaseID string) {
	t.Helper()
	dir := ws.ReleaseDir(wsPath, artistID, releaseID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestFindRelease_Found(t *testing.T) {
	dir := t.TempDir()
	makeArtistDir(t, dir, "artist-a")
	makeReleaseDir(t, dir, "artist-a", "album-1")

	artistID, err := ws.FindRelease(dir, "album-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if artistID != "artist-a" {
		t.Errorf("got artist %q, want %q", artistID, "artist-a")
	}
}

func TestFindRelease_NotFound(t *testing.T) {
	dir := t.TempDir()
	makeArtistDir(t, dir, "artist-a")

	_, err := ws.FindRelease(dir, "nonexistent")
	if err == nil {
		t.Error("expected error for missing release")
	}
}

func TestFindRelease_EmptyWorkspace(t *testing.T) {
	dir := t.TempDir()
	_, err := ws.FindRelease(dir, "anything")
	if err == nil {
		t.Error("expected error for empty workspace")
	}
}

// --- FindTrack ---

func TestFindTrack_Found(t *testing.T) {
	dir := t.TempDir()
	makeArtistDir(t, dir, "artist-a")
	makeReleaseDir(t, dir, "artist-a", "release-x")
	makeTrackFile(t, dir, "artist-a", "release-x", "track-01")

	artistID, releaseID, err := ws.FindTrack(dir, "track-01", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if artistID != "artist-a" || releaseID != "release-x" {
		t.Errorf("got (%q, %q), want (%q, %q)", artistID, releaseID, "artist-a", "release-x")
	}
}

func TestFindTrack_NotFound(t *testing.T) {
	dir := t.TempDir()
	makeArtistDir(t, dir, "artist-a")

	_, _, err := ws.FindTrack(dir, "nonexistent", "")
	if err == nil {
		t.Error("expected error for missing track")
	}
}

func TestFindTrack_Ambiguous(t *testing.T) {
	dir := t.TempDir()
	makeArtistDir(t, dir, "artist-a")
	makeReleaseDir(t, dir, "artist-a", "release-x")
	makeReleaseDir(t, dir, "artist-a", "release-y")
	makeTrackFile(t, dir, "artist-a", "release-x", "track-01")
	makeTrackFile(t, dir, "artist-a", "release-y", "track-01")

	_, _, err := ws.FindTrack(dir, "track-01", "")
	if err == nil {
		t.Error("expected error for ambiguous track")
	}
}

func TestFindTrack_DisambiguateWithHint(t *testing.T) {
	dir := t.TempDir()
	makeArtistDir(t, dir, "artist-a")
	makeReleaseDir(t, dir, "artist-a", "release-x")
	makeReleaseDir(t, dir, "artist-a", "release-y")
	makeTrackFile(t, dir, "artist-a", "release-x", "track-01")
	makeTrackFile(t, dir, "artist-a", "release-y", "track-01")

	artistID, releaseID, err := ws.FindTrack(dir, "track-01", "release-y")
	if err != nil {
		t.Fatalf("unexpected error with hint: %v", err)
	}
	if releaseID != "release-y" {
		t.Errorf("got release %q, want %q", releaseID, "release-y")
	}
	_ = artistID
}

func TestWriteJSON_InvalidValue(t *testing.T) {
	dir := t.TempDir()
	// channels cannot be marshaled to JSON
	ch := make(chan int)
	err := ws.WriteJSON(filepath.Join(dir, "bad.json"), ch)
	if err == nil {
		t.Error("expected marshal error for channel type")
	}
}

// Verify WriteJSON produces valid indented JSON.
func TestWriteJSON_IndentedOutput(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.json")
	if err := ws.WriteJSON(path, map[string]string{"key": "val"}); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	if !json.Valid(data) {
		t.Errorf("output is not valid JSON: %s", data)
	}
}
