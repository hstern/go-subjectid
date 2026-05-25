// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import "regexp"

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
	// 4 to 15 ASCII digits.
	PhoneNumber string
}

// Format returns "phone_number". See [SubjectIdentifier.Format].
func (PhoneNumberID) Format() string { return "phone_number" }

// phoneE164Re matches the basic E.164 shell the library checks:
// a literal "+", then 4 to 15 ASCII digits. ITU-T E.164 allows up
// to 15 digits total in the international number (country code +
// subscriber number); the 4-digit lower bound is the smallest
// realistic country-code-plus-subscriber pair.
var phoneE164Re = regexp.MustCompile(`^\+\d{4,15}$`)

// Validate checks the PhoneNumber member against the basic E.164
// shell: a leading "+" followed by 4 to 15 ASCII digits.
//
// What this deliberately does NOT check:
//
//   - Country-prefix validity. Full E.164 validation would
//     require a libphonenumber-sized country-code database; v0.1
//     stays in the standard library. Consumers that need full
//     conformance can re-validate after [Parse].
//   - Number-portability or assignment status. Beyond E.164's
//     scope.
//   - Separator or formatting characters. E.164 prohibits any
//     punctuation; the library rejects spaces, dashes, parens,
//     and dots in the digit run.
func (p PhoneNumberID) Validate() error {
	if p.PhoneNumber == "" {
		return &ValidationError{
			Rule:   "required",
			Format: "phone_number",
			Reason: `"phone_number" member must be a non-empty E.164 number`,
		}
	}
	if !phoneE164Re.MatchString(p.PhoneNumber) {
		return &ValidationError{
			Rule:   "format:phone_number",
			Format: "phone_number",
			Reason: `value must match the E.164 shell "+<4 to 15 ASCII digits>" with no separators`,
		}
	}
	return nil
}

func (PhoneNumberID) sealed() {}

// Compile-time assertion that PhoneNumberID satisfies SubjectIdentifier.
var _ SubjectIdentifier = (*PhoneNumberID)(nil)
