# pd-m3 Plan — the mover (append → verify → remove + gist write-back)

**Spec:** `specs/progressive-disclosure.md` · **Status:** Ready
**Blocking TBDs resolved:** none block M3.
- Gist format — resolved **in the spec** (`### Gist format`): one bullet
  `- <name> — <gist> → [<name>](<noun>/<slug>.md)`. Consumed here, not
  re-litigated.
- CONVENTIONS `##`-vs-`###` digestion — **deferrable**, graduates to
  blocking at **M4**. The M3 mover handles entry kinds only (learning,
  decision) and **refuses** the convention kind (`ErrApplyNotEntryKind`).
- Suggest-only trigger wiring — **deferrable**, M5.

## Scope & DoD

Milestone 3 delivers **the mover**: the first pass that *writes* canonical
knowledge files. It moves staged entries into per-theme tier-1 files and
folds their gists into the root's `## Themes`, under the spec's guards
(append → verify byte-presence → remove; validate precondition; crash
ordering additive-first, single root rewrite last; fail loud, no
auto-repair). It consumes the `Plan` the `ctx-digest` skill authors from
the M2 `Inspection`.

Entry kinds only. CONVENTIONS digestion stays M4 (its `##`-section model
collides with the region delimiters — pd-m1/pd-m2 amendments).

**The write path is guarded by code, not by agent discipline.** The M1
clobber-bug rationale ("never regenerate from what I recognized") means
the mover must not trust the skill to move files safely — the skill
supplies *semantics* (theme assignment + gist prose); the `Apply`
mechanism enforces *safety* (validate, lossless byte-cut, verify, order).

DoD (confirmed by measurement or by the user — never derived from task
completion):

- [ ] `make lint` = 0 issues, `go test ./...` green, `make audit` passes
- [ ] `disclosure.Apply` moves staged entries to theme files with
      append→verify→remove ordering; on any theme-write/verify failure the
      **root file is byte-identical** to before (T06)
- [ ] **Conservation** on a driven fixture: `entryIDs(staging_before) ==
      moved ⊎ entryIDs(staging_after)` (disjoint union); every moved entry
      byte-present in **exactly one** theme file; zero loss, zero dups (T09)
- [ ] **Idempotency**: `Apply` with an empty plan (or empty staging) is a
      no-op; root byte-identical (T07)
- [ ] Post-apply, all four invariants pass on the result — `Validate`,
      `CheckPairing`, `CheckUniqueness`, `CheckLinks` (T08)
- [ ] **Gist write-back**: a touched/new theme's bullet is
      (re)written in `## Themes`; untouched theme bullets are **byte-preserved**
      (T03, T08); first run on an un-migrated root **creates** `## Themes`
      below staging (T07)
- [ ] `ctx disclosure apply <root> --plan <p.json>` performs the move and
      reports what it did; refuses non-knowledge files and the convention
      kind, leaving the file untouched (T10–T12)
- [ ] The `ctx-digest` skill, driven on a realistic fixture, completes a
      full apply that moves entries and writes gists — the milestone's
      **measurement gate** (T16)
- [ ] Real `.context/LEARNINGS.md` → `.context/DECISIONS.md` rollout is
      performed and verified **only on explicit user sign-off** (T17,
      human-gated)

## Data model & storage

No state file — staging **is** the watermark (DECISION `[2026-07-16-215955]`).
New value types in `internal/disclosure` (`apply.go`/`types.go`):

```go
// Plan is the digest plan the ctx-digest skill authors and Apply
// executes: per target theme, which staged entries move there and the
// gist to write. Entry identity is timestamp+title (IDSeparator-joined),
// matching entryIDs / CheckUniqueness.
type Plan struct {
    Kind        string       `json:"kind"`        // "learning" | "decision"
    Assignments []Assignment `json:"assignments"`
}
type Assignment struct {
    Theme   string   `json:"theme"`   // theme name (bullet label + heading text)
    Slug    string   `json:"slug"`    // theme-file basename stem → <noun>/<slug>.md
    Gist    string   `json:"gist"`    // authored one-line gist (spec ### Gist format)
    Entries []string `json:"entries"` // entry IDs to move here, in file order
}

// ApplyResult reports a successful Apply, for CLI output.
type ApplyResult struct {
    Moved  int      `json:"moved"`  // entries moved
    Themes []string `json:"themes"` // theme slugs created or appended
}
```

**Storage layout** (unchanged from spec): theme files at
`<ctxDir>/<noun>/<slug>.md`, `<noun>` ∈ {`learnings`,`decisions`} via the
`cfgDisc.ThemeDir*` constants. Root rewritten in place.

**Theme-file body format:** on create, `# <Theme>\n\n` + the moved entry
spans; on append, `<existing>\n\n` + the moved spans. The exact separator
is unconstrained by conservation (which checks entry *presence*, not file
bytes); byte-presence verify uses `strings.Contains(fileContent, span)`.

## Contracts

**Disclosure package additions** (`internal/disclosure`):

- `ThemeDir(k Kind) (string, bool)` — Kind → `learnings`/`decisions`
  subdir; `(_, false)` for convention (mover refuses it) and unknown.
- `SplitStaging(staging string, moveIDs []string) (moved map[string]string, remaining string, err error)`
  — the **lossless byte-cut**. Computes each entry's *untrimmed*
  header-to-next-header span over raw `staging` (NOT via
  `heading.ParseEntryBlocks`, which trims trailing blanks and is lossy —
  it is used only for *identity*). Returns each moved id → its verbatim
  span and the remaining staging byte-exact. `ErrEntryNotInStaging` for an
  unknown id; `ErrEntryAssignedTwice` for a repeated id.
- `verifyContains(fileContent, span string) error` — pure predicate;
  `ErrVerifyFailed` when `span` is absent. The byte-presence guard.
- `WriteThemeBullet(themesRaw string, a Assignment, hadThemes bool) string`
  — gist write-back over **raw** `## Themes` text: replace the bullet line
  whose parsed name == `a.Theme`, else append a new bullet; create the
  `## Themes` heading when `!hadThemes`. Untouched bullet lines are
  preserved verbatim (spec: "leave untouched themes alone"). Renders
  `- <theme> — <gist> → [<theme>](<noun>/<slug>.md)` using
  `token.PrefixListDash`, `token.MetaSeparator`, `cfgDisc.ThemeArrow`
  (new const `" → "`), `cfgDisc.LinkOpen`.
- `Apply(rootPath string, plan Plan, ctxDir string) (ApplyResult, error)`
  — the IO mover. Order (spec Guards §1,§3):
  1. Read+`Parse` (kind via `KindFor(basename)`); refuse convention/
     unknown kind (`ErrApplyNotEntryKind`).
  2. `Validate(root)` — refuse malformed, no write (Guard §2).
  3. Validate the plan against staging via `SplitStaging` (fail loud on
     unknown/duplicate ids; `ErrEmptyAssignment` for an assignment with no
     entries).
  4. **Additive first:** for each assignment, append its moved spans to
     `<ctxDir>/<noun>/<slug>.md` (mkdir + create if new).
  5. **Verify:** re-read each written theme file; `verifyContains` every
     moved span. Any miss → abort (`ErrVerifyFailed`), **root untouched**.
  6. Fold gists: `WriteThemeBullet` per assignment over `root.ThemesRaw`.
  7. **Single root rewrite, last:** write `preamble + remaining + newThemes`
     to `rootPath`. This is the only root write; reached only after every
     verify passes, so any earlier failure leaves the root byte-identical.
  Fail-loud, no auto-repair. Worst-case crash = duplicated theme append
  (recoverable), never loss.

**CLI** (`ctx disclosure apply <file> --plan <path|->`):

- New subcommand under the existing `disclosure` group (registered in
  `internal/cli/disclosure/disclosure.go` alongside `inspect.Cmd()`).
  `Use`/`DescKey` constants in `config/embed/cmd`.
- Reads `--plan` (file path, or `-` for stdin) as JSON `Plan`, resolves
  ctxDir, calls `disclosure.Apply`, renders `ApplyResult` (human + `--json`
  via `internal/write/disclosure`). Kind-inference / convention refusal →
  typed sentinel, non-zero exit, **no write**.

**Skill** (`internal/assets/claude/skills/ctx-digest/SKILL.md`):

- Extend from dry-run-only to the **apply** path. After the human approves
  the proposed plan (M2 flow unchanged through step 4): write the `Plan`
  JSON → run `ctx disclosure apply <root> --plan <plan.json>` → re-run
  `ctx disclosure inspect` to confirm staging shrank and themes grew.
  The human gates the apply explicitly (skill presents the plan and asks
  before invoking). States conventions are deferred to M4. Copilot copy
  synced via `make sync-copilot-skills`.

No `ctx agent` change (spec Non-Goals). No `add`-path change.

## Test matrix

| invariant / rule | attempt | expected | task |
|---|---|---|---|
| Validate precondition (Guard §2) | Apply on a two-`## Themes` root | `ErrMultipleThemes`, root untouched | T05 |
| lossless byte-cut | staging with blank-line-separated entries; move a subset | `remaining` byte-exact for un-moved; each `moved` span verbatim (incl. body) | T02 |
| entry not in staging | plan id absent from staging | `ErrEntryNotInStaging`, no write | T02 |
| entry assigned twice | id in two assignments | `ErrEntryAssignedTwice`, no write | T02 |
| empty assignment | assignment with `Entries: []` | `ErrEmptyAssignment`, no write | T02 |
| verify predicate | `verifyContains` on content lacking span | `ErrVerifyFailed` | T06 |
| append→verify→remove; root untouched (Guard §1,§3) | force the *last* theme append to fail (unwritable path) | error returned; earlier theme file got its append; **root byte-identical** | T06 |
| gist write-back — new theme | plan with a new theme | `## Themes` gains `- name — gist → [name](noun/slug.md)`; `parseThemeBullet` round-trips name+link | T03 |
| gist write-back — touched theme | re-touch an existing theme | its bullet line replaced; **other bullet lines byte-preserved**; order kept | T03 |
| first-run themes creation | Apply on `HasThemes=false` root | `## Themes` created below staging; `Validate` passes | T07 |
| idempotency | Apply empty plan / empty staging | no-op; root byte-identical | T07 |
| conservation | Apply moving M of N staged | `entryIDs(before) == moved ⊎ entryIDs(after)`; each moved id present in exactly one theme file | T09 |
| post-apply invariants | after a successful Apply | `Validate`+`CheckPairing`+`CheckUniqueness`+`CheckLinks` all nil | T08 |
| convention refused | Apply on `CONVENTIONS.md` | `ErrApplyNotEntryKind`, no write | T05 |
| CLI apply writes | `apply <fixture> --plan p.json` | root+theme files updated; `ApplyResult` printed | T10 |
| CLI rejects non-knowledge file | `apply README.md` | non-zero, kind sentinel, no write | T11 |
| CLI write-safety on bad root | `apply` on a malformed root | non-zero, file byte-identical | T12 |
| skill drive (measurement) | drive ctx-digest apply on realistic fixture | entries moved, gists written, conservation + all four invariants hold | T16 |

## Task breakdown

| id | st | task | deps | files | [P] | acceptance criterion | spec ref |
|---|---|---|---|---|---|---|---|
| T01 | [x] | `Plan`/`Assignment`/`ApplyResult` types; `ThemeDir(k)`; `ThemeArrow` const | — | `internal/disclosure/types.go`, `apply.go`, `internal/config/disclosure/disclosure.go`, `kind_test.go` | [P] | `go test ./internal/disclosure/ -run TestThemeDir`: learning/decision → dir,true; convention/unknown → _,false | Data model |
| T02 | [x] | `SplitStaging` lossless byte-cut + plan-vs-staging validation (uses T04 sentinels) | T01,T04 | `internal/disclosure/split.go`, `split_test.go` | | `-run TestSplitStaging`: remaining byte-exact for un-moved incl. separators; moved spans verbatim; unknown id → `ErrEntryNotInStaging`; dup id → `ErrEntryAssignedTwice`; empty entries → `ErrEmptyAssignment` | Contracts, Guards §1 |
| T03 | [x] | `WriteThemeBullet` gist write-back over raw `## Themes` | T01 | `internal/disclosure/themes.go`, `themes_test.go` | [P] | `-run TestWriteThemeBullet`: new theme appends a bullet `parseThemeBullet` re-parses (name+link match); touched theme's line replaced, other lines byte-preserved; `!hadThemes` creates the heading | Contracts, Design step 4 |
| T04 | [x] | all new error sentinels + YAML descriptions: `ErrEntryNotInStaging`, `ErrEntryAssignedTwice`, `ErrEmptyAssignment`, `ErrVerifyFailed`, `ErrApplyNotEntryKind` | — | `internal/err/disclosure/*.go`, `commands/text/errors.yaml` | [P] | `go test ./internal/audit/ -run TestDescKeyYAMLLinkage` green (bijection holds); all five sentinels defined | CONVENTIONS |
| T05 | [x] | `Apply` IO mover — validate→append→verify→rewrite; kind/precondition refusal | T02,T03,T04 | `internal/disclosure/apply.go`, `apply_test.go` | | `-run TestApply` (temp dir): moves entries; theme files hold bodies; root = preamble+remaining+themes; `ApplyResult{Moved,Themes}` correct; `CONVENTIONS.md` → `ErrApplyNotEntryKind`; two-themes root → `ErrMultipleThemes`, no write | Contracts, Guards |
| T06 | [x] | abort/ordering guarantee — root untouched on failure | T05 | `internal/disclosure/apply_abort_test.go` | [P] | `-run TestApplyAbort`: (a) `verifyContains` absent span → `ErrVerifyFailed`; (b) force last theme append to fail → error, root file byte-identical, earlier theme append present | Guards §1,§3 |
| T07 | [x] | first-run `## Themes` creation + idempotency | T05 | `internal/disclosure/apply_firstrun_test.go` | [P] | `-run TestApplyFirstRun`: un-migrated root gets `## Themes` below staging, `Validate` nil; `-run TestApplyIdempotent`: empty plan → root byte-identical | Design step 5 |
| T08 | [x] | post-apply invariants hold | T05 | `internal/disclosure/apply_invariants_test.go` | [P] | `-run TestApplyInvariants`: after Apply on a fixture, `Validate`+`CheckPairing`+`CheckUniqueness`+`CheckLinks` all return nil | Invariants |
| T09 | [x] | conservation property test | T05 | `internal/disclosure/apply_conserve_test.go` | [P] | `-run TestApplyConservation`: `entryIDs(before)` == disjoint union of moved ids and `entryIDs(after)`; each moved id present in exactly one theme file; zero loss/dup | Tests (Conservation) |
| T10 | [x] | `ctx disclosure apply` subcommand + result output | T05 | `internal/cli/disclosure/cmd/apply/{cmd.go,run.go}`, `internal/cli/disclosure/disclosure.go`, `internal/config/embed/cmd/*.go`, `internal/write/disclosure/*.go` | | `ctx disclosure apply <fixture> --plan p.json` updates root+theme files, prints result; `--json` decodes to `ApplyResult` | CLI |
| T11 | [x] | CLI rejects non-knowledge file | T10 | `internal/cli/disclosure/cmd/apply/run_test.go` | [P] | `apply README.md` exits non-zero with the kind sentinel; file untouched | Contracts |
| T12 | [x] | CLI write-safety on malformed root | T10 | `internal/cli/disclosure/cmd/apply/run_write_test.go` | [P] | integration: `apply` on a two-`## Themes` root → non-zero, file byte-identical | Scope/DoD |
| T13 | [x] | `doc.go` for the apply cmd + command-wiring guards | T10 | `internal/cli/disclosure/cmd/apply/doc.go` | [P] | `make audit` doc.go/docstring floors pass; `go test ./internal/compliance/ -run TestShippedHooksResolve` and command-tree tests green with the new subcommand | CONVENTIONS |
| T14 | [x] | `ctx-digest` skill — apply path | T10 | `internal/assets/claude/skills/ctx-digest/SKILL.md` | | frontmatter test green; body describes inspect→propose→gist→**build plan JSON→human approve→`ctx disclosure apply`→confirm via inspect**, states the append→verify→remove guard and conventions=M4 | Skill |
| T15 | [x] | copilot skill sync | T14 | `internal/assets/integrations/copilot-cli/skills/ctx-digest/**` | | `make check-copilot-skills` green | Skill |
| T16 | [x] |XX **measurement gate** — drive skill apply on realistic fixture | T14,T08,T09 | scratchpad fixture (copy of real LEARNINGS staging) | | drive ctx-digest apply on N≥4 staged entries → ≥2 themes: entries moved, theme files created with bodies, `## Themes` gists written, conservation holds, all four invariants pass, root bounded | Scope/DoD |
| T17 | [x] |XX real LEARNINGS→DECISIONS rollout (**human-gated**) | T16 | `.context/LEARNINGS.md`, `.context/DECISIONS.md`, `.context/learnings/**`, `.context/decisions/**` | | **only on explicit user approval:** apply on the real roots; conservation + all four invariants verified; committed on user sign-off (own commit) | Spec Phasing 3 |
| T18 | [x] |XX milestone gate | T01–T17 | — | | `make lint` 0, `go test ./...` green, `make audit` pass; changed canonical `.md` pass invariants | Scope/DoD |

**Execution waves:** `T01,T04` ∥ → `T02,T03` ∥ → `T05` →
`T06,T07,T08,T09` ∥ → `T10` → `T11,T12,T13` ∥ → `T14` → `T15` →
`T16` → `T17` → `T18`.

## Epic anchors (TASKS.md projection)

Epics **partition** the task ids:

| Epic | Range | Count |
|---|---|---|
| E1 Mover core (disclosure package) | T01–T09 | 9 |
| E2 `ctx disclosure apply` CLI | T10–T13 | 4 |
| E3 Apply skill, rollout, gate | T14–T18 | 5 |

Arithmetic: 9 + 4 + 5 = **18** = the task count. No id double-counted or
unclaimed. **Completion rule:** an epic is `[x]` in TASKS.md only when
every task in its range is `[x]`/`[o]` here.

## Risks & measurement gates

- **⚠️ T16 is the milestone's measurement gate.** If a full driven apply
  cannot move entries losslessly with gists on a realistic fixture, the
  mover design is wrong and must be reconsidered before touching real
  canonical files (T17).
- **T17 mutates canonical `.context/` files** — the payoff, and the
  clobber-risk class the M1 guards exist for. It is **human-gated**: run
  only on explicit user approval, verify conservation + invariants before
  committing, and land it as its own commit. The code (T01–T16) is
  complete and verifiable **independent of T17**.
- **Lossless byte-cut is the fault line.** `heading.ParseEntryBlocks`
  trims trailing blank lines (`internal/heading/entry.go`), so re-joining
  parsed blocks is lossy. `SplitStaging` must operate on raw
  header-to-next-header spans; T02's byte-exact assertion is what proves
  it. Using `ParseEntryBlocks` for the *cut* (not just identity) would
  silently drop inter-entry whitespace.
- **`parseThemeBullet` folds the `→ [link]` tail into `Theme.Gist`**
  (M2 parse quirk — the parsed `Gist` retains the arrow+link suffix).
  M3 sidesteps it by editing **raw** bullet lines (`WriteThemeBullet`),
  never re-rendering from the parsed `Theme`. Noted, not fixed here (M2
  surface; not blocking).
- **Verify-miss is defensive.** After a successful append, byte-presence
  always holds, so the true miss branch is only reachable via IO
  corruption. T06 proves the guard two ways: the pure `verifyContains`
  predicate (miss → `ErrVerifyFailed`) and the ordering property (append
  failure → root untouched).

## Out of scope (deferred, with pointers)

| Deferred | Milestone | Note |
|---|---|---|
| CONVENTIONS digestion (`##` section model) | M4 | mover refuses `convention` kind (`ErrApplyNotEntryKind`); `##`-vs-`###` TBD is blocking there |
| Suggest-only trigger wiring (growth nudge, /ctx-remember, /ctx-wrap-up) | M5 | nudge-threshold TBD blocking there |
| `--dry-run` on `apply` | — | dry-run already covered by `inspect` + the skill's propose step (M2); not needed |
| Nesting / theme taxonomy | — | structure self-similar; deferred indefinitely (spec Non-Goals) |
| `ctx drift` wiring beyond the existing path check | — | spec does not require it |

## Amendments

| date | what | why |
|---|---|---|
| 2026-07-18 | `Assignment.Entries` changed from `[]string` (NUL-joined ids) to `[]StagedEntry` (`{timestamp,title}`). `FlattenPlan` now builds the ids internally; new `entryID(StagedEntry)` helper. Data-model refinement — no acceptance criterion changed (conservation/byte-cut/guards identical); T02's `SplitStaging` tests are untouched (it still takes `[]string`), only the plan/`FlattenPlan` constructors changed. | The skill authors the plan JSON; a NUL-separated id is a foot-gun to hand-author. `[]StagedEntry` is exactly `inspect --json`'s `staging` shape, so the skill lifts entries verbatim — no transformation, no ` `. Verified end-to-end against the built binary. |
| 2026-07-18 | Split the `movehelpers.go` bag file into semantically-relevant files: `move.go` (the mover's mechanics — append/verify/rewrite), `regions.go` (`renderThemeBullet`, line/offset scanners — beside `parseThemeBullet`/`lineAt`), `collect.go` (`entryID` — beside `entryIDs`). | User convention: "helpers"/"utils" bag files are lazy; functions belong in files named for their concern. [[no-helper-bag-files]] |
| 2026-07-19 | Parser: `firstLinePrefixOffset` and `headingLineOffsets` now skip `## [` / `## Themes` lines inside HTML comments (`htmlCommentSpans`; unterminated open → EOF). Shipped as its own `fix(disclosure)` commit with two regression tests; no knowledge-file edit. | Surfaced by the T17 rollout gate (mirrors pd-m1's afterheader find). DECISIONS.md's `<!-- … -->` format guide carries a `## [YYYY-MM-DD] Decision Title` example; once T17 folded every real entry out, that commented example was the only `## [` above `## Themes`, so `Validate` wrongly returned `ErrStagingUnparsable`. Root-cause fix in the parser — the folded root validates unchanged. |
