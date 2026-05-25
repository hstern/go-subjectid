// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import _ "embed"

//go:embed grammar/rfc3986/absolute-uri.rex
var absoluteURIRegexString string

var absoluteURIRegex = mustCompileAnchored(absoluteURIRegexString)

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

// Validate enforces the RFC 9493 §3.2.7 wire shape: the "uri"
// member must be non-empty and must match the RFC 3986 §4.3
// absolute-URI grammar in its entirety.
func (u URIID) Validate() error {
	if u.URI == "" {
		return MissingFields("uri")
	}
	if !absoluteURIRegex.MatchString(u.URI) {
		return ErrFormatURI
	}
	return nil
}

func (URIID) sealed() {}

// Compile-time assertion that URIID satisfies SubjectIdentifier.
var _ SubjectIdentifier = URIID{}
