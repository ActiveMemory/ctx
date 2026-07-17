# Spec: fix AfterHeader dropping everything after the insert point

## Problem

`insert.AfterHeader` (`internal/cli/add/core/insert/insert.go`) ends with:

```go
return []byte(content[:insertPoint] + entry)
```

It truncates the file at `insertPoint` and appends the entry — the tail
(`content[insertPoint:]`) is **silently discarded**. It is an insert
function that does not insert.

Its sole caller is `beforeFirstEntry`'s fallback
(`internal/cli/add/core/insert/before.go`), taken when a knowledge file
contains no `## [` entry. So `ctx learning add` / `ctx decision add` on
an entry-less LEARNINGS.md/DECISIONS.md destroys any content that sits
below the H1 header and its comment block.

Today this is **masked**, not absent: an entry-less file freshly created
by `ctx init` has nothing after the comment block, so `insertPoint`
lands at EOF and the discarded tail is empty. The bug bites the moment
non-entry content exists below the preamble of a file that has no
entries yet — e.g. a hand-written `## Notes` section, or any structural
section a future feature adds.

This is the same family as the `ctx learning add` clobber bug that
`index.Validate` exists to guard: silent destruction of persisted
memory, recoverable only from git.

Discovered by the layout proof written for
`specs/progressive-disclosure.md` (plan `pd-m1`, T10), which planted a
`## Themes` section below the preamble of an entry-less root and watched
`add` delete it. The fix is specified and shipped **independently** of
that design: it is a defect in current behavior on its own terms.

## Design

Make `AfterHeader` a true insert, preserving the tail — the pattern its
sibling `Task` already follows
(`existingStr[:pendingIdx] + entry + NewlineLF + existingStr[pendingIdx:]`).

When the tail is non-empty, separate the new entry from the following
content with the same delimiter `beforeFirstEntry`'s primary path uses,
so both paths in the same function family produce consistent files:

```
entry + NewlineLF + Separator + NewlineLF + NewlineLF + tail
```

When the tail is **empty**, emit `content[:insertPoint] + entry`
unchanged — byte-identical to today's output. This keeps the fix a
strict superset of current behavior: every file shape that works today
produces identical bytes, and only the shape that currently loses data
changes.

## Implementation

- `internal/cli/add/core/insert/insert.go`: `AfterHeader` returns
  `content[:insertPoint] + entry` when `content[insertPoint:]` is empty,
  and `content[:insertPoint] + entry + sep + tail` otherwise.
- No signature change; one caller; no other call sites.

## Tests

The package had **no tests**. Add:

- `testmain_test.go`: `lookup.Init()` in `TestMain`. Without it
  `desc.Text` returns `""` and `strings.Index(s, "")` is `0`, which
  silently turns every insert anchor into "match at offset 0" — tests
  would exercise a path production never takes and pass for the wrong
  reason. This is a precondition for any honest test of this package.
- Tail preserved: content below the preamble of an entry-less file
  survives an add, and the entry lands above it.
- Empty tail unchanged: an entry-less file with nothing after the
  comment block produces byte-identical output to the pre-fix
  implementation.
- Primary path untouched: a file with `## [` entries still inserts
  before the first one (regression guard for `beforeFirstEntry`).

## Acceptance

- `ctx learning add` against an entry-less file carrying a section below
  the preamble preserves that section.
- Output for every file shape that works today is byte-identical.
- `make lint` clean; full suite green.

## Non-Goals

- No change to `beforeFirstEntry`'s primary (`## [`-anchored) path.
- No change to `Task`, `TaskAfterSection`, `AppendAtEnd`, or `Decision`.
- Not a progressive-disclosure change: that design merely *found* this;
  the fix stands alone and ships alone.
