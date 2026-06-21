// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/hstern/go-subjectid"
)

// TestParseEachBuiltinFormatDispatchesToCorrectType exercises one
// hand-crafted example per built-in format. Spec-example fixtures
// from RFC 9493 §3 land in phase 5 (specfixtures); this test is
// the per-format smoke check.
func TestParseEachBuiltinFormatDispatchesToCorrectType(t *testing.T) {
	tests := []struct {
		check func(t *testing.T, got subjectid.SubjectIdentifier)
		name  string
		raw   string
	}{
		{
			name: "account",
			raw:  `{"format":"account","uri":"acct:example.user@service.example.com"}`,
			check: func(t *testing.T, got subjectid.SubjectIdentifier) {
				a, ok := got.(subjectid.AccountID)
				if !ok {
					t.Fatalf("got %T, want AccountID", got)
				}
				if a.URI != "acct:example.user@service.example.com" {
					t.Errorf("URI = %q", a.URI)
				}
			},
		},
		{
			name: "email",
			raw:  `{"format":"email","email":"user@example.com"}`,
			check: func(t *testing.T, got subjectid.SubjectIdentifier) {
				e, ok := got.(subjectid.EmailID)
				if !ok {
					t.Fatalf("got %T, want EmailID", got)
				}
				if e.Email != "user@example.com" {
					t.Errorf("Email = %q", e.Email)
				}
			},
		},
		{
			name: "iss_sub",
			raw:  `{"format":"iss_sub","iss":"https://issuer.example.com/","sub":"145234573"}`,
			check: func(t *testing.T, got subjectid.SubjectIdentifier) {
				i, ok := got.(subjectid.IssSubID)
				if !ok {
					t.Fatalf("got %T, want IssSubID", got)
				}
				if i.Iss != "https://issuer.example.com/" || i.Sub != "145234573" {
					t.Errorf("Iss=%q Sub=%q", i.Iss, i.Sub)
				}
			},
		},
		{
			name: "opaque",
			raw:  `{"format":"opaque","id":"11112222333344445555"}`,
			check: func(t *testing.T, got subjectid.SubjectIdentifier) {
				o, ok := got.(subjectid.OpaqueID)
				if !ok {
					t.Fatalf("got %T, want OpaqueID", got)
				}
				if o.ID != "11112222333344445555" {
					t.Errorf("ID = %q", o.ID)
				}
			},
		},
		{
			name: "phone_number",
			raw:  `{"format":"phone_number","phone_number":"+12065550100"}`,
			check: func(t *testing.T, got subjectid.SubjectIdentifier) {
				p, ok := got.(subjectid.PhoneNumberID)
				if !ok {
					t.Fatalf("got %T, want PhoneNumberID", got)
				}
				if p.PhoneNumber != "+12065550100" {
					t.Errorf("PhoneNumber = %q", p.PhoneNumber)
				}
			},
		},
		{
			name: "did",
			raw:  `{"format":"did","url":"did:example:123456"}`,
			check: func(t *testing.T, got subjectid.SubjectIdentifier) {
				d, ok := got.(subjectid.DIDID)
				if !ok {
					t.Fatalf("got %T, want DIDID", got)
				}
				if d.URL != "did:example:123456" {
					t.Errorf("URL = %q", d.URL)
				}
			},
		},
		{
			name: "uri",
			raw:  `{"format":"uri","uri":"urn:oasis:names:tc:saml:2.0:nameid-format:transient"}`,
			check: func(t *testing.T, got subjectid.SubjectIdentifier) {
				u, ok := got.(subjectid.URIID)
				if !ok {
					t.Fatalf("got %T, want URIID", got)
				}
				if u.URI != "urn:oasis:names:tc:saml:2.0:nameid-format:transient" {
					t.Errorf("URI = %q", u.URI)
				}
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := subjectid.Parse(json.RawMessage(tc.raw))
			if err != nil {
				t.Fatalf("Parse(%s): err = %v", tc.raw, err)
			}
			if got.Format() != tc.name {
				t.Errorf("Format() = %q, want %q", got.Format(), tc.name)
			}
			tc.check(t, got)
		})
	}
}

// TestParseAliasesIsRecursiveAndHeterogeneous walks an aliases
// payload containing two distinct inner formats, both built-in.
// The inner UnmarshalJSON dispatches go through the same registry
// Parse uses.
func TestParseAliasesIsRecursiveAndHeterogeneous(t *testing.T) {
	raw := `{
	  "format": "aliases",
	  "identifiers": [
	    {"format": "email",   "email": "user@example.com"},
	    {"format": "account", "uri":   "acct:user@example.com"}
	  ]
	}`
	got, err := subjectid.Parse(json.RawMessage(raw))
	if err != nil {
		t.Fatalf("Parse aliases: err = %v", err)
	}
	a, ok := got.(subjectid.AliasesID)
	if !ok {
		t.Fatalf("got %T, want AliasesID", got)
	}
	if len(a.Identifiers) != 2 {
		t.Fatalf("len(Identifiers) = %d, want 2", len(a.Identifiers))
	}
	if _, ok := a.Identifiers[0].(subjectid.EmailID); !ok {
		t.Errorf("Identifiers[0] = %T, want EmailID", a.Identifiers[0])
	}
	if _, ok := a.Identifiers[1].(subjectid.AccountID); !ok {
		t.Errorf("Identifiers[1] = %T, want AccountID", a.Identifiers[1])
	}
}

// TestParseAliasesInnerUnknownFormatBecomesUnknownFormat verifies
// the spec-aligned behavior: an aliases identifier with an inner
// element using an unregistered format does not fail; the inner
// element becomes UnknownFormat and the aliases value is still
// usable.
func TestParseAliasesInnerUnknownFormatBecomesUnknownFormat(t *testing.T) {
	raw := `{
	  "format": "aliases",
	  "identifiers": [
	    {"format": "email", "email": "user@example.com"},
	    {"format": "org.example.future", "data": "opaque-payload"}
	  ]
	}`
	got, err := subjectid.Parse(json.RawMessage(raw))
	if err != nil {
		t.Fatalf("Parse aliases: err = %v", err)
	}
	a := got.(subjectid.AliasesID)
	u, ok := a.Identifiers[1].(subjectid.UnknownFormat)
	if !ok {
		t.Fatalf("Identifiers[1] = %T, want UnknownFormat", a.Identifiers[1])
	}
	if u.FormatName != "org.example.future" {
		t.Errorf("FormatName = %q", u.FormatName)
	}
	if !strings.Contains(string(u.Raw), `"opaque-payload"`) {
		t.Errorf("Raw bytes lost the inner element's data member: %s", u.Raw)
	}
}

func TestParseUnknownTopLevelFormatBecomesUnknownFormat(t *testing.T) {
	raw := `{"format":"org.example.future","custom":42}`
	got, err := subjectid.Parse(json.RawMessage(raw))
	if err != nil {
		t.Fatalf("Parse: err = %v", err)
	}
	u, ok := got.(subjectid.UnknownFormat)
	if !ok {
		t.Fatalf("got %T, want UnknownFormat", got)
	}
	if u.FormatName != "org.example.future" {
		t.Errorf("FormatName = %q", u.FormatName)
	}
	if string(u.Raw) != raw {
		t.Errorf("Raw = %s, want %s", u.Raw, raw)
	}
}

func TestParseUnknownFormatRawIsACopy(t *testing.T) {
	src := []byte(`{"format":"org.example.copy","x":1}`)
	got, err := subjectid.Parse(src)
	if err != nil {
		t.Fatalf("Parse: err = %v", err)
	}
	u := got.(subjectid.UnknownFormat)
	// Mutate the source; the carrier should keep the original
	// bytes because Parse copied them.
	src[0] = '!'
	if u.Raw[0] != '{' {
		t.Errorf("Raw is aliased to input bytes: Raw[0] = %q after src mutation", u.Raw[0])
	}
}

func TestParseMissingFormatMemberIsErrRequired(t *testing.T) {
	cases := []string{
		`{}`,
		`{"email":"x@example.com"}`,
		`{"format":""}`,
	}
	for _, raw := range cases {
		t.Run(raw, func(t *testing.T) {
			_, err := subjectid.Parse(json.RawMessage(raw))
			if err == nil {
				t.Fatal("Parse returned nil error; want non-nil error")
			}
			if !errors.Is(err, subjectid.ErrRequired{}) {
				t.Errorf("errors.Is(err, ErrRequired) = false, want true (err = %v)", err)
			}
		})
	}
}

func TestParseMalformedJSONIsErrJSON(t *testing.T) {
	_, err := subjectid.Parse(json.RawMessage(`{not-json`))
	if err == nil {
		t.Fatal("Parse returned nil error")
	}
	if !errors.Is(err, subjectid.ErrJSON) {
		t.Errorf("errors.Is(err, ErrJSON) = false, want true (err = %v)", err)
	}
}

// extensionID exercises the RegisterFormat → Parse round-trip
// path: an external type embedding Seal, registered as the
// constructor for a custom format string, decoded via Parse.
type extensionID struct {
	subjectid.Seal
	Tag string
}

func (extensionID) Format() string  { return "subjectid.test.parse.extension" }
func (extensionID) Validate() error { return nil }

func (e *extensionID) UnmarshalJSON(data []byte) error {
	var v struct {
		Tag string `json:"tag"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	e.Tag = v.Tag
	return nil
}

func TestParseDispatchesExtensionFormatRegisteredViaRegisterFormat(t *testing.T) {
	const name = "subjectid.test.parse.extension"
	if err := subjectid.RegisterFormat(name, func() subjectid.SubjectIdentifier {
		return &extensionID{}
	}); err != nil {
		t.Fatalf("RegisterFormat: err = %v", err)
	}
	raw := `{"format":"subjectid.test.parse.extension","tag":"hello"}`
	got, err := subjectid.Parse(json.RawMessage(raw))
	if err != nil {
		t.Fatalf("Parse: err = %v", err)
	}
	e, ok := got.(extensionID)
	if !ok {
		t.Fatalf("got %T, want extensionID", got)
	}
	if e.Tag != "hello" {
		t.Errorf("Tag = %q, want %q", e.Tag, "hello")
	}
}

func TestPerTypeUnmarshalIgnoresExtraMembers(t *testing.T) {
	raw := `{"format":"email","email":"u@example.com","extra":42,"nested":{"k":"v"}}`
	var e subjectid.EmailID
	if err := json.Unmarshal([]byte(raw), &e); err != nil {
		t.Fatalf("Unmarshal: err = %v", err)
	}
	if e.Email != "u@example.com" {
		t.Errorf("Email = %q", e.Email)
	}
}

func TestUnknownFormatUnmarshalJSONPopulatesFormatAndRaw(t *testing.T) {
	raw := []byte(`{"format":"org.example.future","x":1}`)
	var u subjectid.UnknownFormat
	if err := json.Unmarshal(raw, &u); err != nil {
		t.Fatalf("Unmarshal: err = %v", err)
	}
	if u.FormatName != "org.example.future" {
		t.Errorf("FormatName = %q", u.FormatName)
	}
	if string(u.Raw) != string(raw) {
		t.Errorf("Raw = %s", u.Raw)
	}
}
