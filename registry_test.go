// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"errors"
	"testing"

	"github.com/hstern/go-subjectid"
)

// fakeTenantID is an external-package SubjectIdentifier used only
// in tests. It exercises the Seal embedding path that consumers
// will use when registering extension formats.
type fakeTenantID struct {
	subjectid.Seal
	Tenant string
}

func (fakeTenantID) Format() string  { return "subjectid.test.tenant" }
func (fakeTenantID) Validate() error { return nil }

// Compile-time assertion that an external type embedding Seal
// satisfies SubjectIdentifier.
var _ subjectid.SubjectIdentifier = (*fakeTenantID)(nil)

func TestRegisterFormatBuiltinCollisionReturnsErrFormatReserved(t *testing.T) {
	for _, name := range []string{
		"account", "email", "iss_sub", "opaque",
		"phone_number", "did", "uri", "aliases",
	} {
		err := subjectid.RegisterFormat(name, func() subjectid.SubjectIdentifier { return nil })
		if err == nil {
			t.Errorf("RegisterFormat(%q): got nil error, want one wrapping ErrFormatReserved", name)
			continue
		}
		if !errors.Is(err, subjectid.ErrFormatReserved) {
			t.Errorf("RegisterFormat(%q): err = %v, want errors.Is(err, ErrFormatReserved)", name, err)
		}
	}
}

func TestRegisterFormatExtensionSucceeds(t *testing.T) {
	const name = "subjectid.test.tenant.registered"
	called := 0
	ctor := func() subjectid.SubjectIdentifier {
		called++
		return &fakeTenantID{}
	}
	if err := subjectid.RegisterFormat(name, ctor); err != nil {
		t.Fatalf("RegisterFormat(%q): err = %v, want nil", name, err)
	}
	// Re-registering the same name silently replaces the prior
	// constructor — documented in the RegisterFormat godoc.
	if err := subjectid.RegisterFormat(name, ctor); err != nil {
		t.Fatalf("re-RegisterFormat(%q): err = %v, want nil", name, err)
	}
	if called != 0 {
		t.Errorf("ctor invoked %d times during registration, want 0 (codec hasn't called it yet)", called)
	}
}

func TestErrFormatReservedMessageMentionsTheName(t *testing.T) {
	err := subjectid.RegisterFormat("email", func() subjectid.SubjectIdentifier { return nil })
	if err == nil {
		t.Fatal("RegisterFormat(\"email\"): got nil error")
	}
	want := `subjectid: format name is reserved for a built-in: "email"`
	if got := err.Error(); got != want {
		t.Errorf("err.Error() = %q, want %q", got, want)
	}
}
