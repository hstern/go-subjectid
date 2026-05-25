// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"encoding/json"
	"testing"

	"github.com/hstern/go-subjectid"
)

func TestUnknownFormatFormatReturnsFormatName(t *testing.T) {
	u := subjectid.UnknownFormat{
		FormatName: "org.example.novel",
		Raw:        json.RawMessage(`{"format":"org.example.novel","x":1}`),
	}
	if got, want := u.Format(), "org.example.novel"; got != want {
		t.Errorf("Format() = %q, want %q (UnknownFormat is per-value, not fixed)", got, want)
	}
	if err := u.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil (UnknownFormat skips validation by design)", err)
	}
}

func TestUnknownFormatZeroValueFormatIsEmpty(t *testing.T) {
	var u subjectid.UnknownFormat
	if got := u.Format(); got != "" {
		t.Errorf("zero-value Format() = %q, want %q", got, "")
	}
}
