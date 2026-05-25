// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import "errors"

// ValidationError reports a single RFC 9493 well-formedness rule
// that an identifier failed.
//
// Callers may branch on the structured Rule and Format fields
// without parsing the human-readable [ValidationError.Error]
// message — those fields are part of the package's compatibility
// surface.
type ValidationError struct {
	// Rule names the spec rule that failed. Values are short,
	// machine-readable identifiers — for example "required",
	// "format:email", or "aliases:no-nested-aliases".
	Rule string

	// Format is the identifier's format discriminator (the same
	// string [SubjectIdentifier.Format] returns). For an element
	// inside an aliases identifier, the value is suffixed with
	// the zero-based element index — "aliases:3" names the
	// fourth element — so a caller can locate the offending
	// element without re-walking the structure.
	Format string

	// Reason is a human-readable explanation of the violation,
	// suitable for inclusion in log lines and user-facing error
	// messages.
	Reason string
}

// Error formats the receiver as
//
//	subjectid: <rule> on <format>: <reason>
//
// The string is stable enough for log inspection but is not part
// of the structured comparison surface. Callers that need to
// branch on a specific rule should compare [ValidationError.Rule]
// directly rather than parsing the message.
func (e *ValidationError) Error() string {
	return "subjectid: " + e.Rule + " on " + e.Format + ": " + e.Reason
}

// ErrFormatReserved is returned by RegisterFormat (defined in a
// later commit) when a caller attempts to register a constructor
// for a format name that the library has already populated as a
// built-in.
//
// The eight built-in formats are populated at package init from
// the IANA "Security Event Identifier Formats" registry as it
// stood at the time of release. Built-ins cannot be overridden;
// consumers that need different behavior for a built-in format
// should wrap the concrete type rather than re-register the
// format.
var ErrFormatReserved = errors.New("subjectid: format name is reserved for a built-in")
