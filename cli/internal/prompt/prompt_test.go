package prompt

import (
	"bufio"
	"strings"
	"testing"
)

// injectInput replaces the package-level scanner with one backed by s,
// and restores it after the test.
func injectInput(t *testing.T, s string) {
	t.Helper()
	old := scanner
	scanner = bufio.NewScanner(strings.NewReader(s))
	t.Cleanup(func() { scanner = old })
}

// --- Ask ---

func TestAsk_UserEntersValue(t *testing.T) {
	injectInput(t, "My Input\n")
	got, err := Ask("Label", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "My Input" {
		t.Errorf("got %q, want %q", got, "My Input")
	}
}

func TestAsk_EmptyInputReturnsDefault(t *testing.T) {
	injectInput(t, "\n")
	got, err := Ask("Label", "default-value")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "default-value" {
		t.Errorf("got %q, want %q", got, "default-value")
	}
}

func TestAsk_WhitespaceOnlyReturnsDefault(t *testing.T) {
	injectInput(t, "   \n")
	got, err := Ask("Label", "fallback")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "fallback" {
		t.Errorf("got %q, want %q", got, "fallback")
	}
}

func TestAsk_EOFReturnsDefault(t *testing.T) {
	injectInput(t, "") // no newline → EOF immediately
	got, err := Ask("Label", "eof-default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "eof-default" {
		t.Errorf("got %q, want %q", got, "eof-default")
	}
}

func TestAsk_NoDefault_ReturnsEmpty(t *testing.T) {
	injectInput(t, "\n")
	got, err := Ask("Label", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

// --- Confirm ---

func TestConfirm_Yes(t *testing.T) {
	for _, input := range []string{"y\n", "Y\n", "yes\n", "YES\n"} {
		injectInput(t, input)
		got, err := Confirm("Proceed?")
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", input, err)
		}
		if !got {
			t.Errorf("input %q: expected true", input)
		}
	}
}

func TestConfirm_No(t *testing.T) {
	for _, input := range []string{"n\n", "N\n", "no\n", "\n"} {
		injectInput(t, input)
		got, err := Confirm("Proceed?")
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", input, err)
		}
		if got {
			t.Errorf("input %q: expected false", input)
		}
	}
}

func TestConfirm_EOF(t *testing.T) {
	injectInput(t, "")
	got, err := Confirm("Proceed?")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got {
		t.Error("expected false on EOF")
	}
}

// --- Select ---

func TestSelect_ByNumber(t *testing.T) {
	injectInput(t, "2\n")
	got, err := Select("Choose", []string{"alpha", "beta", "gamma"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "beta" {
		t.Errorf("got %q, want %q", got, "beta")
	}
}

func TestSelect_ByName(t *testing.T) {
	injectInput(t, "gamma\n")
	got, err := Select("Choose", []string{"alpha", "beta", "gamma"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "gamma" {
		t.Errorf("got %q, want %q", got, "gamma")
	}
}

func TestSelect_ByNameCaseInsensitive(t *testing.T) {
	injectInput(t, "ALPHA\n")
	got, err := Select("Choose", []string{"alpha", "beta"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "alpha" {
		t.Errorf("got %q, want %q", got, "alpha")
	}
}

func TestSelect_InvalidInput(t *testing.T) {
	injectInput(t, "xyz\n")
	_, err := Select("Choose", []string{"alpha", "beta"})
	if err == nil {
		t.Error("expected error for invalid selection")
	}
}

func TestSelect_EOF(t *testing.T) {
	injectInput(t, "")
	_, err := Select("Choose", []string{"alpha", "beta"})
	if err == nil {
		t.Error("expected error on EOF (no selection made)")
	}
}
