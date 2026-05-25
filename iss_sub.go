// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import _ "embed"

//go:embed grammar/rfc7519/iss.rex
var issRegexString string

//go:embed grammar/rfc7519/sub.rex
var subRegexString string

var (
	issRegex = mustCompileAnchored(issRegexString)
	subRegex = mustCompileAnchored(subRegexString)
)

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

// Validate enforces the RFC 9493 §3.2.3 wire shape: both "iss"
// and "sub" must be non-empty; "iss" must be an RFC 3986
// absolute-URI per the iss.abnf grammar; "sub" may be any
// printable issuer-chosen identifier per the sub.abnf grammar.
func (i IssSubID) Validate() error {
	var missing []string
	if i.Iss == "" {
		missing = append(missing, "iss")
	}
	if i.Sub == "" {
		missing = append(missing, "sub")
	}
	if len(missing) > 0 {
		return MissingFields(missing...)
	}
	if !issRegex.MatchString(i.Iss) || !subRegex.MatchString(i.Sub) {
		return ErrFormatIssSub
	}
	return nil
}

func (IssSubID) sealed() {}

// Compile-time assertion that IssSubID satisfies SubjectIdentifier.
var _ SubjectIdentifier = IssSubID{}
