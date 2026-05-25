// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"errors"
	"testing"

	"github.com/hstern/go-subjectid"
)

func TestURIIDFormatIsURI(t *testing.T) {
	var u subjectid.URIID
	if got, want := u.Format(), "uri"; got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
}

// TestURIIDValidateAcceptsAbsoluteURIs locks in the RFC 9493 §3.2.7
// contract: any RFC 3986 URI is acceptable. The library validates
// against the absolute-URI form (RFC 3986 §4.3):
//
//	absolute-URI = scheme ":" hier-part [ "?" query ]
//
// — i.e. a URI without a fragment, with an explicit scheme. RFC
// 9493 §3.2.7 names the format the "URI" fallback for cases where
// no more specific format applies; the URI on the wire is
// expected to be a real identifier, not a same-document reference.
func TestURIIDValidateAcceptsAbsoluteURIs(t *testing.T) {
	cases := []string{
		// RFC 9493 §3.2.7 illustrative example.
		"urn:oasis:names:tc:saml:2.0:nameid-format:transient",

		// Common URN forms (RFC 8141).
		"urn:uuid:6e8bc430-9c3a-11d9-9669-0800200c9a66",
		"urn:isbn:0451450523",
		"urn:ietf:rfc:9493",

		// http / https with various authority and path shapes.
		"https://example.com",
		"https://example.com/",
		"https://example.com/users/42",
		"https://example.com/path?query=value",
		"https://user@example.com:8443/path",
		"http://192.0.2.1/",
		"http://[2001:db8::1]/",

		// Other registered schemes from the IANA URI scheme registry.
		"mailto:user@example.com",
		"ftp://ftp.example.com/pub/file",
		"tel:+12065550100",
		"sip:alice@example.com",
		"ws://example.com/socket",
		"file:///etc/hosts",

		// Tag URI per RFC 4151.
		"tag:example.com,2026:foo",

		// Scheme = ALPHA *( ALPHA / DIGIT / "+" / "-" / "." )
		// per RFC 3986 §3.1 — exercise each non-ALPHA char.
		"a1+b-c.d:rest",
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			if err := (subjectid.URIID{URI: in}.Validate()); err != nil {
				t.Errorf("Validate(%q) = %v, want nil", in, err)
			}
		})
	}
}

// TestURIIDValidateRejectsMalformedURIs locks in the failures.
// These exercise (a) the empty-member rule from RFC 9493 §3, and
// (b) RFC 3986 §3.1 scheme grammar and §4.3 absolute-URI shape
// (no fragment, scheme present).
func TestURIIDValidateRejectsMalformedURIs(t *testing.T) {
	cases := []struct {
		input string
		want  error
	}{
		// RFC 9493 §3 — required members must be non-empty.
		{"", subjectid.ErrRequired{}},

		// No scheme — not even a URI-reference.
		{"not a uri", subjectid.ErrFormatURI},
		{"example.com", subjectid.ErrFormatURI},
		{"justastring", subjectid.ErrFormatURI},

		// Relative-ref forms per RFC 3986 §4.2 — not absolute URIs.
		{"/path/only", subjectid.ErrFormatURI},
		{"//host/path", subjectid.ErrFormatURI}, // network-path reference
		{"./relative", subjectid.ErrFormatURI},
		{"../relative", subjectid.ErrFormatURI},

		// Empty scheme.
		{"://example.com", subjectid.ErrFormatURI},
		{":foo", subjectid.ErrFormatURI},

		// Scheme grammar violations per RFC 3986 §3.1:
		//   scheme = ALPHA *( ALPHA / DIGIT / "+" / "-" / "." )
		// — must start with ALPHA, must contain only listed chars.
		{"1http://example.com", subjectid.ErrFormatURI}, // starts with digit
		{"+http://example.com", subjectid.ErrFormatURI}, // starts with "+"
		{"ht_tp://example.com", subjectid.ErrFormatURI}, // "_" not in scheme chars
		{"ht tp://example.com", subjectid.ErrFormatURI}, // space in scheme

		// Absolute-URI per RFC 3986 §4.3 prohibits fragment.
		// (If implementer interprets RFC 9493 §3.2.7 as the
		// fragment-permitting URI form, drop this row.)
		{"https://example.com/#frag", subjectid.ErrFormatURI},

		// Malformed IP literal.
		{"http://[bad]/", subjectid.ErrFormatURI},
		{"http://[2001:db8::1/", subjectid.ErrFormatURI}, // unclosed bracket

		// Whitespace anywhere in the URI value.
		{"https://example.com/path with space", subjectid.ErrFormatURI},
		{"https://exa mple.com/", subjectid.ErrFormatURI},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			err := subjectid.URIID{URI: tc.input}.Validate()
			if err == nil {
				t.Fatalf("Validate(%q): got nil, want non-nil error", tc.input)
			}
			if !errors.Is(err, tc.want) {
				t.Errorf("errors.Is(err, %v) = false, want true (err = %v)", tc.want, err)
			}
		})
	}
}
