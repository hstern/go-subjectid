// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid_test

import (
	"testing"

	"github.com/hstern/go-subjectid"
)

func TestPhoneNumberIDFormat(t *testing.T) {
	var p subjectid.PhoneNumberID
	if got, want := p.Format(), "phone_number"; got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
	if err := p.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil (no rules yet)", err)
	}
}
