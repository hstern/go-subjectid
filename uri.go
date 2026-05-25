// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

// URIID identifies a subject by any RFC 3986 URI, as defined in
// RFC 9493 §3.2.7. It is the "use when nothing more specific
// applies" fallback within the built-in set.
//
// Wire shape:
//
//	{
//	  "format": "uri",
//	  "uri":    "urn:oasis:names:tc:saml:2.0:nameid-format:transient"
//	}
type URIID struct {
	// URI is the absolute RFC 3986 URI naming the subject. It is
	// the value of the JSON "uri" member; the library validates
	// via net/url.Parse plus a strictness wrapper in a later
	// commit (absolute URI, non-empty scheme).
	URI string
}

// Format returns "uri". See [SubjectIdentifier.Format].
func (URIID) Format() string { return "uri" }

// Validate is a no-op until the RFC 3986 absolute-URI rule from
// RFC 9493 §3.2.7 lands in a later commit. The method exists now
// to satisfy [SubjectIdentifier].
func (URIID) Validate() error { return nil }

func (URIID) sealed() {}

// Compile-time assertion that URIID satisfies SubjectIdentifier.
var _ SubjectIdentifier = (*URIID)(nil)
