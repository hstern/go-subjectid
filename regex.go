// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import "regexp"

// mustCompileAnchored wraps pap-generated patterns in ^(?:...)$ so
// MatchString demands a full-string match. pap emits unanchored
// patterns; without the anchors it would match a valid prefix of an
// invalid input (e.g. "did:example:abc%" would pass because
// "did:example:abc" is a valid prefix).
func mustCompileAnchored(pattern string) *regexp.Regexp {
	return regexp.MustCompile("^(?:" + pattern + ")$")
}
