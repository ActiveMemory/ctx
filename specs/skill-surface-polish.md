---
title: Skill Surface Polish (Phase SK)
date: 2026-05-10
status: ready
owner: jose
scope: CLI flag enforcement + skill file edits + doc additions
design-ref: ideas/002-editorial-pipeline-and-skill-rigor.md §3 "Reframing the wishy-washy skills"
phase: SK (prerequisite for Phase KB; independent of Phase RG)
---

# Spec: Skill Surface Polish (Phase SK)

Tightens the existing capture skills to sibling-project rigor before
the editorial pipeline (Phase KB) lifts that pattern wholesale.
Independent of Phase RG; both can ship in parallel.

## Problem

Three concrete gaps in the current capture surface:

1. **`ctx decision add` and `ctx learning add` accept missing body
   fields.** Today the CLI binds `--context`, `--rationale`,
   `--consequence` (decision) and `--context`, `--lesson`,
   `--application` (learning) as optional strings. Users (and
   agents) can submit decisions with no rationale at all, or with
   placeholder values like `TBD` or `see chat`. The captured entry
   is then lower-fidelity than the skill ceremony implies.

2. **`/ctx-spec` skill has no `--brief <path>` flag.** The sibling
   project ships a path-required spec contract that reads an
   authoritative brief and skips the fresh-template Q&A. Without it,
   every spec session re-runs interview questions even when the
   brief already exists.

3. **`/ctx-plan` skill does not offer to persist the debated brief.**
   Adversarial-interview output evaporates at session end. The
   sibling writes it to `.context/briefs/<TS>-<slug>.md` as the next
   spec's authoritative source.

Plus three smaller alignment items:

4. Skill files lack an "Authority boundary (vs other skills)"
   section, allowing silent promotion (handover→decision,
   learning→convention) without explicit user ask.
5. "Light compression for clarity is allowed; new facts are not"
   wording is inconsistent across capture skills.
6. The `--brief` contract is not documented in `docs/skills.md`.

## Approach

### CLI changes

Add a small helper at `internal/cli/add/core/validate/required.go`:

```go
// RequireBodyFlags marks the listed flags as cobra-required AND
// rejects placeholder values via a PreRunE wrapper. Placeholder
// values are matched case-insensitively against a closed set
// (TBD, see chat, n/a, etc.) plus whitespace-only.
func RequireBodyFlags(c *cobra.Command, flags ...string) error
```

Wire it in:

- `internal/cli/decision/cmd/add/cmd.go` — require `--context`,
  `--rationale`, `--consequence` after `build.Cmd(...)`.
- `internal/cli/learning/cmd/add/cmd.go` — require `--context`,
  `--lesson`, `--application` after `build.Cmd(...)`.

Cobra's `MarkFlagRequired` returns "required flag(s) not set"
errors and exits non-zero. Placeholder rejection returns a typed
error explaining which flag and which value pattern matched.

### Skill changes

- `/ctx-spec` (`internal/assets/claude/skills/ctx-spec/SKILL.md`):
  add `--brief <path>` flag. When present, the skill reads the
  file as authoritative source per the authority order
  (frozen contracts > recorded decisions > debrief > agent
  inference labeled `TBD`); skips the fresh template Q&A.

- `/ctx-plan` (`.../ctx-plan/SKILL.md`): always offer to write the
  debated brief to `.context/briefs/<TS>-<slug>.md` at the end of
  an adversarial-interview pass (create `.context/briefs/` if
  absent).

- `/ctx-decision-add`, `/ctx-learning-add`, `/ctx-task-add`,
  `/ctx-convention-add` (`.../ctx-*/SKILL.md`): add
  "Authority boundary (vs other skills)" section explicitly listing
  what each skill does **not** promote without explicit user ask.

- Capture skills (decide / learn primarily; handover later in
  Phase KB): standardize on the wording
  *"Light compression for clarity is allowed; new facts are not."*

### Documentation

- `docs/skills.md` — document the `--brief` contract.

## Behavior

### Happy path

- `ctx decision add --context X --rationale Y --consequence Z body` —
  works as today.
- `ctx decision add body` — exits non-zero with cobra's
  "required flag(s) ... not set" message.
- `ctx decision add --context TBD --rationale Y --consequence Z body` —
  exits non-zero with placeholder-rejection error.
- `/ctx-spec --brief ideas/003-foo.md` — reads the brief and skips
  template Q&A.

### Edge cases

| Case | Expected |
|------|----------|
| `--context " "` (whitespace-only) | placeholder error |
| `--context "see chat"` | placeholder error |
| `--context "TBD"` (exact case) | placeholder error |
| `--context "tbd"` (lowercase) | placeholder error |
| `--context "rationale to be defined later"` (substring) | accepted — only exact-value matches rejected, not substrings |
| `/ctx-spec --brief nonexistent.md` | skill reports file-not-found and stops |
| `/ctx-plan` ending with no edits | still offers to write the brief |

## Interface

### CLI

`ctx decision add` — same flags as today; now three are required
and reject placeholder values.

`ctx learning add` — same.

`/ctx-spec` skill — new flag-style argument `--brief <path>`.

### Helper

```go
// in internal/cli/add/core/validate/
func RequireBodyFlags(c *cobra.Command, flags ...string) error
```

## Implementation

### Files to create

| File | Purpose |
|------|---------|
| `internal/cli/add/core/validate/required.go` | `RequireBodyFlags` helper |
| `internal/cli/add/core/validate/doc.go` | package doc |
| `internal/cli/add/core/validate/required_test.go` | unit tests |

### Files to modify

| File | Change |
|------|--------|
| `internal/cli/decision/cmd/add/cmd.go` | call `RequireBodyFlags` |
| `internal/cli/learning/cmd/add/cmd.go` | call `RequireBodyFlags` |
| `internal/cli/decision/cmd/add/add_test.go` | tests for required + placeholder rejection |
| `internal/cli/learning/cmd/add/add_test.go` | same |
| `internal/assets/claude/skills/ctx-spec/SKILL.md` | `--brief` contract |
| `internal/assets/claude/skills/ctx-plan/SKILL.md` | brief-save offer |
| `internal/assets/claude/skills/ctx-decision-add/SKILL.md` | authority boundary |
| `internal/assets/claude/skills/ctx-learning-add/SKILL.md` | authority boundary |
| `internal/assets/claude/skills/ctx-task-add/SKILL.md` | authority boundary |
| `internal/assets/claude/skills/ctx-convention-add/SKILL.md` | authority boundary |
| `docs/skills.md` | document `--brief` contract |

### Placeholder set

Exact, case-insensitive matches plus whitespace-only:

- `TBD`
- `tbd`
- `n/a`, `N/A`, `na`
- `see chat`
- `see above`
- `pending`
- whitespace-only (regex `^\s*$`)

Substring matches are NOT placeholders (a legitimate value can
contain the word "TBD" inside a longer sentence).

## Testing

- Unit: `RequireBodyFlags` rejects each placeholder; accepts
  legitimate strings; accepts strings that contain placeholder
  substrings.
- Unit: cobra-level required-flag error fires when the flag is
  missing entirely.
- Integration: `ctx decision add` + `ctx learning add` end-to-end
  with each failure mode.
- No regression in existing add tests.

## Non-Goals

- Changing flag names, short forms, or ordering.
- Requiring body flags on `ctx task add` or `ctx convention add`
  (they have different field semantics).
- Implementing Phase KB closeout/fold logic — that's tracked under
  `specs/kb-editorial-pipeline.md`.
- Auto-rewriting old entries that already contain placeholder
  values — only new entries are gated.

## Open Questions

None — all design decisions are pinned in `ideas/002` §3 and the
debated brief at `ideas/003`.
