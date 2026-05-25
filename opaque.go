// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

// OpaqueID identifies a subject by a case-sensitive opaque string,
// as defined in RFC 9493 §3.2.4.
//
// Wire shape:
//
//	{
//	  "format": "opaque",
//	  "id":     "11112222333344445555"
//	}
type OpaqueID struct {
	// ID is the opaque identifier. It is the value of the JSON
	// "id" member and has no syntax constraints beyond being a
	// non-empty JSON string. Validation that the value is
	// non-empty lands in a later commit.
	ID string
}

// Format returns "opaque". See [SubjectIdentifier.Format].
func (OpaqueID) Format() string { return "opaque" }

// Validate is a no-op until the non-empty-ID rule from RFC 9493
// §3.2.4 lands in a later commit. The method exists now to
// satisfy [SubjectIdentifier].
func (OpaqueID) Validate() error { return nil }

func (OpaqueID) sealed() {}

// Compile-time assertion that OpaqueID satisfies SubjectIdentifier.
var _ SubjectIdentifier = (*OpaqueID)(nil)
