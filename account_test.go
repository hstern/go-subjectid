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

// TestAccountIDValidateAcceptsWellFormedAcctURIs locks in the
// well-formed acct URI shapes RFC 7565 §3 and RFC 3986 §3.2.2
// admit. Each line cites the spec clause that justifies it.
func TestAccountIDValidateAcceptsWellFormedAcctURIs(t *testing.T) {
	cases := []string{
		// RFC 9493 §3.2.1 illustrative example.
		"acct:example.user@service.example.com",

		// Canonical reg-name forms — RFC 3986 §3.2.2 reg-name
		// = *( unreserved / pct-encoded / sub-delims ).
		"acct:bob@example.org",
		"acct:user-name+tag@sub.example.com",
		"acct:1@2",

		// Percent-encoded octets in userpart per RFC 7565 §3 via
		// RFC 3986 pct-encoded; hex digits are case-insensitive
		// per RFC 3986 §6.2.2.1.
		"acct:user%20name@example.com",
		"acct:user%2aname@example.com",
		"acct:user%2Aname@example.com",

		// Sub-delims in userpart per RFC 3986 §2:
		//   "!" / "$" / "&" / "'" / "(" / ")" / "*" / "+" / "," / ";" / "="
		"acct:bob!doe$@example.com",
		"acct:a'b(c)d*e+f,g;h=i@example.com",

		// Unreserved set per RFC 3986 §2.3:
		//   ALPHA / DIGIT / "-" / "." / "_" / "~"
		"acct:a-b.c_d~e@example.com",

		// IPv4 host per RFC 3986 §3.2.2 IPv4address.
		"acct:user@192.0.2.1",

		// IPv6 host per RFC 3986 §3.2.2 IP-literal.
		"acct:user@[2001:db8::1]",
		"acct:user@[::1]",
		"acct:user@[fe80::1%25eth0]", // zone-id pct-encoded "%25" per RFC 6874

		// IPvFuture per RFC 3986 §3.2.2.
		"acct:user@[v1.fe80::a+b]",
		"acct:user@[V7.deadbeef:]",

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

// TestAccountIDValidateRejectsMalformedAcctURIs locks in the
// failures: each entry cites the spec clause that makes the input
// invalid. Callers branch on the wrapped sentinel via errors.Is;
// Format is always "account" for AccountID failures.
func TestAccountIDValidateRejectsMalformedAcctURIs(t *testing.T) {
	cases := []struct {
		input string
		want  error
	}{
		// RFC 9493 §3 — required members must be non-empty.
		{"", subjectid.ErrRequired{}},

		// Wrong scheme — RFC 7565 §3 fixes "acct" as the scheme.
		{"mailto:user@example.com", subjectid.ErrFormatAccount},
		{"http://user@example.com", subjectid.ErrFormatAccount},

		// Uppercase scheme — IANA URI Scheme Registry stores the
		// scheme in lowercase (RFC 7565 §3 / RFC 3986 §3.1).
		// net/url.Parse normalizes case silently; reject explicitly.
		{"ACCT:user@example.com", subjectid.ErrFormatAccount},
		{"Acct:user@example.com", subjectid.ErrFormatAccount},

		// No "@" between userpart and host.
		{"acct:user", subjectid.ErrFormatAccount},
		{"acct:userhost.com", subjectid.ErrFormatAccount},

		// Empty userpart — RFC 7565 §3 userpart = 1*(...) requires
		// at least one character.
		{"acct:@example.com", subjectid.ErrFormatAccount},

		// Two "@" — RFC 7565 §3 grammar admits exactly one.
		{"acct:a@b@c", subjectid.ErrFormatAccount},

		// Whitespace not in unreserved/pct-encoded/sub-delims.
		{"acct:user name@example.com", subjectid.ErrFormatAccount},
		{"acct:user@exa mple.com", subjectid.ErrFormatAccount},

		// Fragment, query, path — none appear in the RFC 7565 ABNF
		// acct-uri production.
		{"acct:user@example.com#frag", subjectid.ErrFormatAccount},
		{"acct:user@example.com?q=1", subjectid.ErrFormatAccount},

		// Bad IP-literal: bracketed but neither IPv6 nor IPvFuture.
		{"acct:user@[bad]", subjectid.ErrFormatAccount},

		// Bracketed IPv4 — RFC 3986 §3.2.2 puts only IPv6 and
		// IPvFuture inside brackets.
		{"acct:user@[192.0.2.1]", subjectid.ErrFormatAccount},

		// Unclosed bracket.
		{"acct:user@[2001:db8::1", subjectid.ErrFormatAccount},

		// Bare "%" not followed by two hex digits — pct-encoded
		// per RFC 3986 §2.1 is exactly "%" HEXDIG HEXDIG.
		{"acct:user%@example.com", subjectid.ErrFormatAccount},
		{"acct:user%2@example.com", subjectid.ErrFormatAccount},
		{"acct:user%gg@example.com", subjectid.ErrFormatAccount},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			err := subjectid.AccountID{URI: tc.input}.Validate()
			if err == nil {
				t.Fatalf("Validate(%q): got nil, want non-nil error", tc.input)
			}
			if !errors.Is(err, tc.want) {
				t.Errorf("errors.Is(err, %v) = false, want true (err = %v)", tc.want, err)
			}
		})
	}
}
