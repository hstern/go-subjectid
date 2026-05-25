// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"errors"
	"strings"
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
		"acct:1@2",
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
		input   string
		rule    string
		message string
	}{
		{"", "required", "non-empty"},
		{"mailto:user@example.com", "format:account", "acct:"},
		{"acct:user", "format:account", `"@"`},
		{"acct:@example.com", "format:account", "userpart is empty"},
		{"acct:user@", "format:account", "host is empty"},
		{"acct:a@b@c", "format:account", "second one"},
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
			if !strings.Contains(ve.Reason, tc.message) {
				t.Errorf("Reason = %q, want substring %q", ve.Reason, tc.message)
			}
		})
	}
}
