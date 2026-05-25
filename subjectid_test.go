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

func TestValidationErrorImplementsErrorAndFormats(t *testing.T) {
	// Compile-time assertion.
	var _ error = (*subjectid.ValidationError)(nil)

	e := &subjectid.ValidationError{
		Rule:   "required",
		Format: "email",
		Reason: "email member must be present",
	}
	got := e.Error()
	want := "subjectid: required on email: email member must be present"
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
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
