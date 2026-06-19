package workspace

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Find locates the workspace directory. If hint is non-empty it is used directly.
// Otherwise walks up from cwd until a workspace/ subdirectory is found.
func Find(hint string) (string, error) {
	if hint != "" {
		abs, err := filepath.Abs(hint)
		if err != nil {
			return "", err
		}
		return abs, nil
	}

	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		candidate := filepath.Join(dir, "workspace")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", errors.New("workspace directory not found; run 'ovl init' to create one")
}

// Root returns the repository root (parent of the workspace directory).
func Root(wsPath string) string {
	return filepath.Dir(wsPath)
}

// ToSlug converts a string to a lowercase hyphen-separated slug safe for use as a directory name.
func ToSlug(s string) string {
	s = strings.ToLower(s)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if len(s) > 64 {
		s = s[:64]
	}
	return s
}

// Path helpers

func LabelFile(wsPath string) string {
	return filepath.Join(wsPath, "label", "profile.json")
}

func StateFile(wsPath string) string {
	return filepath.Join(wsPath, "state", "label-state.md")
}

func ArtistDir(wsPath, artistID string) string {
	return filepath.Join(wsPath, "artists", artistID)
}

func ArtistFile(wsPath, artistID string) string {
	return filepath.Join(ArtistDir(wsPath, artistID), "artist.json")
}

func ReleaseDir(wsPath, artistID, releaseID string) string {
	return filepath.Join(ArtistDir(wsPath, artistID), "releases", releaseID)
}

func ReleaseFile(wsPath, artistID, releaseID string) string {
	return filepath.Join(ReleaseDir(wsPath, artistID, releaseID), "release.json")
}

func TracksDir(wsPath, artistID, releaseID string) string {
	return filepath.Join(ReleaseDir(wsPath, artistID, releaseID), "tracks")
}

func TrackFile(wsPath, artistID, releaseID, trackID string) string {
	return filepath.Join(TracksDir(wsPath, artistID, releaseID), trackID+".json")
}

func MasteringProfileFile(wsPath, artistID, profileID string) string {
	return filepath.Join(ArtistDir(wsPath, artistID), "mastering-profiles", profileID+".json")
}

// Directory listing helpers

func ListArtistIDs(wsPath string) ([]string, error) {
	return listDirs(filepath.Join(wsPath, "artists"))
}

func ListReleaseIDs(wsPath, artistID string) ([]string, error) {
	return listDirs(filepath.Join(ArtistDir(wsPath, artistID), "releases"))
}

func ListTrackIDs(wsPath, artistID, releaseID string) ([]string, error) {
	dir := TracksDir(wsPath, artistID, releaseID)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var ids []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".json") {
			ids = append(ids, strings.TrimSuffix(e.Name(), ".json"))
		}
	}
	return ids, nil
}

func ListMasteringProfileIDs(wsPath, artistID string) ([]string, error) {
	dir := filepath.Join(ArtistDir(wsPath, artistID), "mastering-profiles")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var ids []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".json") {
			ids = append(ids, strings.TrimSuffix(e.Name(), ".json"))
		}
	}
	return ids, nil
}

func listDirs(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var ids []string
	for _, e := range entries {
		if e.IsDir() {
			ids = append(ids, e.Name())
		}
	}
	return ids, nil
}

// Search helpers

// FindRelease searches all artists for a release ID. Returns (artistID, error).
func FindRelease(wsPath, releaseID string) (string, error) {
	artistIDs, err := ListArtistIDs(wsPath)
	if err != nil {
		return "", err
	}
	for _, aID := range artistIDs {
		rIDs, err := ListReleaseIDs(wsPath, aID)
		if err != nil {
			continue
		}
		if slices.Contains(rIDs, releaseID) {
			return aID, nil
		}
	}
	return "", fmt.Errorf("release %q not found", releaseID)
}

// FindTrack searches for a track ID across all releases.
// Returns (artistID, releaseID, error). Use --release to disambiguate duplicates.
func FindTrack(wsPath, trackID, releaseHint string) (artistID, releaseID string, err error) {
	artistIDs, err := ListArtistIDs(wsPath)
	if err != nil {
		return "", "", err
	}
	type match struct{ artistID, releaseID string }
	var matches []match
	for _, aID := range artistIDs {
		rIDs, err := ListReleaseIDs(wsPath, aID)
		if err != nil {
			continue
		}
		for _, rID := range rIDs {
			if releaseHint != "" && rID != releaseHint {
				continue
			}
			tIDs, err := ListTrackIDs(wsPath, aID, rID)
			if err != nil {
				continue
			}
			for _, tID := range tIDs {
				if tID == trackID {
					matches = append(matches, match{aID, rID})
				}
			}
		}
	}
	if len(matches) == 0 {
		return "", "", fmt.Errorf("track %q not found", trackID)
	}
	if len(matches) > 1 {
		return "", "", fmt.Errorf("track %q exists in multiple releases; use --release to disambiguate", trackID)
	}
	return matches[0].artistID, matches[0].releaseID, nil
}

// I/O helpers

// WriteJSON marshals v as indented JSON and writes it to path, creating parent dirs.
func WriteJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// ReadJSON reads and unmarshals a JSON file into v.
func ReadJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
