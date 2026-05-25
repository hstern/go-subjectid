// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hstern/go-subjectid"
)

func TestSpecVersionIsRFC9493(t *testing.T) {
	if got, want := subjectid.SpecVersion, "RFC 9493"; got != want {
		t.Fatalf("SpecVersion = %q, want %q", got, want)
	}
}

func TestErrRequiredCarriesFieldsAndMatchesErrorsIs(t *testing.T) {
	// Compile-time assertion.
	var _ error = subjectid.ErrRequired{}

	e := subjectid.MissingFields("iss", "sub")

	// Categorical errors.Is against the zero value matches any
	// ErrRequired regardless of Fields.
	if !errors.Is(e, subjectid.ErrRequired{}) {
		t.Errorf("errors.Is(e, ErrRequired{}) = false, want true")
	}

	// errors.As recovers the named fields.
	var req subjectid.ErrRequired
	if !errors.As(e, &req) {
		t.Fatalf("errors.As did not extract ErrRequired from %v", e)
	}
	if got, want := req.Fields, []string{"iss", "sub"}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("Fields = %v, want %v", got, want)
	}

	// Error message names the missing fields.
	if got := e.Error(); got != "subjectid: required fields iss, sub are missing or empty" {
		t.Errorf("Error() = %q", got)
	}
}

func TestErrFormatReservedIsComparableViaErrorsIs(t *testing.T) {
	if subjectid.ErrFormatReserved == nil {
		t.Fatal("ErrFormatReserved is nil")
	}
	wrapped := fmt.Errorf("registering %q: %w", "email", subjectid.ErrFormatReserved)
	if !errors.Is(wrapped, subjectid.ErrFormatReserved) {
		t.Fatal("errors.Is did not match the wrapped sentinel")
	}
}
