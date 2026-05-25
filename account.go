// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import _ "embed"

//go:embed grammar/rfc7565/account-uri.rex
var accountURIRegexString string

var accountURIRegex = mustCompileAnchored(accountURIRegexString)

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
	if !accountURIRegex.MatchString(a.URI) {
		return ErrFormatAccount
	}
	return nil
}

func (AccountID) sealed() {}

// Compile-time assertion that AccountID satisfies SubjectIdentifier.
var _ SubjectIdentifier = AccountID{}
