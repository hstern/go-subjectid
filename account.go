// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import "strings"

// AccountID identifies a subject by an "acct" URI per RFC 7565, as
// defined in RFC 9493 §3.2.1.
//
// Wire shape:
//
//	{
//	  "format": "account",
//	  "uri":    "acct:example.user@service.example.com"
//	}
type AccountID struct {
	// URI is the acct URI naming the account. It is the value of
	// the JSON "uri" member; the scheme is required to be "acct"
	// and the userpart / host follow the RFC 7565 §3 grammar.
	URI string
}

// Format returns "account". See [SubjectIdentifier.Format].
func (AccountID) Format() string { return "account" }

// Validate checks the URI member against RFC 7565 §3 syntax: the
// scheme must be "acct" (lowercase, per the registration), and
// the scheme-specific part must be a single userpart "@" host
// pair with both sides non-empty.
//
// No DNS lookup is performed and no per-component grammar beyond
// "non-empty separated by a single @" is enforced. RFC 7565
// allows percent-encoded octets and a richer userpart grammar
// that the library does not parse here; a v0.2 release could
// tighten this with a dedicated parser if interop demand
// appears.
func (a AccountID) Validate() error {
	if a.URI == "" {
		return &ValidationError{
			Rule:   "required",
			Format: "account",
			Reason: `"uri" member must be a non-empty acct URI`,
		}
	}
	const scheme = "acct:"
	if !strings.HasPrefix(a.URI, scheme) {
		return &ValidationError{
			Rule:   "format:account",
			Format: "account",
			Reason: `acct URI must begin with "acct:" (RFC 7565 §3)`,
		}
	}
	body := a.URI[len(scheme):]
	at := strings.IndexByte(body, '@')
	if at < 0 {
		return &ValidationError{
			Rule:   "format:account",
			Format: "account",
			Reason: `acct URI must be "acct:userpart@host" — missing "@"`,
		}
	}
	user, host := body[:at], body[at+1:]
	if user == "" {
		return &ValidationError{
			Rule:   "format:account",
			Format: "account",
			Reason: `acct URI userpart is empty`,
		}
	}
	if host == "" {
		return &ValidationError{
			Rule:   "format:account",
			Format: "account",
			Reason: `acct URI host is empty`,
		}
	}
	if strings.IndexByte(host, '@') >= 0 {
		return &ValidationError{
			Rule:   "format:account",
			Format: "account",
			Reason: `acct URI must contain exactly one "@" — host contains a second one`,
		}
	}
	return nil
}

func (AccountID) sealed() {}

// Compile-time assertion that AccountID satisfies SubjectIdentifier.
var _ SubjectIdentifier = (*AccountID)(nil)
