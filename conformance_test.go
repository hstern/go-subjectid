// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hstern/go-subjectid"
	"github.com/hstern/go-subjectid/internal/specfixtures"
)

// TestSpecFixturesRoundTripByteStable iterates every embedded RFC 9493
// §3 example payload and asserts the three-step contract:
//
//  1. [subjectid.Parse] decodes the wire bytes to a built-in
//     concrete [subjectid.SubjectIdentifier] (no [UnknownFormat]
//     fall-through).
//  2. Validate() returns nil on the parsed value, since every spec
//     example is, by definition, well-formed.
//  3. json.Marshal of the parsed value produces bytes identical to
//     the original wire payload — modulo nothing, since the
//     embedded fixtures are stored as the canonical compact form.
//
// This is the library's conformance claim: every illustrative
// payload from the spec round-trips byte-stably through the codec
// + validator.
func TestSpecFixturesRoundTripByteStable(t *testing.T) {
	fixtures := specfixtures.All()
	if len(fixtures) == 0 {
		t.Fatal("specfixtures.All() returned no fixtures; embed directive may be broken")
	}

	// Every built-in format from the IANA registry must be
	// represented in the fixture set. Sanity check before we walk
	// the rows so a missing fixture surfaces as a clear failure
	// rather than as an incomplete iteration.
	want := map[string]bool{
		"account": false, "email": false, "iss_sub": false, "opaque": false,
		"phone_number": false, "did": false, "uri": false, "aliases": false,
	}
	for _, fx := range fixtures {
		if _, ok := want[fx.Format]; ok {
			want[fx.Format] = true
		}
	}
	for name, seen := range want {
		if !seen {
			t.Errorf("missing fixture for built-in format %q", name)
		}
	}

	for _, fx := range fixtures {
		t.Run(fx.Format, func(t *testing.T) {
			id, err := subjectid.Parse(fx.Wire)
			if err != nil {
				t.Fatalf("Parse(%s): %v", fx.Wire, err)
			}
			if got, want := id.Format(), fx.Format; got != want {
				t.Errorf("parsed Format() = %q, want %q", got, want)
			}
			if _, isUnknown := id.(subjectid.UnknownFormat); isUnknown {
				t.Errorf("parsed into UnknownFormat; built-in dispatch did not fire")
			}
			if verr := id.Validate(); verr != nil {
				t.Errorf("Validate() on spec example: %v", verr)
			}
			got, merr := json.Marshal(id)
			if merr != nil {
				t.Fatalf("Marshal: %v", merr)
			}
			if !bytes.Equal(got, fx.Wire) {
				t.Errorf("byte-stability broken:\n got:  %s\n want: %s", got, fx.Wire)
			}
		})
	}
}
