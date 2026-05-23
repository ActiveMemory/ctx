# i18n.Fold Helper + AST Ban on strings.ToLower

`strings.ToLower` produces incorrect output for non-ASCII
text (Turkish İ→i̇ vs i, German ß→ß vs ss, etc.). Every
case-insensitive comparison in ctx using `strings.ToLower`
on potentially-non-ASCII input is a latent i18n bug. Today's
codebase has 45 production callsites + 4 in tests; some
operate on ASCII-bounded input (URL schemes, Go identifiers)
and are correct by accident, others operate on user-typed
strings (slugs, search queries, classifications) and are
real bugs waiting for the first non-English user.

This is the "fix tooling first" prerequisite for
`specs/placeholder-i18n.md` (which adds non-ASCII Turkish
placeholders and would activate the latent bug in the
placeholder validator).

## Problem

No AST or compliance test in ctx flags `strings.ToLower`
for case-insensitive matching. Discovery happens only by
observation. The `audit` and `compliance` packages even
use `strings.ToLower` themselves (operating on Go
identifiers — ASCII-safe, but indistinguishable to a
reviewer from the unsafe uses elsewhere).

A per-call annotation scheme (`// case-fold: ascii-only`
comments) would add visual noise everywhere and rot the
moment a contributor forgets. A package-level allowlist
("these packages are grandfathered") legitimizes broken
windows. Neither acceptable.

## Solution

**Make the correct path the only path.** New
`internal/i18n/` package with a single exported function:

```go
// Fold returns the Unicode case-folded form of s,
// suitable for case-insensitive comparison. Replaces
// strings.ToLower in all comparison contexts.
func Fold(s string) string {
    return foldCaser.String(s)
}
```

Backed by `cases.Fold(language.Und)` from
`golang.org/x/text/cases` (already a transitive dependency
at v0.37.0). For ASCII input, output is byte-identical to
`strings.ToLower` — so the swap is behavior-preserving on
every callsite that was already safe; on the unsafe sites
it fixes the latent bug.

**AST ban.** New compliance test
(`internal/compliance/no_strings_tolower_test.go`) walks
every .go file under the project root (skipping
vendor/.git/dist/site/node_modules per the existing
`allGoFiles` helper) and fails on any `strings.ToLower`
call expression. The single allowed callsite is
`internal/i18n/fold.go` itself, which uses the upstream
`cases` API and cannot use its own helper (chicken-egg).
This exception is enforced **structurally**, not via
allowlist: the test checks the package import path, not a
file list, and there is exactly one package
(`internal/i18n`) where the call is legitimate.

**Sweep.** Replace all ~49 existing callsites in a single
commit with the gate. No allowlist, no "grandfathered" —
they are all latent bugs of varying severity, the swap is
safe for the currently-safe ones (ASCII-bounded input
folds identically), and the all-at-once sweep prevents
the AST ban from becoming aspirational.

## Out of Scope

- `cases.Title`, `cases.Upper`, or other case operations.
  `Fold` is the right primitive for *comparison*; the
  others have different correctness contracts.
- Locale-aware folding. The helper uses `language.Und`
  (Unicode default folding) which is correct for
  comparison purposes regardless of input locale. Adding
  per-locale folding would require a richer API and is
  not motivated by any current call site.
- A linter rule for `strings.EqualFold`. That function is
  ASCII-only by Go stdlib documentation; if a callsite
  uses it on non-ASCII input that's a separate bug class
  worth a future sweep, but bundling it into this commit
  conflates two ban surfaces.
- Replacing the `slug` package's `strings.ToLower` with
  more than just `i18n.Fold`. Slug generation is doing
  more than case-fold (it strips diacritics, normalizes
  to NFC, etc.) — that's a different correctness story.
  For this commit, `slug.go` swaps ToLower→Fold and
  preserves its other transforms.

## Verification

- `internal/compliance/no_strings_tolower_test.go` passes
  (zero direct `strings.ToLower` calls outside
  `internal/i18n/`).
- `make lint` clean.
- `go test ./...` clean — the helper is a drop-in for
  ASCII input, and the unsafe sites had no test coverage
  exercising non-ASCII input today, so behavior on the
  ASCII paths is unchanged. The Fold semantics for
  non-ASCII are documented in `internal/i18n/fold.go`'s
  doc comment.
