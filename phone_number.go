// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

// PhoneNumberID identifies a subject by an ITU-T E.164 phone
// number, as defined in RFC 9493 §3.2.5.
//
// Wire shape:
//
//	{
//	  "format":       "phone_number",
//	  "phone_number": "+12065550100"
//	}
type PhoneNumberID struct {
	// PhoneNumber is the E.164 number. It is the value of the
	// JSON "phone_number" member: a leading "+" followed by
	// ASCII digits, 4 to 15 digits total. The library validates
	// only this shell in a later commit; full country-prefix
	// validation is out of scope for v0.1 (libphonenumber-class
	// data).
	PhoneNumber string
}

// Format returns "phone_number". See [SubjectIdentifier.Format].
func (PhoneNumberID) Format() string { return "phone_number" }

// Validate is a no-op until the basic E.164 shell rule from
// RFC 9493 §3.2.5 lands in a later commit. The method exists now
// to satisfy [SubjectIdentifier].
func (PhoneNumberID) Validate() error { return nil }

func (PhoneNumberID) sealed() {}

// Compile-time assertion that PhoneNumberID satisfies SubjectIdentifier.
var _ SubjectIdentifier = (*PhoneNumberID)(nil)
