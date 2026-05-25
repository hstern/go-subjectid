// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

// Package subjectid implements RFC 9493 Subject Identifiers for
// Security Event Tokens.
//
// The implementation provides a sealed SubjectIdentifier interface
// plus per-format concrete types for each format in the IANA
// "Security Event Identifier Formats" registry, a JSON codec that
// round-trips byte-stably against every example figure in RFC 9493 §3,
// an opt-in Validate method per format, and an extension mechanism via
// RegisterFormat for formats outside the built-in set.
//
// This file is currently a stub. The type surface, codec, and
// validation rules arrive in subsequent commits.
package subjectid

// SpecVersion identifies the RFC this package implements. RFCs have
// no minor or patch numbers; errata to RFC 9493 are absorbed into
// Go-minor releases of this module without changing the value of
// this constant.
const SpecVersion = "RFC 9493"
