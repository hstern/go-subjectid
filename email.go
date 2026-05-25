// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import (
	_ "embed"
	"sync/atomic"
)

//go:embed grammar/rfc5322/addr-spec.rex
var addrSpecRegexString string

var addrSpecRegex = mustCompileAnchored(addrSpecRegexString)

// emailValidator holds the optional custom validator installed via
// [WithEmailValidator]. nil means use the built-in addr-spec regex.
var emailValidator atomic.Pointer[func(EmailID) error]

// WithEmailValidator installs a custom syntax check for the email
// format and returns the previously-installed validator (or nil).
//
// When set, the custom validator REPLACES the built-in RFC 5322
// addr-spec regex check; the RFC 9493 §3 required-non-empty rule
// still runs first, so the custom validator may assume e.Email is
// non-empty. Pass nil to restore the built-in.
//
// Use this when v0.1's narrowed addr-spec grammar (which rejects
// quoted-string local-parts and the RFC 5322 obs-* productions) is
// too strict or too loose: plug in [net/mail.ParseAddress], a
// compliance-grade parser, or a permissive corporate-directory
// matcher.
//
// The hook is process-global. Tests must restore the previous
// value, typically via defer:
//
//	defer subjectid.WithEmailValidator(subjectid.WithEmailValidator(myFn))
//
// Safe for concurrent use; the underlying value is an
// [atomic.Pointer].
func WithEmailValidator(fn func(EmailID) error) func(EmailID) error {
	var prev func(EmailID) error
	if p := emailValidator.Load(); p != nil {
		prev = *p
	}
	if fn == nil {
		emailValidator.Store(nil)
	} else {
		emailValidator.Store(&fn)
	}
	return prev
}

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
	if v := emailValidator.Load(); v != nil {
		return (*v)(e)
	}
	if !addrSpecRegex.MatchString(e.Email) {
		return ErrFormatEmail
	}
	return nil
}

func (EmailID) sealed() {}

// Compile-time assertion that EmailID satisfies SubjectIdentifier.
var _ SubjectIdentifier = EmailID{}
