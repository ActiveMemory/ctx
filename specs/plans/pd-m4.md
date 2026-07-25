# pd-m4 Plan — CONVENTIONS digestion (the convention kind)

**Spec:** `specs/progressive-disclosure.md` · **Status:** Ready
**Blocking TBDs resolved:** convention model (`### `-under-`## Recent`
vs curated `## `-section) — decided **curated `## `-section, unified
into the entry-kind path** in the spec's revised Layout/Invariants and
`DECISIONS.md [2026-07-19-100259]`.

## Scope & DoD

Fold `CONVENTIONS.md` (18 stable, curated `## ` sections) into a bounded
root under progressive disclosure, by teaching the M3 mover the
`convention` kind — same `preamble | staging | ## Themes` layout as
entry kinds, parametrized by a per-kind entry prefix (`## ` vs `## [`)
and identity (section title vs timestamp). Retire the `### `-under-`##
Recent` model. Convention add-path prepends above `## Themes`.

DoD:
- [x] `CONVENTIONS.md` folded into ~5 theme files under
  `.context/conventions/`; root bounded to `## Themes` gists + links.
- [x] `ctx convention add` prepends above `## Themes`; the
  entry-below-themes invariant holds on the folded root.
- [x] All invariants pass on the folded root: per-kind entry-below-themes,
  gist↔file pairing, one-place uniqueness, **duplicate-title fail-loud**.
- [x] Entry-kind (LEARNINGS/DECISIONS) path **unchanged** — M3 suite
  green; `heading.ParseEntryBlocks` untouched.
- [x] `make lint` 0, `go test ./...` green, `make audit` pass.

## Data model & storage

- `Kind` gains no new value (`KindConvention` already exists).
- `StagedEntry.Timestamp` becomes **optional**: empty for conventions.
  Identity `Timestamp + IDSeparator + Title` collapses to title-only.
- Plan JSON for conventions: `{"title": "..."}` (no `timestamp`).
- New per-kind accessor `EntryPrefix(Kind)`: `## [` for learning/decision,
  `## ` for convention. Consumed by `parseEntryKind`, `Validate`, and the
  add-path anchor.
- `ThemeDir(KindConvention)` → `("conventions", true)` (constant
  `ThemeDirConvention` already exists; the case is currently missing).
- Retire constants `ConventionLinePrefix (= "### ")` and
  `HeadingRecent (= "## Recent")`.
- New sentinel `ErrDuplicateStagedTitle` (+ embedded desc text).
- Theme files created under `.context/conventions/<slug>.md`.

## Contracts

- `disclosure.Parse(content, KindConvention)` → `Root{Preamble, Staging,
  ThemesRaw, HasThemes}` with `Staging` = the `## ` sections (excluding
  the structural `## Themes`). `parseConvention` deleted; `parseEntryKind`
  is the single path, prefix-parametrized.
- Convention staging enumeration: a `## `-boundary splitter in the
  `disclosure` package (NOT `heading.ParseEntryBlocks`), returning
  `[]StagedEntry{Title, Timestamp:""}`.
- `disclosure.Apply(root, plan, KindConvention, …)` — same guarded
  append→verify→remove→single-root-rewrite as entry kinds.
- `disclosure.Validate(Root)`: entry-below-themes generalized to
  `EntryPrefix(kind)`; new duplicate-title rule for convention staging.
- `ctx disclosure apply .context/CONVENTIONS.md --plan …` — accepted
  (no longer `ErrApplyNotEntryKind`).
- `ctx convention add "<text>"` — insert anchor
  `beforeFirstEntry(## , skip ## Themes)`, fallback `AfterHeader`; lands
  newest-first above `## Themes`.

## Test matrix

| # | Invariant / behavior | Violation attempt | Expected failure/behavior | Task |
|---|---|---|---|---|
| C1 | Parse round-trips conventions byte-exact | fold + reconstruct | `Reconstruct == input` | T07 |
| C2 | Convention enumeration | 18-section fixture | 18 `StagedEntry{Title, Timestamp:""}` | T05 |
| C3 | Title identity is unique | two `## Naming` | `ErrDuplicateStagedTitle`, no move | T11 |
| C4 | `SplitStaging` on `## ` | cut N sections | byte-exact spans; conservation | T08 |
| C5 | Apply convention | fold N→M themes | conservation, all invariants, bounded root, root rewritten once | T14 |
| C6 | Kind gate | apply on convention / on unknown kind | convention accepted; unknown → `ErrApplyNotEntryKind` | T13 |
| C7 | Comment-skip (M3 carryover) | `## ` inside preamble `<!-- -->` | not an entry | T09 |
| C8 | Entry-kind regression | run M3 suite | LEARNINGS/DECISIONS unchanged; green | T22 |
| C9 | Measurement gate | drive digest on real CONVENTIONS staging | ≥3 themes, conservation, invariants, bounded | T20 |
| C10 | Post-fold add | `ctx convention add` after fold | lands above `## Themes`; entry-below-themes holds | T17 |
| C11 | Entry-below-themes (per-kind) | `## X` below `## Themes` | `ErrEntryBelowThemes` | T10 |

## Task breakdown

| id | st | task | deps | files | [P] | acceptance criterion | spec ref |
|---|---|---|---|---|---|---|---|
| T01 | [x] | `EntryPrefix(Kind)` accessor (`## [` / `## `) | — | `internal/disclosure/kind.go` | [P] | unit: `EntryPrefix(KindConvention)=="## "`, `EntryPrefix(KindLearning)=="## ["` | Layout |
| T02 | [x] | retire `ConventionLinePrefix`, `HeadingRecent` + their config tests | T01 | `internal/config/disclosure/*.go` | | `grep -rn 'HeadingRecent\|ConventionLinePrefix' internal/` → no non-test hits; config test green | Layout |
| T03 | [x] | sentinel `ErrDuplicateStagedTitle` + embedded desc text | — | `internal/err/disclosure/disclosure.go`, embed text yaml | [P] | `make audit` message-registry test green; sentinel resolves non-empty | Invariants |
| T04 | [x] | `ThemeDir(KindConvention)` → `("conventions", true)` | — | `internal/disclosure/kind.go` | [P] | unit: `ThemeDir(KindConvention)==("conventions",true)` | Tiers |
| T05 | [x] | convention `## `-boundary enumerator (title identity, `Timestamp:""`) | T01 | `internal/disclosure/collect.go` | | unit: 18-section fixture → 18 `StagedEntry{Title,""}` (**C2**) | Layout |
| T06 | [x] | kind-aware `StagedEntries`/`entryIDs` dispatch (heading vs convention enum) | T05 | `internal/disclosure/inspect.go`, `collect.go` | | inspect on convention fixture = title-only entries; entry-kind inspect unchanged | Layout |
| T07 | [x] | `parseEntryKind` prefix-parametrized; route conventions; **delete `parseConvention`** | T01,T04 | `internal/disclosure/parse.go`, `regions.go` | | Parse(convention) round-trips byte-exact (**C1**); `grep parseConvention` empty | Layout |
| T08 | [x] | `SplitStaging` cuts `## ` sections by title id | T05,T07 | `internal/disclosure/split.go` | | split test byte-exact + conservation on convention fixture (**C4**) | Guards |
| T09 | [x] | comment-skip regression for conventions | T07 | `internal/disclosure/parse_test.go` | [P] | `## ` in preamble `<!-- -->` not an entry (**C7**) | Guards |
| T10 | [x] | generalize Validate entry-below-themes to `EntryPrefix`; delete `## Recent`/`### ` branches | T07 | `internal/disclosure/validate.go` | | `## X` below `## Themes` → `ErrEntryBelowThemes`; entry-kinds unchanged (**C11**) | Invariants |
| T11 | [x] | dup-title fail-loud in Validate | T03,T05 | `internal/disclosure/validate.go` | | two `## Naming` → `ErrDuplicateStagedTitle`, no move (**C3**) | Invariants |
| T12 | [x] | rule 3 (staging enumerates) via kind-aware enumerator | T05,T07 | `internal/disclosure/validate.go` |  | non-empty convention staging, 0 parseable → `ErrStagingUnparsable` | Guards |
| T13 | [x] | lift `ErrApplyNotEntryKind` for `KindConvention` | T04,T07 | `internal/disclosure/apply.go` | | apply(convention) not refused; unknown kind still refused (**C6**) | Guards |
| T14 | [x] | Apply convention end-to-end (split+bullet+rewrite, prefix/title) | T08,T13 | `internal/disclosure/apply.go`, `move.go` | | N→M fold: conservation, invariants, bounded root, one rewrite (**C5**) | The pass |
| T15 | [x] | title-only identity through `FlattenPlan`/`entryID` | T05 | `internal/disclosure/collect.go`, `apply.go` |  | title-only plan flattens; conservation | Layout |
| T16 | [x] | `ctx convention add` prepend anchor (`beforeFirstEntry(## , skip ## Themes)`/`AfterHeader`) | T01 | `internal/cli/add/core/insert/append.go`, mode selection | | convention add on folded fixture lands above `## Themes`, newest-first; empty-staging via AfterHeader | Layout |
| T17 | [x] | post-fold add invariant test | T16,T10 | add insert test | [P] | `ctx convention add` after fold keeps entry-below-themes clean (**C10**) | Invariants |
| T18 | [x] | `ctx-digest` skill — convention path (title identity, no timestamps) | T14 | `internal/assets/claude/skills/ctx-digest/SKILL.md` | | frontmatter test green; body covers inspect→propose→gist→apply for conventions + dup-title guard | Skill |
| T19 | [x] | copilot skill sync | T18 | `internal/assets/integrations/copilot-cli/skills/ctx-digest/**` | [P] | `make check-copilot-skills` green | Skill |
| T20 | [x] | **measurement gate** — drive digest on realistic CONVENTIONS fixture | T14,T16,T18 | scratchpad fixture (copy of real CONVENTIONS) | | ≥8 sections → ≥3 themes: entries moved, gists written, conservation, all invariants, bounded root (**C9**) | Tests |
| T21 | [x] | real CONVENTIONS rollout (**human-gated**) | T20 | `.context/CONVENTIONS.md`, `.context/conventions/**` | | **only on explicit user approval:** fold real root into ~5 themes; conservation + invariants verified; own commit on sign-off | Phasing 4 |
| T22 | [x] | milestone gate | T01–T21 | — | | `make lint` 0, `go test ./...` green, `make audit` pass; changed canonical `.md` pass invariants; M3 entry-kind suite green (**C8**) | Scope/DoD |

**Execution waves** (each `∥` group is file-disjoint): `T01 ∥ T03` → `T02 ∥ T04 ∥ T05` → `T06 ∥ T07` → `T08 ∥ T09 ∥ T10 ∥ T16` → `T11` → `T12` → `T13` → `T14` → `T15` → `T17 ∥ T18` → `T19` → `T20` → `T21` → `T22`.
- **Dumb-scanner false boundary:** a `## ` inside a fenced code block is a
  *unique* string, not a duplicate, so C3 won't catch it. Safety is
  byte-conservation + the human-gated plan review (spurious entry visible
  at inspect). No fence detection (spec decision). Authoring fix if it
  ever occurs: rewrite that one line.

## Out of scope (deferred, with pointers)

| Deferred | Milestone | Note |
|---|---|---|
| Suggest-only trigger wiring (growth nudge, /ctx-remember, /ctx-wrap-up) | M5 | nudge-threshold TBD blocking there |
| Tier-2 recursion (overgrown theme file becomes a root) | — | structure self-similar; deferred indefinitely (spec Non-Goals) |
| CONSTITUTION.md / TASKS.md digestion | — | out of scope by spec (small / auto-archived) |

## Amendments

| date | what | why |
|---|---|---|
| 2026-07-23 | **T02 resequenced after T07** (was wave 2) | T02's acceptance (`grep -rn 'HeadingRecent\|ConventionLinePrefix' internal/` → no non-test hits) cannot pass while `parseConvention` still consumes both constants (`regions.go:52,64`); T07 deletes it. Wave order only — no criterion change. |
| 2026-07-23 | **Explicit convention gate added to `Apply`**, to be removed by T13 | The mover's refusal of conventions was *derived* from `ThemeDir` returning false, so T04 alone silently opened a destructive mover to CONVENTIONS.md with no convention parse/validate/dup-title path in place (caught by `TestApply_Refusals/convention_kind_refused`). T04 and T13 are coupled through this line; the gate is now explicit so T13 opens it deliberately. |
| 2026-07-23 | `TestThemeDir` (inspect_test.go) convention case updated `(false,"")` → `(true,"conventions")` | M3 test pinned the old refusal; T04 changes it by design. Entry-kind cases untouched (C8 holds). |
| 2026-07-23 | **`Reconstruct` unified** (dropped its `KindConvention` branch) | Not called out in T07, but required by it: the old region order for conventions was `preamble + themes + staging`. With staging above themes for every kind, one order serves all three and C1 round-trip holds. |
| 2026-07-23 | **`SplitStaging` gained a `Kind` parameter** (exported signature change) | T08 needs the enumerator to vary by kind. Internal-only API (`internal/disclosure`); all callers updated. |
| 2026-07-23 | New unexported `stagedBlock` type placed in `types.go`, not `collect.go` | `TestTypeFileConvention` (internal/compliance) enforces type definitions living in `types.go`; the first placement failed it. |
| 2026-07-23 | `parse_test.go` convention fixtures rewritten to the `## `-section layout | `conventionMigrated`/`conventionUnmigrated` encoded the retired `### `-under-`## Recent` model, so they pinned exactly what T02/T07 remove. |
| 2026-07-23 | **T16 redesigned** — `ctx convention add` inserts a bullet at the top of an H2 **section** (honoring the existing `--section` flag), instead of prepending above `## Themes`. **User-approved.** | The plan treated the `## ` section as the entry, but a convention entry is a *bullet*; the section is the grouping (as in ctx's own CONVENTIONS.md, where each `## ` heading groups related bullets across its 18 sections). Since `Parse` starts staging at the first `## ` line, a bullet prepended above it lands in the **preamble** — unenumerable, undigestible, and growing without bound. The original T16 would have passed its acceptance while defeating the milestone. |
| 2026-07-23 | `--section` was already registered for `ctx convention add` and silently ignored | `build.Cmd` threads `Section` into `AppendEntry` for every kind; the convention branch fell through to `AppendAtEnd`. T16 wires it rather than adding a flag. |
| 2026-07-23 | New `normalize.ConventionSection` (H2) + `cfgDisc.SectionDefaultConvention` (`"General"`) | Section names normalize to `## `, not `### ` (which would not be a staged entry). The default covers the post-fold case where no section survives: the add-path creates one above `## Themes` rather than appending at EOF (which would violate entry-below-themes). |
| 2026-07-23 | **`--section` is now REQUIRED for `ctx convention add`; placeholders refused.** **User-directed.** | A default section is a catch-all candy store: an agent that has not decided where a convention belongs dumps everything there. Choosing the section *is* the thinking, so the CLI refuses to do it. `build.requireSection` routes through the existing `validate.RejectPlaceholder`, which covers empty/whitespace (flag-empty) and the shipped placeholder set — `tbd`, `n/a`, `none`, `pending`, … (flag-placeholder). Pinned by `TestConventionAddRequiresSection` (7 cases); a refused add leaves the root byte-identical. |
| 2026-07-23 | **`--section Themes` declares a theme**, on all three canonical kinds. **User-directed.** | Content is `"<name> — <gist>"`, split on the same em-dash separator `parseThemeBullet` reads back, so the write round-trips. Slug via `slug.FromTitle`. New `internal/write/theme` + `disclosure.AddTheme`. It writes **both** halves — the gist bullet *and* the theme file — because a bullet without its file fails the gist↔file pairing invariant and would leave a root `ctx disclosure` refuses to touch. Re-declaring revises the gist and preserves the file's bodies. |
| 2026-07-23 | Theme adds bypass the per-kind body-flag gate | `--context`/`--rationale`/`--consequence`/`--lesson`/`--application` describe an *entry*; a theme has only a name and a gist. The decision/learning `PreRunE` blocks were duplicated, so both now call the shared `build.RequiredBodyFlags`, which exempts theme targets. `entry.Validate` returns early for the same reason. Ordinary adds still require their body flags (verified). |
| 2026-07-23 | `SectionDefaultConvention` (`"General"`) replaced by `SectionUnfiled` (`"Unfiled"`) | The CLI can no longer reach it. It survives only for direct `AppendEntry` callers with no section and a root with no sections; the name is deliberately unattractive so a human notices rather than settles into it. |
| 2026-07-23 | **Bug fixed in the shipped `ctx-digest` skill: the plan-JSON example was unusable.** | It documented `"entries": ["<ts> <title>"]` (strings), but `Assignment.Entries` is `[]StagedEntry`, so the CLI answers `json: cannot unmarshal string into Go struct field`. Any agent following the skill verbatim failed at the apply step. Verified both forms against the built binary; corrected to objects (`{"timestamp","title"}`) lifted verbatim from `inspect`, and added to the quality checklist. Pre-existing (M3), surfaced by writing the convention variant. |
| 2026-07-25 | **T20 measurement gate PASSED on the real 390-line, 18-section `CONVENTIONS.md`** (scratchpad copy, nothing in-repo touched) | Drove the full pipeline via the built binary: inspect → 18 title-identity entries, empty timestamps (C2 on real data); plan grouped into 5 themes; `apply` moved 18/18. Verified: **conservation** (all 18 section spans byte-present in exactly one theme file, none missing/duplicated/leaked), **bounded root** (390→22 lines, preamble + `## Themes` only), **preamble byte-exact**, **gist↔file pairing** (no broken links, no orphans), empty-plan re-apply a byte-exact no-op. Post-fold `ctx convention add --section Naming` landed a recreated section above `## Themes` (C10 on real folded root). |
| 2026-07-25 | **T21 real rollout DONE — `.context/CONVENTIONS.md` folded, user-approved grouping.** | 18 sections → 5 themes (code-style·7, layout·2, workflow·3+Refactoring=4, docs·3, cli·2). Root 390 → 22 lines. Same rigor as T20 re-run on the real result: conservation (18/18 in exactly one theme file), preamble byte-exact, bounded root (`## Themes` only), staged 0 / themes 5, no broken links. Typography linter excludes `.context/` (line 109), so theme-bullet em-dashes are clean. Refactoring moved layout→workflow vs the T20 grouping (change-discipline, not on-disk). **Not yet committed** — wants its own commit, which needs the branch carved (see session's separable change-sets). |
| 2026-07-25 | Gate limitation noted: the real file has **zero** fenced `## ` lines, so it does not exercise the dumb-scanner false-boundary path | Built a separate adversarial fixture (a `## ` inside a ```` ```md ```` fence) to close the gap: the scanner does mis-see it as a section (as documented), but folding both into one theme **conserved every byte** — mis-grouped, never corrupted. Confirms the plan's "byte-conservation is the safety net" claim on the exact shape it was written for. |
| 2026-07-23 | T18 also documents the new add-path surface | The skill gained a "Related" note for `ctx <kind> add --section Themes`, since declaring a theme and folding entries into one are now two different operations and an agent needs to know which is which. |
| 2026-07-23 | T16 acceptance leg "empty-staging via `AfterHeader`" not met as written | `AfterHeader` lands in the preamble. Empty staging now creates a default section above `## Themes` instead. Behavior is strictly safer; the criterion's mechanism is superseded by the redesign above. |
| 2026-07-23 | `layout_convention_test.go` rewritten; helpers split into `conventionsection.go` | The old file was the pd-m1 proof of the retired model. The split satisfies `TestNoMixedVisibility` (one file may not mix exported and unexported funcs); `TestNoMagicStrings` also forced `token.Whitespace` over a literal `" \t"`. |
