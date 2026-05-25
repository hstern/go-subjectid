// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/hstern/go-subjectid"
)

// TestStrictMarshalDefaultsOff locks in the Postel default: an
// invalid identifier marshals fine when the toggle is off, so
// callers that round-trip-and-mutate aren't punished for momentarily
// invalid intermediate states.
func TestStrictMarshalDefaultsOff(t *testing.T) {
	defer subjectid.StrictMarshal(subjectid.StrictMarshal(false))

	// Empty email is invalid per RFC 9493 §3 but must marshal
	// while the toggle is off.
	if _, err := json.Marshal(subjectid.EmailID{Email: ""}); err != nil {
		t.Errorf("default-off Marshal(invalid): err = %v, want nil", err)
	}
}

// TestStrictMarshalEnablesValidation locks in the documented
// behavior: with the toggle on, MarshalJSON runs Validate first and
// surfaces the validation error.
func TestStrictMarshalEnablesValidation(t *testing.T) {
	defer subjectid.StrictMarshal(subjectid.StrictMarshal(true))

	// Each format's invalid value should fail marshal under
	// StrictMarshal. We exercise one per format to confirm every
	// MarshalJSON honors the toggle.
	cases := []struct {
		name string
		val  subjectid.SubjectIdentifier
		want error
	}{
		{"account empty", subjectid.AccountID{URI: ""}, subjectid.ErrRequired{}},
		{"account bad", subjectid.AccountID{URI: "not-acct"}, subjectid.ErrFormatAccount},
		{"email empty", subjectid.EmailID{Email: ""}, subjectid.ErrRequired{}},
		{"email bad", subjectid.EmailID{Email: "not-email"}, subjectid.ErrFormatEmail},
		{"iss_sub empty", subjectid.IssSubID{}, subjectid.ErrRequired{}},
		{"iss_sub bad", subjectid.IssSubID{Iss: "not-uri", Sub: "x"}, subjectid.ErrFormatIssSub},
		{"opaque empty", subjectid.OpaqueID{ID: ""}, subjectid.ErrOpaqueEmpty},
		{"phone empty", subjectid.PhoneNumberID{PhoneNumber: ""}, subjectid.ErrRequired{}},
		{"phone bad", subjectid.PhoneNumberID{PhoneNumber: "not-phone"}, subjectid.ErrFormatPhoneNumber},
		{"did empty", subjectid.DIDID{URL: ""}, subjectid.ErrRequired{}},
		{"did bad", subjectid.DIDID{URL: "not-did"}, subjectid.ErrFormatDID},
		{"uri empty", subjectid.URIID{URI: ""}, subjectid.ErrRequired{}},
		{"uri bad", subjectid.URIID{URI: "not-uri"}, subjectid.ErrFormatURI},
		{"aliases empty", subjectid.AliasesID{}, subjectid.ErrRequired{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := json.Marshal(tc.val)
			if err == nil {
				t.Fatalf("strict Marshal(%s): err = nil, want non-nil", tc.name)
			}
			if !errors.Is(err, tc.want) {
				t.Errorf("strict Marshal(%s): errors.Is(err, %v) = false, want true (err = %v)",
					tc.name, tc.want, err)
			}
		})
	}
}

// TestStrictMarshalValidValueMarshals locks in that StrictMarshal
// does not break the happy path: a valid identifier still marshals
// to its spec-order bytes under the toggle.
func TestStrictMarshalValidValueMarshals(t *testing.T) {
	defer subjectid.StrictMarshal(subjectid.StrictMarshal(true))

	got, err := json.Marshal(subjectid.EmailID{Email: "user@example.com"})
	if err != nil {
		t.Fatalf("strict Marshal(valid): err = %v, want nil", err)
	}
	want := `{"format":"email","email":"user@example.com"}`
	if string(got) != want {
		t.Errorf("strict Marshal(valid): got %s, want %s", got, want)
	}
}

// TestStrictMarshalReturnsPreviousValue pins the swap-semantics
// documented on [StrictMarshal] — callers can save/restore via
//
//	defer StrictMarshal(StrictMarshal(true))
//
// which only works if the function returns the prior value.
func TestStrictMarshalReturnsPreviousValue(t *testing.T) {
	defer subjectid.StrictMarshal(subjectid.StrictMarshal(false))

	subjectid.StrictMarshal(false)
	if prev := subjectid.StrictMarshal(true); prev {
		t.Errorf("StrictMarshal(true) when off: returned %v, want false", prev)
	}
	if prev := subjectid.StrictMarshal(true); !prev {
		t.Errorf("StrictMarshal(true) when on: returned %v, want true", prev)
	}
	if prev := subjectid.StrictMarshal(false); !prev {
		t.Errorf("StrictMarshal(false) when on: returned %v, want true", prev)
	}
}

// TestStrictMarshalNestedAliasesFails locks in that aliases'
// composite Validate (which recursively checks each inner element)
// is honored at marshal time too.
func TestStrictMarshalNestedAliasesFails(t *testing.T) {
	defer subjectid.StrictMarshal(subjectid.StrictMarshal(true))

	id := subjectid.AliasesID{
		Identifiers: []subjectid.SubjectIdentifier{
			subjectid.EmailID{Email: "user@example.com"},
			subjectid.AliasesID{Identifiers: []subjectid.SubjectIdentifier{
				subjectid.EmailID{Email: "x@y.z"},
			}},
		},
	}
	_, err := json.Marshal(id)
	if err == nil {
		t.Fatal("strict Marshal(nested aliases): err = nil, want non-nil")
	}
	if !errors.Is(err, subjectid.ErrNestedAliases) {
		t.Errorf("errors.Is(err, ErrNestedAliases) = false, want true (err = %v)", err)
	}
}

// TestWithEmailValidatorReplacesBuiltin locks in the per-format
// hook contract documented on every With*Validator: when set, the
// custom fn REPLACES the built-in regex; the required-non-empty
// rule still runs first; nil restores the built-in; the setter
// returns the previous validator for save/restore via defer.
func TestWithEmailValidatorReplacesBuiltin(t *testing.T) {
	// Save and restore — every other test must see the built-in.
	defer subjectid.WithEmailValidator(subjectid.WithEmailValidator(nil))

	// Sentinel returned by the custom validator.
	errCustom := errors.New("custom email validator failed")
	subjectid.WithEmailValidator(func(e subjectid.EmailID) error {
		if e.Email == "blocked@example.com" {
			return errCustom
		}
		return nil
	})

	// Custom fn is consulted: a value the built-in regex would
	// REJECT now passes, and a value the built-in would accept
	// fails because the custom fn rejects it.
	cases := []struct {
		name  string
		email string
		want  error
	}{
		{"empty still fails required", "", subjectid.ErrRequired{}},
		{"custom accepts what built-in would reject", "not a valid addr-spec", nil},
		{"custom rejects what built-in would accept", "blocked@example.com", errCustom},
		{"custom accepts what built-in accepts", "user@example.com", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := subjectid.EmailID{Email: tc.email}.Validate()
			if tc.want == nil {
				if err != nil {
					t.Errorf("Validate(%q) = %v, want nil", tc.email, err)
				}
				return
			}
			if !errors.Is(err, tc.want) {
				t.Errorf("Validate(%q): errors.Is(err, %v) = false, want true (err = %v)",
					tc.email, tc.want, err)
			}
		})
	}
}

// TestWithEmailValidatorNilRestoresBuiltin locks in the "pass nil
// to restore the built-in" contract.
func TestWithEmailValidatorNilRestoresBuiltin(t *testing.T) {
	defer subjectid.WithEmailValidator(subjectid.WithEmailValidator(nil))

	subjectid.WithEmailValidator(func(subjectid.EmailID) error {
		return errors.New("custom")
	})
	// Custom is in effect now.
	if err := (subjectid.EmailID{Email: "user@example.com"}).Validate(); err == nil {
		t.Fatal("expected custom validator to fire, got nil")
	}
	// Restore.
	subjectid.WithEmailValidator(nil)
	if err := (subjectid.EmailID{Email: "user@example.com"}).Validate(); err != nil {
		t.Errorf("after WithEmailValidator(nil): Validate(valid) = %v, want nil", err)
	}
}

// TestWithEmailValidatorReturnsPrevious pins the save-and-restore
// idiom documented on the setter.
func TestWithEmailValidatorReturnsPrevious(t *testing.T) {
	defer subjectid.WithEmailValidator(subjectid.WithEmailValidator(nil))

	fn1 := func(subjectid.EmailID) error { return errors.New("fn1") }
	fn2 := func(subjectid.EmailID) error { return errors.New("fn2") }

	// Start with built-in.
	if prev := subjectid.WithEmailValidator(fn1); prev != nil {
		t.Errorf("first WithEmailValidator: prev = %p, want nil", prev)
	}
	// fn1 → fn2: prev should return fn1's value (we can't compare
	// function identity portably across Go releases, but invoking
	// the returned fn must observe fn1's behavior).
	prev := subjectid.WithEmailValidator(fn2)
	if prev == nil {
		t.Fatal("second WithEmailValidator: prev = nil, want fn1")
	}
	if got := prev(subjectid.EmailID{}); got == nil || got.Error() != "fn1" {
		t.Errorf("prev validator: got %v, want fn1", got)
	}
}

// TestAllPerFormatValidatorHooksFire smokes the remaining five
// per-format hooks (account, iss_sub, phone_number, did, uri) with
// the same shape as the email-specific tests above: install fn,
// confirm it replaces the built-in, restore via the returned
// previous value.
//
// Each row uses a value the built-in WOULD reject and a custom
// validator that ACCEPTS it; observing nil from Validate confirms
// the custom validator is in effect and the built-in is bypassed.
func TestAllPerFormatValidatorHooksFire(t *testing.T) {
	t.Run("account", func(t *testing.T) {
		defer subjectid.WithAccountValidator(subjectid.WithAccountValidator(nil))
		subjectid.WithAccountValidator(func(subjectid.AccountID) error { return nil })
		if err := (subjectid.AccountID{URI: "not-acct"}.Validate()); err != nil {
			t.Errorf("custom account validator should bypass built-in: err = %v", err)
		}
	})
	t.Run("iss_sub", func(t *testing.T) {
		defer subjectid.WithIssSubValidator(subjectid.WithIssSubValidator(nil))
		subjectid.WithIssSubValidator(func(subjectid.IssSubID) error { return nil })
		if err := (subjectid.IssSubID{Iss: "not-uri", Sub: "x"}.Validate()); err != nil {
			t.Errorf("custom iss_sub validator should bypass built-in: err = %v", err)
		}
	})
	t.Run("phone_number", func(t *testing.T) {
		defer subjectid.WithPhoneNumberValidator(subjectid.WithPhoneNumberValidator(nil))
		subjectid.WithPhoneNumberValidator(func(subjectid.PhoneNumberID) error { return nil })
		if err := (subjectid.PhoneNumberID{PhoneNumber: "not-phone"}.Validate()); err != nil {
			t.Errorf("custom phone_number validator should bypass built-in: err = %v", err)
		}
	})
	t.Run("did", func(t *testing.T) {
		defer subjectid.WithDIDValidator(subjectid.WithDIDValidator(nil))
		subjectid.WithDIDValidator(func(subjectid.DIDID) error { return nil })
		if err := (subjectid.DIDID{URL: "not-did"}.Validate()); err != nil {
			t.Errorf("custom did validator should bypass built-in: err = %v", err)
		}
	})
	t.Run("uri", func(t *testing.T) {
		defer subjectid.WithURIValidator(subjectid.WithURIValidator(nil))
		subjectid.WithURIValidator(func(subjectid.URIID) error { return nil })
		if err := (subjectid.URIID{URI: "not-uri"}.Validate()); err != nil {
			t.Errorf("custom uri validator should bypass built-in: err = %v", err)
		}
	})
}

// TestPerFormatHooksRequiredCheckStillRuns locks in the contract
// that the RFC 9493 §3 required-non-empty rule runs BEFORE the
// custom validator, so a custom fn never sees an empty value (or
// for iss_sub, never sees either empty member).
func TestPerFormatHooksRequiredCheckStillRuns(t *testing.T) {
	defer subjectid.WithEmailValidator(subjectid.WithEmailValidator(nil))

	called := false
	subjectid.WithEmailValidator(func(subjectid.EmailID) error {
		called = true
		return nil
	})
	err := subjectid.EmailID{Email: ""}.Validate()
	if !errors.Is(err, subjectid.ErrRequired{}) {
		t.Errorf("empty-check should return ErrRequired before consulting custom; got %v", err)
	}
	if called {
		t.Error("custom validator was called for an empty value; required-check must run first")
	}
}
