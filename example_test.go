// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hstern/go-subjectid"
)

// ExampleParse shows the discriminator-driven decode path:
// [subjectid.Parse] peeks at the wire's "format" member and
// dispatches to the matching concrete type, returning a value
// that satisfies [subjectid.SubjectIdentifier].
func ExampleParse() {
	wire := []byte(`{"format":"email","email":"user@example.com"}`)

	id, err := subjectid.Parse(json.RawMessage(wire))
	if err != nil {
		panic(err)
	}

	email := id.(subjectid.EmailID)
	fmt.Println(id.Format(), email.Email)
	// Output: email user@example.com
}

// ExampleParse_unknownFormat shows the forward-compatibility
// path: a format the library does not recognize parses into
// [subjectid.UnknownFormat] carrying the original bytes verbatim,
// so the payload still round-trips.
func ExampleParse_unknownFormat() {
	wire := []byte(`{"format":"org.example.future","widget_id":42}`)

	id, err := subjectid.Parse(json.RawMessage(wire))
	if err != nil {
		panic(err)
	}

	unk := id.(subjectid.UnknownFormat)
	out, _ := json.Marshal(unk)
	fmt.Println(unk.FormatName)
	fmt.Println(string(out))
	// Output:
	// org.example.future
	// {"format":"org.example.future","widget_id":42}
}

// ExampleSubjectIdentifier_Format shows the per-format
// discriminator: every built-in type returns its IANA registry
// name from [subjectid.SubjectIdentifier.Format], the same value
// the JSON "format" member carries on the wire.
func ExampleSubjectIdentifier_Format() {
	ids := []subjectid.SubjectIdentifier{
		subjectid.AccountID{URI: "acct:user@example.com"},
		subjectid.EmailID{Email: "user@example.com"},
		subjectid.IssSubID{Iss: "https://issuer.example.com/", Sub: "1"},
		subjectid.OpaqueID{ID: "abc123"},
	}
	for _, id := range ids {
		fmt.Println(id.Format())
	}
	// Output:
	// account
	// email
	// iss_sub
	// opaque
}

// ExampleAccountID_Validate shows the per-format validator
// returning nil for a well-formed acct: URI and a wrapped
// [subjectid.ErrFormatAccount] sentinel for a malformed one.
func ExampleAccountID_Validate() {
	ok := subjectid.AccountID{URI: "acct:user@example.com"}
	bad := subjectid.AccountID{URI: "not-an-acct-uri"}

	fmt.Println(ok.Validate())
	fmt.Println(errors.Is(bad.Validate(), subjectid.ErrFormatAccount))
	// Output:
	// <nil>
	// true
}

// ExampleEmailID_Validate shows the addr-spec syntax check
// firing for malformed input. The empty-string case returns
// [subjectid.ErrRequired] (a struct type recoverable via
// [errors.As]) rather than the format sentinel — "missing
// member" is a different rule than "syntax violation".
func ExampleEmailID_Validate() {
	bad := subjectid.EmailID{Email: ""}
	err := bad.Validate()

	var req subjectid.ErrRequired
	if errors.As(err, &req) {
		fmt.Println("missing:", req.Fields)
	}
	// Output: missing: [email]
}

// ExampleIssSubID_Validate shows the per-member check on the
// iss_sub composite: both "iss" and "sub" must be non-empty, and
// "iss" must be an RFC 3986 absolute URI.
func ExampleIssSubID_Validate() {
	id := subjectid.IssSubID{
		Iss: "https://issuer.example.com/",
		Sub: "145234573",
	}
	fmt.Println(id.Validate())
	// Output: <nil>
}

// ExampleOpaqueID_Validate shows the only failure mode:
// [subjectid.ErrOpaqueEmpty] when ID is the empty string. The
// opaque format imposes no syntax beyond non-emptiness.
func ExampleOpaqueID_Validate() {
	fmt.Println(errors.Is(subjectid.OpaqueID{}.Validate(), subjectid.ErrOpaqueEmpty))
	// Output: true
}

// ExamplePhoneNumberID_Validate shows the E.164 shell check —
// "+" followed by 4 to 15 ASCII digits, with no separators. The
// library does NOT run full per-country plan validation; install
// a [WithPhoneNumberValidator]-style hook (see godoc) for that.
func ExamplePhoneNumberID_Validate() {
	fmt.Println(subjectid.PhoneNumberID{PhoneNumber: "+12065550100"}.Validate())
	// Output: <nil>
}

// ExampleDIDID_Validate shows the W3C DID Core URL shell check.
// Per-method validation (did:web, did:key, did:plc, …) is out
// of scope — the library validates the universal shape only.
func ExampleDIDID_Validate() {
	fmt.Println(subjectid.DIDID{URL: "did:example:123456"}.Validate())
	// Output: <nil>
}

// ExampleURIID_Validate shows the RFC 3986 absolute-URI check.
// "uri" is the fallback format for identifiers without a more
// specific format in the IANA registry.
func ExampleURIID_Validate() {
	id := subjectid.URIID{URI: "urn:oasis:names:tc:saml:2.0:nameid-format:transient"}
	fmt.Println(id.Validate())
	// Output: <nil>
}

// ExampleAliasesID_Validate shows the composite rule from
// RFC 9493 §3.2.8: the identifiers slice must be non-empty AND
// no element may itself be an aliases identifier. Inner-element
// errors join into the returned error so [errors.Is] surfaces
// any inner-format sentinel.
func ExampleAliasesID_Validate() {
	id := subjectid.AliasesID{
		Identifiers: []subjectid.SubjectIdentifier{
			subjectid.EmailID{Email: "user@example.com"},
			subjectid.AccountID{URI: "acct:user@example.com"},
		},
	}
	fmt.Println(id.Validate())
	// Output: <nil>
}

// ExampleAliasesID_Validate_nestedRejected shows the no-nested-
// aliases rule firing: the outer aliases parses and validates
// individually-valid inner elements, but a nested AliasesID is
// rejected with [subjectid.ErrNestedAliases].
func ExampleAliasesID_Validate_nestedRejected() {
	id := subjectid.AliasesID{
		Identifiers: []subjectid.SubjectIdentifier{
			subjectid.AliasesID{Identifiers: []subjectid.SubjectIdentifier{
				subjectid.EmailID{Email: "user@example.com"},
			}},
		},
	}
	fmt.Println(errors.Is(id.Validate(), subjectid.ErrNestedAliases))
	// Output: true
}

// exampleTenantID is a minimal extension type for the
// ExampleRegisterFormat demo. Production consumers would carry
// the same shape: embed [subjectid.Seal], implement Format,
// Validate, MarshalJSON, UnmarshalJSON, then register a
// constructor at consumer init time.
type exampleTenantID struct {
	subjectid.Seal
	Tenant string
}

func (exampleTenantID) Format() string  { return "org.example.exampletenant" }
func (exampleTenantID) Validate() error { return nil }

func (o exampleTenantID) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Format string `json:"format"`
		Tenant string `json:"tenant"`
	}{Format: o.Format(), Tenant: o.Tenant})
}

func (o *exampleTenantID) UnmarshalJSON(data []byte) error {
	var v struct {
		Tenant string `json:"tenant"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	o.Tenant = v.Tenant
	return nil
}

// ExampleRegisterFormat shows the consumer extension path. A
// custom format name is associated with a constructor at init
// time; [subjectid.Parse] then dispatches to it automatically for
// any payload carrying that format discriminator.
func ExampleRegisterFormat() {
	_ = subjectid.RegisterFormat("org.example.exampletenant",
		func() subjectid.SubjectIdentifier { return &exampleTenantID{} })

	wire := []byte(`{"format":"org.example.exampletenant","tenant":"acme-prod"}`)
	id, err := subjectid.Parse(json.RawMessage(wire))
	if err != nil {
		panic(err)
	}

	out, _ := json.Marshal(id)
	fmt.Println(string(out))
	// Output: {"format":"org.example.exampletenant","tenant":"acme-prod"}
}

// ExampleRegisterFormat_reservedName shows the IANA-name guard:
// re-registering one of the eight built-in format names returns
// an error wrapping [subjectid.ErrFormatReserved], so a consumer
// cannot accidentally override a built-in.
func ExampleRegisterFormat_reservedName() {
	err := subjectid.RegisterFormat("email",
		func() subjectid.SubjectIdentifier { return &exampleTenantID{} })
	fmt.Println(errors.Is(err, subjectid.ErrFormatReserved))
	// Output: true
}
