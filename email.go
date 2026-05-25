// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import (
	"fmt"
	"net/mail"
	"strings"
)

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

// Validate checks the Email member against RFC 5322 §3.4.1
// addr-spec syntax via [net/mail.ParseAddress], plus three rules
// the library applies on top:
//
//  1. The value must be non-empty (RFC 9493 §3 — required
//     members are non-null and non-empty).
//  2. The value must be a bare addr-spec, not a name-addr.
//     mail.ParseAddress also accepts "Name <addr@host>", but
//     the wire shape carries the addr-spec on its own; a
//     name-addr is a syntax error here.
//  3. Quoted-string local parts are rejected. They are valid
//     per RFC 5322 but rarely intended in modern use and are
//     more often a sign of malformed input than of a real
//     deliverable mailbox. Consumers that genuinely need them
//     can re-validate after taking the EmailID through
//     [Parse]; this method does not provide a per-instance
//     override hook in v0.1.
//
// Validation does not lowercase the local-part, strip plus-
// addressing, or otherwise canonicalize: RFC 9493 §3.2.2 is
// explicit that canonicalization is the recipient's choice.
func (e EmailID) Validate() error {
	if e.Email == "" {
		return &ValidationError{
			Rule:   "required",
			Format: "email",
			Reason: `"email" member must be a non-empty addr-spec`,
		}
	}
	if strings.HasPrefix(e.Email, `"`) {
		return &ValidationError{
			Rule:   "format:email",
			Format: "email",
			Reason: "quoted-string local parts are not supported in v0.1",
		}
	}
	addr, err := mail.ParseAddress(e.Email)
	if err != nil {
		return &ValidationError{
			Rule:   "format:email",
			Format: "email",
			Reason: fmt.Sprintf("not a valid RFC 5322 addr-spec: %v", err),
		}
	}
	if addr.Address != e.Email {
		return &ValidationError{
			Rule:   "format:email",
			Format: "email",
			Reason: "value must be a bare addr-spec, not a name-addr",
		}
	}
	return nil
}

func (EmailID) sealed() {}

// Compile-time assertion that EmailID satisfies SubjectIdentifier.
var _ SubjectIdentifier = (*EmailID)(nil)
