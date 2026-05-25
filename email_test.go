// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"errors"
	"testing"

	"github.com/hstern/go-subjectid"
)

func TestEmailIDFormatIsEmail(t *testing.T) {
	var e subjectid.EmailID
	if got, want := e.Format(), "email"; got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
}

// TestEmailIDValidateAcceptsBareAddrSpec locks in the
// addr-spec shapes RFC 5322 §3.4.1 admits. The library performs
// no canonicalization per RFC 9493 §3.2.2 (recipient's choice).
func TestEmailIDValidateAcceptsBareAddrSpec(t *testing.T) {
	cases := []string{
		// RFC 9493 §3.2.2 illustrative example.
		"user@example.com",

		// RFC 5322 §3.4.1 addr-spec = local-part "@" domain.
		"user.name+tag@sub.example.org",
		"a@b",
		"john.doe@example.com",
		"a.b.c@d.e.f.g",

		// Domain literal per RFC 5322 §3.4.1 — domain may be a
		// domain-literal "[" *(...) "]".
		"user@[127.0.0.1]",
		"user@[IPv6:2001:db8::1]",

		// Mixed case preserved verbatim — RFC 9493 §3.2.2 leaves
		// canonicalization to the recipient.
		"User.Name+Promotions@Example.COM",
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

// TestEmailIDValidateRejectsInvalidInputs locks in the failures.
// Quoted-string local parts are valid RFC 5322 but rejected here
// per the library's v0.1 policy: they are almost always a sign of
// malformed input.
func TestEmailIDValidateRejectsInvalidInputs(t *testing.T) {
	cases := []struct {
		input string
		want  error
	}{
		// RFC 9493 §3 — required members must be non-empty.
		{"", subjectid.ErrRequired{}},

		// Quoted-string local-part — RFC 5322 §3.4.1 admits, but
		// the library rejects in v0.1 (rare in modern use; almost
		// always a sign of malformed input).
		{`"quoted local"@example.com`, subjectid.ErrFormatEmail},
		{`"a b"@example.com`, subjectid.ErrFormatEmail},

		// Not an addr-spec at all.
		{"not-an-email", subjectid.ErrFormatEmail},
		{"plainstring", subjectid.ErrFormatEmail},

		// Name-addr per RFC 5322 §3.4 (the wire shape carries a
		// bare addr-spec only).
		{`User <user@example.com>`, subjectid.ErrFormatEmail},
		{`"User" <user@example.com>`, subjectid.ErrFormatEmail},

		// Missing local-part or domain.
		{"user@", subjectid.ErrFormatEmail},
		{"@example.com", subjectid.ErrFormatEmail},
		{"@", subjectid.ErrFormatEmail},

		// Two "@".
		{"a@b@c", subjectid.ErrFormatEmail},

		// Whitespace (RFC 5322 admits FWS only in specific
		// contexts; not in a bare addr-spec on the wire).
		{"user @example.com", subjectid.ErrFormatEmail},
		{"user@ example.com", subjectid.ErrFormatEmail},

		// Bare comment per RFC 5322 §3.2.2 — not part of the
		// addr-spec wire value.
		{"user(comment)@example.com", subjectid.ErrFormatEmail},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			err := subjectid.EmailID{Email: tc.input}.Validate()
			if err == nil {
				t.Fatalf("Validate(%q): got nil, want non-nil error", tc.input)
			}
			if !errors.Is(err, tc.want) {
				t.Errorf("errors.Is(err, %v) = false, want true (err = %v)", tc.want, err)
			}
		})
	}
}

// TestEmailIDValidateDoesNotCanonicalize pins the contract that
// Validate is a pure read — it does not mutate the receiver,
// downcase the local-part, strip plus-addressing, or otherwise
// transform the value. RFC 9493 §3.2.2 is explicit that
// canonicalization is the recipient's choice.
func TestEmailIDValidateDoesNotCanonicalize(t *testing.T) {
	const in = "User.Name+Promotions@Example.COM"
	e := subjectid.EmailID{Email: in}
	if err := e.Validate(); err != nil {
		t.Fatalf("Validate(%q) = %v, want nil", in, err)
	}
	if e.Email != in {
		t.Errorf("Email mutated by Validate: got %q, want %q", e.Email, in)
	}
}
