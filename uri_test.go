// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"testing"

	"github.com/hstern/go-subjectid"
)

func TestURIIDFormat(t *testing.T) {
	var u subjectid.URIID
	if got, want := u.Format(), "uri"; got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
	if err := u.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil (no rules yet)", err)
	}
}
