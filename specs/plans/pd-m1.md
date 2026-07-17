# pd-m1 Plan — guards, invariants, and structural vocabulary

**Spec:** `specs/progressive-disclosure.md` · **Status:** Ready
**Blocking TBDs resolved:** B1 (validate vs first-run contradiction) —
resolved **in the spec**, Guards §2: validate accepts **zero or one**
`## Themes`; zero = not-yet-migrated (the pass creates it), two+ =
malformed → refuse. Invariants need no carve-out (vacuously true on an
un-migrated root).

## Scope & DoD

Milestone 1 builds the **guards, invariants, and structural vocabulary**
and proves the layout premise. **Nothing moves**: no entry body is
relocated, no gist is authored, no `.context` knowledge content changes.
This ordering is deliberate — the pass (M2+) moves entry bodies, which is
the clobber/data-loss risk class, so the refusal machinery and the
layout proof land first.

DoD (confirmed by measurement or by the user — never derived from task
completion):

- [ ] `make lint` reports 0 issues and `go test ./...` is green
- [ ] The **layout proof** passes for all three kinds, with populated
      *and* empty staging: `add` lands in the staging zone with a
      `## Themes` section present (T10–T12)
- [ ] `Validate` refuses every malformed shape in the test matrix with
      the correct sentinel (T06)
- [ ] LEARNINGS.md / DECISIONS.md / CONVENTIONS.md bodies are
      **byte-identical** to milestone start (`git diff --stat` shows no
      change to them)

## Data model & storage

No persistence, no migrations, no DDL. In-memory only:

```go
type Kind int // KindLearning | KindDecision | KindConvention

type Theme struct {
    Name string // heading text under ## Themes
    Gist string // the "just enough" line(s)
    Link string // relative path to the theme file
}

type Root struct {
    Preamble string   // everything before the staging zone
    Staging  string    // raw staging region (entry bodies live here)
    Themes   []Theme   // parsed ## Themes section (empty when un-migrated)
    HasThemes bool     // false = not yet migrated (first run)
}
```

`Staging` stays **raw** in M1: the mover (M2) needs byte-exact bodies,
and re-serializing parsed entries is how content gets silently
normalized. Parsing into blocks uses the existing
`heading.ParseEntryBlocks`; M1 never writes a root back.

## Contracts

New package `internal/disclosure`:

- `Parse(content string, k Kind) (Root, error)` — splits a root into
  preamble / staging / themes. Round-trips: `Preamble+Staging+Themes`
  reconstructs the input byte-for-byte.
- `Validate(r Root) error` — the precondition (spec Guards §2).
- `CheckPairing(root Root, themeDir string) error` — gists ↔ theme files 1:1.
- `CheckUniqueness(root Root, themeDir string) error` — entry in exactly one place.
- `CheckLinks(root Root, ctxDir string) error` — every theme link resolves.

Constants → `internal/config/disclosure`: `HeadingThemes = "## Themes"`,
`HeadingRecent = "## Recent"`.
Errors → `internal/err/disclosure` (constructors in `internal/err` per
CONVENTIONS; sentinels are `entity.Sentinel` consts with text in
`commands/text/errors.yaml` keyed `err.disclosure.*` — never English in
Go).

No CLI surface in M1. No `ctx agent` change (spec Non-Goals).

## Test matrix

| Invariant / rule | Violation attempt | Expected failure | Task |
|---|---|---|---|
| zero-or-one `## Themes` | two `## Themes` headings | `Validate` → `ErrMultipleThemes` | T06 |
| un-migrated root is valid | zero `## Themes` | `Validate` → `nil` (first run) | T06 |
| no `## [` below `## Themes` | entry planted below Themes | `Validate` → `ErrEntryBelowThemes` | T06 |
| staging parses into discrete entries | malformed staging region | `Validate` → `ErrStagingUnparsable` | T06 |
| gists ↔ theme files 1:1 | theme file with no gist | `CheckPairing` → `ErrOrphanThemeFile` | T07 |
| gists ↔ theme files 1:1 | gist with no theme file | `CheckPairing` → `ErrMissingThemeFile` | T07 |
| pairing vacuous when un-migrated | 0 gists, 0 theme files | `CheckPairing` → `nil` | T07 |
| entry in exactly one place | same entry in staging *and* a theme file | `CheckUniqueness` → `ErrDuplicateEntry` | T08 |
| entry in exactly one place | same entry in two theme files | `CheckUniqueness` → `ErrDuplicateEntry` | T08 |
| theme links resolve | link to nonexistent path | `CheckLinks` → `ErrBrokenThemeLink` | T09 |
| add lands in staging (learning) | Themes present, staging populated | new entry index `<` Themes index | T10 |
| add lands in staging (learning) | Themes present, **staging empty** | new entry index `<` Themes index (AfterHeader fallback) | T10 |
| add lands in staging (decision) | both above cases | new entry index `<` Themes index | T11 |
| add lands in staging (convention) | Themes present | new section lands inside `## Recent`, after Themes | T12 |
| invariants tolerate the real tree | run against today's `.context` | vacuous pass | T14 |

## Task breakdown

| id | st | task | deps | files | [P] | acceptance criterion | spec ref |
|---|---|---|---|---|---|---|---|
| T01 | [ ] | Structural vocabulary constants | — | `internal/config/disclosure/disclosure.go` | [P] | `go test ./internal/config/disclosure/` asserts `HeadingThemes == "## Themes"` and `HeadingRecent == "## Recent"` | Design/Layout |
| T02 | [ ] | Error sentinels + i18n text | — | `internal/err/disclosure/*.go`, `commands/text/errors.yaml` | [P] | `go test ./internal/err/...` green; test asserts every `err.disclosure.*` sentinel's `Error()` resolves non-empty via `desc.Text` (no English literal in Go) | CONVENTIONS/Error Handling |
| T03 | [ ] | `Root`/`Theme`/`Kind` types | — | `internal/disclosure/types.go` | [P] | `go build ./...` green; `make audit` type-file report shows 0 violations | Data model |
| T04 | [ ] | `Parse` for entry kinds (LEARNINGS/DECISIONS) | T01,T03 | `internal/disclosure/parse.go` | | `go test ./internal/disclosure/ -run TestParse_EntryKind`: fixture splits into preamble/staging/themes **and** round-trips byte-for-byte | Design/Layout |
| T05 | [ ] | `Parse` for CONVENTIONS kind (`## Themes` then `## Recent`) | T04 | `internal/disclosure/parse.go` | | `-run TestParse_ConventionKind`: round-trips byte-for-byte; `Staging` == the `## Recent` section | Design/Layout |
| T06 | [ ] | `Validate` precondition | T02,T04,T05 | `internal/disclosure/validate.go` | | `-run TestValidate` table test: every matrix row T06 returns its named sentinel; valid + un-migrated return `nil` | Guards §2 |
| T07 | [ ] | Invariant: gists ↔ theme files 1:1 | T04 | `internal/disclosure/invariant.go` | | `-run TestInvariant_Pairing`: orphan file → `ErrOrphanThemeFile`; gist w/o file → `ErrMissingThemeFile`; 1:1 → nil; 0↔0 → nil | Invariants |
| T08 | [ ] | Invariant: entry in exactly one place | T07 | `internal/disclosure/invariant.go` | | `-run TestInvariant_Uniqueness`: dup across staging+theme → `ErrDuplicateEntry`; dup across two themes → `ErrDuplicateEntry`; single → nil | Invariants |
| T09 | [ ] | Invariant: theme links resolve | T08 | `internal/disclosure/invariant.go` | | `-run TestInvariant_Links`: broken link → `ErrBrokenThemeLink`; resolving link → nil | Invariants |
| T10 | [ ] | **Layout proof — LEARNINGS** (populated + empty staging) | T01 | `internal/cli/add/core/insert/layout_learning_test.go` | [P] | `-run TestAdd_LearningLandsAboveThemes`: after `AppendEntry` on a Themes-bearing fixture, `strings.Index(out,newEntry) < strings.Index(out,"## Themes")` — for **both** populated and empty staging | Design/Layout |
| T11 | [ ] | **Layout proof — DECISIONS** (populated + empty staging) | T01 | `internal/cli/add/core/insert/layout_decision_test.go` | [P] | `-run TestAdd_DecisionLandsAboveThemes`: same assertion, both cases | Design/Layout |
| T12 | [ ] | **Layout proof — CONVENTIONS** (`## Recent`) | T01 | `internal/cli/add/core/insert/layout_convention_test.go` | [P] | `-run TestAdd_ConventionLandsInRecent`: after `AppendEntry`, new section index `>` `## Recent` index and `>` `## Themes` index | Design/Layout |
| T13 | [ ] | `doc.go` for `internal/disclosure` | T04–T09 | `internal/disclosure/doc.go` | | `make audit` green (doc.go quality floor: behavior-grounded, ~25–100 body lines, related-packages section) | CONVENTIONS/Documentation |
| T14 | [ ] | Compliance wiring — invariants vs the real tree | T06–T09 | `internal/compliance/disclosure_test.go` | | `go test ./internal/compliance/ -run TestDisclosureInvariants` green on today's tree (vacuous pass); proven-both-ways via a planted-violation temp fixture that fails it | Invariants |
| T15 | [ ] | Milestone gate | T01–T14 | — | | `make lint` = 0 issues; `go test ./...` green; `git diff --stat` shows **no** change to `.context/{LEARNINGS,DECISIONS,CONVENTIONS}.md` | Scope & DoD |

**Execution waves** (topological order): `T01,T02,T03` ∥ → `T04` → `T05`
→ `T06` → `T07` → `T08` → `T09` → `T13,T14` → `T15`.
`T10,T11,T12` are `[P]` from T01 onward (distinct test files, no shared
edges/sequences) and may run alongside any wave after T01.

## Epic anchors (TASKS.md projection)

TASKS.md carries **epic-level anchors only**; this plan is the single
source of truth for milestone progress. The epics **partition** the task
ids — every id in exactly one epic:

| Epic | Range | Count |
|---|---|---|
| E1 Vocabulary, types, errors | T01–T03 | 3 |
| E2 Root parser + validate | T04–T06 | 3 |
| E3 Cross-file invariants | T07–T09 | 3 |
| E4 Layout proofs (add-path de-risking) | T10–T12 | 3 |
| E5 doc.go, compliance wiring, milestone gate | T13–T15 | 3 |

Arithmetic: 3 + 3 + 3 + 3 + 3 = **15** = the task count. No id is
double-counted; no id is unclaimed.

**Completion rule:** an epic is checked `[x]` in TASKS.md only when
every task in its range is `[x]` or `[o]` here.

## Risks & measurement gates

- **⚠️ Measurement gate — T10–T12 are the milestone's load-bearing
  result.** The entire design rests on "`add` needs zero change because
  its anchor always lands above `## Themes`." If any layout proof fails,
  the premise is wrong and the spec's Layout section must be revisited
  **through `/ctx-plan`** — not patched here. Everything downstream
  (M2–M5) reshapes if this fires.
- **CONVENTIONS is the least-verified path.** `AppendAtEnd` was read but
  its exact EOF behavior against a trailing `## Recent` section (and any
  trailing newline handling) is unproven — T12 is where that lands.
- **`Staging` kept raw** to avoid silent normalization; if M2 later needs
  structured staging, that is an amendment, not a redesign.

## Out of scope (deferred, with pointers)

| Deferred | Milestone | Note |
|---|---|---|
| The pass itself; gist authoring | M2 | *Gist format ("just enough") TBD graduates to **blocking** here* |
| Real rollout: LEARNINGS → DECISIONS | M3 | Requires M1 guards green |
| CONVENTIONS prose rollout; edits-behind-a-link UX | M4 | *CONVENTIONS staging-parse TBD graduates to blocking here* |
| Suggest-only trigger wiring | M5 | *Growth-nudge threshold TBD graduates to blocking here* |
| Nesting / taxonomy | — | Deferred indefinitely; structure is self-similar (spec Non-Goals) |
| `ctx drift` wiring beyond the existing path check | — | Spec does not require it |

## Amendments

| date | what | why |
|---|---|---|
| — | — | — |
