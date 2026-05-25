// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

// EmailID identifies a subject by an email address, as defined in
// RFC 9493 §3.2.2.
//
// Wire shape:
//
//	{
//	  "format": "email",
//	  "email":  "user@example.com"
//	}
type EmailID struct {
	// Email is an RFC 5322 §3.4.1 addr-spec. It is the value of
	// the JSON "email" member.
	//
	// Per RFC 9493 §3.2.2, the recipient SHOULD apply its own
	// canonicalization (lowercasing the local-part, stripping
	// plus-addressing, etc.) — the library performs no
	// canonicalization on its own. Syntax validation lands in a
	// later commit.
	Email string
}

// Format returns "email". See [SubjectIdentifier.Format].
func (EmailID) Format() string { return "email" }

// Validate is a no-op until the RFC 5322 addr-spec validation
// rule for RFC 9493 §3.2.2 lands in a later commit. The method
// exists now to satisfy [SubjectIdentifier].
func (EmailID) Validate() error { return nil }

func (EmailID) sealed() {}

// Compile-time assertion that EmailID satisfies SubjectIdentifier.
var _ SubjectIdentifier = (*EmailID)(nil)
