package output_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/open-virtual-label/ovl/internal/output"
)

// captureStdout replaces os.Stdout, runs f, then restores it and returns what was written.
func captureStdout(t *testing.T, f func()) string {
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

// --- JSON mode ---

func TestSetJSON_IsJSON(t *testing.T) {
	output.SetJSON(false)
	if output.IsJSON() {
		t.Error("expected IsJSON=false after SetJSON(false)")
	}
	output.SetJSON(true)
	if !output.IsJSON() {
		t.Error("expected IsJSON=true after SetJSON(true)")
	}
	output.SetJSON(false) // restore
}

// --- JSON output ---

func TestJSON_Valid(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
	}
	got := captureStdout(t, func() {
		if err := output.JSON(payload{Name: "test"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	var result payload
	if err := json.Unmarshal([]byte(got), &result); err != nil {
		t.Errorf("output is not valid JSON: %v\noutput: %s", err, got)
	}
	if result.Name != "test" {
		t.Errorf("got name %q, want %q", result.Name, "test")
	}
}

func TestJSON_Invalid(t *testing.T) {
	// channels cannot be marshaled to JSON — should return an error.
	ch := make(chan int)
	err := output.JSON(ch)
	if err == nil {
		t.Error("expected error marshaling channel")
	}
}

// --- Quiet mode ---

func TestQuiet_SuppressesSuccess(t *testing.T) {
	output.SetQuiet(true)
	got := captureStdout(t, func() {
		output.Success("should not appear")
	})
	output.SetQuiet(false)
	if got != "" {
		t.Errorf("expected no output in quiet mode, got %q", got)
	}
}

func TestQuiet_SuppressesPrint(t *testing.T) {
	output.SetQuiet(true)
	got := captureStdout(t, func() {
		output.Print("silent")
	})
	output.SetQuiet(false)
	if got != "" {
		t.Errorf("expected no output in quiet mode, got %q", got)
	}
}

func TestNonQuiet_ShowsSuccess(t *testing.T) {
	output.SetQuiet(false)
	got := captureStdout(t, func() {
		output.Success("created %s", "test.json")
	})
	if !strings.Contains(got, "test.json") {
		t.Errorf("expected output to contain %q, got %q", "test.json", got)
	}
}

func TestNonQuiet_ShowsPrint(t *testing.T) {
	output.SetQuiet(false)
	got := captureStdout(t, func() {
		output.Print("hello %s", "world")
	})
	if !strings.Contains(got, "hello world") {
		t.Errorf("expected %q in output, got %q", "hello world", got)
	}
}

// Fail writes regardless of quiet mode.
func TestFail_NotSuppressedByQuiet(t *testing.T) {
	output.SetQuiet(true)
	got := captureStdout(t, func() {
		output.Fail("something broke")
	})
	output.SetQuiet(false)
	if !strings.Contains(got, "something broke") {
		t.Errorf("expected Fail output even in quiet mode, got %q", got)
	}
}

// Table should render without panicking.
func TestTable_Renders(t *testing.T) {
	headers := []string{"ID", "Name"}
	rows := [][]string{{"1", "Alice"}, {"2", "Bob"}}
	got := captureStdout(t, func() {
		output.Table(headers, rows)
	})
	if !strings.Contains(got, "Alice") || !strings.Contains(got, "Bob") {
		t.Errorf("table output missing expected rows: %q", got)
	}
}

func TestTable_Empty(t *testing.T) {
	// Should not panic with empty rows.
	captureStdout(t, func() {
		output.Table([]string{"ID"}, [][]string{})
	})
}

// Print uses fmt.Printf-style formatting.
func TestPrint_Formatting(t *testing.T) {
	output.SetQuiet(false)
	got := captureStdout(t, func() {
		output.Print("count: %d", 42)
	})
	if !strings.Contains(got, "count: 42") {
		t.Errorf("got %q", got)
	}
}

// Info writes to stderr; we just verify it doesn't panic.
func TestInfo_NoPanic(t *testing.T) {
	output.SetQuiet(false)
	output.Info("info %s", "message")
	output.SetQuiet(true)
	output.Info("suppressed")
	output.SetQuiet(false)
}

// Error writes to stderr regardless of quiet mode.
func TestError_NoPanic(t *testing.T) {
	output.SetQuiet(true)
	output.Error("error %s", "detail")
	output.SetQuiet(false)
	output.Error("another %s", "error")
}

