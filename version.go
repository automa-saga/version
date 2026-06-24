// SPDX-License-Identifier: Apache-2.0

// Package version exposes a binary's build metadata (version, commit, build
// date, Go version) in a single, dependency-free place that can be shared
// across every binary in an organisation.
//
// Values are injected at build time via the Go linker's -X flag and fall back
// to development placeholders for `go run` / `go build` invocations that don't
// stamp them — so an unstamped build correctly reports itself as "dev".
//
// Inject the values by targeting the exported package variables below. Because
// `-X` can write to any package linked into the binary, every consuming binary
// stamps these same variables at its own build step, each with its own values:
//
//	go build -ldflags "\
//	  -X github.com/automa-saga/version.Version=v1.2.3 \
//	  -X github.com/automa-saga/version.Commit=$(git rev-parse HEAD) \
//	  -X github.com/automa-saga/version.Date=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
//
// Read the running binary's metadata through Get:
//
//	info := version.Get()
//	fmt.Println(info.Text())
//
// The package intentionally has zero third-party dependencies. It renders text
// and JSON using only the standard library; any other format (e.g. YAML) is
// left to the caller, which can marshal the exported Info fields directly.
package version

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

// Build metadata, set via -ldflags "-X github.com/automa-saga/version.<var>=...".
// The defaults are deliberate, human-readable markers for unstamped builds.
var (
	// Version is the semantic version (e.g. v1.2.3), or "dev" for unstamped builds.
	Version = "dev"
	// Commit is the git commit hash the binary was built from, or "none".
	Commit = "none"
	// Date is the build timestamp (RFC 3339), or "unknown".
	Date = "unknown"
)

// Output formats supported by Format.
const (
	FormatJSON = "json"
	FormatText = "text"
)

// Info is a snapshot of the running binary's build metadata.
type Info struct {
	Version   string `json:"version" yaml:"version"`
	Commit    string `json:"commit" yaml:"commit"`
	Date      string `json:"date,omitempty" yaml:"date,omitempty"`
	GoVersion string `json:"goVersion" yaml:"goVersion"`
}

// Get returns a snapshot of the running binary's build metadata, including the
// Go runtime version it was compiled with.
func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		Date:      Date,
		GoVersion: runtime.Version(),
	}
}

// Format renders the Info in the requested format. An empty format defaults to
// JSON. Only "json" and "text" are supported by this package; for other formats
// (e.g. YAML) marshal the Info fields directly with your own dependency.
func (i Info) Format(format string) (string, error) {
	switch strings.ToLower(format) {
	case "", FormatJSON:
		out, err := json.Marshal(i)
		if err != nil {
			return "", fmt.Errorf("marshal version to json: %w", err)
		}
		return string(out), nil
	case FormatText:
		return i.Text(), nil
	default:
		return "", fmt.Errorf("unsupported version format %q (want %q or %q)", format, FormatJSON, FormatText)
	}
}

// Text returns a human-readable, multi-line version string. The build date line
// is omitted when the date was never stamped.
func (i Info) Text() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Version:    %s\n", i.Version)
	fmt.Fprintf(&b, "Commit:     %s\n", i.Commit)
	if i.Date != "" && i.Date != "unknown" {
		fmt.Fprintf(&b, "Built:      %s\n", i.Date)
	}
	fmt.Fprintf(&b, "Go Version: %s", i.GoVersion)
	return b.String()
}

// String implements fmt.Stringer, returning the same value as Text.
func (i Info) String() string { return i.Text() }

// ShortCommit returns the commit hash truncated to n characters (useful for
// compact display). It returns the full commit if it is shorter than n.
func (i Info) ShortCommit(n int) string {
	if n <= 0 || len(i.Commit) <= n {
		return i.Commit
	}
	return i.Commit[:n]
}
