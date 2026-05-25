// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

// IssSubID identifies a subject by a JWT-style (issuer, subject)
// pair, as defined in RFC 9493 §3.2.3. The pair is also the iss
// and sub claim pair from RFC 7519.
//
// Wire shape:
//
//	{
//	  "format": "iss_sub",
//	  "iss":    "https://issuer.example.com/",
//	  "sub":    "145234573"
//	}
type IssSubID struct {
	// Iss is the JWT issuer URI. It is the value of the JSON
	// "iss" member and a JSON-string-encoded URI.
	Iss string

	// Sub is the issuer-scoped subject identifier. It is the
	// value of the JSON "sub" member; its grammar is the
	// issuer's choice.
	Sub string
}

// Format returns "iss_sub". See [SubjectIdentifier.Format].
func (IssSubID) Format() string { return "iss_sub" }

// Validate is a no-op until the per-member rules for RFC 9493
// §3.2.3 (iss must be a URI, sub must be non-empty) land in a
// later commit. The method exists now to satisfy
// [SubjectIdentifier].
func (IssSubID) Validate() error { return nil }

func (IssSubID) sealed() {}

// Compile-time assertion that IssSubID satisfies SubjectIdentifier.
var _ SubjectIdentifier = (*IssSubID)(nil)
