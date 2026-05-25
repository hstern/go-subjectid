// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/hstern/go-subjectid"
)

func TestEmailIDFormatIsEmail(t *testing.T) {
	var e subjectid.EmailID
	if got, want := e.Format(), "email"; got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
}

func TestEmailIDValidateAcceptsBareAddrSpec(t *testing.T) {
	cases := []string{
		"user@example.com",
		"user.name+tag@sub.example.org",
		"a@b",
		"user@[127.0.0.1]",
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			e := subjectid.EmailID{Email: in}
			if err := e.Validate(); err != nil {
				t.Errorf("Validate(%q) = %v, want nil", in, err)
			}
		})
	}
}

func TestEmailIDValidateRejectsInvalidInputs(t *testing.T) {
	cases := []struct {
		input   string
		rule    string
		message string
	}{
		{"", "required", "non-empty"},
		{`"quoted local"@example.com`, "format:email", "quoted-string"},
		{"not-an-email", "format:email", "RFC 5322"},
		{`User <user@example.com>`, "format:email", "name-addr"},
		{"user@", "format:email", "RFC 5322"},
		{"@example.com", "format:email", "RFC 5322"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			err := subjectid.EmailID{Email: tc.input}.Validate()
			if err == nil {
				t.Fatalf("Validate(%q): got nil, want a *ValidationError", tc.input)
			}
			var ve *subjectid.ValidationError
			if !errors.As(err, &ve) {
				t.Fatalf("err type = %T, want *ValidationError", err)
			}
			if ve.Rule != tc.rule {
				t.Errorf("Rule = %q, want %q", ve.Rule, tc.rule)
			}
			if ve.Format != "email" {
				t.Errorf("Format = %q, want %q", ve.Format, "email")
			}
			if !strings.Contains(ve.Reason, tc.message) {
				t.Errorf("Reason = %q, want substring %q", ve.Reason, tc.message)
			}
		})
	}
}

func TestEmailIDValidateDoesNotCanonicalize(t *testing.T) {
	// Spec: recipient SHOULD canonicalize; the library does not.
	// Validation must accept the input verbatim — mixed case,
	// plus-addressing, etc.
	const in = "User.Name+Promotions@Example.COM"
	e := subjectid.EmailID{Email: in}
	if err := e.Validate(); err != nil {
		t.Fatalf("Validate(%q) = %v, want nil", in, err)
	}
	if e.Email != in {
		t.Errorf("Email mutated by Validate: got %q, want %q", e.Email, in)
	}
}
