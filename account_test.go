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
		// Canonical reg-name forms.
		"acct:example.user@service.example.com",
		"acct:bob@example.org",
		"acct:user-name+tag@sub.example.com",
		"acct:1@2",
		// Percent-encoded octets in userpart (RFC 7565 §3 via
		// RFC 3986 pct-encoded), both cases per RFC 3986 §6.2.2.1.
		"acct:user%20name@example.com",
		"acct:user%2aname@example.com",
		"acct:user%2Aname@example.com",
		// Sub-delims in userpart per RFC 3986 §2:
		//   "!" / "$" / "&" / "'" / "(" / ")" / "*" / "+" / "," / ";" / "="
		"acct:bob!doe$@example.com",
		// IPv4 host per RFC 3986 §3.2.2 IPv4address.
		"acct:user@192.0.2.1",
		// IPv6 host per RFC 3986 §3.2.2 IP-literal.
		"acct:user@[2001:db8::1]",
		"acct:user@[::1]",
		// IPvFuture per RFC 3986 §3.2.2.
		"acct:user@[v1.fe80::a+b]",
		// Empty reg-name is permitted by RFC 3986 §3.2.2:
		//   reg-name = *( unreserved / pct-encoded / sub-delims )
		// The library does not impose a stricter rule on top.
		"acct:user@",
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
		// Wrong scheme.
		{"mailto:user@example.com", "format:account"},
		// Uppercase scheme — net/url.Parse normalizes silently;
		// the explicit literal-prefix check rejects.
		{"ACCT:user@example.com", "format:account"},
		// No @.
		{"acct:user", "format:account"},
		// Empty userpart.
		{"acct:@example.com", "format:account"},
		// Two @.
		{"acct:a@b@c", "format:account"},
		// Whitespace not in the unreserved/pct-encoded/sub-delims
		// character classes.
		{"acct:user name@example.com", "format:account"},
		{"acct:user@exa mple.com", "format:account"},
		// Fragment, query, path — none in the RFC 7565 ABNF.
		{"acct:user@example.com#frag", "format:account"},
		{"acct:user@example.com?q=1", "format:account"},
		// Bad IP-literal: bracketed but neither IPv6 nor IPvFuture.
		{"acct:user@[bad]", "format:account"},
		// Bracketed IPv4 — not allowed by RFC 3986 §3.2.2
		// (only IPv6 and IPvFuture go inside brackets).
		{"acct:user@[192.0.2.1]", "format:account"},
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
