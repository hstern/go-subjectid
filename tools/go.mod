// Separate module so the runtime library's go.mod stays stdlib-only.
// Tools here run via `make gen-grammars` (which cd's into this directory
// to use this module); the generated output is checked in.
module github.com/hstern/go-subjectid/tools

go 1.26

require github.com/pandatix/go-abnf v0.4.3

require github.com/hashicorp/go-uuid v1.0.3 // indirect

// Pinned to the fork while
//   https://github.com/pandatix/go-abnf/pull/214
// is in review. The patch teaches (Repetition).regex to emit
// Go-RE2-compatible quantifiers — fixes `{,m}` -> `{0,m}` and
// adds the missing case for the general bounded form m*nRULE
// (which silently dropped the quantifier in v0.4.3). Drop this
// replace once the upstream release with the fix is tagged and
// bump the require above to it.
replace github.com/pandatix/go-abnf => github.com/hstern/go-abnf v0.0.0-20260525091812-a38271311938
