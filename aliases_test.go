// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"errors"
	"testing"

	"github.com/hstern/go-subjectid"
)

func TestAliasesIDFormatIsAliases(t *testing.T) {
	var a subjectid.AliasesID
	if got, want := a.Format(), "aliases"; got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
}

// TestAliasesIDIdentifiersAcceptsSubjectIdentifier documents the
// type-level contract: AliasesID.Identifiers is a slice of the
// sealed SubjectIdentifier interface, so it accepts any built-in
// (or registered extension) at compile time. A nested AliasesID
// is type-legal but spec-illegal — RFC 9493 §3.2.8 forbids it —
// and Validate enforces the prohibition at the well-formedness
// layer rather than at the type layer (see the reject table).
func TestAliasesIDIdentifiersAcceptsSubjectIdentifier(t *testing.T) {
	inner := subjectid.AliasesID{}
	outer := subjectid.AliasesID{
		Identifiers: []subjectid.SubjectIdentifier{inner},
	}
	if got, want := len(outer.Identifiers), 1; got != want {
		t.Fatalf("len(Identifiers) = %d, want %d", got, want)
	}
	if got, want := outer.Identifiers[0].Format(), "aliases"; got != want {
		t.Errorf("Identifiers[0].Format() = %q, want %q", got, want)
	}
}

// TestAliasesIDValidateAccepts locks in the well-formed shapes:
// a non-empty list of non-aliases inner identifiers, each itself
// valid per its own Validate.
func TestAliasesIDValidateAccepts(t *testing.T) {
	cases := []struct {
		name string
		ids  []subjectid.SubjectIdentifier
	}{
		{
			// RFC 9493 §3.2.8 illustrative example.
			name: "spec example: email + account",
			ids: []subjectid.SubjectIdentifier{
				subjectid.EmailID{Email: "user@example.com"},
				subjectid.AccountID{URI: "acct:user@example.com"},
			},
		},
		{
			name: "single inner email",
			ids: []subjectid.SubjectIdentifier{
				subjectid.EmailID{Email: "user@example.com"},
			},
		},
		{
			name: "four formats mixed",
			ids: []subjectid.SubjectIdentifier{
				subjectid.EmailID{Email: "user@example.com"},
				subjectid.AccountID{URI: "acct:user@example.com"},
				subjectid.PhoneNumberID{PhoneNumber: "+12065550100"},
				subjectid.OpaqueID{ID: "abc"},
			},
		},
		{
			name: "iss_sub + opaque",
			ids: []subjectid.SubjectIdentifier{
				subjectid.IssSubID{Iss: "https://issuer.example.com/", Sub: "1"},
				subjectid.OpaqueID{ID: "x"},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := subjectid.AliasesID{Identifiers: tc.ids}
			if err := a.Validate(); err != nil {
				t.Errorf("Validate(%s) = %v, want nil", tc.name, err)
			}
		})
	}
}

// TestAliasesIDValidateRejects locks in the failure modes per
// RFC 9493 §3.2.8:
//
//   - identifiers MUST be present and non-empty (RFC 9493 §3).
//   - identifiers MUST NOT contain a nested aliases identifier.
//   - each inner identifier MUST itself be valid; the error
//     bubbles up as an [errors.Join] of the per-element
//     sentinel errors; callers recover them with [errors.Is]
//     or [errors.As].
func TestAliasesIDValidateRejects(t *testing.T) {
	cases := []struct {
		name       string
		ids        []subjectid.SubjectIdentifier
		wantErr    error
		wantFormat string
	}{
		{
			name:       "nil identifiers",
			ids:        nil,
			wantErr:    subjectid.ErrRequired{},
			wantFormat: "aliases",
		},
		{
			name:       "empty identifiers",
			ids:        []subjectid.SubjectIdentifier{},
			wantErr:    subjectid.ErrRequired{},
			wantFormat: "aliases",
		},
		{
			name: "single nested aliases",
			ids: []subjectid.SubjectIdentifier{
				subjectid.AliasesID{Identifiers: []subjectid.SubjectIdentifier{
					subjectid.EmailID{Email: "user@example.com"},
				}},
			},
			wantErr:    subjectid.ErrNestedAliases,
			wantFormat: "aliases:0",
		},
		{
			name: "nested aliases at index 1",
			ids: []subjectid.SubjectIdentifier{
				subjectid.EmailID{Email: "user@example.com"},
				subjectid.AliasesID{Identifiers: []subjectid.SubjectIdentifier{
					subjectid.EmailID{Email: "user@example.com"},
				}},
			},
			wantErr:    subjectid.ErrNestedAliases,
			wantFormat: "aliases:1",
		},
		{
			name: "invalid inner email at index 0",
			ids: []subjectid.SubjectIdentifier{
				subjectid.EmailID{Email: ""},
			},
			wantErr:    subjectid.ErrRequired{},
			wantFormat: "aliases:0",
		},
		{
			name: "invalid inner account at index 1",
			ids: []subjectid.SubjectIdentifier{
				subjectid.EmailID{Email: "user@example.com"},
				subjectid.AccountID{URI: "not-an-acct-uri"},
			},
			wantErr:    subjectid.ErrFormatAccount,
			wantFormat: "aliases:1",
		},
		{
			name: "invalid inner phone at index 2",
			ids: []subjectid.SubjectIdentifier{
				subjectid.EmailID{Email: "user@example.com"},
				subjectid.OpaqueID{ID: "x"},
				subjectid.PhoneNumberID{PhoneNumber: "not-a-phone"},
			},
			wantErr:    subjectid.ErrFormatPhoneNumber,
			wantFormat: "aliases:2",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := subjectid.AliasesID{Identifiers: tc.ids}.Validate()
			if err == nil {
				t.Fatalf("Validate(%s): got nil, want non-nil error", tc.name)
			}
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("errors.Is(err, %v) = false, want true (err = %v)", tc.wantErr, err)
			}
		})
	}
}
