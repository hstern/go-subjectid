// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hstern/go-subjectid"
)

// TestForwardCompatUnknownFormatRoundTrip pins the forward-
// compatibility contract on [subjectid.UnknownFormat]: a payload
// whose "format" discriminator is outside the library's built-in
// set and unregistered as an extension parses cleanly into an
// UnknownFormat carrier rather than erroring, and the carrier
// round-trips byte-stably so the unknown payload survives any
// pass through the codec.
//
// This is the rule that makes the library safe to ship in a
// long-lived consumer: the IANA Subject Identifier Formats
// registry may grow new entries after the consumer's build is
// frozen, and the library must not turn those into hard errors.
func TestForwardCompatUnknownFormatRoundTrip(t *testing.T) {
	// A synthesized payload using a format name guaranteed to be
	// outside the built-in set. The body carries an arbitrary
	// extension shape ("widget_id" + an int) so we can confirm
	// the entire wire payload survives.
	wire := []byte(`{"format":"org.example.future","widget_id":42}`)

	id, err := subjectid.Parse(wire)
	if err != nil {
		t.Fatalf("Parse(unknown format): %v", err)
	}

	unk, ok := id.(subjectid.UnknownFormat)
	if !ok {
		t.Fatalf("Parse returned %T, want subjectid.UnknownFormat", id)
	}
	if got, want := unk.Format(), "org.example.future"; got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
	if got, want := unk.FormatName, "org.example.future"; got != want {
		t.Errorf("FormatName = %q, want %q", got, want)
	}
	if verr := unk.Validate(); verr != nil {
		t.Errorf("Validate() on UnknownFormat: %v, want nil", verr)
	}

	got, merr := json.Marshal(unk)
	if merr != nil {
		t.Fatalf("Marshal: %v", merr)
	}
	if !bytes.Equal(got, wire) {
		t.Errorf("byte-stability broken:\n got:  %s\n want: %s", got, wire)
	}
}

// TestForwardCompatUnknownFormatInsideAliases verifies the
// composite case: an aliases payload whose inner identifiers
// include an unknown format does not fail the surrounding parse,
// and the resulting AliasesID's Identifiers slice carries the
// unknown alongside the known.
//
// This guards against a regression where an over-eager
// validation-on-parse path would reject the whole aliases for a
// single unknown inner element. The Postel reading: the wire is
// well-formed JSON with a recognizable structure; the unknown
// inner is preserved for the caller to inspect.
func TestForwardCompatUnknownFormatInsideAliases(t *testing.T) {
	wire := []byte(`{"format":"aliases","identifiers":[` +
		`{"format":"email","email":"user@example.com"},` +
		`{"format":"org.example.future","widget_id":42}]}`)

	id, err := subjectid.Parse(wire)
	if err != nil {
		t.Fatalf("Parse(aliases with unknown inner): %v", err)
	}
	aliases, ok := id.(subjectid.AliasesID)
	if !ok {
		t.Fatalf("Parse returned %T, want subjectid.AliasesID", id)
	}
	if got, want := len(aliases.Identifiers), 2; got != want {
		t.Fatalf("len(Identifiers) = %d, want %d", got, want)
	}
	if _, ok := aliases.Identifiers[0].(subjectid.EmailID); !ok {
		t.Errorf("Identifiers[0] = %T, want subjectid.EmailID", aliases.Identifiers[0])
	}
	if _, ok := aliases.Identifiers[1].(subjectid.UnknownFormat); !ok {
		t.Errorf("Identifiers[1] = %T, want subjectid.UnknownFormat", aliases.Identifiers[1])
	}
}
