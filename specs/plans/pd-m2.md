# pd-m2 Plan — the dry-run pass (inspect + propose, move nothing)

**Spec:** `specs/progressive-disclosure.md` · **Status:** Ready
**Blocking TBDs resolved:** gist format — resolved **in the spec**
(`### Gist format`): one bullet `- <name> — <gist> → [<name>](<noun>/<slug>.md)`,
gist one line ≤~140 chars describing coverage, authored by the pass.

## Scope & DoD

Milestone 2 delivers the digesting pass **in dry-run only**: a read-only
CLI that reports a root's staging entries and current themes, and a
skill that reads it, proposes a theme per staging entry (agent-semantic,
human-overridable), authors proposed gists, and shows the plan of what
*would* move — **moving nothing**. The mover (append→verify→remove, gist
write-back) is M3.

This ordering keeps the risky write path (M3) behind a fully-exercised
read+plan path: the same structured data the mover will consume is built
and tested here first.

DoD (confirmed by measurement or by the user — never derived from task
completion):

- [ ] `make lint` = 0 issues, `go test ./...` green, `make audit` passes
- [ ] `ctx disclosure inspect <file>` reports, for each kind, the staging
      entries and current themes; `--json` emits the same structured
      (E2)
- [ ] Running `ctx disclosure inspect` against any file **writes
      nothing** — the file is byte-identical after (E2 acceptance)
- [ ] The dry-run skill, exercised on a fixture root, produces a
      theme→entries plan with proposed gists and moves/writes nothing
      (E3, verified by driving it)
- [ ] `git diff` shows **no change** to `.context/{LEARNINGS,DECISIONS,
      CONVENTIONS}.md` across the whole milestone

## Data model & storage

No persistence. The CLI is a pure read/report over M1's `disclosure.Parse`.
One new value type for structured output:

```go
// internal/disclosure (types.go)
type Inspection struct {
    Kind    string        // "learning" | "decision" | "convention"
    Staging []StagedEntry // un-digested entries, file order
    Themes  []Theme       // current ## Themes (reused from M1)
}
type StagedEntry struct {
    Timestamp string
    Title     string
}
```

Kind is inferred from the file's basename (LEARNINGS.md → learning, …)
via a new `disclosure.KindFor(basename) (Kind, bool)`.

## Contracts

**Disclosure package additions** (`internal/disclosure`):

- `KindFor(basename string) (Kind, bool)` — maps a canonical filename to
  its Kind; false for anything else.
- `StagedEntries(root Root) []StagedEntry` — the staging zone's entries,
  via `heading.ParseEntryBlocks(root.Staging)`. Empty for conventions
  (no `## [` entries) — an accepted M2 limitation, see Risks.
- `Inspect(content string, k Kind) Inspection` — Parse + assemble the
  Inspection. Total (mirrors Parse).

**CLI** (`ctx disclosure inspect <file> [--json]`):

- New command group `ctx disclosure` (hidden=false), one subcommand
  `inspect`, registered in `internal/bootstrap/group.go`. `Use` strings
  are `config/embed/cmd` constants.
- Reads the file, infers Kind (error if not a canonical knowledge file),
  prints a human summary (kind, N staged, list; M themes, list) or
  `--json` (the Inspection). **Never writes.**
- Kind-inference failure → a typed error in `internal/err/disclosure`.

**Skill** (`internal/assets/claude/skills/ctx-digest/SKILL.md`, shipped):

- Dry-run pass. Steps: run `ctx disclosure inspect --json` on a root →
  propose a theme per staged entry (semantic; human can rename/merge/
  override) → author a proposed gist per theme (spec `### Gist format`)
  → present the plan (theme → its entries → proposed gist + link).
  **Explicitly moves/writes nothing in M2**; ends by naming M3 as the
  apply step. Copilot copy synced via `make sync-copilot-skills`.

No `ctx agent` change (spec Non-Goals). No mover, no write path.

## Test matrix

| Behavior / rule | Attempt | Expected | Task |
|---|---|---|---|
| kind inference | "LEARNINGS.md"/"DECISIONS.md"/"CONVENTIONS.md" | KindFor → (kind, true) | T01 |
| kind inference | "README.md" | KindFor → (_, false) | T01 |
| staged entries listed | entry-kind root with N staged | StagedEntries → N in file order | T02 |
| staged entries: empty staging | migrated root, empty staging | StagedEntries → nil | T02 |
| Inspect assembles | fixture root | kind+staging+themes match Parse | T03 |
| inspect reports | `inspect LEARNINGS.md` on a fixture | human output lists staged + themes | T05 |
| inspect --json | `inspect --json` | valid JSON decodes to Inspection | T06 |
| inspect writes nothing | `inspect` on a real root | file byte-identical after | T07 |
| inspect rejects non-knowledge file | `inspect README.md` | non-zero exit, typed error | T08 |
| dry-run moves nothing | drive the skill on a fixture | fixture byte-identical; plan produced | T12 |

## Task breakdown

| id | st | task | deps | files | [P] | acceptance criterion | spec ref |
|---|---|---|---|---|---|---|---|
| T01 | [x] | `KindFor` filename→Kind | — | `internal/disclosure/kind.go`, `kind_test.go` | [P] | `go test ./internal/disclosure/ -run TestKindFor`: three canonical names map true; others false | Contracts |
| T02 | [x] | `StagedEntry` type + `StagedEntries` | — | `internal/disclosure/types.go`, `inspect.go`, `inspect_test.go` | | `-run TestStagedEntries`: N entries in order; empty staging → nil | Data model |
| T03 | [x] | `Inspection` type + `Inspect` | T01,T02 | `internal/disclosure/types.go`, `inspect.go` | | `-run TestInspect`: kind/staging/themes match a Parse of the same fixture | Contracts |
| T04 | [x] | `Use` constants + kind-inference error sentinel | — | `internal/config/embed/cmd/*.go`, `internal/err/disclosure/*.go`, `commands/text/errors.yaml` | [P] | `go test ./internal/audit/ -run TestDescKeyYAMLLinkage` green (bijection holds) | CONVENTIONS |
| T05 | [x] | `ctx disclosure inspect` group + human output | T03,T04 | `internal/cli/disclosure/**`, `internal/bootstrap/group.go` | | `ctx disclosure inspect <fixture>` prints kind, staged list, theme list | CLI |
| T06 | [x] | `--json` output | T05 | `internal/cli/disclosure/cmd/inspect/run.go` | | `inspect --json <fixture>` output `json.Unmarshal`s into Inspection with expected values | CLI |
| T07 | [x] | write-nothing guarantee | T05 | `internal/cli/disclosure/cmd/inspect/run_test.go` | | integration test: file bytes before == after `inspect` | Scope/DoD |
| T08 | [x] | reject non-knowledge file | T05 | `internal/cli/disclosure/cmd/inspect/run.go`, `run_test.go` | | `inspect README.md` exits non-zero with the kind-inference sentinel | Contracts |
| T09 | [x] | `doc.go` for the CLI group + disclosure additions | T05 | `internal/cli/disclosure/doc.go`, `cmd/**/doc.go` | | `make audit` doc.go + docstring floors pass | CONVENTIONS |
| T10 | [x] | command wiring guard stays green | T05 | — | | `go test ./internal/compliance/ -run TestShippedHooksResolve` and command-tree tests green with the new group | CLI |
| T11 | [ ] | `ctx-digest` skill (dry-run) SKILL.md | T06 | `internal/assets/claude/skills/ctx-digest/SKILL.md` | | skill frontmatter valid (`go test ./internal/assets/... -run Frontmatter`); body describes inspect→propose→gist→plan and states "moves nothing (M2); apply is M3" | Skill |
| T12 | [ ] | **dry-run drive** (measurement) | T11 | scratchpad fixture | | drive the skill on a fixture root: a theme→entries plan with gists is produced AND the fixture is byte-identical after | Scope/DoD |
| T13 | [ ] | copilot skill sync | T11 | `internal/assets/integrations/copilot-cli/skills/ctx-digest/**` | | `make check-copilot-skills` green | Skill |
| T14 | [ ] | milestone gate | T01–T13 | — | | `make lint` 0, `go test ./...` green, `make audit` pass; canonical `.md` byte-identical | Scope/DoD |

**Execution waves:** `T01,T02,T04` ∥ → `T03` → `T05` → `T06,T07,T08` →
`T09,T10` → `T11` → `T12,T13` → `T14`.

## Epic anchors (TASKS.md projection)

Epics **partition** the task ids:

| Epic | Range | Count |
|---|---|---|
| E1 Disclosure inspect model | T01–T04 | 4 |
| E2 `ctx disclosure inspect` CLI | T05–T10 | 6 |
| E3 Dry-run skill | T11–T14 | 4 |

Arithmetic: 4 + 6 + 4 = **14** = the task count. No id double-counted or
unclaimed. **Completion rule:** an epic is `[x]` in TASKS.md only when
every task in its range is `[x]`/`[o]` here.

## Risks & measurement gates

- **⚠️ T12 is the milestone's measurement gate.** If the dry-run cannot
  produce a coherent theme→entries plan from the inspect output, the
  skill's design is wrong and should be reconsidered before M3 builds the
  mover on top of it.
- **CONVENTIONS staging is empty under the entry-based model.** Real
  CONVENTIONS uses `##` sections, not `## [` entries, so `StagedEntries`
  returns nil for it — inspect reports zero staged conventions. This is
  the same `##`-vs-`###` gap flagged in pd-m1's amendments; conventions
  digestion is genuinely an **M4** problem. M2 inspect still works for
  conventions (kind + themes), it just lists no staged entries.
- **New user-facing CLI surface.** `ctx disclosure` is a new command;
  the out-of-band `_ctx-surface-audit` will want docs/cli coverage —
  tracked for the release pass, not M2 (spec adds no docs requirement).

## Out of scope (deferred, with pointers)

| Deferred | Milestone | Note |
|---|---|---|
| The mover: append→verify→remove, gist write-back | M3 | consumes the same Inspection this milestone builds |
| First real rollout LEARNINGS→DECISIONS | M3 | requires the mover |
| CONVENTIONS digestion (`##` model) | M4 | *`##`-vs-`###` TBD is blocking here* |
| Theme taxonomy / naming discipline | M3 | agent proposes at runtime in M2 dry-run; matters when files are created (M3) |
| Suggest-only trigger wiring | M5 | *nudge-threshold TBD blocking here* |

## Amendments

| date | what | why |
|---|---|---|
| — | — | — |
