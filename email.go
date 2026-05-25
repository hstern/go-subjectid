// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import _ "embed"

//go:embed grammar/rfc5322/addr-spec.rex
var addrSpecRegexString string

var addrSpecRegex = mustCompileAnchored(addrSpecRegexString)

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
	// canonicalization on its own.
	Email string
}

// Format returns "email". See [SubjectIdentifier.Format].
func (EmailID) Format() string { return "email" }

// Validate enforces the RFC 9493 §3.2.2 wire shape: the "email"
// member must be non-empty and must match the RFC 5322 §3.4.1
// addr-spec grammar (narrowed per v0.1 policy — see addr-spec.abnf)
// in its entirety.
func (e EmailID) Validate() error {
	if e.Email == "" {
		return MissingFields("email")
	}
	if !addrSpecRegex.MatchString(e.Email) {
		return ErrFormatEmail
	}
	return nil
}

func (EmailID) sealed() {}

// Compile-time assertion that EmailID satisfies SubjectIdentifier.
var _ SubjectIdentifier = EmailID{}
