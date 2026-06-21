// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"errors"
	"testing"

	"github.com/hstern/go-subjectid"
)

// TestParseReturnsValueForm pins the canonical-form decision:
// Parse returns the value form of an identifier, so the
// dynamic type a caller reads back from Parse is identical to the
// value literal they would write by hand (subjectid.IssSubID, not
// *subjectid.IssSubID). A direct value assertion must succeed and the
// pointer assertion must fail.
func TestParseReturnsValueForm(t *testing.T) {
	id, err := subjectid.Parse([]byte(
		`{"format":"iss_sub","iss":"https://idp.example.com/","sub":"alice"}`))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if _, ok := id.(subjectid.IssSubID); !ok {
		t.Fatalf("Parse returned %T, want value form subjectid.IssSubID", id)
	}
	if _, ok := id.(*subjectid.IssSubID); ok {
		t.Fatalf("Parse returned pointer form %T, want value-canonical", id)
	}
}

// TestParseUnknownReturnsValueForm checks that the forward-compat
// UnknownFormat carrier also escapes Parse in value form, consistent
// with the canonical-form rule.
func TestParseUnknownReturnsValueForm(t *testing.T) {
	id, err := subjectid.Parse([]byte(
		`{"format":"org.example.unknown","opaque_id":"xyz"}`))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if _, ok := id.(subjectid.UnknownFormat); !ok {
		t.Fatalf("Parse returned %T, want value form subjectid.UnknownFormat", id)
	}
}

// TestParseNestedAliasesFromWireRejected guards the latent bug the
// pointer/value split caused: a nested aliases identifier arriving
// from the wire (which Parse populates recursively) must be rejected
// by Validate with ErrNestedAliases, exactly as a hand-built nested
// alias is. Before the canonical-form fix the value-typed nested
// check missed the pointer-form wire element and this slipped through.
func TestParseNestedAliasesFromWireRejected(t *testing.T) {
	raw := []byte(`{"format":"aliases","identifiers":[` +
		`{"format":"aliases","identifiers":[` +
		`{"format":"email","email":"user@example.com"}]}]}`)
	id, err := subjectid.Parse(raw)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if err := id.Validate(); !errors.Is(err, subjectid.ErrNestedAliases) {
		t.Fatalf("Validate() of wire-parsed nested aliases = %v, want ErrNestedAliases", err)
	}
}
