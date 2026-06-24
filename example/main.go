// SPDX-License-Identifier: Apache-2.0

// Command example demonstrates consuming github.com/automa-saga/version.
//
// Build it with stamped metadata:
//
//	go build -ldflags "\
//	  -X github.com/automa-saga/version.Version=v1.2.3 \
//	  -X github.com/automa-saga/version.Commit=$(git rev-parse HEAD) \
//	  -X github.com/automa-saga/version.Date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
//	  -o example ./example
//
// Then run it with a format argument (text, json):
//
//	./example          # text
//	./example json
package main

import (
	"fmt"
	"os"

	"github.com/automa-saga/version"
)

func main() {
	format := "text"
	if len(os.Args) > 1 {
		format = os.Args[1]
	}

	out, err := version.Get().Format(format)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(out)
}
