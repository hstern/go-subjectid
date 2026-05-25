// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"testing"

	"github.com/hstern/go-subjectid"
)

func TestOpaqueIDFormat(t *testing.T) {
	var o subjectid.OpaqueID
	if got, want := o.Format(), "opaque"; got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
	if err := o.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil (no rules yet)", err)
	}
}
