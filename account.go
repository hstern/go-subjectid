// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

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
	// "userpart @ host" grammar.
	URI string
}

// Format returns "account". See [SubjectIdentifier.Format].
func (AccountID) Format() string { return "account" }

// Validate checks the URI member against the RFC 7565 §3 grammar
// for an acct URI: scheme "acct" (lowercase, per the IANA
// registration), then a userpart and a host separated by a single
// "@", each non-empty and built from RFC 3986 §2 unreserved /
// pct-encoded / sub-delims characters. No fragment, query, or
// path is allowed by the ABNF, and this is enforced by the regex
// anchoring on ^ and $.
//
// The regex itself ([acctURIRe]) is generated from
// tools/genabnf/abnf/acct.abnf by `make gen-grammars`; the
// generator path keeps this validator in lock-step with the
// spec ABNF rather than a hand-paraphrase.
//
// What this deliberately does NOT check:
//
//   - DNS resolvability of the host. RFC 7565 is explicit that
//     an acct URI does not imply a network resource exists.
//   - IP-literal hosts (RFC 3986 §3.2.2 bracketed IPv6 or
//     IPvFuture). The ABNF in abnf/acct.abnf accepts only
//     reg-name-shaped hosts; broaden the grammar if interop
//     demand surfaces.
//   - IDN handling. The library matches A-label hosts on the
//     wire; clients that want U-label support should normalize
//     before passing the URI in.
func (a AccountID) Validate() error {
	if a.URI == "" {
		return &ValidationError{
			Rule:   "required",
			Format: "account",
			Reason: `"uri" member must be a non-empty acct URI`,
		}
	}
	if !acctURIRe.MatchString(a.URI) {
		return &ValidationError{
			Rule:   "format:account",
			Format: "account",
			Reason: `value does not match the RFC 7565 §3 acct URI ABNF`,
		}
	}
	return nil
}

func (AccountID) sealed() {}

// Compile-time assertion that AccountID satisfies SubjectIdentifier.
var _ SubjectIdentifier = (*AccountID)(nil)
