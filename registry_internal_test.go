// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import "testing"

func TestBuiltinFormatsHasEightEntries(t *testing.T) {
	got := builtinFormats()
	if want := 8; len(got) != want {
		t.Errorf("len(builtinFormats()) = %d, want %d", len(got), want)
	}
}

func TestBuiltinFormatsReturnsFreshMap(t *testing.T) {
	m1 := builtinFormats()
	m2 := builtinFormats()
	delete(m1, "email")
	if _, ok := m2[" email"]; ok {
		t.Fatalf("unreachable")
	}
	if _, ok := m2["email"]; !ok {
		t.Error("mutating one builtinFormats() map affected another — must be fresh per call")
	}
}

func TestBuiltinFormatsReturnsCorrectConcreteTypes(t *testing.T) {
	tests := []struct {
		want SubjectIdentifier
		name string
	}{
		{name: "account", want: &AccountID{}},
		{name: "email", want: &EmailID{}},
		{name: "iss_sub", want: &IssSubID{}},
		{name: "opaque", want: &OpaqueID{}},
		{name: "phone_number", want: &PhoneNumberID{}},
		{name: "did", want: &DIDID{}},
		{name: "uri", want: &URIID{}},
		{name: "aliases", want: &AliasesID{}},
	}
	m := builtinFormats()
	for _, tc := range tests {
		ctor, ok := m[tc.name]
		if !ok {
			t.Errorf("builtinFormats()[%q] missing", tc.name)
			continue
		}
		got := ctor()
		if gotFormat, wantFormat := got.Format(), tc.want.Format(); gotFormat != wantFormat {
			t.Errorf("builtinFormats()[%q]().Format() = %q, want %q", tc.name, gotFormat, wantFormat)
		}
	}
}

func TestLookupReturnsBuiltinAndExtension(t *testing.T) {
	if got := lookup("email"); got == nil {
		t.Error("lookup(\"email\") = nil, want a constructor")
	}
	if got := lookup("subjectid.test.lookup.unregistered"); got != nil {
		t.Errorf("lookup(unregistered) = %T, want nil", got)
	}

	const name = "subjectid.test.lookup.extension"
	if err := RegisterFormat(name, func() SubjectIdentifier { return &UnknownFormat{FormatName: name} }); err != nil {
		t.Fatalf("RegisterFormat(%q): err = %v, want nil", name, err)
	}
	ctor := lookup(name)
	if ctor == nil {
		t.Fatalf("lookup(%q) = nil after RegisterFormat, want a constructor", name)
	}
	got := ctor()
	if got == nil {
		t.Fatal("ctor() = nil")
	}
	if gotName := got.Format(); gotName != name {
		t.Errorf("ctor().Format() = %q, want %q", gotName, name)
	}
}
