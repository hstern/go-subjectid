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

func TestPhoneNumberIDValidateAcceptsBasicE164(t *testing.T) {
	cases := []string{
		"+12065550100",     // RFC 9493 §3.2.5 example
		"+1234",            // shortest legal: exactly 4 digits
		"+123456789012345", // longest legal: exactly 15 digits
		"+447911123456",    // UK mobile-shaped
		"+33123456789",     // FR landline-shaped
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			if err := (subjectid.PhoneNumberID{PhoneNumber: in}.Validate()); err != nil {
				t.Errorf("Validate(%q) = %v, want nil", in, err)
			}
		})
	}
}

func TestPhoneNumberIDValidateRejectsMalformedNumbers(t *testing.T) {
	cases := []struct {
		input string
		rule  string
	}{
		{"", "required"},
		{"12065550100", "format:phone_number"},       // missing +
		{"+", "format:phone_number"},                 // no digits
		{"+123", "format:phone_number"},              // 3 digits (below 4)
		{"+1234567890123456", "format:phone_number"}, // 16 digits (above 15)
		{"+1 206 555 0100", "format:phone_number"},   // spaces
		{"+1-206-555-0100", "format:phone_number"},   // hyphens
		{"+1(206)5550100", "format:phone_number"},    // parens
		{"+12.06.555.0100", "format:phone_number"},   // dots
		{"++12065550100", "format:phone_number"},     // double +
		{"+abc1234", "format:phone_number"},          // non-digit
		{"+1206555X100", "format:phone_number"},      // non-digit in middle
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			err := subjectid.PhoneNumberID{PhoneNumber: tc.input}.Validate()
			if err == nil {
				t.Fatalf("Validate(%q): got nil, want *ValidationError", tc.input)
			}
			var ve *subjectid.ValidationError
			if !errors.As(err, &ve) {
				t.Fatalf("err type = %T, want *ValidationError", err)
			}
			if ve.Rule != tc.rule {
				t.Errorf("Rule = %q, want %q", ve.Rule, tc.rule)
			}
			if ve.Format != "phone_number" {
				t.Errorf("Format = %q, want %q", ve.Format, "phone_number")
			}
		})
	}
}
