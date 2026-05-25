// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import (
	"errors"
	"fmt"
	"strings"
)

// Err is the top-level umbrella sentinel for every error this
// package returns. Every other sentinel ([ErrJSON], [ErrFormatDID],
// [ErrFormatReserved], the per-validator format errors built via
// [FormatErr], and the [ErrRequired] struct type) is detectable as
// a [subjectid.Err] via [errors.Is].
//
// Use this when handling errors from multiple sources to branch on
// "did this come from subjectid?" without naming each leaf sentinel:
//
//	if errors.Is(err, subjectid.Err) {
//	    // handle as a subjectid validation/codec failure
//	} else {
//	    // unrelated to subjectid
//	}
var Err = errors.New("subjectid")

// FormatErr builds a sentinel error from a descriptive message,
// wrapping [Err] so the result matches errors.Is(returned, Err).
//
// Pass the full message the caller wants in the error text — the
// helper adds no "invalid" prefix or other framing. The result reads
// "subjectid: <msg>".
//
// Built-in format sentinels (ErrFormatAccount, ErrFormatEmail, etc.)
// and ErrFormatReserved are constructed via this helper. Consumers
// that register their own format via RegisterFormat can use FormatErr
// to mint a matching sentinel so callers can branch with errors.Is
// in the same idiom that works for the built-ins:
//
//	var ErrFormatMyCustom = subjectid.FormatErr("invalid my-custom format")
//
//	func (m MyCustomID) Validate() error {
//	    if !myRegex.MatchString(m.Value) {
//	        return ErrFormatMyCustom
//	    }
//	    return nil
//	}
//
// The returned error is a fresh value; two FormatErr calls with the
// same message produce distinct sentinels that do NOT compare equal
// under errors.Is. Store the result in a package-level var if callers
// need a stable comparison target.
func FormatErr(msg string) error {
	return fmt.Errorf("%w: %s", Err, msg)
}

// ErrRequired reports that one or more required members of the wire
// shape were missing or empty. RFC 9493 §3 requires every format's
// defining members to be present and non-empty.
//
// Fields names the specific member(s) that were missing — values like
// "phone_number", "url", or []string{"iss", "sub"}. Use [errors.As]
// to recover them:
//
//	var req subjectid.ErrRequired
//	if errors.As(err, &req) {
//	    log.Printf("missing fields: %v", req.Fields)
//	}
//
// [errors.Is] also matches against the zero value, so categorical
// checks work without naming the fields:
//
//	if errors.Is(err, subjectid.ErrRequired{}) { ... }
//
// And against the top-level package umbrella:
//
//	if errors.Is(err, subjectid.Err) { ... }
type ErrRequired struct {
	// Fields names the wire-shape members that were missing or
	// empty. An empty slice is used when the failing field set
	// is not known.
	Fields []string
}

// Error implements the error interface.
func (e ErrRequired) Error() string {
	switch len(e.Fields) {
	case 0:
		return "subjectid: required field is missing or empty"
	case 1:
		return "subjectid: required field " + e.Fields[0] + " is missing or empty"
	default:
		return "subjectid: required fields " + strings.Join(e.Fields, ", ") + " are missing or empty"
	}
}

// Is matches any ErrRequired regardless of Fields (so callers can
// write errors.Is(err, ErrRequired{}) for a categorical check) and
// also matches the package umbrella [Err].
func (ErrRequired) Is(target error) bool {
	if target == Err {
		return true
	}
	_, ok := target.(ErrRequired)
	return ok
}

// MissingFields constructs an ErrRequired naming the specific wire-
// shape member(s) that were missing or empty. Pass one name for the
// common single-field case, or several for composite formats like
// iss_sub where multiple members can be missing at once.
func MissingFields(names ...string) error {
	return ErrRequired{Fields: names}
}

// Sentinel errors categorizing the kinds of well-formedness failure
// a Subject Identifier can have. Each is the [errors.Is] target a
// caller uses to branch on the failure category. Each also wraps
// [Err] so caller code that only cares whether the error is a
// subjectid problem can use errors.Is(err, subjectid.Err).
var (
	// ErrJSON reports that the bytes handed to Parse are not a
	// JSON object exposing a "format" string member.
	ErrJSON = fmt.Errorf("%w: input is not a Subject Identifier JSON object", Err)

	// ErrNestedAliases reports that an aliases identifier contains
	// another aliases identifier, which RFC 9493 §3.2.7 forbids.
	ErrNestedAliases = fmt.Errorf("%w: aliases identifier must not contain a nested aliases identifier", Err)

	// ErrAliasesEmpty reports that an aliases identifier's
	// "identifiers" array was nil or zero-length. RFC 9493 §3.2.7
	// requires at least one element.
	ErrAliasesEmpty = fmt.Errorf("%w: aliases identifier must contain at least one element", Err)

	// ErrOpaqueEmpty reports that an OpaqueID's "id" member was
	// the empty string. RFC 9493 §3.2.4 imposes no internal
	// structure on opaque identifiers but §3 still forbids empty
	// values. This sentinel exists alongside ErrRequired so
	// callers can distinguish "the wire shape was missing the
	// member" (ErrRequired, raised by Parse before the member is
	// inspected) from "the member was present but the empty
	// string" (ErrOpaqueEmpty, raised by OpaqueID.Validate).
	ErrOpaqueEmpty = fmt.Errorf("%w: opaque identifier is the empty string", Err)

	// Per-format invalid-value sentinels, one per RFC 9493 format
	// whose value is structurally checked. Built via FormatErr so
	// the surface stays uniform and each wraps Err.
	ErrFormatAccount     = FormatErr("invalid acct URI")
	ErrFormatEmail       = FormatErr("invalid email addr-spec")
	ErrFormatIssSub      = FormatErr("invalid iss/sub identifier")
	ErrFormatPhoneNumber = FormatErr("invalid E.164 phone number")
	ErrFormatDID         = FormatErr("invalid DID URL")
	ErrFormatURI         = FormatErr("invalid absolute URI")

	// ErrFormatReserved is returned by RegisterFormat (defined in
	// a later commit) when a caller attempts to register a
	// constructor for a format name that the library has already
	// populated as a built-in.
	ErrFormatReserved = FormatErr("format name is reserved for a built-in")
)
