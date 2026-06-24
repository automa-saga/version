// SPDX-License-Identifier: Apache-2.0

package version

import (
	"encoding/json"
	"runtime"
	"strings"
	"testing"
)

func TestGetIncludesGoVersion(t *testing.T) {
	got := Get()
	if got.GoVersion != runtime.Version() {
		t.Errorf("GoVersion = %q, want %q", got.GoVersion, runtime.Version())
	}
	// Defaults apply for unstamped test builds.
	if got.Version != "dev" {
		t.Errorf("Version = %q, want %q", got.Version, "dev")
	}
}

func TestFormatJSON(t *testing.T) {
	info := Info{Version: "v1.2.3", Commit: "abc123", Date: "2026-06-24T00:00:00Z", GoVersion: "go1.21"}

	for _, format := range []string{"", "json", "JSON"} {
		out, err := info.Format(format)
		if err != nil {
			t.Fatalf("Format(%q) error: %v", format, err)
		}
		var round Info
		if err := json.Unmarshal([]byte(out), &round); err != nil {
			t.Fatalf("Format(%q) produced invalid json %q: %v", format, out, err)
		}
		if round != info {
			t.Errorf("Format(%q) round-trip = %+v, want %+v", format, round, info)
		}
	}
}

func TestFormatText(t *testing.T) {
	info := Info{Version: "v1.2.3", Commit: "abc123", Date: "2026-06-24T00:00:00Z", GoVersion: "go1.21"}
	out, err := info.Format("text")
	if err != nil {
		t.Fatalf("Format(text) error: %v", err)
	}
	for _, want := range []string{"v1.2.3", "abc123", "2026-06-24T00:00:00Z", "go1.21"} {
		if !strings.Contains(out, want) {
			t.Errorf("Format(text) = %q, missing %q", out, want)
		}
	}
}

func TestTextOmitsUnknownDate(t *testing.T) {
	for _, date := range []string{"", "unknown"} {
		out := Info{Version: "dev", Commit: "none", Date: date, GoVersion: "go1.21"}.Text()
		if strings.Contains(out, "Built:") {
			t.Errorf("Text() with date %q should omit Built line, got %q", date, out)
		}
	}
}

func TestFormatUnsupported(t *testing.T) {
	if _, err := (Info{}).Format("xml"); err == nil {
		t.Error("Format(xml) expected error, got nil")
	}
}

func TestShortCommit(t *testing.T) {
	info := Info{Commit: "0123456789abcdef"}
	if got := info.ShortCommit(8); got != "01234567" {
		t.Errorf("ShortCommit(8) = %q, want %q", got, "01234567")
	}
	if got := info.ShortCommit(64); got != info.Commit {
		t.Errorf("ShortCommit(64) = %q, want full commit", got)
	}
	if got := info.ShortCommit(0); got != info.Commit {
		t.Errorf("ShortCommit(0) = %q, want full commit", got)
	}
}

func TestStringMatchesText(t *testing.T) {
	info := Info{Version: "v1.0.0", Commit: "deadbeef", GoVersion: "go1.21"}
	if info.String() != info.Text() {
		t.Error("String() should equal Text()")
	}
}
