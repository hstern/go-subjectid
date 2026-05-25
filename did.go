// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import (
	_ "embed"
	"sync/atomic"
)

//go:embed grammar/w3c-did-core/did-url.rex
var didURLRegexString string

var didURLRegex = mustCompileAnchored(didURLRegexString)

// didValidator holds the optional custom validator installed via
// [WithDIDValidator]. nil means use the built-in W3C DID Core
// shell regex.
var didValidator atomic.Pointer[func(DIDID) error]

// WithDIDValidator installs a custom syntax check for the did
// format and returns the previously-installed validator (or nil).
//
// When set, the custom validator REPLACES the built-in W3C DID
// Core URL shell regex (which does not enforce per-method rules);
// the RFC 9493 §3 required-non-empty rule still runs first, so
// the custom validator may assume d.URL is non-empty. Pass nil to
// restore the built-in. See [WithEmailValidator] for the contract
// details.
//
// Typical use: layer per-method validation (did:web, did:key,
// did:plc, did:peer, …) on top of the shell. The library only
// validates the universal shape; method-specific rigor is the
// consumer's choice.
func WithDIDValidator(fn func(DIDID) error) func(DIDID) error {
	var prev func(DIDID) error
	if p := didValidator.Load(); p != nil {
		prev = *p
	}
	if fn == nil {
		didValidator.Store(nil)
	} else {
		didValidator.Store(&fn)
	}
	return prev
}

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

// Validate enforces the RFC 9493 §3.2.6 wire shape: the "url"
// member must be non-empty and must match the W3C DID Core ABNF
// from did-url.abnf in its entirety. Per-method validation (did:web,
// did:key, etc.) is out of scope.
func (d DIDID) Validate() error {
	if d.URL == "" {
		return MissingFields("url")
	}
	if v := didValidator.Load(); v != nil {
		return (*v)(d)
	}
	if !didURLRegex.MatchString(d.URL) {
		return ErrFormatDID
	}
	return nil
}

func (DIDID) sealed() {}

// Compile-time assertion that DIDID satisfies SubjectIdentifier.
var _ SubjectIdentifier = DIDID{}
