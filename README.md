# go-subjectid

A Go library implementing
[RFC 9493 — Subject Identifiers for Security Event Tokens][rfc9493].

[rfc9493]: https://www.rfc-editor.org/rfc/rfc9493.html

> **Status: pre-release.** API surface is being shaped against RFC 9493
> §3 with the goal of byte-stable round-trip on every spec example.
> First tag will be `v0.1.0`. The module path
> `github.com/hstern/go-subjectid` is stable from the first commit; do
> not depend on the repo until `v0.1.0` is tagged.

## What it is

`go-subjectid` is the Go ecosystem's reference implementation of the
RFC 9493 Subject Identifier types — the discriminated-union JSON
values used by SET / SSF / CAEP / RISC to identify the subject of a
security event.

The library handles:

- **The eight built-in formats** from the IANA "Security Event
  Identifier Formats" registry: `account`, `email`, `iss_sub`,
  `opaque`, `phone_number`, `did`, `uri`, `aliases`.
- **JSON codec** — discriminator-driven `Unmarshal` dispatch plus
  spec-order, byte-stable `Marshal` output.
- **Validation** — opt-in `Validate()` per format, returning sentinel
  errors (`ErrFormatXxx`, `ErrRequired`, `ErrOpaqueEmpty`,
  `ErrNestedAliases`, …). Callers branch with `errors.Is`; every
  sentinel also matches the umbrella `subjectid.Err`. `ErrRequired`
  is a struct type — use `errors.As` to recover the missing field
  names.
- **Forward compatibility** — unknown formats parse into an
  `UnknownFormat` carrier that preserves the wire bytes verbatim, so
  the identifier round-trips even when the library can't fully parse
  it.
- **Extension** — `RegisterFormat` lets downstream code add a new
  format constructor for non-registry types.

## Install

```sh
go get github.com/hstern/go-subjectid@latest
```

Requires Go 1.26 or newer.

## Build prerequisites (contributors only)

Per-format validation regexes are generated from the relevant RFCs'
ABNF grammars in `grammar/<rfc>/*.abnf` and committed to the tree as
`*.rex` files alongside, so the runtime tree depends only on the Go
standard library. Regenerating those files needs the
[pandatix/go-abnf](https://github.com/pandatix/go-abnf) `pap` CLI.

`go install` does not work because the `cmd/pap` go.mod uses a local
`replace` directive; build from source:

```sh
git clone --depth=1 https://github.com/pandatix/go-abnf.git /tmp/go-abnf
go build -C /tmp/go-abnf/cmd/pap -o "$(go env GOPATH)/bin/pap"
```

Then regenerate with `make generate` (or `make -B generate` to force a
full rebuild). CI runs the same incantation in the `pap` job below and
fails the build on uncommitted drift, so always commit the `.rex`
files alongside any `.abnf` change.

## Quickstart

_To be filled in once Phase 3 (codec) lands._

## Design

The library's design rationale, especially the wire-fidelity choices
that bite every RFC 9493 implementer once, will be summarized in
`design.md` ahead of the `v0.1.0` tag. The headline decisions:

- Sealed `SubjectIdentifier` interface with per-format concrete
  types — no `map[string]any`, no zero-valued union-fields struct.
- `encoding/json` stdlib with custom `UnmarshalJSON` dispatch on the
  `format` discriminator.
- **Lenient on unmarshal, strict on marshal.** Postel's law: decode
  whatever the wire gave us, validate at the marshal boundary.
- Byte-stable output: `format` first, format-specific members in the
  order RFC 9493 §3 defines them.
- Open-extension fields are `json.RawMessage`, not `map[string]any`
  — interop scenarios pin exact JSON bytes, and `map` reorders keys.

## Compatibility

- **Spec version**: `const SpecVersion = "RFC 9493"`. RFCs don't
  have minor/patch numbers; errata are absorbed into Go-minor
  releases without changing the constant.
- **Go version**: 1.26+.
- **Dependencies**: standard library only, at runtime. Test
  dependencies, if any, are listed in `go.mod`.
- **Library SemVer** is independent of the spec version. Major-
  version handling follows the `go-jose` branch pattern (no
  versioned subdirectories — `vN` lives in `go.mod` on a `vN`
  branch).

## Contributing

Contributor conventions are in [`AGENTS.md`](AGENTS.md): commit
message style, code review expectations, the per-file SPDX header,
and the local pre-PR checks the CI also runs.

Bugs and feature ideas are welcome via the project's issue tracker.

## License

Apache-2.0. See [`LICENSE`](LICENSE) for the full text.

Every source file carries an SPDX identifier:

```go
// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0
```
