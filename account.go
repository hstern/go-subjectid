// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import (
	"fmt"
	"net"
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
	// "userpart @ host" grammar, with host per RFC 3986 §3.2.2.
	URI string
}

// Format returns "account". See [SubjectIdentifier.Format].
func (AccountID) Format() string { return "account" }

// acctUserpartRe matches the userpart per RFC 7565 §3:
//
//	userpart = 1*( unreserved / pct-encoded / sub-delims )
//
// — equivalent regex form of the spec's left-recursive
// production. Character classes from RFC 3986 §2:
// unreserved (ALPHA / DIGIT / "-" / "." / "_" / "~"),
// pct-encoded ("%" HEXDIG HEXDIG, case-insensitive per
// RFC 3986 §6.2.2.1), sub-delims ("!" / "$" / "&" / "'" /
// "(" / ")" / "*" / "+" / "," / ";" / "=").
var acctUserpartRe = regexp.MustCompile(
	`^(?:[A-Za-z0-9._~!$&'()*+,;=-]|%[0-9A-Fa-f]{2})+$`)

// regNameRe matches RFC 3986 §3.2.2 reg-name:
//
//	reg-name = *( unreserved / pct-encoded / sub-delims )
//
// Note the star: an empty reg-name is permitted by the spec
// (RFC 3986 §3.2.2 does not require non-empty).
var regNameRe = regexp.MustCompile(
	`^(?:[A-Za-z0-9._~!$&'()*+,;=-]|%[0-9A-Fa-f]{2})*$`)

// ipvFutureRe matches RFC 3986 §3.2.2 IPvFuture:
//
//	IPvFuture = "v" 1*HEXDIG "." 1*( unreserved / sub-delims / ":" )
//
// (No pct-encoded inside IPvFuture per the spec; only the bare
// unreserved + sub-delims character class plus ":".)
var ipvFutureRe = regexp.MustCompile(
	`^[vV][0-9A-Fa-f]+\.[A-Za-z0-9._~!$&'()*+,;=:-]+$`)

// Validate checks the URI member against the RFC 7565 §3 grammar
// for an acct URI:
//
//	acct-uri = "acct" ":" userpart "@" host
//
// with userpart per RFC 7565 §3 and host per RFC 3986 §3.2.2
// (IP-literal / IPv4address / reg-name). The implementation
// layers stdlib parsers ([net/url.Parse], [net.ParseIP]) under
// spec-conformance checks so violations the lenient stdlib
// would normalize away (case-folded scheme, fragment / query
// / path on what should be a pure opaque URI, etc.) are
// rejected explicitly.
//
// What this deliberately does NOT check:
//
//   - DNS resolvability of the host. RFC 7565 is explicit that
//     an acct URI does not imply a network resource exists.
//   - IDN normalization. The library matches A-label hosts on
//     the wire; clients that want U-label support should
//     normalize before passing the URI in.
func (a AccountID) Validate() error {
	if a.URI == "" {
		return &ValidationError{
			Rule:   "required",
			Format: "account",
			Reason: `"uri" member must be a non-empty acct URI`,
		}
	}
	// Literal-prefix scheme check. net/url.Parse lowercases the
	// scheme into u.Scheme silently, so "ACCT:..." would otherwise
	// slip past a u.Scheme == "acct" check.
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
	// The RFC 7565 ABNF for acct URI has no fragment, query, or
	// path. net/url.Parse pulls those components into their own
	// fields (so they would not appear in u.Opaque); we have to
	// reject them explicitly or "acct:user@host#frag" would
	// round-trip silently.
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
	// Split the opaque on the first "@". RFC 7565 §3 grammar
	// admits exactly one "@" between userpart and host. The
	// userpart and host character classes (unreserved /
	// pct-encoded / sub-delims) do not include "@", so any
	// additional "@" is a syntax error.
	userpart, host, ok := strings.Cut(u.Opaque, "@")
	if !ok {
		return &ValidationError{
			Rule:   "format:account",
			Format: "account",
			Reason: `acct URI scheme-specific part must contain "@" separating userpart from host`,
		}
	}
	if strings.ContainsRune(host, '@') {
		return &ValidationError{
			Rule:   "format:account",
			Format: "account",
			Reason: `acct URI must contain exactly one "@" — host contains another`,
		}
	}
	if !acctUserpartRe.MatchString(userpart) {
		return &ValidationError{
			Rule:   "format:account",
			Format: "account",
			Reason: `userpart must match RFC 7565 §3 1*( unreserved / pct-encoded / sub-delims ) and be non-empty`,
		}
	}
	if err := validateRFC3986Host(host); err != nil {
		return &ValidationError{
			Rule:   "format:account",
			Format: "account",
			Reason: fmt.Sprintf("host: %v", err),
		}
	}
	return nil
}

// validateRFC3986Host returns nil iff s matches the RFC 3986
// §3.2.2 host production:
//
//	host = IP-literal / IPv4address / reg-name
//
// reg-name accepts the empty string per the spec ABNF; the
// caller is responsible for any "non-empty" policy on top.
//
// Returned error is a bare error with a single-line message
// describing which alternative failed; callers wrap it into
// their format-specific ValidationError.
func validateRFC3986Host(s string) error {
	// IP-literal = "[" ( IPv6address / IPvFuture ) "]"
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		inner := s[1 : len(s)-1]
		if strings.HasPrefix(inner, "v") || strings.HasPrefix(inner, "V") {
			if !ipvFutureRe.MatchString(inner) {
				return fmt.Errorf(
					`bracketed IPvFuture %q does not match RFC 3986 §3.2.2 grammar`, s)
			}
			return nil
		}
		ip := net.ParseIP(inner)
		if ip == nil || ip.To4() != nil {
			// net.ParseIP accepts IPv4-mapped-into-IPv6
			// representations and returns a 4-byte slice via
			// To4 for them; we want canonical IPv6 only inside
			// brackets.
			return fmt.Errorf(
				`bracketed host %q is not a valid IPv6 address`, s)
		}
		return nil
	}
	// IPv4address — try net.ParseIP, accept only the 4-byte form.
	if ip := net.ParseIP(s); ip != nil && ip.To4() != nil {
		return nil
	}
	// reg-name (may be empty per RFC 3986 §3.2.2).
	if !regNameRe.MatchString(s) {
		return fmt.Errorf(
			`%q is not a valid IPv4address, IP-literal, or reg-name`, s)
	}
	return nil
}

func (AccountID) sealed() {}

// Compile-time assertion that AccountID satisfies SubjectIdentifier.
var _ SubjectIdentifier = (*AccountID)(nil)
