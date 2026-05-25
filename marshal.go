// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import (
	"encoding/json"
	"errors"
)

// MarshalJSON implements [json.Marshaler] for AccountID. The
// output is the spec-order JSON object
//
//	{"format":"account","uri":"<URI>"}
//
// — format member first, then the format-specific members in the
// order RFC 9493 §3.2.1 defines them. Output is canonical
// (no whitespace) and byte-stable for a given input.
func (a AccountID) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Format string `json:"format"`
		URI    string `json:"uri"`
	}{Format: "account", URI: a.URI})
}

// MarshalJSON implements [json.Marshaler] for EmailID. See
// RFC 9493 §3.2.2 for the wire shape.
func (e EmailID) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Format string `json:"format"`
		Email  string `json:"email"`
	}{Format: "email", Email: e.Email})
}

// MarshalJSON implements [json.Marshaler] for IssSubID. See
// RFC 9493 §3.2.3 — iss precedes sub, matching the JWT claim
// order.
func (i IssSubID) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Format string `json:"format"`
		Iss    string `json:"iss"`
		Sub    string `json:"sub"`
	}{Format: "iss_sub", Iss: i.Iss, Sub: i.Sub})
}

// MarshalJSON implements [json.Marshaler] for OpaqueID. See
// RFC 9493 §3.2.4 for the wire shape.
func (o OpaqueID) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Format string `json:"format"`
		ID     string `json:"id"`
	}{Format: "opaque", ID: o.ID})
}

// MarshalJSON implements [json.Marshaler] for PhoneNumberID. See
// RFC 9493 §3.2.5 for the wire shape.
func (p PhoneNumberID) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Format      string `json:"format"`
		PhoneNumber string `json:"phone_number"`
	}{Format: "phone_number", PhoneNumber: p.PhoneNumber})
}

// MarshalJSON implements [json.Marshaler] for DIDID. See
// RFC 9493 §3.2.6 for the wire shape.
func (d DIDID) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Format string `json:"format"`
		URL    string `json:"url"`
	}{Format: "did", URL: d.URL})
}

// MarshalJSON implements [json.Marshaler] for URIID. See
// RFC 9493 §3.2.7 for the wire shape.
func (u URIID) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Format string `json:"format"`
		URI    string `json:"uri"`
	}{Format: "uri", URI: u.URI})
}

// MarshalJSON implements [json.Marshaler] for AliasesID. The
// inner Identifiers slice is marshaled element-by-element; each
// element's own MarshalJSON is invoked, so a heterogeneous
// slice round-trips with each element in its own format's
// spec-order shape.
func (a AliasesID) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Format      string              `json:"format"`
		Identifiers []SubjectIdentifier `json:"identifiers"`
	}{Format: "aliases", Identifiers: a.Identifiers})
}

// errUnknownFormatEmpty is returned by [UnknownFormat.MarshalJSON]
// when both Raw and FormatName are empty — the zero value cannot
// produce a valid Subject Identifier object on its own.
var errUnknownFormatEmpty = errors.New(
	"subjectid: UnknownFormat has neither Raw bytes nor a FormatName to emit")

// MarshalJSON implements [json.Marshaler] for UnknownFormat by
// returning the Raw bytes verbatim when populated, so the wire
// payload round-trips through the codec. Note that
// [encoding/json.Marshal] runs [json.Compact] over a
// [json.Marshaler]'s output, so insignificant whitespace inside
// Raw is normalized on the way out — the surrounding payload
// sees canonical JSON. Round-trip byte stability therefore
// holds modulo [json.Compact], matching the conformance test
// suite's expectation.
//
// Two edge cases:
//
//   - Raw is non-empty but not a valid JSON object: an error is
//     returned. UnknownFormat's contract is "round-trip valid
//     JSON the library does not recognize"; emitting invalid
//     JSON would silently corrupt the surrounding payload.
//   - Raw is empty: a minimal {"format":"<FormatName>"} object
//     is emitted. If FormatName is also empty an error is
//     returned — there is no Subject Identifier to write.
//
// Mutating FormatName after a Parse does not change the emitted
// bytes (Raw is authoritative when present); to change the wire
// format string, set Raw to nil before marshaling.
func (u UnknownFormat) MarshalJSON() ([]byte, error) {
	if len(u.Raw) > 0 {
		if !json.Valid(u.Raw) {
			return nil, errors.New("subjectid: UnknownFormat.Raw is not valid JSON")
		}
		return u.Raw, nil
	}
	if u.FormatName == "" {
		return nil, errUnknownFormatEmpty
	}
	return json.Marshal(struct {
		Format string `json:"format"`
	}{Format: u.FormatName})
}
