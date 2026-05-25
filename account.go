// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

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
	// and the user-part / host-part follow the RFC 7565 grammar.
	// Syntax validation is implemented in a later commit.
	URI string
}

// Format returns "account". See [SubjectIdentifier.Format].
func (AccountID) Format() string { return "account" }

// Validate is a no-op until the per-format validation rules for
// RFC 9493 §3.2.1 and the RFC 7565 acct URI syntax land in a
// later commit. The method exists now to satisfy
// [SubjectIdentifier] so callers can program against the
// interface from day one.
func (AccountID) Validate() error { return nil }

func (AccountID) sealed() {}

// Compile-time assertion that AccountID satisfies SubjectIdentifier.
var _ SubjectIdentifier = (*AccountID)(nil)
