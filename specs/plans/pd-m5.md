# pd-m5 Plan — suggest-only triggers (two-signal knowledge health)

**Spec:** `specs/progressive-disclosure.md` · **Status:** Ready
**Blocking TBDs resolved:** nudge-threshold — decided **two-signal
health (foldability count + weight bytes)** in the spec's revised
`### Triggers` and `DECISIONS.md` (2026-07-25: two-signal health; the
Markdown-linter ceiling).

## Scope & DoD

Wire the three existing session surfaces — the `check-knowledge` growth
nudge, `/ctx-remember`, `/ctx-wrap-up` — to progressive disclosure via
one shared `knowledge.Health`, emitting two suggestion-only signals:
**foldable root** (staging count → `/ctx-digest`) and **heavy page**
(bytes over root + theme files → split / extract-to-tooling). No
auto-fold, no state file.

DoD:
- [ ] `knowledge.Health` emits foldable + heavy findings from one scan;
  both consumed by the hook and the skill report path.
- [ ] Foldability is staging-based across all three kinds via
  `disclosure.StagedEntries`; `convention_line_count` retired for
  `convention_section_count`.
- [ ] Weight scan covers the root **and** every theme file, in bytes.
- [ ] The growth nudge suggests `/ctx-digest` (foldable) and
  split/extract (heavy); `/ctx-remember` and `/ctx-wrap-up` surface both.
- [ ] Suggestion-only; existing daily-throttle/snooze reused; no new
  state file.
- [ ] `make lint` 0, `go test ./...` green, `make audit` pass.

## Data model & storage

- No persistence. Signals derive from on-disk content each scan.
- `finding` (types.go) gains a **signal kind** (`foldable` | `heavy`)
  and a **Path** (for heavy theme-file findings, not just a canonical
  basename); `Unit` extends to `bytes`/`sections`.
- rc (`internal/rc/types.go`, accessors `rc.go`): **retire**
  `ConventionLineCount`; **add** `ConventionSectionCount` and
  `ThemePageByteCeiling`.
- runtime defaults (`internal/config/runtime/runtime.go`): retire
  `DefaultConventionLineCount (200)`; add
  `DefaultConventionSectionCount = 12` and
  `DefaultThemePageByteCeiling = 65536` (64 KB). Both tunable; `0`
  disables (existing convention).

## Contracts

- `knowledge.Health(ctxDir string, cfg Thresholds) []finding` — the
  single source; foldable via `len(disclosure.StagedEntries(Parse(...)))`,
  heavy via `len(bytes)` over root + `disclosure.ThemeDir(kind)` files.
- Both-fire on one root → foldable finding ordered first (fold reduces
  both).
- `ctx system check-knowledge` gains a report path emitting the findings
  (text) the skills relay; the hook path is unchanged in shape.
- Warning text (`hooks.yaml`, `check-knowledge/warning.txt`): foldable →
  `/ctx-digest` primary, consolidate/drift secondary; heavy → "split the
  theme, or extract to tooling (a context file this heavy is a linter in
  prose)".

## Test matrix

| # | Invariant / behavior | Violation attempt | Expected | Task |
|---|---|---|---|---|
| M1 | Foldability counts staging, not whole file | un-migrated root, 40 entries | foldable finding, count 40 | T15 |
| M2 | Folded root quiets | migrated root, 3 staged + 12 themes | no foldable finding | T15 |
| M3 | Weight scans theme files | migrated root + one 80 KB theme file | heavy finding on that file | T16 |
| M4 | Both signals on one root | large un-migrated root | foldable ordered before heavy | T17 |
| M5 | Convention measure is sections | 5-section / 250-line CONVENTIONS | no foldable nudge | T18 |
| M6 | Threshold fires at `> N` | count == N exactly | no finding; `N+1` → finding | T19 |
| M7 | `0` disables a check | threshold 0 | that signal never fires | T19 |
| M8 | Surface parity | same root via hook + report | identical findings | T20 |

## Task breakdown

| id | st | task | deps | files | [P] | acceptance criterion | spec ref |
|---|---|---|---|---|---|---|---|
| T01 | [ ] | runtime: retire `DefaultConventionLineCount`; add `DefaultConventionSectionCount=12`, `DefaultThemePageByteCeiling=65536` | — | `internal/config/runtime/runtime.go`, `doc.go` | [P] | build; `grep DefaultConventionLineCount` → none | Triggers/Config |
| T02 | [ ] | rc types+accessors: retire `ConventionLineCount`; add `ConventionSectionCount`, `ThemePageByteCeiling` | T01 | `internal/rc/types.go`, `rc.go`, `default.go` | | unit: accessors return defaults; `rc.ConventionSectionCount()==12` | Config |
| T03 | [ ] | retire `ConventionLineCount` refs project-wide | T02 | `internal/**`, validate_test yaml keys | | `grep -rn ConventionLineCount internal/` → no non-test hits; config test green | Config |
| T04 | [ ] | `finding` type: add `Kind` (foldable\|heavy) + `Path`; extend `Unit` | — | `internal/cli/system/core/knowledge/types.go` | [P] | unit: zero-value + both kinds construct | Data model |
| T05 | [ ] | foldable signal: count staging via `disclosure.StagedEntries` per kind (replaces ParseEntryBlocks + line-count) | T02,T04 | `internal/cli/system/core/knowledge/knowledge.go` | | unit: un-migrated 40-entry → foldable(40); folded → none (**M1,M2**) | Signal 1 |
| T06 | [ ] | heavy signal: byte scan of the root file | T04 | `knowledge.go` | | unit: root > ceiling → heavy finding on root path | Signal 2 |
| T07 | [ ] | heavy signal: byte scan of theme files via `disclosure.ThemeDir(kind)` | T04,T06 | `knowledge.go` | | unit: 80 KB theme file → heavy finding on that file (**M3**) | Signal 2 |
| T08 | [ ] | `knowledge.Health`: combine signals; both-fire → foldable first | T05,T06,T07 | `knowledge.go` | | unit: large un-migrated root → foldable ordered before heavy (**M4**) | Contracts |
| T09 | [ ] | warning text: foldable→/ctx-digest primary (consolidate/drift secondary); heavy→split/extract | T08 | `internal/assets/commands/text/hooks.yaml`, `internal/assets/hooks/messages/check-knowledge/warning.txt` | | build embeds; message-registry test green; text names /ctx-digest | Contracts |
| T10 | [ ] | `FormatWarnings`/`EmitWarning` route by `Kind` | T08,T09 | `knowledge.go` | | unit: foldable vs heavy render distinct remedy lines | Contracts |
| T11 | [ ] | `ctx system check-knowledge` report path the skills call | T08 | `internal/cli/system/cmd/checkknowledge/*.go` | | e2e: command prints findings for an oversized fixture root | Contracts |
| T12 | [ ] | hook wiring: `CheckHealth` uses `Health`; throttle/log-first unchanged | T08,T10 | `knowledge.go` | | e2e: hook nudge names /ctx-digest for a foldable root; throttled once/day | Triggers |
| T13 | [ ] | `/ctx-remember` skill: surface foldable/heavy at session start (read-only) | T11 | `internal/assets/claude/skills/ctx-remember/SKILL.md` | [P] | frontmatter test green; body calls the report path, relays suggestions | Triggers |
| T14 | [ ] | `/ctx-wrap-up` skill: surface at session end, suggest-only never inline | T11 | `internal/assets/claude/skills/ctx-wrap-up/SKILL.md` | [P] | frontmatter test green; body surfaces + explicitly does not fold | Triggers |
| T15 | [ ] | `Health` unit tests: foldable, folded-quiet fixtures | T05 | `internal/cli/system/core/knowledge/*_test.go` | | **M1,M2** pass | Tests |
| T16 | [ ] | `Health` unit tests: heavy root + heavy theme file | T07 | `knowledge` test | [P] | **M3** pass | Tests |
| T17 | [ ] | `Health` unit test: both-fire ordering | T08 | `knowledge` test | [P] | **M4** pass | Tests |
| T18 | [ ] | convention measure test: sections not lines | T05 | `knowledge` test | [P] | **M5** pass | Tests |
| T19 | [ ] | boundary + disable tests (`>N`, `0` disables) | T08 | `knowledge` test | [P] | **M6,M7** pass | Tests |
| T20 | [ ] | surface parity test: hook findings == report findings | T11,T12 | `knowledge`/checkknowledge test | | **M8** pass | Tests |
| T21 | [ ] | copilot skill sync | T13,T14 | `internal/assets/integrations/copilot-cli/skills/**` | [P] | `make check-copilot-skills` green | Skill |
| T22 | [ ] |XX **measurement gate**: drive Health on a realistic .context (oversized root + fat theme file) | T12,T11 | scratchpad fixture | | foldable + heavy both surface; suggestions correct; throttle holds | Risks |
| T23 | [ ] |XX milestone gate | T01–T22 | — | | `make lint` 0, `go test ./...` green, `make audit` pass; DoD confirmed | Scope/DoD |

**Execution waves** (each `∥` group is file-disjoint): `T01 ∥ T04` →
`T02` → `T03 ∥ T05 ∥ T06` → `T07` → `T08` → `T09` → `T10 ∥ T11` →
`T12` → `T13 ∥ T14 ∥ T15 ∥ T16 ∥ T17 ∥ T18 ∥ T19` → `T20 ∥ T21` →
`T22` → `T23`.

## Risks & measurement gates

- **T22 measurement gate**: driving `Health` on a real oversized
  `.context` (not unit fixtures) is where a signal-ordering or
  theme-scan-path bug hides — expect it to find something (every prior
  disclosure milestone's gate did). Verify both signals surface with
  correct remedies before T23.
- **Default values are guesses** (`section_count=12`, `byte_ceiling=64KB`):
  the gate may show them too eager or too lax on real content; tune in an
  amendment, not by editing acceptance criteria.
- **Throttle coverage**: confirm the `check-knowledge` daily throttle
  actually suppresses the fold nudge across a session (assumption from
  the stress-test); if it does not, that is a new task, not a silent gap.

## Out of scope (deferred, with pointers)

| Deferred | Milestone | Note |
|---|---|---|
| Tier-2 theme recursion (overgrown theme file → its own root) | — | deferred indefinitely (spec Non-Goals); heavy-page advises split/extract instead |
| Auto-fold at any trigger | — | precluded by design (suggest-only) |
| CONSTITUTION.md / TASKS.md health | — | out of scope by spec |

## Amendments

| date | what | why |
|---|---|---|
| — | — | — |
