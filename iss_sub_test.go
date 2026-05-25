// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"testing"

	"github.com/hstern/go-subjectid"
)

func TestIssSubIDFormat(t *testing.T) {
	var i subjectid.IssSubID
	if got, want := i.Format(), "iss_sub"; got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
	if err := i.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil (no rules yet)", err)
	}
}
