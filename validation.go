// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import "sync/atomic"

// strictMarshal is the process-wide toggle consulted by every
// built-in MarshalJSON method. See [StrictMarshal].
var strictMarshal atomic.Bool

// StrictMarshal sets the process-wide marshal-validation toggle and
// returns the previous setting. When true, every MarshalJSON method
// on a built-in [SubjectIdentifier] type runs [SubjectIdentifier.Validate]
// first and propagates any failure as the returned error.
//
// Default is false — Postel's law on outbound: the producer knows
// what it is emitting, and a strict-marshal default would punish
// the round-trip-and-mutate workflow common in test harnesses.
// Enable it for end-to-end test harnesses where emitting invalid
// wire output is itself a bug worth surfacing, or for paranoid
// production paths that want belt-and-suspenders enforcement.
//
// The toggle is process-global. Library consumers embedded in
// larger applications should generally NOT enable it from library
// code — the setting will affect every Marshal call elsewhere in
// the process. Prefer an explicit
//
//	if err := v.Validate(); err != nil { ... }
//	json.Marshal(v)
//
// at the producer boundary. Tests that DO toggle it must restore
// the prior value, typically via defer:
//
//	defer subjectid.StrictMarshal(subjectid.StrictMarshal(true))
//
// Safe for concurrent use; the underlying value is an
// [atomic.Bool].
func StrictMarshal(enabled bool) bool {
	return strictMarshal.Swap(enabled)
}

// maybeValidate is the marshal-side hook: when [StrictMarshal] is
// enabled, run v.Validate and surface any failure; otherwise emit
// nil so the caller falls through to the unconditional encode path.
// Inlined into every built-in MarshalJSON.
func maybeValidate(v SubjectIdentifier) error {
	if strictMarshal.Load() {
		return v.Validate()
	}
	return nil
}
