// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"errors"
	"testing"

	"github.com/hstern/go-subjectid"
)

func TestOpaqueIDFormatIsOpaque(t *testing.T) {
	var o subjectid.OpaqueID
	if got, want := o.Format(), "opaque"; got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
}

// TestOpaqueIDValidateAccepts locks in the RFC 9493 §3.2.4
// contract: "id" is a case-sensitive opaque string with no
// syntax constraints beyond being a non-empty JSON string. The
// validator must accept anything non-empty.
func TestOpaqueIDValidateAccepts(t *testing.T) {
	cases := []string{
		// RFC 9493 §3.2.4 illustrative example.
		"11112222333344445555",

		// Minimal — single character.
		"a",
		" ", // single space is a non-empty string; spec imposes no charset
		"0",

		// Arbitrary printable ASCII.
		"ANY case-sensitive string %%% !",
		"abc-DEF_123.~",

		// Whitespace-bearing values — spec imposes no constraint;
		// canonicalization is the recipient's choice.
		"  leading and trailing spaces  ",
		"tabs\tand\nnewlines",

		// Non-ASCII (JSON strings carry arbitrary Unicode).
		"日本語",
		"émigré", // émigré
		"🔑",
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			if err := (subjectid.OpaqueID{ID: in}.Validate()); err != nil {
				t.Errorf("Validate(%q) = %v, want nil", in, err)
			}
		})
	}
}

// TestOpaqueIDValidateRejectsEmpty locks in the only failure mode:
// the empty string violates RFC 9493 §3's required-non-empty rule.
func TestOpaqueIDValidateRejectsEmpty(t *testing.T) {
	err := subjectid.OpaqueID{ID: ""}.Validate()
	if err == nil {
		t.Fatalf("Validate(\"\"): got nil, want non-nil error")
	}
	if !errors.Is(err, subjectid.ErrOpaqueEmpty) {
		t.Errorf("errors.Is(err, ErrOpaqueEmpty) = false, want true (err = %v)", err)
	}
}
