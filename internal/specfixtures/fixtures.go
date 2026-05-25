// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

// Package specfixtures embeds the eight illustrative Subject
// Identifier payloads from RFC 9493 §3 as canonical (compact) JSON
// so the parent package's conformance tests can exercise every
// built-in format against the actual spec wording.
//
// The fixtures are stored as ready-to-compare bytes — each file is
// the exact byte sequence the library's MarshalJSON emits for the
// corresponding parsed value, with no surrounding whitespace. This
// lets a round-trip test compare bytes directly without running
// [encoding/json.Compact] on the input or pinning a "modulo
// whitespace" exemption in the assertion.
//
// The package is internal-only by design — its sole consumer is
// the parent package's `_test.go` files. Promotion to a stable
// public surface would invite extension by external consumers,
// who would then re-anchor their tests on the library's choice of
// spec-example wording; that is a maintenance burden the library
// should not take on for v0.1.
package specfixtures

import (
	"embed"
	"io/fs"
	"path"
	"sort"
	"strings"
)

// raw holds every *.json file in this directory, embedded at build
// time. Going through embed.FS rather than per-file string vars
// keeps the list of formats and the on-disk files in sync — adding
// a new fixture is a single new .json file with no Go-side
// bookkeeping.
//
//go:embed *.json
var raw embed.FS

// A Fixture pairs the format-name discriminator from RFC 9493 §3
// with the exact wire bytes the spec example produces. Format is
// the value of the JSON "format" member (e.g. "account",
// "iss_sub"); Wire is the canonical compact-JSON encoding the
// library's MarshalJSON emits for the parsed value.
type Fixture struct {
	// Format is the RFC 9493 format-name discriminator —
	// the basename of the embedded fixture file.
	Format string

	// Wire is the canonical compact-JSON byte sequence: byte-
	// identical to what MarshalJSON produces for the Format
	// type's representation of the spec's example. Conformance
	// tests compare the result of unmarshal-then-marshal against
	// this slice directly.
	Wire []byte
}

// All returns every embedded fixture, ordered alphabetically by
// Format so iteration is deterministic across runs. Each call
// returns a freshly-allocated slice with freshly-copied Wire
// slices, so callers cannot mutate the embedded payload by
// holding a reference.
func All() []Fixture {
	entries, err := fs.ReadDir(raw, ".")
	if err != nil {
		// Embedded FS is read-only at build time; a read error
		// here would mean the embed directive itself was wrong,
		// which would have failed compilation.
		panic("specfixtures: embedded FS unreadable: " + err.Error())
	}
	out := make([]Fixture, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		body, err := fs.ReadFile(raw, e.Name())
		if err != nil {
			panic("specfixtures: read " + e.Name() + ": " + err.Error())
		}
		cp := make([]byte, len(body))
		copy(cp, body)
		out = append(out, Fixture{
			Format: strings.TrimSuffix(path.Base(e.Name()), ".json"),
			Wire:   cp,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Format < out[j].Format })
	return out
}
