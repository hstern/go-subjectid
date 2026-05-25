// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"testing"

	"github.com/hstern/go-subjectid"
)

func TestDIDIDFormat(t *testing.T) {
	var d subjectid.DIDID
	if got, want := d.Format(), "did"; got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
	if err := d.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil (no rules yet)", err)
	}
}
