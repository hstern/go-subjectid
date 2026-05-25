// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

// DIDID identifies a subject by a W3C Decentralized Identifier
// (DID) URL, as defined in RFC 9493 §3.2.6.
//
// Wire shape:
//
//	{
//	  "format": "did",
//	  "url":    "did:example:123456"
//	}
type DIDID struct {
	// URL is the DID URL. It is the value of the JSON "url"
	// member; the grammar is the W3C DID Core URL —
	// did:method:method-specific-id[?query][#fragment]. The
	// library validates the top-level shell only in a later
	// commit; per-method validation (did:web, did:key, etc.) is
	// out of scope.
	URL string
}

// Format returns "did". See [SubjectIdentifier.Format].
func (DIDID) Format() string { return "did" }

// Validate is a no-op until the DID URL shell validation rule
// from RFC 9493 §3.2.6 lands in a later commit. The method
// exists now to satisfy [SubjectIdentifier].
func (DIDID) Validate() error { return nil }

func (DIDID) sealed() {}

// Compile-time assertion that DIDID satisfies SubjectIdentifier.
var _ SubjectIdentifier = (*DIDID)(nil)
