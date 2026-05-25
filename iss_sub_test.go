// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"errors"
	"testing"

	"github.com/hstern/go-subjectid"
)

func TestIssSubIDFormatIsIssSub(t *testing.T) {
	var i subjectid.IssSubID
	if got, want := i.Format(), "iss_sub"; got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
}

// TestIssSubIDValidateAccepts locks in the per-member contract
// from RFC 9493 §3.2.3: "iss" is the JWT issuer URI (RFC 7519
// §4.1.1 — StringOrURI, in this context required to be a URI),
// "sub" is the issuer-scoped subject identifier (non-empty
// string; the grammar is the issuer's choice).
func TestIssSubIDValidateAccepts(t *testing.T) {
	cases := []struct {
		name string
		iss  string
		sub  string
	}{
		// RFC 9493 §3.2.3 illustrative example.
		{"spec example", "https://issuer.example.com/", "145234573"},

		// Other URI issuer shapes per RFC 7519 §4.1.1.
		{"https issuer", "https://op.example.com", "abc"},
		{"https issuer with port", "https://op.example.com:8443", "user-1"},
		{"https issuer with path", "https://op.example.com/realms/main", "u"},
		{"urn issuer", "urn:example:issuer", "1"},
		{"did issuer", "did:example:issuer", "alice"},
		{"http issuer", "http://op.example.com", "x"},

		// Sub may contain any non-empty value the issuer chose.
		{"opaque digits sub", "https://issuer.example.com/", "11112222333344445555"},
		{"opaque base64-ish sub", "https://issuer.example.com/", "abc-XYZ_123.tilde~"},
		{"sub with space", "https://issuer.example.com/", "subject with space"},
		{"single-char sub", "https://issuer.example.com/", "x"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			i := subjectid.IssSubID{Iss: tc.iss, Sub: tc.sub}
			if err := i.Validate(); err != nil {
				t.Errorf("Validate({iss:%q, sub:%q}) = %v, want nil",
					tc.iss, tc.sub, err)
			}
		})
	}
}

// TestIssSubIDValidateRejects locks in the failures. The wrapped
// sentinel is [subjectid.ErrRequired{}] for empty members (RFC 9493
// §3) and [subjectid.ErrFormatIssSub] for iss values that are not
// URIs (RFC 7519 §4.1.1 / RFC 3986 §4.3).
func TestIssSubIDValidateRejects(t *testing.T) {
	cases := []struct {
		name string
		iss  string
		sub  string
		want error
	}{
		// Empty iss or sub or both — RFC 9493 §3.
		{"empty iss", "", "x", subjectid.ErrRequired{}},
		{"empty sub", "https://issuer.example.com/", "", subjectid.ErrRequired{}},
		{"both empty", "", "", subjectid.ErrRequired{}},

		// iss not a URI — RFC 7519 §4.1.1 requires StringOrURI;
		// the library tightens to URI per the design table.
		{"iss plain word", "issuer", "x", subjectid.ErrFormatIssSub},
		{"iss relative ref", "/issuer", "x", subjectid.ErrFormatIssSub},
		{"iss network-path", "//issuer.example.com", "x", subjectid.ErrFormatIssSub},
		{"iss empty scheme", "://issuer.example.com", "x", subjectid.ErrFormatIssSub},
		{"iss scheme starts digit", "1http://issuer.example.com", "x", subjectid.ErrFormatIssSub},
		{"iss with whitespace", "https:// issuer.example.com", "x", subjectid.ErrFormatIssSub},
		{"iss bare colon", ":", "x", subjectid.ErrFormatIssSub},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := subjectid.IssSubID{Iss: tc.iss, Sub: tc.sub}.Validate()
			if err == nil {
				t.Fatalf("Validate({iss:%q, sub:%q}): got nil, want non-nil error",
					tc.iss, tc.sub)
			}
			if !errors.Is(err, tc.want) {
				t.Errorf("errors.Is(err, %v) = false, want true (err = %v)", tc.want, err)
			}
		})
	}
}
