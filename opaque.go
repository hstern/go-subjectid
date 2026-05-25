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
	// non-empty JSON string. The non-empty rule is enforced by
	// Validate and reported as [ErrOpaqueEmpty].
	ID string
}

// Format returns "opaque". See [SubjectIdentifier.Format].
func (OpaqueID) Format() string { return "opaque" }

// Validate enforces the RFC 9493 §3.2.4 non-empty-ID rule. The
// returned error is [ErrOpaqueEmpty] rather than [ErrRequired] so
// callers can distinguish "id was present but empty" from "id
// member was missing entirely from the JSON wire shape".
func (o OpaqueID) Validate() error {
	if o.ID == "" {
		return ErrOpaqueEmpty
	}
	return nil
}

func (OpaqueID) sealed() {}

// Compile-time assertion that OpaqueID satisfies SubjectIdentifier.
var _ SubjectIdentifier = (*OpaqueID)(nil)
