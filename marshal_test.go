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

// TestMarshalEachBuiltinFormatIsSpecOrderAndByteStable pins the
// canonical output for each built-in. The expected bytes are
// what RFC 9493 §3 figures show after canonical-JSON
// whitespace normalization.
func TestMarshalEachBuiltinFormatIsSpecOrderAndByteStable(t *testing.T) {
	tests := []struct {
		value subjectid.SubjectIdentifier
		want  string
		name  string
	}{
		{
			name:  "account",
			value: subjectid.AccountID{URI: "acct:example.user@service.example.com"},
			want:  `{"format":"account","uri":"acct:example.user@service.example.com"}`,
		},
		{
			name:  "email",
			value: subjectid.EmailID{Email: "user@example.com"},
			want:  `{"format":"email","email":"user@example.com"}`,
		},
		{
			name:  "iss_sub",
			value: subjectid.IssSubID{Iss: "https://issuer.example.com/", Sub: "145234573"},
			want:  `{"format":"iss_sub","iss":"https://issuer.example.com/","sub":"145234573"}`,
		},
		{
			name:  "opaque",
			value: subjectid.OpaqueID{ID: "11112222333344445555"},
			want:  `{"format":"opaque","id":"11112222333344445555"}`,
		},
		{
			name:  "phone_number",
			value: subjectid.PhoneNumberID{PhoneNumber: "+12065550100"},
			want:  `{"format":"phone_number","phone_number":"+12065550100"}`,
		},
		{
			name:  "did",
			value: subjectid.DIDID{URL: "did:example:123456"},
			want:  `{"format":"did","url":"did:example:123456"}`,
		},
		{
			name:  "uri",
			value: subjectid.URIID{URI: "urn:oasis:names:tc:saml:2.0:nameid-format:transient"},
			want:  `{"format":"uri","uri":"urn:oasis:names:tc:saml:2.0:nameid-format:transient"}`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.value)
			if err != nil {
				t.Fatalf("Marshal(%T): err = %v", tc.value, err)
			}
			if string(got) != tc.want {
				t.Errorf("Marshal output:\n  got:  %s\n  want: %s", got, tc.want)
			}
		})
	}
}

func TestMarshalAliasesIsSpecOrderAndDelegatesToElementMarshalers(t *testing.T) {
	v := subjectid.AliasesID{
		Identifiers: []subjectid.SubjectIdentifier{
			subjectid.EmailID{Email: "user@example.com"},
			subjectid.AccountID{URI: "acct:user@example.com"},
		},
	}
	got, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Marshal: err = %v", err)
	}
	want := `{"format":"aliases","identifiers":[` +
		`{"format":"email","email":"user@example.com"},` +
		`{"format":"account","uri":"acct:user@example.com"}` +
		`]}`
	if string(got) != want {
		t.Errorf("Marshal output:\n  got:  %s\n  want: %s", got, want)
	}
}

// TestMarshalUnknownFormatReturnsRawBytesCompacted exercises the
// round-trip contract: an UnknownFormat with populated Raw emits
// the Raw bytes verbatim from MarshalJSON, which the
// encoding/json package then runs through json.Compact before
// writing — so the surrounding output sees canonicalized
// whitespace. The byte-stability invariant holds modulo this
// compaction, which matches the conformance phase's expectation.
func TestMarshalUnknownFormatReturnsRawBytesCompacted(t *testing.T) {
	raw := json.RawMessage(`{"format": "org.example.future", "extra": [1, 2, 3]}`)
	u := subjectid.UnknownFormat{
		FormatName: "org.example.future",
		Raw:        raw,
	}
	got, err := json.Marshal(u)
	if err != nil {
		t.Fatalf("Marshal: err = %v", err)
	}
	want := `{"format":"org.example.future","extra":[1,2,3]}`
	if string(got) != want {
		t.Errorf("Marshal output:\n  got:  %s\n  want: %s", got, want)
	}
}

func TestMarshalUnknownFormatEmptyRawEmitsFormatOnlyObject(t *testing.T) {
	u := subjectid.UnknownFormat{FormatName: "org.example.future"}
	got, err := json.Marshal(u)
	if err != nil {
		t.Fatalf("Marshal: err = %v", err)
	}
	want := `{"format":"org.example.future"}`
	if string(got) != want {
		t.Errorf("Marshal output:\n  got:  %s\n  want: %s", got, want)
	}
}

func TestMarshalUnknownFormatZeroValueErrors(t *testing.T) {
	_, err := json.Marshal(subjectid.UnknownFormat{})
	if err == nil {
		t.Fatal("Marshal of zero-value UnknownFormat returned nil error; want a refusal")
	}
	// The json package wraps marshaler errors in *json.MarshalerError;
	// unwrap to find ours.
	for err != nil {
		if strings.Contains(err.Error(), "UnknownFormat") {
			return
		}
		err = errors.Unwrap(err)
	}
	t.Fatal("error chain did not mention UnknownFormat")
}

func TestMarshalUnknownFormatInvalidRawErrors(t *testing.T) {
	u := subjectid.UnknownFormat{
		FormatName: "org.example.future",
		Raw:        json.RawMessage(`{not-json`),
	}
	if _, err := json.Marshal(u); err == nil {
		t.Fatal("Marshal of invalid-Raw UnknownFormat returned nil error")
	}
}

// TestParseMarshalRoundTripIsByteStableForCanonicalInput
// verifies the round-trip invariant the spec-fixtures phase
// will lean on: canonical-input → Parse → Marshal returns the
// exact same bytes.
func TestParseMarshalRoundTripIsByteStableForCanonicalInput(t *testing.T) {
	inputs := []string{
		`{"format":"account","uri":"acct:user@example.com"}`,
		`{"format":"email","email":"a@b.example"}`,
		`{"format":"iss_sub","iss":"https://i.example/","sub":"1"}`,
		`{"format":"opaque","id":"id-1"}`,
		`{"format":"phone_number","phone_number":"+15555550100"}`,
		`{"format":"did","url":"did:web:example.com"}`,
		`{"format":"uri","uri":"urn:example:1"}`,
		`{"format":"aliases","identifiers":[` +
			`{"format":"email","email":"a@b.example"},` +
			`{"format":"account","uri":"acct:a@b.example"}` +
			`]}`,
	}
	for _, in := range inputs {
		t.Run(in, func(t *testing.T) {
			v, err := subjectid.Parse(json.RawMessage(in))
			if err != nil {
				t.Fatalf("Parse: err = %v", err)
			}
			out, err := json.Marshal(v)
			if err != nil {
				t.Fatalf("Marshal: err = %v", err)
			}
			if string(out) != in {
				t.Errorf("round-trip drift:\n  in:  %s\n  out: %s", in, out)
			}
		})
	}
}
