# Contributor conventions

This file orients new contributors — human or otherwise — to the
conventions this repository expects. The `README.md` answers "what is
this library and how do I use it"; this file answers "how do I work
on it".

## Goal

A Go reference implementation of
[RFC 9493 — Subject Identifiers for Security Event Tokens][rfc9493].
Wire fidelity over ergonomic shortcuts; byte-stable round-trip on
every spec example; zero non-test runtime dependencies.

[rfc9493]: https://www.rfc-editor.org/rfc/rfc9493.html

## Quick checks before pushing

```sh
gofmt -l .                       # must be empty
go vet ./...
go mod tidy -diff                # no drift
go test -race -shuffle=on ./...
golangci-lint run ./...
```

All five are run in CI. Failing any of them locally is a hard stop.

## Code style

- **`gofmt` is canonical**, `gofumpt` is acceptable for new files.
- **Linter config is in `.golangci.yml`.** Don't introduce `nolint`
  pragmas without a one-line comment naming the rule and the reason.
- **Naming follows Go conventions** — `MixedCaps` for exported,
  `mixedCaps` for unexported, short receiver names, no ALL_CAPS
  constants, no `Get`-prefix on getters.
- **Errors**: sentinel `ErrFoo` for stable comparison targets,
  custom error types when the consumer needs more than `Is`/`As`,
  `fmt.Errorf(... %w ...)` to wrap.
- **No `panic` outside `init()` and clearly-bug-only invariants.**

## Per-file SPDX header

Every `.go` file (including tests) starts with exactly two lines
before the `package` declaration:

```go
// Copyright 2026 The go-subjectid Authors
// SPDX-License-Identifier: Apache-2.0

package subjectid
```

The full Apache-2.0 license text lives in `LICENSE`. Per-file SPDX
tags are the machine-readable reference for license-detection tools.
Don't add long per-file preambles, `// originally written by` lines,
or `// authored-by:` comments — authorship lives in git history.

## Commit messages

Detailed by default. Short single-line commits are a code smell.

- **Title**: imperative, ≤72 chars, no trailing period, lowercase
  first letter (unless a proper noun or literal identifier).
- **Body**: prose paragraph(s) explaining *why* the change exists —
  the spec clause honored, the bug observed, the constraint
  satisfied. Wrap at ~72 chars.
- **References**: spec section numbers (§3.2.1), public RFC numbers
  (RFC 5322), public commit SHAs. Anything verifiable by a
  stranger.
- **No conventional-commits prefixes** (`feat:`, `fix:`, `chore:`).
  The body is where the explanation goes.

Exceptions where a one-line title is acceptable: typo / whitespace
fixes, comment-only changes, dependency bumps with no API impact,
the initial scaffold commit.

Trailers go at the bottom, separated by a blank line. If a tool
produced the commit, include the appropriate co-author trailer.

## Pull requests

- **Branch from current `origin/main`.** Branch protection blocks
  direct push to `main` once it's enabled.
- **One topic per PR.** A PR that does two unrelated things gets
  asked to split.
- **CI must be green before merge.** The `static`, `test`, and
  `lint` jobs all run on every PR and are required checks.
- **Squash on merge** is the default; the squash commit's message
  follows the same shape as a regular commit message.

## Tests

- **Table-driven tests by default.** The stdlib `testing` package is
  the standard; `testify` is not currently a dependency, and adding
  one needs a discussion.
- **`-race` and `-shuffle=on` are always on.** Tests must pass
  deterministically under both.
- **Fixtures live alongside their tests** unless they're spec-
  example payloads, in which case they go under
  `internal/specfixtures/`.
- **Example functions** (`func Example…`) double as godoc rendering
  and executable verification. Add one for each public type as the
  API stabilizes.

## Dependencies

- **Runtime: standard library only** through v0.x. Adding a runtime
  dependency needs a justification in the PR description.
- **Tests: standard library only by default.** Test-only deps are
  more flexible but still need a one-line justification.
- **`go.mod`**: keep the `module` path stable at
  `github.com/hstern/go-subjectid` (no `/vN` suffix for v0.x/v1.x —
  Go SemVer rule). Major-version bumps follow the `go-jose` branch
  pattern.

## Spec fidelity

The library exists to round-trip RFC 9493 Subject Identifier values
exactly as the spec describes them. When the spec and Go idiom
conflict, the spec wins:

- **`format` member always first** in marshal output. RFC 9493
  doesn't mandate ordering but every published example puts it
  first, and the library matches byte-for-byte.
- **Aliases MUST NOT be nested.** RFC 9493 §3.2.8. Enforced in
  `Validate()`, not at the Go type level.
- **Unknown formats parse into `UnknownFormat`**, not an error.
  Forward-compatibility.
- **Extra fields are silently dropped on unmarshal, rejected on
  marshal.** Postel's law.
- **Required members are non-null and non-empty.** A value whose
  zero state passes `Validate()` is a bug.
