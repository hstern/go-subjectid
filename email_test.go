// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"testing"

	"github.com/hstern/go-subjectid"
)

func TestEmailIDFormat(t *testing.T) {
	var e subjectid.EmailID
	if got, want := e.Format(), "email"; got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
	if err := e.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil (no rules yet)", err)
	}
}
