# version

[![Go Reference](https://pkg.go.dev/badge/github.com/automa-saga/version.svg)](https://pkg.go.dev/github.com/automa-saga/version)

A tiny, **dependency-free** Go package for exposing a binary's build metadata —
version, commit, build date, and Go runtime version — in one place that can be
shared across every binary in your organisation.

Values are injected at build time with the Go linker's `-X` flag and fall back
to development placeholders when unstamped, so a plain `go run` / `go build`
correctly reports itself as `dev`.

## Why

Most projects re-implement the same version package over and over. This package
standardises it:

- **Zero third-party dependencies** — stdlib only (`encoding/json`, `runtime`).
- **One injection mechanism everywhere** — `-ldflags -X`, the de-facto Go idiom.
  No generated files, no `go generate`, no per-OS shell scripts.
- **Works across many binaries** — because `-X` can stamp any package linked
  into a binary, each binary injects these same variables at its own build step
  with its own values. No registry, no per-binary subpackages.
- **Framework-agnostic** — no CLI framework baked in. Wire it into Cobra, urfave,
  flag, or anything else yourself (see below).

## Install

```bash
go get github.com/automa-saga/version
```

Requires Go 1.26+.

## Stamping the build

Inject values by targeting the exported package variables `Version`, `Commit`,
and `Date`:

```bash
go build -ldflags "\
  -X github.com/automa-saga/version.Version=v1.2.3 \
  -X github.com/automa-saga/version.Commit=$(git rev-parse HEAD) \
  -X github.com/automa-saga/version.Date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -o mybinary ./cmd/mybinary
```

### With [Task](https://taskfile.dev)

```yaml
vars:
  VERSION:
    sh: git describe --tags --always
  COMMIT:
    sh: git rev-parse HEAD
  DATE:
    sh: date -u +%Y-%m-%dT%H:%M:%SZ
  LDFLAGS: >-
    -X github.com/automa-saga/version.Version={{.VERSION}}
    -X github.com/automa-saga/version.Commit={{.COMMIT}}
    -X github.com/automa-saga/version.Date={{.DATE}}

tasks:
  build:
    cmds:
      - go build -ldflags "{{.LDFLAGS}}" -o bin/mybinary ./cmd/mybinary
```

### Unstamped builds

Without `-ldflags`, the defaults apply:

```
Version:    dev
Commit:     none
Go Version: go1.21.x
```

(The `Built:` line is omitted when the date is unset/`unknown`.)

## Reading the version

```go
package main

import (
	"fmt"

	"github.com/automa-saga/version"
)

func main() {
	info := version.Get()

	fmt.Println(info.Version)          // v1.2.3
	fmt.Println(info.ShortCommit(8))   // abc12345
	fmt.Println(info.Text())           // multi-line human-readable block

	json, _ := info.Format("json")     // {"version":"v1.2.3",...}
	fmt.Println(json)
}
```

`Get()` returns an `Info` snapshot:

```go
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date,omitempty"`
	GoVersion string `json:"goVersion"`
}
```

Supported `Format` values are `"json"` (default) and `"text"`.

## Integrating with a CLI

The package is CLI-framework-agnostic. Build the command yourself and call into
`version`. Example with [Cobra](https://github.com/spf13/cobra):

```go
func newVersionCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information and exit",
		RunE: func(cmd *cobra.Command, _ []string) error {
			out, err := version.Get().Format(format)
			if err != nil {
				return err
			}
			cmd.Println(out)
			return nil
		},
	}
	cmd.Flags().StringVarP(&format, "output", "o", "text", "Output format: text or json")
	return cmd
}
```

## Other output formats (e.g. YAML)

To keep this package dependency-free, only `json` and `text` are built in. The
`Info` fields carry both `json` and `yaml` struct tags, so a consumer that wants
YAML can marshal `Info` directly with its own YAML library:

```go
import "gopkg.in/yaml.v3"

out, _ := yaml.Marshal(version.Get())
```

## Using it across multiple binaries

Each binary stamps the same variables at its own build step, so several binaries
— even within one module — can carry different versions:

```bash
# CLI binary
go build -ldflags "-X github.com/automa-saga/version.Version=v1.2.3" ./cmd/cli

# Daemon binary, released on its own cadence
go build -ldflags "-X github.com/automa-saga/version.Version=daemon-v0.4.0" ./cmd/daemon
```

No shared registry or per-binary wrapper packages are needed — the linker writes
the right value into each binary independently.

## License

[Apache 2.0](LICENSE)
