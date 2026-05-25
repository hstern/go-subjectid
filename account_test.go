// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"errors"
	"testing"

	"github.com/hstern/go-subjectid"
)

func TestAccountIDFormatIsAccount(t *testing.T) {
	var a subjectid.AccountID
	if got, want := a.Format(), "account"; got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
}

func TestAccountIDValidateAcceptsWellFormedAcctURIs(t *testing.T) {
	cases := []string{
		"acct:example.user@service.example.com",
		"acct:bob@example.org",
		"acct:user-name+tag@sub.example.com",
		// Single-character userpart and host: the ABNF lower
		// bound is one character on each side.
		"acct:1@2",
		// Percent-encoded octets in userpart (RFC 7565 §3 via
		// RFC 3986 pct-encoded).
		"acct:user%20name@example.com",
		// Sub-delims in userpart per RFC 3986: ! $ & ' ( ) * + , ; =
		"acct:bob!doe$@example.com",
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			if err := (subjectid.AccountID{URI: in}.Validate()); err != nil {
				t.Errorf("Validate(%q) = %v, want nil", in, err)
			}
		})
	}
}

func TestAccountIDValidateRejectsMalformedAcctURIs(t *testing.T) {
	cases := []struct {
		input string
		rule  string
	}{
		{"", "required"},
		{"mailto:user@example.com", "format:account"},
		{"ACCT:user@example.com", "format:account"},
		{"acct:user", "format:account"},
		{"acct:@example.com", "format:account"},
		{"acct:user@", "format:account"},
		{"acct:a@b@c", "format:account"},
		{"acct:user name@example.com", "format:account"},
		{"acct:user@exa mple.com", "format:account"},
		{"acct:user@example.com#frag", "format:account"},
		{"acct:user@example.com?q=1", "format:account"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			err := subjectid.AccountID{URI: tc.input}.Validate()
			if err == nil {
				t.Fatalf("Validate(%q): got nil, want *ValidationError", tc.input)
			}
			var ve *subjectid.ValidationError
			if !errors.As(err, &ve) {
				t.Fatalf("err type = %T, want *ValidationError", err)
			}
			if ve.Rule != tc.rule {
				t.Errorf("Rule = %q, want %q", ve.Rule, tc.rule)
			}
			if ve.Format != "account" {
				t.Errorf("Format = %q, want %q", ve.Format, "account")
			}
		})
	}
}
