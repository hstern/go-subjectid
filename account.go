// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

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

// acctOpaqueRe encodes the userpart-and-host portion of the
// RFC 7565 §3 acct URI ABNF:
//
//	userpart = 1*(unreserved / pct-encoded / sub-delims)
//	host     = reg-name = 1*(unreserved / pct-encoded / sub-delims)
//
// The character classes are from RFC 3986 §2; unreserved
// (ALPHA / DIGIT / "-" / "." / "_" / "~") and sub-delims
// ("!" / "$" / "&" / "'" / "(" / ")" / "*" / "+" / "," / ";"
// / "=") are folded into one bracket expression. IP-literal
// hosts (RFC 3986 §3.2.2) are not accepted in v0.1.
//
// The outer "acct:" scheme and overall URI structure are
// validated by [net/url.Parse] before this regex sees the
// opaque part, so the scheme alphabet, percent-encoding
// well-formedness, etc. do not need to be re-checked here.
var acctOpaqueRe = regexp.MustCompile(
	`^(?:[A-Za-z0-9._~!$&'()*+,;=-]|%[0-9A-Fa-f]{2})+` +
		`@(?:[A-Za-z0-9._~!$&'()*+,;=-]|%[0-9A-Fa-f]{2})+$`)

// Validate checks the URI member against the RFC 7565 §3 grammar
// for an acct URI:
//
//  1. The string must parse as an RFC 3986 URI ([net/url.Parse]).
//  2. The scheme must be exactly "acct" — lowercase, per the
//     IANA scheme registration.
//  3. The scheme-specific part must be a userpart and a host
//     separated by a single "@", both non-empty, each made up
//     of unreserved / pct-encoded / sub-delims characters from
//     RFC 3986 §2 (acctOpaqueRe).
//
// Delegating step 1 to [net/url.Parse] catches percent-encoding
// well-formedness, invalid characters, and other URI-level
// problems for free; the local regex only carries the acct-
// specific shape.
//
// What this deliberately does NOT check:
//
//   - DNS resolvability of the host. RFC 7565 is explicit that
//     an acct URI does not imply a network resource exists.
//   - IP-literal hosts (RFC 3986 §3.2.2 bracketed IPv6 or
//     IPvFuture). Subject identifiers in practice use reg-name
//     hosts; v0.2 could broaden the grammar if interop demand
//     surfaces.
//   - IDN handling. The library matches A-label hosts on the
//     wire; clients that want U-label support should normalize
//     before passing the URI in.
//   - Fragment, query, path, or authority components. The
//     RFC 7565 ABNF does not allow them; the regex anchors on
//     ^ and $ so anything beyond "userpart@host" in the opaque
//     part fails.
func (a AccountID) Validate() error {
	if a.URI == "" {
		return &ValidationError{
			Rule:   "required",
			Format: "account",
			Reason: `"uri" member must be a non-empty acct URI`,
		}
	}
	// Pre-check the literal scheme prefix. net/url.Parse lowercases
	// the scheme into u.Scheme silently, so "ACCT:..." would
	// otherwise slip past a u.Scheme == "acct" check.
	if !strings.HasPrefix(a.URI, "acct:") {
		return &ValidationError{
			Rule:   "format:account",
			Format: "account",
			Reason: `scheme must be exactly "acct" — lowercase, per the IANA scheme registration (RFC 7565 §3)`,
		}
	}
	u, err := url.Parse(a.URI)
	if err != nil {
		return &ValidationError{
			Rule:   "format:account",
			Format: "account",
			Reason: fmt.Sprintf("not a valid RFC 3986 URI: %v", err),
		}
	}
	// The RFC 7565 ABNF for acct URI has no fragment, no query,
	// no path. net/url.Parse pulls those components into their
	// own fields (so they would not appear in u.Opaque); we have
	// to reject them explicitly or a value like
	// "acct:user@host#frag" would round-trip silently.
	if u.Fragment != "" || u.RawFragment != "" {
		return &ValidationError{
			Rule:   "format:account",
			Format: "account",
			Reason: `acct URI must not contain a fragment (RFC 7565 §3)`,
		}
	}
	if u.RawQuery != "" || u.ForceQuery {
		return &ValidationError{
			Rule:   "format:account",
			Format: "account",
			Reason: `acct URI must not contain a query (RFC 7565 §3)`,
		}
	}
	if !acctOpaqueRe.MatchString(u.Opaque) {
		return &ValidationError{
			Rule:   "format:account",
			Format: "account",
			Reason: `scheme-specific part must match the RFC 7565 §3 "userpart @ host" grammar`,
		}
	}
	return nil
}

func (AccountID) sealed() {}

// Compile-time assertion that AccountID satisfies SubjectIdentifier.
var _ SubjectIdentifier = (*AccountID)(nil)
