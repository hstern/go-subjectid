// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import "encoding/json"

// UnknownFormat carries a Subject Identifier whose format is not
// in this build's built-in set and was not registered via
// RegisterFormat. It is the package's forward-compatibility
// carrier: a library compiled today, decoding a payload that uses
// a format the IANA registry adds in 2027, returns an
// UnknownFormat rather than an error so the value can still
// round-trip and the caller can branch on whatever subset of
// formats it actually understands.
//
// The wire bytes are preserved verbatim in Raw so that re-encoding
// produces the original payload byte-for-byte (modulo JSON
// whitespace canonicalization). Raw is intentionally a
// [json.RawMessage] rather than a map[string]any: interop
// scenarios often pin exact JSON bytes, and a map reorders its
// keys on every encode.
type UnknownFormat struct {
	// FormatName is the value of the JSON "format" member as
	// decoded — the discriminator the library did not recognize.
	FormatName string

	// Raw is the entire JSON object that was decoded, byte-for-
	// byte as it appeared on the wire. Re-encoding returns these
	// bytes unchanged.
	Raw json.RawMessage
}

// Format returns the unrecognized format discriminator. The
// returned value is the FormatName field, not a fixed string —
// UnknownFormat is the one type whose Format method varies per
// value. See [SubjectIdentifier.Format].
func (u UnknownFormat) Format() string { return u.FormatName }

// Validate is a no-op. UnknownFormat exists precisely to round-
// trip values whose well-formedness rules the library does not
// know; performing validation on it would defeat the purpose of
// the carrier.
func (UnknownFormat) Validate() error { return nil }

func (UnknownFormat) sealed() {}

// Compile-time assertion that UnknownFormat satisfies SubjectIdentifier
// in value form — the canonical dynamic form Parse returns.
var _ SubjectIdentifier = UnknownFormat{}
