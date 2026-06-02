# Fix: `ctx learning add` Destroys Bodies Trapped in the Index Block

`ctx learning add` (and `ctx decision add`, and every `reindex`
command) silently destroys entry bodies that live between the
`<!-- INDEX:START -->` / `<!-- INDEX:END -->` markers. This phase
makes the index-regeneration path **fail loud and touch nothing**
when it cannot regenerate without data loss.

## Problem

`index.Update` (`internal/index/index.go:121-138`) treats the entire
span between `INDEX:START` and `INDEX:END` as a disposable, machine-
owned index and replaces it wholesale with a freshly generated table:

```go
before := content[:startIdx+len(marker.IndexStart)]
after := content[endIdx:]
return before + nl + indexContent + after // span destroyed
```

That assumption is false for files where real entry bodies sit inside
the marker span. This is exactly the shape an older dash-bullet
`ctx init` produced (the marker block doubled as the entry list), and
the shape a `ctx hub` file drifts into. Because `ParseHeaders` scans
the *whole* file, the regenerated table still looks complete — which
**masks** that every body was just deleted.

Verified live on `main` (scratch repo, `fix/learning-add-index-data-loss`):
a `LEARNINGS.md` with two hand-authored bodies between the markers
collapsed, after a single `ctx learning add`, to *just the index
table* — both existing bodies **and** the newly added body were gone.
Field-observed in `things-wtf-hub` session aa32f065 (commit 2dc4d1a,
−44 lines; recovered only via `git show`).

A second, related failure lives in the same function: the
"no valid markers" branch (`index.go:140-164`) inserts a fresh block
via `marker.IndexBlockFmt` / `IndexBlockAppendFmt`, both of which emit
a **new** `INDEX:START`/`INDEX:END` pair (`marker/index_fmt.go:15-23`).
When the existing markers are duplicated, single, or out of order,
this produces a *second* `INDEX:START` rather than failing.

Severity: HIGH — silent destruction of persisted memory, the one
thing ctx promises to protect. Only git made it recoverable.

## Approach

Fail loud, touch nothing. Do **not** auto-repair the file (an opt-in
`reindex --repair` is deferred to a follow-up task). Add a precondition
guard that runs *before any write* and refuses to proceed when
regenerating the index would lose data or duplicate a marker.

New function `index.Validate(content, fileName string) error`:

1. Count `INDEX:START` and `INDEX:END` occurrences.
   - `0` and `0` → OK (legitimate fresh-index creation; `Update`'s
     insert path preserves all content).
   - any other count `!= 1` each → malformed (duplicate or missing
     one) → error. Refusing here also closes the duplicate-marker
     branch.
2. With exactly one of each: if `INDEX:END` precedes `INDEX:START`
   → malformed → error.
3. Else inspect the span between the markers with `regex.EntryHeader`.
   Any `## [timestamp] Title` match → trapped bodies → error.
   Otherwise → OK.

The guard never mutates; it only reads and classifies.

### Wiring (two choke points)

Every index-regenerating path funnels through one of these, and both
read the file before mutating it:

- `entry.Write` (`internal/entry/write.go`) — the add path (covers
  `learning/decision add`, `memory.Promote`, `watch addEntry`,
  `add Run`). Call `index.Validate` on the *existing* content right
  after read, gated to the indexed types (`Decision`, `Learning`),
  before `AppendEntry` / the first `SafeWriteFile`. On error, return
  it — nothing is written.
- `index.Reindex` (`internal/index/index.go`) — all three reindex
  commands. Call `index.Validate` on the read content before
  `updateFunc` / write.

`index.Update`'s signature is left unchanged (`func(string) string`)
to avoid rippling through `Reindex`'s `updateFunc func(string) string`
and the reindex call sites — the CRITICAL blast radius gitnexus
flagged. The guard makes malformed content unreachable by `Update`.

### Errors

New `internal/err/index` package, two constructors backed by i18n desc
keys in `internal/assets/commands/text/errors.yaml`:

- `err.index.entries-in-block` — `%s` is the file name. Explains that
  entry content sits between the markers, that regeneration is refused
  to avoid deletion, and that the fix is to move `INDEX:END` above the
  entries.
- `err.index.malformed-markers` — `%s` is the file name. Markers are
  missing, duplicated, or out of order; restore a single well-formed
  pair and retry.

## Guard (round-trip test, BOTH formats)

Per the task: a round-trip test for **both** index formats that
asserts existing bodies survive an add and exactly one marker pair
remains.

- `index` package: table tests for `Validate` — empty index (OK),
  populated table between markers (OK), zero markers (OK), `##`
  header between markers (error), duplicate `INDEX:START` (error),
  `END` before `START` (error).
- `entry`/add level: well-formed dash-bullet-seeded and table-seeded
  files both survive an add with all prior bodies intact and exactly
  one marker pair; a malformed (bodies-in-block) file is refused with
  **no write** (file byte-identical after the failed add).

## Out of Scope

- `ctx <type> reindex --repair`: an opt-in self-heal that would
  relocate `INDEX:END` above trapped bodies instead of refusing. Not
  built and not filed — the chosen behavior is fail-loud + manual fix.
  Recorded here only as the obvious extension if demand appears.

## Settled Decisions

1. Behavior is fail-loud + touch-nothing, not auto-repair (user
   decision, this session).
2. `index.Update` signature stays `func(string) string`; the guard is
   a separate precondition, not a return-value change — keeps the
   CRITICAL-risk call graph untouched.
3. Zero markers is allowed (fresh creation), not treated as malformed.
