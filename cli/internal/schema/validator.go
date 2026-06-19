package schema

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

//go:embed data
var schemaFS embed.FS

// compiled caches compiled schemas by name.
var compiled = map[string]*jsonschema.Schema{}

// schemaFiles maps schema names to their embedded paths.
var schemaFiles = map[string]string{
	"label":             "data/label.schema.json",
	"artist":            "data/artist.schema.json",
	"release":           "data/release.schema.json",
	"track":             "data/track.schema.json",
	"mastering-profile": "data/mastering-profile.schema.json",
	"opportunity":       "data/opportunity.schema.json",
	"finance-entry":     "data/finance-entry.schema.json",
}

func getSchema(name string) (*jsonschema.Schema, error) {
	if s, ok := compiled[name]; ok {
		return s, nil
	}
	path, ok := schemaFiles[name]
	if !ok {
		return nil, fmt.Errorf("unknown schema: %q", name)
	}
	data, err := schemaFS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading schema %s: %w", name, err)
	}
	c := jsonschema.NewCompiler()
	if err := c.AddResource(name+".json", strings.NewReader(string(data))); err != nil {
		return nil, fmt.Errorf("loading schema %s: %w", name, err)
	}
	s, err := c.Compile(name + ".json")
	if err != nil {
		return nil, fmt.Errorf("compiling schema %s: %w", name, err)
	}
	compiled[name] = s
	return s, nil
}

// Validate validates rawJSON against the named schema.
// Returns a slice of human-readable error strings (empty slice means valid).
func Validate(schemaName string, rawJSON []byte) ([]string, error) {
	s, err := getSchema(schemaName)
	if err != nil {
		return nil, err
	}
	var v any
	if err := json.Unmarshal(rawJSON, &v); err != nil {
		return []string{fmt.Sprintf("invalid JSON: %v", err)}, nil
	}
	if err := s.Validate(v); err != nil {
		var ve *jsonschema.ValidationError
		if errors.As(err, &ve) {
			return flattenErrors(ve), nil
		}
		return []string{err.Error()}, nil
	}
	return nil, nil
}

// ValidateFile reads a file and validates it against the given schema.
func ValidateFile(schemaName, filePath string) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return Validate(schemaName, data)
}

// InferSchema returns the schema name for a workspace-relative file path.
func InferSchema(wsPath, filePath string) (string, error) {
	rel, err := filepath.Rel(wsPath, filePath)
	if err != nil {
		return "", err
	}
	rel = filepath.ToSlash(rel)

	switch {
	case rel == "label/profile.json":
		return "label", nil
	case matchGlob(rel, "artists/*/artist.json"):
		return "artist", nil
	case matchGlob(rel, "artists/*/releases/*/release.json"):
		return "release", nil
	case matchGlob(rel, "artists/*/releases/*/tracks/*.json"):
		return "track", nil
	case matchGlob(rel, "artists/*/mastering-profiles/*.json"):
		return "mastering-profile", nil
	case matchGlob(rel, "opportunities/*.json"):
		return "opportunity", nil
	case matchGlob(rel, "finance/*.json"):
		return "finance-entry", nil
	default:
		return "", fmt.Errorf("cannot infer schema for path %q", rel)
	}
}

func flattenErrors(ve *jsonschema.ValidationError) []string {
	if len(ve.Causes) == 0 {
		loc := ve.InstanceLocation
		if loc == "" {
			loc = "/"
		}
		return []string{fmt.Sprintf("  %s: %s", loc, ve.Message)}
	}
	var out []string
	for _, cause := range ve.Causes {
		out = append(out, flattenErrors(cause)...)
	}
	return out
}

// matchGlob matches a slash-separated path against a simple glob pattern.
// Each path segment is matched against the corresponding pattern segment using
// filepath.Match, so patterns like "*.json" work correctly within a segment.
func matchGlob(path, pattern string) bool {
	pp := strings.Split(path, "/")
	pg := strings.Split(pattern, "/")
	if len(pp) != len(pg) {
		return false
	}
	for i, seg := range pg {
		ok, err := filepath.Match(seg, pp[i])
		if err != nil || !ok {
			return false
		}
	}
	return true
}
