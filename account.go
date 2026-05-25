// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import (
	_ "embed"
	"sync/atomic"
)

//go:embed grammar/rfc7565/account-uri.rex
var accountURIRegexString string

var accountURIRegex = mustCompileAnchored(accountURIRegexString)

// accountValidator holds the optional custom validator installed
// via [WithAccountValidator]. nil means use the built-in RFC 7565
// regex.
var accountValidator atomic.Pointer[func(AccountID) error]

// WithAccountValidator installs a custom syntax check for the
// account format and returns the previously-installed validator
// (or nil).
//
// When set, the custom validator REPLACES the built-in RFC 7565
// acct-URI regex check; the RFC 9493 §3 required-non-empty rule
// still runs first, so the custom validator may assume a.URI is
// non-empty. Pass nil to restore the built-in. See
// [WithEmailValidator] for the rationale, gotchas, and
// concurrency notes — the contract is uniform across all six
// per-format validator hooks.
func WithAccountValidator(fn func(AccountID) error) func(AccountID) error {
	var prev func(AccountID) error
	if p := accountValidator.Load(); p != nil {
		prev = *p
	}
	if fn == nil {
		accountValidator.Store(nil)
	} else {
		accountValidator.Store(&fn)
	}
	return prev
}

// AccountID identifies a subject by an "acct" URI per RFC 7565, as
// defined in RFC 9493 §3.2.1.
//
// Wire shape:
//
//	{
//	  "format": "account",
//	  "uri":    "acct:example.user@service.example.com"
//	}
type AccountID struct {
	// URI is the acct URI naming the account. It is the value of
	// the JSON "uri" member; the scheme is required to be "acct"
	// and the scheme-specific part follows the RFC 7565 §3
	// "userpart @ host" grammar, with host per RFC 3986 §3.2.2.
	URI string
}

// Format returns "account". See [SubjectIdentifier.Format].
func (AccountID) Format() string { return "account" }

// Validate enforces the RFC 9493 §3.2.1 wire shape: the "uri"
// member must be non-empty and must match the RFC 7565 + RFC 3986
// grammar from account-uri.abnf in its entirety.
func (a AccountID) Validate() error {
	if a.URI == "" {
		return MissingFields("uri")
	}
	if v := accountValidator.Load(); v != nil {
		return (*v)(a)
	}
	if !accountURIRegex.MatchString(a.URI) {
		return ErrFormatAccount
	}
	return nil
}

func (AccountID) sealed() {}

// Compile-time assertion that AccountID satisfies SubjectIdentifier.
var _ SubjectIdentifier = AccountID{}
