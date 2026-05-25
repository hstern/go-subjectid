// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

// Package subjectid implements RFC 9493 Subject Identifiers for
// Security Event Tokens.
//
// The package provides:
//
//   - The sealed [SubjectIdentifier] interface implemented by every
//     identifier value, defined in this file.
//   - Per-format concrete types for each format in the IANA
//     "Security Event Identifier Formats" registry, each in its
//     own file (account.go, email.go, …).
//   - An UnknownFormat carrier that preserves the wire bytes of
//     formats the library does not natively recognize, so forward-
//     compatibility round-trips correctly.
//   - An opt-in [SubjectIdentifier.Validate] method per format,
//     returning the package's sentinel errors (per-format
//     [FormatErr]-built sentinels, [ErrRequired] for missing
//     members, etc.). Callers branch with [errors.Is]; every
//     sentinel also matches the umbrella [Err].
package subjectid

// SpecVersion identifies the RFC this package implements. RFCs have
// no minor or patch numbers; errata to RFC 9493 are absorbed into
// Go-minor releases of this module without changing the value of
// this constant.
const SpecVersion = "RFC 9493"

// SubjectIdentifier is the sealed interface implemented by every
// RFC 9493 Subject Identifier value the package exposes.
//
// Implementations are confined to this package by the unexported
// sealed marker method. Downstream code that needs a format
// outside the built-in set calls RegisterFormat (defined in a
// later commit) rather than implementing the interface directly:
// registration feeds the codec dispatch table, whereas a direct
// implementation would be invisible to it.
type SubjectIdentifier interface {
	// Format returns the IANA "format" discriminator for this
	// identifier — "email", "iss_sub", "aliases", and so on. The
	// returned value is also the value of the "format" member in
	// the JSON wire shape.
	Format() string

	// Validate reports whether the identifier satisfies the
	// RFC 9493 well-formedness rules for its format. A nil return
	// means valid; a non-nil return is one of the package's
	// sentinel errors — the per-format [FormatErr]-built sentinel
	// (e.g. [ErrFormatPhoneNumber]), [ErrRequired] (recoverable
	// via [errors.As] to learn which field was missing), or a
	// format-specific sentinel like [ErrOpaqueEmpty] or
	// [ErrNestedAliases]. Every return value matches the umbrella
	// [Err] under [errors.Is].
	Validate() error

	// sealed is the unexported marker that confines
	// implementations of SubjectIdentifier to this package. It
	// has no behavior.
	sealed()
}
