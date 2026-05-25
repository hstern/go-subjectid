// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import (
	"regexp"
	"sync/atomic"
)

// e164Shape is the RFC 9493 §3.2.5 wire shape: a leading "+" followed
// by 4 to 15 ASCII digits, with no separators. phonenumbers.Parse is
// lenient about spaces, dashes, parens, and dots, so we re-check the
// raw input against the spec before delegating to the parser.
var e164Shape = regexp.MustCompile(`^\+\d{4,15}$`)

// phoneNumberValidator holds the optional custom validator installed
// via [WithPhoneNumberValidator]. nil means use the built-in E.164
// shell regex.
var phoneNumberValidator atomic.Pointer[func(PhoneNumberID) error]

// WithPhoneNumberValidator installs a custom syntax check for the
// phone_number format and returns the previously-installed
// validator (or nil).
//
// When set, the custom validator REPLACES the built-in E.164
// shell regex check (which only enforces "+" + 4-15 ASCII digits);
// the RFC 9493 §3 required-non-empty rule still runs first, so
// the custom validator may assume p.PhoneNumber is non-empty.
// Pass nil to restore the built-in. See [WithEmailValidator] for
// the contract details.
//
// Typical use: wire in libphonenumber-grade validation
// (github.com/nyaruka/phonenumbers and friends) for callers who
// need real per-country plan / portability checks rather than the
// shape-only default.
func WithPhoneNumberValidator(fn func(PhoneNumberID) error) func(PhoneNumberID) error {
	var prev func(PhoneNumberID) error
	if p := phoneNumberValidator.Load(); p != nil {
		prev = *p
	}
	if fn == nil {
		phoneNumberValidator.Store(nil)
	} else {
		phoneNumberValidator.Store(&fn)
	}
	return prev
}

// PhoneNumberID identifies a subject by an ITU-T E.164 phone
// number, as defined in RFC 9493 §3.2.5.
//
// The Phone Number Identifier Format identifies a subject using a
// telephone number.  Subject Identifiers in this format MUST contain a
// "phone_number" member whose value is a string containing the full
// telephone number of the subject, including an international dialing
// prefix, formatted according to E.164 [E164].  The "phone_number"
// member is REQUIRED and MUST NOT be null or empty.  The Phone Number
// Identifier Format is identified by the name "phone_number".
//
// Below is a non-normative example Subject Identifier in the Phone
// Number Identifier Format:
//
//	{
//	  "format": "phone_number",
//	  "phone_number": "+12065550100"
//	}
//
// Figure 8: Example: Subject Identifier in the Phone Number
// Identifier Format
type PhoneNumberID struct {
	// PhoneNumber is the E.164 number. It is the value of the
	// JSON "phone_number" member: a leading "+" followed by
	// 4 to 15 ASCII digits.
	PhoneNumber string
}

// Format returns "phone_number". See [SubjectIdentifier.Format].
func (PhoneNumberID) Format() string { return "phone_number" }

// Validate enforces the RFC 9493 §3.2.5 wire shape: "phone_number"
// must be a non-empty string matching the basic E.164 shell — "+"
// followed by 4 to 15 ASCII digits with no separators. Returns
// [MissingFields]("phone_number") for the empty value and
// [ErrFormatPhoneNumber] for shape violations.
func (p PhoneNumberID) Validate() error {
	if p.PhoneNumber == "" {
		return MissingFields("phone_number")
	}
	if v := phoneNumberValidator.Load(); v != nil {
		return (*v)(p)
	}
	if !e164Shape.MatchString(p.PhoneNumber) {
		return ErrFormatPhoneNumber
	}
	return nil
}

func (PhoneNumberID) sealed() {}

// Compile-time assertion that PhoneNumberID satisfies SubjectIdentifier.
var _ SubjectIdentifier = PhoneNumberID{}
