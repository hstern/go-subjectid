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

// phoneE164BasicRe matches the basic E.164 international form:
// a leading "+", then 4 to 15 ASCII digits, no separators.
// Encoded as a literal regex of the simplified spec grammar
//
//	phone = "+" 4*15DIGIT
//
// Full country-prefix validity (assigned country codes,
// per-country length rules) is intentionally not enforced here;
// that lands in a follow-up PR backed by nyaruka/phonenumbers
// (Google libphonenumber port).
var phoneE164BasicRe = regexp.MustCompile(`^\+\d{4,15}$`)

// Validate checks the PhoneNumber member against the basic E.164
// shell: a leading "+" followed by 4 to 15 ASCII digits.
//
// What this deliberately does NOT check (pending the follow-up
// nyaruka/phonenumbers integration):
//
//   - Country-prefix validity (which 1- to 3-digit prefixes are
//     assigned per the ITU-T E.164 country-code numbering plan).
//   - Per-country subscriber-number length rules.
//   - Formatting / canonicalization.
//
// What this does enforce:
//
//   - Non-empty.
//   - Leading "+".
//   - ASCII digits only — no spaces, dashes, parens, or dots.
//   - Length 4 to 15 digits (the E.164 lower and upper bounds).
func (p PhoneNumberID) Validate() error {
	if p.PhoneNumber == "" {
		return &ValidationError{
			Rule:   "required",
			Format: "phone_number",
			Reason: `"phone_number" member must be a non-empty E.164 number`,
		}
	}
	if !phoneE164BasicRe.MatchString(p.PhoneNumber) {
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
