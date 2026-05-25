// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"testing"

	"github.com/hstern/go-subjectid"
)

func TestAccountIDFormat(t *testing.T) {
	var a subjectid.AccountID
	if got, want := a.Format(), "account"; got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
	if err := a.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil (no rules yet)", err)
	}
}
