// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/hstern/go-subjectid"
)

// TestFormatErrYieldsSentinelComparableViaErrorsIs exercises the
// consumer-facing contract documented on [subjectid.FormatErr]: the
// returned error compares equal under errors.Is to itself (so a
// package-level var can act as the branch target callers test
// against), wraps cleanly through fmt.Errorf with %w, and carries a
// message that names the format.
func TestFormatErrYieldsSentinelComparableViaErrorsIs(t *testing.T) {
	// Mint a sentinel exactly the way a consumer registering a
	// custom format would.
	errFormatCustom := subjectid.FormatErr("my-custom format")

	// errors.Is matches the sentinel against itself.
	if !errors.Is(errFormatCustom, errFormatCustom) {
		t.Errorf("errors.Is(sentinel, sentinel) = false, want true")
	}

	// Wrapping with %w preserves identity through errors.Is.
	wrapped := errors.New("validation context")
	wrapped = errorsJoin(wrapped, errFormatCustom)
	if !errors.Is(wrapped, errFormatCustom) {
		t.Errorf("errors.Is on wrapped chain = false, want true")
	}

	// Two FormatErr calls with the same name yield distinct
	// values that do NOT compare equal — exactly the property the
	// doc comment promises, so consumers know to store the
	// sentinel in a package-level var rather than re-minting it
	// per call site.
	a := subjectid.FormatErr("dup")
	b := subjectid.FormatErr("dup")
	if errors.Is(a, b) {
		t.Errorf("errors.Is(distinct FormatErr calls) = true, want false")
	}

	// Message names the format.
	if got := errFormatCustom.Error(); !strings.Contains(got, "my-custom format") {
		t.Errorf("Error() = %q; want it to name the format", got)
	}
}

// TestErrRequiredConsumerUsage shows the two idiomatic ways a
// consumer recovers information about a required-field failure:
// categorical (errors.Is against the zero value) and structural
// (errors.As to recover the Fields slice).
func TestErrRequiredConsumerUsage(t *testing.T) {
	// Drive a validator that should return ErrRequired.
	err := subjectid.IssSubID{Iss: "", Sub: ""}.Validate()
	if err == nil {
		t.Fatal("Validate(empty iss, empty sub): got nil, want error")
	}

	// Categorical match against the zero value: works without
	// caring which fields are missing.
	if !errors.Is(err, subjectid.ErrRequired{}) {
		t.Errorf("errors.Is(err, ErrRequired{}) = false, want true")
	}

	// Structural extraction: recover the field names.
	var req subjectid.ErrRequired
	if !errors.As(err, &req) {
		t.Fatalf("errors.As did not extract ErrRequired from %v", err)
	}
	if len(req.Fields) != 2 || req.Fields[0] != "iss" || req.Fields[1] != "sub" {
		t.Errorf("req.Fields = %v, want [iss sub]", req.Fields)
	}
}

// errorsJoin wraps two errors so the chain carries both, using
// errors.Join from the stdlib. Inlined as a tiny helper so the test
// body reads top-to-bottom without an import alias.
func errorsJoin(a, b error) error {
	return errors.Join(a, b)
}

// TestErrUmbrellaMatchesEveryPackageSentinel verifies the contract
// on the top-level [subjectid.Err]: every error this package can
// return — sentinel or struct — matches errors.Is(err, subjectid.Err)
// so a consumer with mixed error sources can branch on "subjectid
// problem at all?" without naming the leaf sentinels.
func TestErrUmbrellaMatchesEveryPackageSentinel(t *testing.T) {
	cases := []struct {
		name string
		err  error
	}{
		{"Err itself", subjectid.Err},
		{"ErrJSON", subjectid.ErrJSON},
		{"ErrNestedAliases", subjectid.ErrNestedAliases},
		{"ErrAliasesEmpty", subjectid.ErrAliasesEmpty},
		{"ErrOpaqueEmpty", subjectid.ErrOpaqueEmpty},
		{"ErrFormatAccount", subjectid.ErrFormatAccount},
		{"ErrFormatEmail", subjectid.ErrFormatEmail},
		{"ErrFormatIssSub", subjectid.ErrFormatIssSub},
		{"ErrFormatPhoneNumber", subjectid.ErrFormatPhoneNumber},
		{"ErrFormatDID", subjectid.ErrFormatDID},
		{"ErrFormatURI", subjectid.ErrFormatURI},
		{"ErrFormatReserved", subjectid.ErrFormatReserved},
		{"ErrRequired{} (zero value)", subjectid.ErrRequired{}},
		{"MissingFields(...) (typed)", subjectid.MissingFields("iss", "sub")},
		{"FormatErr (consumer-built)", subjectid.FormatErr("invalid consumer format")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if !errors.Is(tc.err, subjectid.Err) {
				t.Errorf("errors.Is(%v, Err) = false, want true", tc.err)
			}
		})
	}

	// Sanity: a non-subjectid error must NOT match.
	other := errors.New("some other package: oops")
	if errors.Is(other, subjectid.Err) {
		t.Errorf("errors.Is(unrelated, Err) = true, want false")
	}
}
