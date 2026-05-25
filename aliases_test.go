// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"testing"

	"github.com/hstern/go-subjectid"
)

func TestAliasesIDFormat(t *testing.T) {
	var a subjectid.AliasesID
	if got, want := a.Format(), "aliases"; got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
	if err := a.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil (no rules yet)", err)
	}
}

// TestAliasesIDIdentifiersAcceptsSubjectIdentifier documents the
// type-level contract that lets the codec carry a heterogeneous
// mix of inner identifiers. The element type is the sealed
// interface, so once all eight built-ins land any combination
// fits. Here we exercise the contract with the only
// SubjectIdentifier-satisfying type defined in this PR (an inner
// AliasesID); broader cross-format coverage lands with the codec
// in phase 3.
//
// A nested AliasesID is spec-illegal — RFC 9493 §3.2.8 forbids
// it — but is type-legal, and Validate (phase 4) will reject it
// at the well-formedness layer rather than at the type layer.
func TestAliasesIDIdentifiersAcceptsSubjectIdentifier(t *testing.T) {
	inner := subjectid.AliasesID{}
	outer := subjectid.AliasesID{
		Identifiers: []subjectid.SubjectIdentifier{inner},
	}
	if got, want := len(outer.Identifiers), 1; got != want {
		t.Fatalf("len(Identifiers) = %d, want %d", got, want)
	}
	if got, want := outer.Identifiers[0].Format(), "aliases"; got != want {
		t.Errorf("Identifiers[0].Format() = %q, want %q", got, want)
	}
}
