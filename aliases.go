// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import (
	"errors"
)

// AliasesID is a composite identifier representing the same
// subject under multiple formats, as defined in RFC 9493 §3.2.8.
//
// Wire shape:
//
//	{
//	  "format": "aliases",
//	  "identifiers": [
//	    {"format": "email",   "email": "user@example.com"},
//	    {"format": "account", "uri":   "acct:user@example.com"}
//	  ]
//	}
//
// Per RFC 9493 §3.2.8, an aliases identifier MUST NOT itself
// contain a nested aliases identifier. The Go type cannot express
// "[]SubjectIdentifier minus AliasesID" without losing the
// uniformity that makes the dispatch-and-recurse codec work, so
// the prohibition is enforced in [AliasesID.Validate] (a later
// commit) rather than at the type level.
type AliasesID struct {
	// Identifiers is the heterogeneous slice of inner Subject
	// Identifiers. It is the value of the JSON "identifiers"
	// member. The slice MUST be non-empty and MUST NOT contain
	// any element whose Format returns "aliases"; both rules
	// are enforced in Validate.
	Identifiers []SubjectIdentifier
}

// Format returns "aliases". See [SubjectIdentifier.Format].
func (AliasesID) Format() string { return "aliases" }

// Validate enforces the two RFC 9493 §3.2.8 rules on aliases:
// the identifiers array must be non-empty, and no element may
// itself be an aliases identifier (no nesting). Inner errors from
// per-element Validate calls are joined; callers branch via
// [errors.Is] against the package sentinels.
func (a AliasesID) Validate() error {
	if len(a.Identifiers) == 0 {
		return MissingFields("identifiers")
	}
	errs := make([]error, 0, len(a.Identifiers))
	for _, id := range a.Identifiers {
		if _, ok := id.(AliasesID); ok {
			errs = append(errs, ErrNestedAliases)
			continue
		}
		if err := id.Validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (AliasesID) sealed() {}

// Compile-time assertion that AliasesID satisfies SubjectIdentifier.
var _ SubjectIdentifier = AliasesID{}
