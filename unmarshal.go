// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// envelope is the wire shape Parse peeks at first to discover
// the "format" discriminator. It exists only to read that one
// member without binding the rest of the JSON to any concrete
// type yet.
type envelope struct {
	Format string `json:"format"`
}

// deref normalizes a freshly-parsed identifier to its canonical
// value form. Registry constructors allocate a pointer so the codec
// can populate it via UnmarshalJSON; Parse returns the dereferenced
// value so the dynamic type a caller reads back from Parse matches
// the value literal they would write by hand (subjectid.IssSubID,
// not *subjectid.IssSubID). See [SubjectIdentifier] for the
// canonical-form contract.
//
// Reflection keeps this uniform across the built-in formats and any
// extension type registered via [RegisterFormat]; the constraint it
// imposes — extension types must satisfy [SubjectIdentifier] with
// value receivers, so the dereferenced value still implements the
// interface — is documented on RegisterFormat. A non-pointer is
// returned unchanged.
//
// An extension type registered with pointer-receiver methods would
// not satisfy the interface in dereferenced form; rather than panic
// or drop it, deref returns such a value unchanged. Normalization is
// thus best-effort, and value-canonical is guaranteed only for the
// built-ins and conforming extensions.
func deref(id SubjectIdentifier) SubjectIdentifier {
	rv := reflect.ValueOf(id)
	if rv.Kind() != reflect.Pointer {
		return id
	}
	if v, ok := rv.Elem().Interface().(SubjectIdentifier); ok {
		return v
	}
	return id
}

// Parse decodes a Subject Identifier from raw JSON, dispatching
// on the value of the "format" member to the appropriate
// concrete type. The returned identifier is always in value form
// (e.g. [IssSubID], never *IssSubID) — the canonical dynamic form
// for every value this package produces, matching what a caller
// constructs as a struct literal. See [SubjectIdentifier].
//
//   - For a built-in format (account, email, iss_sub, opaque,
//     phone_number, did, uri, aliases), the matching per-format
//     type is returned with its members populated from the wire.
//   - For an extension format registered via [RegisterFormat],
//     the registered constructor is invoked and json.Unmarshal
//     delegates to its UnmarshalJSON.
//   - For a format string the registry does not know, an
//     [UnknownFormat] is returned carrying the original bytes
//     verbatim. The library never errors on an unrecognized
//     format — that is the forward-compatibility contract.
//
// Errors are reserved for two conditions only:
//
//   - The bytes are not valid JSON, or the "format" member is
//     not present or not a string. In the first case the error
//     wraps [ErrJSON]; in the second it is [ErrRequired] with
//     "format" in its Fields. Both match [Err] under [errors.Is].
//   - A registered constructor's UnmarshalJSON returns an
//     error. The error is returned verbatim.
//
// Extra members in the JSON object that are not defined for the
// format are silently dropped, per RFC 9493's "be liberal in what
// you accept" reading of the wire shape. Strict checking happens
// at marshal time, not at unmarshal.
func Parse(raw json.RawMessage) (SubjectIdentifier, error) {
	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrJSON, err)
	}
	if env.Format == "" {
		return nil, MissingFields("format")
	}
	if ctor := lookup(env.Format); ctor != nil {
		target := ctor()
		if err := json.Unmarshal(raw, target); err != nil {
			return nil, err
		}
		return deref(target), nil
	}
	// Unknown format — copy the bytes so the caller cannot mutate
	// our internal state by holding a reference to raw.
	bytes := make(json.RawMessage, len(raw))
	copy(bytes, raw)
	return UnknownFormat{
		FormatName: env.Format,
		Raw:        bytes,
	}, nil
}

// UnmarshalJSON implements [json.Unmarshaler] for AccountID. It
// decodes the "uri" member; the "format" member and any other
// extra members are ignored (Postel's law on receive). Use Parse
// to dispatch on "format" automatically.
func (a *AccountID) UnmarshalJSON(data []byte) error {
	var v struct {
		URI string `json:"uri"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	a.URI = v.URI
	return nil
}

// UnmarshalJSON implements [json.Unmarshaler] for EmailID. See
// [AccountID.UnmarshalJSON] for the general contract.
func (e *EmailID) UnmarshalJSON(data []byte) error {
	var v struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	e.Email = v.Email
	return nil
}

// UnmarshalJSON implements [json.Unmarshaler] for IssSubID. See
// [AccountID.UnmarshalJSON] for the general contract.
func (i *IssSubID) UnmarshalJSON(data []byte) error {
	var v struct {
		Iss string `json:"iss"`
		Sub string `json:"sub"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	i.Iss = v.Iss
	i.Sub = v.Sub
	return nil
}

// UnmarshalJSON implements [json.Unmarshaler] for OpaqueID. See
// [AccountID.UnmarshalJSON] for the general contract.
func (o *OpaqueID) UnmarshalJSON(data []byte) error {
	var v struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	o.ID = v.ID
	return nil
}

// UnmarshalJSON implements [json.Unmarshaler] for PhoneNumberID.
// See [AccountID.UnmarshalJSON] for the general contract.
func (p *PhoneNumberID) UnmarshalJSON(data []byte) error {
	var v struct {
		PhoneNumber string `json:"phone_number"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	p.PhoneNumber = v.PhoneNumber
	return nil
}

// UnmarshalJSON implements [json.Unmarshaler] for DIDID. See
// [AccountID.UnmarshalJSON] for the general contract.
func (d *DIDID) UnmarshalJSON(data []byte) error {
	var v struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	d.URL = v.URL
	return nil
}

// UnmarshalJSON implements [json.Unmarshaler] for URIID. See
// [AccountID.UnmarshalJSON] for the general contract.
func (u *URIID) UnmarshalJSON(data []byte) error {
	var v struct {
		URI string `json:"uri"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	u.URI = v.URI
	return nil
}

// UnmarshalJSON implements [json.Unmarshaler] for AliasesID by
// reading the "identifiers" array as raw JSON elements and
// recursively dispatching each through [Parse]. The element
// types are arbitrary — heterogeneous slices fall out of the
// uniform dispatch.
//
// An element whose "format" is not recognized by the registry
// becomes an [UnknownFormat] in the resulting slice rather than
// failing the whole AliasesID. An element that is not valid JSON
// or is missing the "format" member fails Parse and is returned
// as the AliasesID-level error — the spec rule that aliases
// contain non-null subject identifiers leaves no room for a
// totally unparseable inner element.
func (a *AliasesID) UnmarshalJSON(data []byte) error {
	var v struct {
		Identifiers []json.RawMessage `json:"identifiers"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	a.Identifiers = make([]SubjectIdentifier, 0, len(v.Identifiers))
	for n, raw := range v.Identifiers {
		inner, err := Parse(raw)
		if err != nil {
			return fmt.Errorf("aliases element %d: %w", n, err)
		}
		a.Identifiers = append(a.Identifiers, inner)
	}
	return nil
}

// UnmarshalJSON implements [json.Unmarshaler] for UnknownFormat.
// Stores the "format" discriminator in FormatName and the entire
// input bytes verbatim in Raw so MarshalJSON (a later commit) can
// re-emit byte-for-byte.
//
// UnknownFormat is normally produced by [Parse] for unrecognized
// formats; this method is provided so an external caller can
// also json.Unmarshal directly into an UnknownFormat target
// when its purpose is to round-trip the bytes whole.
func (u *UnknownFormat) UnmarshalJSON(data []byte) error {
	var env envelope
	if err := json.Unmarshal(data, &env); err != nil {
		return err
	}
	bytes := make(json.RawMessage, len(data))
	copy(bytes, data)
	u.FormatName = env.Format
	u.Raw = bytes
	return nil
}
