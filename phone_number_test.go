// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"errors"
	"testing"

	"github.com/hstern/go-subjectid"
)

func TestPhoneNumberIDFormatIsPhoneNumber(t *testing.T) {
	var p subjectid.PhoneNumberID
	if got, want := p.Format(), "phone_number"; got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
}

// TestPhoneNumberIDValidateAcceptsBasicE164 locks in the
// ITU-T E.164 international form: a leading "+", then ASCII
// digits with no separators, with a real country prefix and a
// libphonenumber-accepted national subscriber number. The v0.1
// contract delegates country-prefix and length-per-country rules
// to libphonenumber.
func TestPhoneNumberIDValidateAcceptsBasicE164(t *testing.T) {
	cases := []string{
		// RFC 9493 §3.2.5 illustrative example.
		"+12065550100",

		// Length boundaries libphonenumber accepts: the shortest
		// valid international number is 7 digits (e.g. Austria),
		// and the E.164 ceiling is 15 digits.
		"+4312",            // shortest libphonenumber-valid
		"+431234567890123", // E.164 maximum of 15 digits

		// Country-shaped examples (not validating prefixes here).
		"+447911123456",  // UK mobile-shaped
		"+33123456789",   // FR landline-shaped
		"+819012345678",  // JP mobile-shaped
		"+5511987654321", // BR mobile-shaped
		"+861012345678",  // CN-shaped
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			if err := (subjectid.PhoneNumberID{PhoneNumber: in}.Validate()); err != nil {
				t.Errorf("Validate(%q) = %v, want nil", in, err)
			}
		})
	}
}

// TestPhoneNumberIDValidateRejectsMalformedNumbers locks in the
// failures. Separators, ASCII-digit violations, and length
// boundaries are exercised explicitly.
func TestPhoneNumberIDValidateRejectsMalformedNumbers(t *testing.T) {
	cases := []struct {
		input string
		want  error
	}{
		// RFC 9493 §3 — required members must be non-empty.
		{"", subjectid.ErrRequired{}},

		// Missing leading "+".
		{"1234", subjectid.ErrFormatPhoneNumber},
		{"12065550100", subjectid.ErrFormatPhoneNumber},
		{"206 555 0100", subjectid.ErrFormatPhoneNumber},

		// "+" alone or with too few digits — below E.164 minimum.
		{"+", subjectid.ErrFormatPhoneNumber},
		{"+1", subjectid.ErrFormatPhoneNumber},
		{"+12", subjectid.ErrFormatPhoneNumber},
		{"+123", subjectid.ErrFormatPhoneNumber},

		// Above E.164 maximum of 15 digits.
		{"+1234567890123456", subjectid.ErrFormatPhoneNumber},

		// Separators not in the basic E.164 shell.
		{"+1 206 555 0100", subjectid.ErrFormatPhoneNumber}, // spaces
		{"+1-206-555-0100", subjectid.ErrFormatPhoneNumber}, // hyphens
		{"+1(206)5550100", subjectid.ErrFormatPhoneNumber},  // parens
		{"+12.06.555.0100", subjectid.ErrFormatPhoneNumber}, // dots
		{"+1\t2065550100", subjectid.ErrFormatPhoneNumber},  // tab

		// Doubled or interior "+".
		{"++12065550100", subjectid.ErrFormatPhoneNumber},
		{"+1206+5550100", subjectid.ErrFormatPhoneNumber},

		// Non-digit characters anywhere.
		{"+abc1234", subjectid.ErrFormatPhoneNumber},
		{"+1206555X100", subjectid.ErrFormatPhoneNumber},
		{"+1206555100x", subjectid.ErrFormatPhoneNumber},

		// Leading whitespace.
		{" +12065550100", subjectid.ErrFormatPhoneNumber},

		// Non-ASCII digits (Arabic-Indic 4) — RFC 9493 §3.2.5
		// is ASCII per E.164.
		{"+۱۲۳۴", subjectid.ErrFormatPhoneNumber},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			err := subjectid.PhoneNumberID{PhoneNumber: tc.input}.Validate()
			if err == nil {
				t.Fatalf("Validate(%q): got nil, want non-nil error", tc.input)
			}
			if !errors.Is(err, tc.want) {
				t.Errorf("errors.Is(err, %v) = false, want true (err = %v)", tc.want, err)
			}
		})
	}
}
