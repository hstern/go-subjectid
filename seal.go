// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

// Seal is the embeddable marker that lets a type defined outside
// this package satisfy [SubjectIdentifier] when registered via
// [RegisterFormat]. Embed it as an anonymous field in your custom
// format type to inherit the unexported sealed marker method.
//
//	type OrgTenantID struct {
//	    subjectid.Seal
//	    Tenant string
//	}
//
//	func (OrgTenantID) Format() string { return "org.example.tenant" }
//	func (OrgTenantID) Validate() error { return nil }
//
//	func init() {
//	    _ = subjectid.RegisterFormat("org.example.tenant", func() subjectid.SubjectIdentifier {
//	        return &OrgTenantID{}
//	    })
//	}
//
// The unexported sealed method is promoted to the embedding
// type, so it satisfies the interface's seal without allowing
// arbitrary external implementations — any extension type must
// opt in by embedding Seal.
//
// The built-in format types do not use Seal; they implement
// sealed directly because they live inside this package.
type Seal = sealMarker

// sealMarker is the concrete struct that carries the unexported
// sealed method. [Seal] is a type alias for sealMarker; the
// alias keeps "subjectid.Seal" as the spelling consumers see in
// embedding declarations while leaving the underlying type
// renamable without an API break.
type sealMarker struct{}

// sealed satisfies the unexported marker on [SubjectIdentifier].
// Promoted into any type that embeds [Seal].
//
// it is invoked only via interface dispatch on external types
// that embed Seal (typically registered via RegisterFormat). No
// in-package caller exists by design — the built-in format types
// declare their own sealed methods because they live in this
// package and do not need the embedding escape hatch.
//
//nolint:unused // The unused linter cannot see this method's role:
func (sealMarker) sealed() {}
