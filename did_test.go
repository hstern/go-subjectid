// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"errors"
	"testing"

	"github.com/hstern/go-subjectid"
)

func TestDIDIDFormatIsDID(t *testing.T) {
	var d subjectid.DIDID
	if got, want := d.Format(), "did"; got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
}

// TestDIDIDValidateAcceptsWellFormedDIDURLs locks in the DID URL
// shell from W3C DID Core (https://www.w3.org/TR/did-core/) §3.1:
//
//	did-url            = did path-abempty [ "?" query ] [ "#" fragment ]
//	did                = "did:" method-name ":" method-specific-id
//	method-name        = 1*method-char
//	method-char        = %x61-7A / DIGIT     ; lowercase a-z / 0-9
//	method-specific-id = *( *idchar ":" ) 1*idchar
//	idchar             = ALPHA / DIGIT / "." / "-" / "_" / pct-encoded
//
// RFC 9493 §3.2.6 names DID Core as the grammar reference and
// scopes the library to the shell only (no per-method rules).
func TestDIDIDValidateAcceptsWellFormedDIDURLs(t *testing.T) {
	cases := []string{
		// RFC 9493 §3.2.6 illustrative example.
		"did:example:123456",

		// DID Core appendix examples (representative methods).
		"did:web:example.com",
		"did:web:example.com%3A8443",
		"did:key:z6MkpTHR8VNsBxYAAWHut2Geadd9jSrUEAGmGGwhB7nXgRdT",
		"did:peer:0z6MkpTHR8VNsBxYAAWHut2Geadd9jSrUEAGmGGwhB7nXgRdT",

		// Method-specific-id with multiple colon segments.
		"did:example:123:456:789",
		"did:web:example.com:user:alice",

		// idchar set per DID Core: ALPHA / DIGIT / "." / "-" / "_"
		// / pct-encoded.
		"did:example:abc.def-ghi_jkl",
		"did:example:abc%20def",

		// Optional path-abempty per DID Core §3.1.
		"did:example:123/path",
		"did:example:123/path/to/resource",

		// Optional query per DID Core §3.2.1 (service parameter).
		"did:example:123?service=files",
		"did:example:123?service=files&relativeRef=/resume.pdf",

		// Optional fragment per DID Core §3.2.1 (verification
		// method or service reference).
		"did:example:123#key-1",
		"did:example:123#agent",

		// Combined path + query + fragment.
		"did:example:123/path?service=x#frag",

		// Degenerate-but-well-formed trailing separators per RFC
		// 3986: §3.3 path-abempty = *( "/" segment ) and §3.5
		// fragment = *(...). Both productions admit a separator
		// followed by zero chars.
		"did:example:abc/",
		"did:example:abc#",

		// Method-name with digits.
		"did:m1:abc",
		"did:web2:example.com",
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			if err := (subjectid.DIDID{URL: in}.Validate()); err != nil {
				t.Errorf("Validate(%q) = %v, want nil", in, err)
			}
		})
	}
}

// TestDIDIDValidateRejectsMalformedDIDURLs locks in the failures
// from the same DID Core §3.1 grammar.
func TestDIDIDValidateRejectsMalformedDIDURLs(t *testing.T) {
	cases := []struct {
		input string
		want  error
	}{
		// RFC 9493 §3 — required members must be non-empty.
		{"", subjectid.ErrRequired{}},

		// Wrong scheme.
		{"http://example.com/did/123", subjectid.ErrFormatDID},
		{"urn:example:123", subjectid.ErrFormatDID},
		{"acct:user@example.com", subjectid.ErrFormatDID},

		// Uppercase scheme — DID Core scheme is "did" lowercase.
		{"DID:example:123", subjectid.ErrFormatDID},
		{"Did:example:123", subjectid.ErrFormatDID},

		// Missing components.
		{"did:", subjectid.ErrFormatDID},         // no method or msid
		{"did:example", subjectid.ErrFormatDID},  // missing msid colon
		{"did::123", subjectid.ErrFormatDID},     // empty method-name
		{"did:example:", subjectid.ErrFormatDID}, // empty msid (must end ≥1 idchar)

		// method-name violates 1*method-char (lowercase ALPHA / DIGIT).
		{"did:Example:123", subjectid.ErrFormatDID},  // uppercase method
		{"did:exam_ple:123", subjectid.ErrFormatDID}, // underscore not in method-char
		{"did:exam-ple:123", subjectid.ErrFormatDID}, // hyphen not in method-char
		{"did:exam.ple:123", subjectid.ErrFormatDID}, // dot not in method-char
		{"did:exam ple:123", subjectid.ErrFormatDID}, // space
		{"did:exam$ple:123", subjectid.ErrFormatDID}, // sub-delim, not in method-char

		// method-specific-id has an idchar outside ALPHA / DIGIT /
		// "." / "-" / "_" / pct-encoded.
		{"did:example:abc def", subjectid.ErrFormatDID}, // space
		// Note: "did:example:abc#" with an empty fragment and
		// "did:example:abc/" with an empty trailing segment are both
		// well-formed per RFC 3986 (§3.5 fragment = *(...) and §3.3
		// path-abempty = *( "/" segment ) with segment = *pchar), so we
		// do not reject either.

		// Bad pct-encoding in idchar.
		{"did:example:abc%2", subjectid.ErrFormatDID},
		{"did:example:abc%gg", subjectid.ErrFormatDID},
		{"did:example:abc%", subjectid.ErrFormatDID},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			err := subjectid.DIDID{URL: tc.input}.Validate()
			if err == nil {
				t.Fatalf("Validate(%q): got nil, want non-nil error", tc.input)
			}
			if !errors.Is(err, tc.want) {
				t.Errorf("errors.Is(err, %v) = false, want true (err = %v)", tc.want, err)
			}
		})
	}
}
