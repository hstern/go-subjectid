# Changelog

All notable changes to `go-subjectid` are documented here. The format
is loosely based on [Keep a Changelog](https://keepachangelog.com/),
and the project adheres to [Semantic Versioning](https://semver.org/).

The library's version is independent of the spec it implements
(RFC 9493). Errata to the RFC are absorbed into Go-minor releases
without bumping the spec-version constant.

## [Unreleased]

## [0.2.0]

### Changed

- **Breaking:** `Parse` now returns Subject Identifier values in their
  value form (e.g. `IssSubID`) rather than pointer form (`*IssSubID`),
  so a value built as a struct literal and a value read back from
  `Parse` share the same dynamic type. Code that type-asserts or
  type-switches on a `Parse` result must match the value forms. The
  `SubjectIdentifier` and `Parse` documentation state this canonical
  form, and extension types registered via `RegisterFormat` must
  satisfy the interface with value receivers.

### Fixed

- `AliasesID.Validate` now rejects a nested `aliases` identifier
  regardless of how it was produced. A nested aliases parsed from JSON
  previously escaped the `ErrNestedAliases` check; a hand-built one was
  already caught.

## [0.1.0]

- Initial release: the eight RFC 9493 Subject Identifier formats, a
  discriminator-driven JSON codec (`Parse` plus spec-order, byte-stable
  `Marshal`), opt-in per-format `Validate` with sentinel errors, a
  forward-compatible `UnknownFormat` carrier that preserves unknown
  wire bytes verbatim, and the `RegisterFormat` extension hook.
