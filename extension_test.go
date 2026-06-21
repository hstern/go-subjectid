// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/hstern/go-subjectid"
)

// orgTenantID is the test-only extension format used by
// TestRegisterFormatExtensionSmoke. It lives in the test binary
// only — production consumers register their own format types.
type orgTenantID struct {
	subjectid.Seal
	Tenant string
}

// Format returns the extension format's IANA-style discriminator.
// Names below the "org." prefix are conventionally private; the
// library does not police the namespace, but reserving a private
// prefix is the lowest-friction way to avoid collisions with
// future IANA additions.
func (orgTenantID) Format() string { return "org.example.tenant" }

// MarshalJSON emits the extension's spec-order wire shape.
// Mirroring the built-ins' layout — "format" first, then the
// format-specific members — keeps a consumer's wire output
// indistinguishable from the library's own.
func (o orgTenantID) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Format string `json:"format"`
		Tenant string `json:"tenant"`
	}{Format: o.Format(), Tenant: o.Tenant})
}

// UnmarshalJSON reads the "tenant" member. "format" is ignored —
// [subjectid.Parse] has already dispatched by the time this method
// runs, so re-reading it would be redundant.
func (o *orgTenantID) UnmarshalJSON(data []byte) error {
	var v struct {
		Tenant string `json:"tenant"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	o.Tenant = v.Tenant
	return nil
}

// errOrgTenantEmpty is the extension's per-format sentinel. Built
// via [subjectid.FormatErr] so it participates in the package's
// uniform errors.Is contract — consumers can branch on
// errors.Is(err, errOrgTenantEmpty), errors.Is(err, subjectid.Err),
// or both.
var errOrgTenantEmpty = subjectid.FormatErr("invalid org.example.tenant")

// Validate enforces the format's well-formedness — here, a
// non-empty Tenant.
func (o orgTenantID) Validate() error {
	if o.Tenant == "" {
		return errOrgTenantEmpty
	}
	return nil
}

// TestRegisterFormatExtensionSmoke walks the end-to-end consumer
// path for [subjectid.RegisterFormat]: register a custom format,
// parse a payload using it, validate, marshal back, observe
// byte-stable round-trip. This is the contract every external
// extension is built against.
//
// The test does not deregister the format afterwards — the
// library exposes no deregister API. The choice is intentional:
// extension registration is a process-lifecycle event (consumer
// init), and a deregister door would invite consumers to treat
// the registry as mutable per request. Re-registering an existing
// name silently replaces the constructor, which keeps test
// re-runs deterministic without a cleanup step.
func TestRegisterFormatExtensionSmoke(t *testing.T) {
	if err := subjectid.RegisterFormat("org.example.tenant", func() subjectid.SubjectIdentifier {
		return &orgTenantID{}
	}); err != nil {
		t.Fatalf("RegisterFormat: %v", err)
	}

	wire := []byte(`{"format":"org.example.tenant","tenant":"acme-prod"}`)

	id, err := subjectid.Parse(wire)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	tenant, ok := id.(orgTenantID)
	if !ok {
		t.Fatalf("Parse returned %T, want orgTenantID", id)
	}
	if got, want := tenant.Tenant, "acme-prod"; got != want {
		t.Errorf("Tenant = %q, want %q", got, want)
	}
	if verr := id.Validate(); verr != nil {
		t.Errorf("Validate: %v", verr)
	}

	got, merr := json.Marshal(id)
	if merr != nil {
		t.Fatalf("Marshal: %v", merr)
	}
	if !bytes.Equal(got, wire) {
		t.Errorf("byte-stability broken:\n got:  %s\n want: %s", got, wire)
	}

	// Validate failure path: empty Tenant should surface the
	// per-format sentinel and also match the package umbrella.
	if err := (orgTenantID{}).Validate(); err == nil {
		t.Error("Validate(empty Tenant): err = nil, want non-nil")
	} else {
		if !errors.Is(err, errOrgTenantEmpty) {
			t.Errorf("errors.Is(err, errOrgTenantEmpty) = false, want true (err = %v)", err)
		}
		if !errors.Is(err, subjectid.Err) {
			t.Errorf("errors.Is(err, subjectid.Err) = false, want true (err = %v)", err)
		}
	}
}

// TestRegisterFormatRejectsBuiltinNames locks in the
// no-overrides-of-built-ins rule from [subjectid.RegisterFormat]'s
// godoc: every IANA-registered built-in name is reserved, and
// re-registering one returns an error wrapping
// [subjectid.ErrFormatReserved] for programmatic dispatch.
func TestRegisterFormatRejectsBuiltinNames(t *testing.T) {
	builtins := []string{
		"account", "email", "iss_sub", "opaque",
		"phone_number", "did", "uri", "aliases",
	}
	for _, name := range builtins {
		t.Run(name, func(t *testing.T) {
			err := subjectid.RegisterFormat(name, func() subjectid.SubjectIdentifier {
				return &orgTenantID{}
			})
			if err == nil {
				t.Fatalf("RegisterFormat(%q): err = nil, want ErrFormatReserved", name)
			}
			if !errors.Is(err, subjectid.ErrFormatReserved) {
				t.Errorf("errors.Is(err, ErrFormatReserved) = false, want true (err = %v)", err)
			}
		})
	}
}
