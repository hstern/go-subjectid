// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid

import (
	"fmt"
	"sync"
)

// Constructor is the per-format factory the registry holds. Each
// Constructor returns a freshly-allocated zero-valued
// [SubjectIdentifier] of the matching format type; the codec
// invokes UnmarshalJSON on the returned value to fill it in.
//
// Constructors registered for non-built-in formats must return a
// type that embeds [Seal] so it satisfies the sealed interface.
//
// The returned value is conventionally a pointer so the codec can
// populate it via UnmarshalJSON, but [Parse] normalizes the result
// to its canonical value form before returning it. The extension
// type must therefore satisfy [SubjectIdentifier] with value
// receivers, so the dereferenced value still implements the
// interface; a type whose methods use pointer receivers would not
// survive that normalization.
type Constructor func() SubjectIdentifier

// formatRegistry is the package-global dispatch table. It is
// initialized from [builtinFormats] at package load and amended
// only via [RegisterFormat]. Concurrent reads use the RLock; the
// rare write path is [RegisterFormat], typically called once per
// extension format at consumer init time.
//
// Field order — map first, mutex second — is chosen so the GC
// scans only the 8 pointer bytes of the map header rather than
// the full struct: sync.RWMutex contains no pointers, so placing
// it after the map lets the GC stop scanning earlier.
var formatRegistry = struct {
	m  map[string]Constructor
	mu sync.RWMutex
}{m: builtinFormats()}

// builtinFormats returns a freshly-allocated map populated with
// the eight formats from the IANA "Security Event Identifier
// Formats" registry as it stood at the time of release. A new
// map is returned per call so callers cannot mutate the
// package's shared registry by holding a reference to the
// initial value.
func builtinFormats() map[string]Constructor {
	return map[string]Constructor{
		"account":      func() SubjectIdentifier { return &AccountID{} },
		"email":        func() SubjectIdentifier { return &EmailID{} },
		"iss_sub":      func() SubjectIdentifier { return &IssSubID{} },
		"opaque":       func() SubjectIdentifier { return &OpaqueID{} },
		"phone_number": func() SubjectIdentifier { return &PhoneNumberID{} },
		"did":          func() SubjectIdentifier { return &DIDID{} },
		"uri":          func() SubjectIdentifier { return &URIID{} },
		"aliases":      func() SubjectIdentifier { return &AliasesID{} },
	}
}

// builtinFormatNames is the set of format names [builtinFormats]
// populates. [RegisterFormat] checks against it to detect
// built-in collisions without needing to know whether a registry
// entry was placed by builtinFormats or by an earlier
// RegisterFormat call.
var builtinFormatNames = map[string]struct{}{
	"account":      {},
	"email":        {},
	"iss_sub":      {},
	"opaque":       {},
	"phone_number": {},
	"did":          {},
	"uri":          {},
	"aliases":      {},
}

// RegisterFormat registers a [Constructor] for a Subject
// Identifier format outside the IANA built-in set.
//
// Call once per extension format, typically at consumer init
// time. Re-registering an already-registered extension format
// silently replaces the prior constructor — a single consumer
// init owns each extension format by convention. Concurrent
// callers are serialized by an internal mutex.
//
// Returns an error wrapping [ErrFormatReserved] if name matches
// one of the eight built-in formats: those cannot be overridden;
// consumers needing different per-built-in behavior should wrap
// the concrete type rather than re-register the format. Compare
// with [errors.Is].
func RegisterFormat(name string, ctor Constructor) error {
	if _, builtin := builtinFormatNames[name]; builtin {
		return fmt.Errorf("%w: %q", ErrFormatReserved, name)
	}
	formatRegistry.mu.Lock()
	defer formatRegistry.mu.Unlock()
	formatRegistry.m[name] = ctor
	return nil
}

// lookup returns the constructor for the named format, or nil
// if no constructor is registered. The codec (a later commit)
// uses lookup to dispatch; a nil return causes the decoder to
// fall back to [UnknownFormat].
func lookup(name string) Constructor {
	formatRegistry.mu.RLock()
	defer formatRegistry.mu.RUnlock()
	return formatRegistry.m[name]
}
