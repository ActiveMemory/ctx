# Computed Index Projection Plan — `ctx index`

**Spec:** specs/computed-index-projection.md · **Status:** Ready
**Blocking TBDs resolved:** none — the spec's Open Questions were all resolved
before commit (package rename → `internal/heading`; KB indexing out of scope;
`--depth`/`--json` ship in v1; no "malformed heading" concept; `ctx agent` needs
no rewire). No TBD remains that any task below embeds an assumption about.

Single-milestone spec: the whole spec is one milestone (this plan). No
`--milestone` split.

## Scope & DoD

Replace the stored `<!-- INDEX -->` blocks with a computed `ctx index <file>`
heading projector; delete the block-maintenance machinery; rename the surviving
parser package. Net deletion + one new command.

Definition of Done (confirmed by measurement/user, not derived from task ticks):

- [x] `ctx index <file>` projects L2 (and `--depth 3` L3) headings for any
      knowledge file, including TASKS.md `## Phase` sections.
- [x] `--depth` and `--json` both work; missing/unreadable path exits non-zero
      with a path-bearing message.
- [x] All real `<!-- INDEX -->` blocks removed from `.context/{DECISIONS,LEARNINGS,TASKS}.md`
      and none regenerated when an entry is added.
- [x] `ctx reindex`, `ctx decision reindex`, `ctx learning reindex` no longer exist
      (unknown-command, exit ≠ 0). KB indexing untouched.
- [x] `internal/index` renamed to `internal/heading`; `GenerateTable`/`Update*`/
      `Reindex`/`Validate` deleted; `ParseHeaders`/`ParseEntryBlocks`/`EntryBlock`
      retained and their consumers green.
- [x] `ctx agent` output is unaffected by block removal (regression passes).
- [x] `make build` green; `make lint` (`golangci-lint run` v2.11.4) **0 issues**;
      `make test` (whole `./...`) green; AST `audit` + `compliance` suites green;
      `lint-drift` clean; docstrings clean.

## Data model & storage

No persisted state. New in-memory type:

```go
// internal/entity
type Heading struct {
    Level int    // 2 for ##, 3 for ###
    Text  string // heading text after the # markers
}
```

Retained parser types (unchanged, moved to `internal/heading`): `IndexEntry`,
`EntryBlock`. No migrations; the "migration" is a one-shot strip of the blocks
(T19), tolerant of already-absent blocks in downstream repos (they strip on next
write, which no longer regenerates).

## Contracts

- `heading.ParseHeaders(content string) []entity.IndexEntry` — retained, entry-specific (`## [ts] Title`).
- `heading.ParseEntryBlocks(content string) []heading.EntryBlock` — retained; `ctx agent` scoring.
- `heading.Headings(content string, maxDepth int) []entity.Heading` — **new**, generic ATX matcher, code-fence-aware.
- CLI: `ctx index <file> [--depth N] [--json]` — top-level command.
  - default output: one heading per line, file order.
  - `--depth` default 2 (L2 only); 3 includes L3.
  - `--json`: array of `{level, text}`, matching `ctx status/drift --json`.

## Test matrix

| Invariant / rule (spec) | Violation attempt | Expected failure/behavior | Task |
|---|---|---|---|
| Generic matcher accepts any ATX heading (no "malformed") | Heading without timestamp (`## Phase 1`) | Included in projection | T04, T13 |
| Code-fence-aware extraction | `## x` inside ```` ``` ```` fence | Not projected | T03, T04 |
| Depth default = 2 | file with `###` sub-headings, no `--depth` | L3 omitted; `--depth 3` includes them | T16 |
| Empty / no-heading file | project a file with no `##` | empty output, exit 0 | T04, T18 |
| Missing/unreadable path | `ctx index /nope` | exit ≠ 0, path-bearing message | T18 |
| No arg | `ctx index` | usage error | T18 |
| Residual INDEX block inert | project a file still holding a stale block | block ignored, not double-counted | T04 |
| `--json` well-formed | `ctx index --json f \| jq .` | valid JSON array | T17 |
| Deleted symbol still referenced | leftover `heading.Validate` call | compile failure | T11, T12, T24 |
| No INDEX block remains under `.context/` | grep after strip | zero matches | T19, T22 |
| Reindex commands gone | `ctx reindex` / `ctx decision reindex` / `ctx learning reindex` | unknown command, exit ≠ 0 | T05–T07, T22 |
| `ctx agent` unaffected by block removal | run agent with blocks stripped | packet still lists decisions/learnings | T21 |
| add stops touching index | `ctx decision add` after strip | entry appended, no block written | T09 |

## Task breakdown

| id | st | task | deps | files | [P] | acceptance criterion | spec ref |
|----|----|------|------|-------|-----|----------------------|----------|
| **Epic A — Rename & generic parser primitive** |||||||
| T01 | [x] | `git mv internal/index → internal/heading`; update package decl + all 10 non-test importers + tests | — | internal/index/** → internal/heading/**; the 10 importer files | | `make build` green; `grep -rn 'ctx/internal/index"' internal/` empty | Impl (rename) |
| T02 | [x] | Add `entity.Heading{Level,Text}` type | — | internal/entity/heading.go | [P] | `make build` green | Data model |
| T03 | [x] | Add code-fence-aware generic ATX heading regex/util | — | internal/regex/*.go | [P] | unit test: matches `## a`,`### b`; skips fenced `## c` and `#### d` | Edge (fences) |
| T04 | [x] | Implement `heading.Headings(content,maxDepth)` + unit tests | T01,T02,T03 | internal/heading/headings.go(+_test) | | `go test ./internal/heading/ -run Headings` green over: default L2, depth3 L3, empty, no-heading, fenced-ignored, stale-block-ignored | Approach |
| **Epic B — Remove reindex command surfaces** |||||||
| T05 | [x] | Delete `ctx reindex` tree; deregister at bootstrap/group.go:109; remove `UseReindex`/`DescKeyReindex` | T01 | internal/cli/reindex/**, internal/bootstrap/group.go, internal/config/embed/cmd/base.go | | `ctx reindex` → unknown command (exit ≠ 0); `make build` green | Impl |
| T06 | [x] | Delete `ctx decision reindex`; remove its AddCommand + `DescKeyDecisionReindex` | T01 | internal/cli/decision/cmd/reindex/**, decision command wiring, config/embed/cmd/decision.go | | `ctx decision reindex` → unknown; build green | Impl |
| T07 | [x] | Delete `ctx learning reindex`; remove AddCommand + `DescKeyLearningReindex` | T01 | internal/cli/learning/cmd/reindex/**, learning command wiring, config/embed/cmd/learning.go | | `ctx learning reindex` → unknown; build green | Impl |
| T08 | [x] | Remove orphaned reindex description YAML; verify Go↔YAML linkage (lint-drift check 5) | T05,T06,T07 | internal/assets/commands/*.yaml | | `make lint` green (no orphaned DescKey) | Impl |
| **Epic C — Detach write path & delete block API** |||||||
| T09 | [x] | Remove index-mutation block from `entry/write.go` (drop `Validate`+`UpdateDecisions`/`UpdateLearnings`, ~L65–125); entries append only; drop now-unused `errAdd.IndexUpdate` | T01 | internal/entry/write.go, internal/err/(add) | | `ctx decision add`/`ctx learning add` append entry, write no INDEX block; `go test ./internal/entry/` green | Impl |
| T10 | [x] | Delete `TestWrite_RefusesEntriesTrappedInIndexBlock` + related block tests | T09 | internal/entry/write_test.go | | `go test ./internal/entry/` green | Impl |
| T11 | [x] | Delete `GenerateTable`/`Update`/`UpdateDecisions`/`UpdateLearnings`/`Reindex`/`Validate` + their tests from `internal/heading` | T05,T06,T07,T09,T10 | internal/heading/index.go,entry.go,*_test.go | | `grep -rnE 'func (GenerateTable\|Update\|Reindex\|Validate)\(' internal/heading` empty; `make build` green | Impl |
| T12 | [x] | Audit retained-parser consumers; confirm only `ParseHeaders`/`ParseEntryBlocks` used; remove `drift/check` block/`Validate` reliance if any | T11 | internal/trace/resolve_entry.go, internal/cli/agent/core/{budget/parse,score/score,score/types}.go, internal/memory/extract.go, internal/drift/check.go, internal/cli/system/core/knowledge/knowledge.go | | `make build && make test` green; `grep -rn 'heading.Validate\|heading.Update' internal/` empty | Impl |
| **Epic D — `ctx index` command** |||||||
| T13 | [x] | Scaffold `internal/cli/index/` (cmd/root/cmd.go,run.go; core projector; doc.go) calling `heading.Headings`; positional `<file>` | T04 | internal/cli/index/** | | `ctx index .context/DECISIONS.md` prints `##` headings | Interface |
| T14 | [x] | Add `UseIndex`/`DescKeyIndex` + description YAML | — | internal/config/embed/cmd/base.go, internal/assets/commands/*.yaml | [P] | `make build` green; lint linkage green | Interface |
| T15 | [x] | Register `ctx index` top-level in bootstrap | T13,T14 | internal/bootstrap/group.go (or bootstrap.go) | | `ctx --help` lists `index`; `ctx index <file>` runs | Interface |
| T16 | [x] | `--depth` flag (default 2; 3 adds `###`) | T13 | internal/cli/index/** | | `ctx index --depth 3 f` includes L3; default omits | Interface |
| T17 | [x] | `--json` flag via `internal/write` | T13 | internal/cli/index/**, internal/write/(index) | | `ctx index --json f` emits valid JSON array | Interface |
| T18 | [x] | Error handling: unreadable/missing path → path-bearing err exit ≠ 0; no arg → usage err | T13 | internal/cli/index/**, internal/err/(index) | | `ctx index /nope` exit ≠ 0 + message; `ctx index` (no arg) usage err | Error Handling |
| **Epic E — Strip blocks, marker cleanup, guards** |||||||
| T19 | [x] | Strip `<!-- INDEX -->` blocks: DECISIONS ×1, LEARNINGS ×2, TASKS ×6 | T09 | .context/DECISIONS.md, .context/LEARNINGS.md, .context/TASKS.md | | `grep -c INDEX:START` on all three = 0 | Impl |
| T20 | [x] | Remove `IndexBlockFmt`/`IndexBlockAppendFmt`/`INDEX:START`/`END` marker constants | T09,T11 | internal/config/marker/index_fmt.go, marker.go | | `grep -rn 'IndexBlock\|INDEX:START' internal/config/marker` empty; build green | Impl |
| T21 | [x] | Regression: `ctx agent` projects DECISIONS/LEARNINGS with blocks stripped | T15,T19 | (test) internal/cli/agent/** or integration | | `ctx agent` output lists decisions & learnings, no error | Approach |
| T22 | [x] | Audit test: no `<!-- INDEX -->` under `.context/`; no reindex command registered | T15,T19 | internal/audit/*_test.go | | `go test ./internal/audit/` green with new check | Testing |
| **Epic F — Docs & final gate** |||||||
| T23 | [x] | Docs: remove `reindex` from docs/cli/**; add `ctx index` reference; fix any CLAUDE.md/skill mentions | T15 | docs/cli/**, docs/**, CLAUDE.md template | | `grep -rn 'reindex' docs/` only historical/release-notes; `ctx index` documented | Interface |
| T24 | [x] | Final gate: `make build && make lint && make test`; run command-audit for stale reindex refs | ALL | — | | all three make targets exit 0; no stale `reindex` reference | DoD |

**Epic → id partition (must sum to 24):** A=T01–T04 (4) · B=T05–T08 (4) ·
C=T09–T12 (4) · D=T13–T18 (6) · E=T19–T22 (4) · F=T23–T24 (2). 4+4+4+6+4+2 = **24**.
Each id in exactly one epic.

**Completion rule:** an epic anchor in TASKS.md is `[x]` only when every task in
its id-range is `[x]` or `[o]` in this plan. This plan is the single source of
truth for milestone progress; TASKS.md epics project it.

**Execution waves (topological check):**
1. T01, T02, T03 (roots; T02/T03 `[P]`)
2. T04 · T05 · T06 · T07 · T09 · T14 (all depend only on T01/roots)
3. T08 · T10 · T13
4. T11 · T15 · T16 · T17 · T18
5. T12 · T19 · T20 · T23
6. T21 · T22
7. T24

## Risks & measurement gates

- **`drift/check.go` may embed block-`Validate` logic** (not just parse). If T12
  finds it validates the stored index, that check is obsoleted by the removal —
  delete or repurpose it, and note here. Measurement: does `ctx drift` still pass
  meaningfully after T12? If it loses a real check, surface as a spec gap.
- **`ctx agent` scoring depends on `ParseEntryBlocks`.** If the rename/trim
  changes its output ordering or scores, T21 catches it. Gate: T21 must pass
  before T24.
- **YAML↔Go linkage (lint-drift check 5)** is strict; orphaned reindex DescKeys
  fail lint. T08 owns clearing them; gate is `make lint` in T24.

## Out of scope

- **Richer `ctx list/search`** (filtering, full-text) — separate TASKS.md → Misc
  task; layers on this index primitive later.
- **KB indexing / `ctx kb reindex`** — different machinery, untouched (spec Non-Goal).
- **Time-sharding / cold-bucket** — killed in `/ctx-plan` (brief), not revisited.

## Amendments

- **2026-07-14 · T06/T07 acceptance corrected (measurement gate).** Original
  criterion: "`ctx decision reindex` / `ctx learning reindex` → unknown command,
  exit ≠ 0." Reality: `decision`/`learning` are built on `parent.Cmd`, which
  intentionally keeps cobra's default (help + exit 0) on an unknown subcommand —
  documented by `internal/cli/system/system_test.go:TestParentCmdScopeUnchanged`.
  Only `ctx system` opts into the loud `unknown.HandlerFor` handler. So the
  exit-code expectation was wrong about the codebase, not a regression I
  introduced. **Corrected criterion (verified):** the reindex subcommand is
  deregistered — absent from `ctx decision --help` / `ctx learning --help` (only
  `add` remains) — and top-level `ctx reindex` exits 1 (cobra default errors at
  root). The feature is fully removed. **Not minted here:** making the
  decision/learning parents fail loud on unknown subcommands (like `ctx system`)
  would contradict the documented `TestParentCmdScopeUnchanged` scope and is a
  separate design decision — route through `/ctx-plan` if wanted; filed as a
  follow-up thought, not a task.
- **2026-07-14 · T03 satisfied by reuse, no new regex.** `regex.MarkdownHeading`
  (`^(#{1,6}) (.+)$`) and `regex.CodeFenceLine` already existed; the fence-aware
  logic lives in `heading.Headings` (T04). No file added under `internal/regex`.
- **2026-07-14 · T04 hardened: HTML-comment-aware (bug found by driving the real
  file).** Verifying `ctx index .context/DECISIONS.md` (Epic D) surfaced that the
  matcher projected the `## Quick Format` / `## Full Format` example headings
  inside the `<!-- DECISION FORMATS ... -->` legend comment. Same class as code
  fences. `heading.Headings` now skips multi-line HTML comment blocks (inline
  `<!-- x -->` comments do not swallow following headings); two tests added. Not
  a criterion change — the T04 implementation was incomplete; this is the fix.
- **2026-07-14 · T19 scope corrected (verify-before-strip caught a data hazard).**
  Plan said "DECISIONS ×1, LEARNINGS ×2, TASKS ×6" from a raw `grep -c INDEX:START`.
  Inspecting actual marker positions: LEARNINGS has **one** real block (L17–112);
  its "2nd" hit is prose ("INDEX:START/END markers:") inside an entry body. TASKS
  has **zero** real blocks — all 6 hits are prose inside task descriptions about
  this very index-drop work. A naive regex strip would have corrupted TASKS bodies.
  Corrected: stripped only exact marker-delimited blocks in DECISIONS (70 entries
  preserved) and LEARNINGS (96 preserved) via an entry-count-guarded awk; TASKS
  left untouched. **Corrected DoD/matrix reading:** "no real `<!-- INDEX -->`
  marker LINES in DECISIONS/LEARNINGS" (met); TASKS retains its prose mentions by
  design.
- **2026-07-14 · T22 scoped to the durable invariant.** The plan's "no INDEX
  block under `.context/`" half is a one-time migration state, not a perpetual
  code invariant (downstream repos' `.context` is out of ctx's control), so it is
  not asserted by a test. The durable guard IS a test: `TestReindexCommandRemoved`
  in bootstrap_test.go walks the command tree and fails if any `reindex` command
  reappears (the `kb` subtree is exempt — KB indexing is a spec Non-Goal). `index`
  added to `TestInitialize`'s expected-commands list.
- **2026-07-14 · T24: `golangci-lint` installed and run (0 issues).** Installed
  the CI-pinned `golangci-lint v2.11.4` (from `.github/workflows/ci.yml`) and ran
  `make lint` over the whole project. It caught one real issue in new code —
  gosec G602 (false-positive out-of-range on a hand-rolled slice-compare helper
  in `headings_test.go`); replaced with `slices.Equal`. Re-run: **0 issues.** The
  AST `internal/audit` + `internal/compliance` suites also exercised the removal
  heavily — they caught 12 orphaned dead exports + 3 style nits that greps
  missed, all fixed.
- **2026-07-14 · Orphan sweep beyond the plan (audit-driven).** Deleting the
  block-maintenance API orphaned 12 exports the plan did not enumerate: the index
  table formatters (`marker.TableRowFmt/TableSepFmt`), warn labels
  (`warn.IndexHeader/IndexSeparator/IndexRow`), column DescKeys (`colummn.go`,
  whole file), drift-index writers (`write/drift/reindex.go`, whole file, +
  `DescKeyDriftCleared/Regenerated`), and journal reindex errors
  (`err/journal.ReindexFile{NotFound,Read,Write}` + DescKeys). All removed with
  their YAML, confirmed zero live callers first. `TestNoDeadExports` was the
  forcing function.
