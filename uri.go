// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import (
	_ "embed"
	"sync/atomic"
)

//go:embed grammar/rfc3986/absolute-uri.rex
var absoluteURIRegexString string

var absoluteURIRegex = mustCompileAnchored(absoluteURIRegexString)

// uriValidator holds the optional custom validator installed via
// [WithURIValidator]. nil means use the built-in RFC 3986
// absolute-URI regex.
var uriValidator atomic.Pointer[func(URIID) error]

// WithURIValidator installs a custom syntax check for the uri
// format and returns the previously-installed validator (or nil).
//
// When set, the custom validator REPLACES the built-in RFC 3986
// §4.3 absolute-URI regex check; the RFC 9493 §3 required-non-
// empty rule still runs first, so the custom validator may assume
// u.URI is non-empty. Pass nil to restore the built-in. See
// [WithEmailValidator] for the contract details.
//
// Typical use: tighten to a closed set of allowed schemes
// (https-only for a SaaS deployment), check against IANA's URI
// Scheme registry, or relax the absolute-URI requirement to allow
// fragment-bearing forms when the consumer's protocol uses them.
func WithURIValidator(fn func(URIID) error) func(URIID) error {
	var prev func(URIID) error
	if p := uriValidator.Load(); p != nil {
		prev = *p
	}
	if fn == nil {
		uriValidator.Store(nil)
	} else {
		uriValidator.Store(&fn)
	}
	return prev
}

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
	if v := uriValidator.Load(); v != nil {
		return (*v)(u)
	}
	if !absoluteURIRegex.MatchString(u.URI) {
		return ErrFormatURI
	}
	return nil
}

func (URIID) sealed() {}

// Compile-time assertion that URIID satisfies SubjectIdentifier.
var _ SubjectIdentifier = URIID{}
