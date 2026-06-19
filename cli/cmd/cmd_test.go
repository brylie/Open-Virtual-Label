package cmd

import (
	"regexp"
	"testing"
)

// --- ExitError ---

func TestExitError_Error(t *testing.T) {
	e := &ExitError{Code: 1, Msg: "something failed"}
	if e.Error() != "something failed" {
		t.Errorf("got %q, want %q", e.Error(), "something failed")
	}
}

func TestExitError_ZeroCode(t *testing.T) {
	e := &ExitError{Code: 0, Msg: "ok"}
	if e.Error() != "ok" {
		t.Errorf("got %q, want %q", e.Error(), "ok")
	}
}

// --- ISRC pattern ---

var isrcValidCases = []string{
	"GB-EMI-21-00001",  // with dashes
	"GBEMI2100001",     // no dashes
	"US-S1Z-99-12345",
	"FI-A1B-03-00042",
}

var isrcInvalidCases = []string{
	"",
	"GBEMI210000",      // too short (missing last digit)
	"GBEMI21000001",    // too long (extra digit)
	"12-EMI-21-00001",  // country code must be uppercase letters only
	"gb-emi-21-00001",  // must be uppercase
	"GB-EM-21-00001",   // publisher code must be 3 chars
}

func TestIsrcPattern_ValidCases(t *testing.T) {
	for _, tc := range isrcValidCases {
		if !isrcPattern.MatchString(tc) {
			t.Errorf("expected %q to match ISRC pattern", tc)
		}
	}
}

func TestIsrcPattern_InvalidCases(t *testing.T) {
	for _, tc := range isrcInvalidCases {
		if isrcPattern.MatchString(tc) {
			t.Errorf("expected %q to NOT match ISRC pattern", tc)
		}
	}
}

// Verify the pattern compiles to the expected regex.
func TestIsrcPattern_Regex(t *testing.T) {
	expected := regexp.MustCompile(`^[A-Z]{2}-?[A-Z0-9]{3}-?\d{2}-?\d{5}$`)
	if isrcPattern.String() != expected.String() {
		t.Errorf("pattern string mismatch:\n  got  %q\n  want %q", isrcPattern.String(), expected.String())
	}
}

// --- Constants ---

func TestConstants(t *testing.T) {
	// Verify that string constants used in pipeline logic have expected values.
	// These are referenced by agents and state-machine logic; changes would be breaking.
	cases := map[string]string{
		"statusInProduction": statusInProduction,
		"statusMastering":    statusMastering,
		"statusQC":           statusQC,
		"statusReady":        statusReady,
		"statusSubmitted":    statusSubmitted,
		"statusLive":         statusLive,
		"msgCanceled":        msgCanceled,
		"cmdUseList":         cmdUseList,
	}
	want := map[string]string{
		"statusInProduction": "in-production",
		"statusMastering":    "mastering",
		"statusQC":           "qc",
		"statusReady":        "ready",
		"statusSubmitted":    "submitted",
		"statusLive":         "live",
		"msgCanceled":        "canceled",
		"cmdUseList":         "list",
	}
	for k, got := range cases {
		if got != want[k] {
			t.Errorf("constant %s: got %q, want %q", k, got, want[k])
		}
	}
}
