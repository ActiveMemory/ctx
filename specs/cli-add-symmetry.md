---
title: CLI add Symmetry — Noun-First Canonical Form
status: proposed
date: 2026-05-03
owner: jose
scope: behavioral — CLI surface, slash-skill alignment, docs sweep
related:
  - specs/cli-namespace-cleanup.md
  - specs/ctx-init-overwrite-safety.md
---

# Spec: CLI add Symmetry

## Problem

The CLI today mixes two organizing principles:

| Operation               | Form         | Lives under  |
|-------------------------|--------------|--------------|
| `ctx task complete`     | noun-first   | `ctx task`   |
| `ctx task archive`      | noun-first   | `ctx task`   |
| `ctx task snapshot`     | noun-first   | `ctx task`   |
| `ctx decision reindex`  | noun-first   | `ctx decision` |
| `ctx learning reindex`  | noun-first   | `ctx learning` |
| `ctx task add`          | **missing**  | (must use `ctx add task`) |
| `ctx decision add`      | **missing**  | (must use `ctx add decision`) |
| `ctx learning add`      | **missing**  | (must use `ctx add learning`) |
| `ctx convention add`    | **missing**  | (`ctx convention` doesn't exist at all — must use `ctx add convention`) |

Every other verb in the artifact namespace is noun-first. Only
`add` lives under a verb-first parent. A user who learns `ctx
task complete` reasonably expects `ctx task add` to exist. It
doesn't, and the asymmetry is invisible until they hit
"unknown command".

The asymmetry is also self-inflicted from the slash-skill side:
the user-facing skills are already noun-first
(`ctx-task-add`, `ctx-decision-add`, `ctx-learning-add`,
`ctx-convention-add`). The CLI is the odd one out.

The 2026-05-03 incident that surfaced this: assistant referenced
`/ctx-decision-add` and `/ctx-learning-add` (the skill names);
user reasonably checked `ctx decision add` / `ctx task add` on
the CLI; user found nothing and concluded the schema-enforcing
commands had been deleted. They hadn't — they live at `ctx add
<noun>` — but the dissonance between CLI naming and
skill/docstring naming is real and recurring.

Stale documentation in three places already implies the noun-first
form is canonical:

- `internal/write/add/doc.go:8-9` — "entry addition commands
  (ctx task add, ctx decision add, ctx learning add, ctx
  convention add)"
- `specs/explicit-context-dir.md:267` — `ctx task add "foo"`
- 2026-04-19 journal entry — `ctx task add`

Three sources of drift in the same direction is not a typo
pattern; it's an organizing intuition that the CLI never matched.

## Goals

1. Single, consistent organizing principle: **noun-first** for
   every operation on a context artifact.
2. Make CLI naming match slash-skill naming
   (`ctx-<noun>-add` ↔ `ctx <noun> add`).
3. Remove the verb-first `ctx add` parent (or reduce it to a
   thin deprecation alias — see Open Questions).
4. Ensure `ctx <noun> --help` lists *every* operation on that
   noun, including `add`. No "must use a different parent"
   surprises.

## Non-goals

- Changing the schema (`--context`, `--rationale`,
  `--consequence`, `--lesson`, `--application`, etc. unchanged).
- Changing the underlying file formats.
- Renaming slash skills (already in the target form).
- Touching unrelated noun-first asymmetries elsewhere in the CLI
  (separate audit, separate spec).

## Design

### A. New canonical commands

Create the four noun-first add subcommands, each backed by the
existing `internal/cli/add/core/` machinery (which becomes
shared library, not a CLI parent's private core):

```
ctx task add        → internal/cli/task/cmd/add/
ctx decision add    → internal/cli/decision/cmd/add/
ctx learning add    → internal/cli/learning/cmd/add/
ctx convention add  → internal/cli/convention/cmd/add/  (new parent)
```

Flag set, validation, error messages, and exit codes match the
current `ctx add <noun>` exactly. Each new subcommand is a thin
adapter that calls the shared core.

### B. Create `ctx convention` parent

`ctx convention` is currently not a registered command (`ctx
convention` returns "unknown command, did you mean
'connection'?"). The flip requires a `ctx convention` parent
that hosts `add`.

**Decided: ship `add` only.** CONVENTIONS.md doesn't have an
INDEX block today, so `reindex` would mean designing one — out
of scope for this spec. Initial subcommand set:

```
ctx convention add  (port from `ctx add convention`)
```

If CONVENTIONS.md ever grows an index, add `ctx convention
reindex` as a follow-up. The `ctx convention` parent stands on
its own with one subcommand for now.

### C. Refactor shared core

`internal/cli/add/` is currently both:
- A CLI parent package (registers `ctx add`, owns
  `cmd/root/cmd.go`).
- A shared library (`core/append.go`, `core/insert.go`,
  `core/index.go`, etc., consumed by the parent's run.go).

After the flip, the CLI-parent role goes away (or shrinks to a
deprecation alias). The shared library should keep working from
the same import path so the new noun-first commands import it
cleanly.

Two options:

1. **In place**: keep `internal/cli/add/core/` as the shared
   library; new noun commands import it directly. The
   `internal/cli/add/cmd/root/` shrinks to either nothing
   (option D1) or a thin deprecation shim (option D2).
2. **Relocate**: move `internal/cli/add/core/` to
   `internal/add/` (no `cli/` prefix, signaling it's
   non-CLI shared logic). New noun commands import the new
   path. `internal/cli/add/` package goes away entirely (option
   D1) or becomes the deprecation shim (option D2).

Recommendation: **option 1** for the migration commit (low
churn, no import-path-rewrite blast); follow up with option 2
in a separate commit if `internal/cli/add/` ends up genuinely
empty after deprecation.

### D. Deprecation strategy: hard cut

**Decided: hard cut.** Remove `ctx add <noun>` in the same
commit that adds `ctx <noun> add`. Matches the project's prior
namespace-cleanup pattern (`78fbdf7d` and `f4117b87` were both
hard renames). Release-notes call-out plus a one-line migration
sed in the breaking-change section is the only mitigation
needed.

After the flip, `ctx add` either disappears entirely or returns
"unknown command" with a Cobra `did you mean` suggestion
pointing at the new noun-first form.

### E. Slash-skill alignment

The four `ctx-<noun>-add` skills under
`internal/assets/claude/skills/` already use noun-first naming
externally, but they likely invoke `ctx add <noun>` internally.
After the flip:

- Audit each skill's SKILL.md / scripts for the verb-first
  invocation; replace with noun-first.
- If D2 is chosen, the skills can be updated immediately (the
  new form is preferred); if D1, they MUST be updated in the
  same commit.

### F. Documentation sweep

20+ files reference the verb-first form (from grep across
`docs/`, `internal/assets/`):

- `docs/home/{common-workflows,repeated-mistakes,first-session,context-files,is-ctx-right,joining-a-project}.md`
- `docs/cli/{connection,context,connect}.md`
- `docs/operations/{autonomous-loop,runbooks/sanitize-permissions,runbooks/hub-deployment}.md`
- `docs/blog/{2026-02-17-context-as-infrastructure,2026-02-03-the-attention-budget,2026-02-01-ctx-v0.2.0-the-archaeology-release,2026-02-07-the-anatomy-of-a-skill-that-works}.md`
- `docs/recipes/{task-management,hub-team,external-context,hub-personal}.md`
- `internal/write/add/doc.go` (the godoc that already implied
  noun-first — becomes correct after the flip)
- `specs/explicit-context-dir.md` (line 267)
- AGENT_PLAYBOOK.md and CLAUDE.md template snippets (audit)

**Blog posts**: do not rewrite history. Add an editor's note
("CLI surface changed in v0.X — see [link]") rather than
silently rewriting published commands.

### G. Tests

- Cobra-tree test: assert each of the four nouns has an `add`
  subcommand; assert flag set matches the old `ctx add <noun>`
  flag set exactly (catch flag drift during the migration).
- If D2: alias test — `ctx add <noun> ...` still produces the
  same output as `ctx <noun> add ...` and prints the deprecation
  notice on stderr.
- If D1: removed-command test — `ctx add <noun>` returns
  "unknown command" with a helpful "did you mean `ctx <noun>
  add`?" suggestion (Cobra's built-in fuzzy match should handle
  this; verify it does).
- Drift test: lint-drift or doc-drift check for the
  `internal/write/add/doc.go` claim, so it stays accurate after
  the flip.

## Open questions

1. **Shared-core relocation (option 1 vs option 2 in section
   C)**: in-place for the migration commit, follow-up move
   later? Default proposal: in-place now, separate refactor
   commit later if warranted.
2. **"Did you mean" text**: Cobra's default is acceptable but
   a custom one-liner explaining the principle is friendlier:
   "operations on a noun live under that noun: `ctx <noun>
   <verb>`". Resolve at implementation time; minor.

## Decisions baked in

- **Deprecation horizon**: hard cut (D1). `ctx add <noun>`
  removed in the same commit that adds `ctx <noun> add`.
- **`ctx convention reindex`**: deferred. Ship `ctx convention
  add` only.

## Risks

- **Breaking change without deprecation period (D1)**: any
  external user/CI calling `ctx add <noun>` breaks. Mitigation:
  release notes call this out explicitly; provide a one-line
  migration sed.
- **Forgotten doc reference**: 20+ docs is a lot. Mitigation:
  the drift test in section G; also a final `rg "ctx add "
  --type md` sweep gated in CI for the release commit.
- **Slash-skill regression**: skills currently work; the audit
  in section E catches the invocation drift, but a missed skill
  silently breaks for users invoking it. Mitigation: explicit
  test that each `ctx-<noun>-add` skill resolves to a working
  CLI invocation post-flip.
- **`ctx convention` parent name collision**: low risk;
  `convention` is unused as a top-level today.

## Tasks

Implementation order: shared core stays put → noun commands
land → deprecation choice executed → docs swept →
release notes.

### Phase A: New noun commands (CLI surface)

- [ ] Add `ctx task add` subcommand under
      `internal/cli/task/cmd/add/`; thin adapter calling
      `internal/cli/add/core/`
- [ ] Add `ctx decision add` subcommand under
      `internal/cli/decision/cmd/add/`
- [ ] Add `ctx learning add` subcommand under
      `internal/cli/learning/cmd/add/`
- [ ] Create `ctx convention` parent command in
      `internal/cli/convention/cmd/root/`; add `ctx convention
      add` subcommand under `internal/cli/convention/cmd/add/`
- [ ] Verify each new subcommand's flag set matches the old
      `ctx add <noun>` exactly (Cobra-tree test)

### Phase B: Hard-cut removal of `ctx add <noun>`

- [ ] Remove `ctx add` parent registration from the root
      command tree
- [ ] Remove `internal/cli/add/cmd/root/` (CLI-parent role
      gone); keep `internal/cli/add/core/` as shared library
- [ ] Verify Cobra returns a useful "unknown command" with a
      "did you mean `ctx <noun> add`?" suggestion; if Cobra's
      default is unhelpful, add a custom suggestion
- [ ] Removed-command test: `ctx add task` returns non-zero
      with the suggestion in stderr

### Phase C: Slash-skill alignment

- [ ] Audit `internal/assets/claude/skills/ctx-task-add/`
      contents (SKILL.md, any helper scripts); replace any
      `ctx add task` invocation with `ctx task add`
- [ ] Same for `ctx-decision-add`, `ctx-learning-add`,
      `ctx-convention-add`
- [ ] Add a test that each skill's invocation resolves on the
      post-flip CLI

### Phase D: Documentation, docstrings, comments sweep (full surface)

This is the "as usual" sweep. Source-tree comments and godocs
matter as much as user-facing docs — the
`internal/write/add/doc.go` drift is exactly what fooled both
user and assistant in the incident that prompted this spec.

- [ ] **Source-tree godocs and comments** — `rg "ctx add
      (decision|task|learning|convention)\b" --type go`;
      every hit needs review. Notably:
  - `internal/write/add/doc.go:8-9` (already noun-first;
    confirm it's still accurate post-flip)
  - Any `// Example: ctx add ...` comments
  - Any cobra `Example:` field strings in
    `internal/cli/**/cmd/root/cmd.go`
  - Embed-level help text in
    `internal/config/embed/text/**` if any references exist
- [ ] **Spec files** — `specs/explicit-context-dir.md:267`
      and any other spec hits from `rg "ctx add " specs/`
- [ ] **User-facing docs** — files identified in section F:
      `docs/home/{common-workflows,repeated-mistakes,first-session,context-files,is-ctx-right,joining-a-project}.md`,
      `docs/cli/{connection,context,connect}.md`,
      `docs/operations/{autonomous-loop,runbooks/sanitize-permissions,runbooks/hub-deployment}.md`,
      `docs/recipes/{task-management,hub-team,external-context,hub-personal}.md`
- [ ] **Generated CLI reference** — if `docs/cli/index.md` or
      similar is generated from Cobra metadata, regenerate.
      Otherwise update by hand.
- [ ] **Blog posts** — `docs/blog/{2026-02-17-context-as-infrastructure,2026-02-03-the-attention-budget,2026-02-01-ctx-v0.2.0-the-archaeology-release,2026-02-07-the-anatomy-of-a-skill-that-works}.md`
      — do NOT rewrite the body; add an editor's note at the
      top: "*Editor's note (YYYY-MM-DD): the CLI surface
      changed in v0.X. `ctx add <noun>` is now `ctx <noun>
      add`; see [link to release notes].*" Published commands
      stay as published.
- [ ] **CLAUDE.md and AGENT_PLAYBOOK templates** — audit
      `internal/assets/claude/CLAUDE*.md`,
      `AGENT_PLAYBOOK*.md`, the project root CLAUDE.md, and
      any embedded snippets. Replace verb-first invocations
      with noun-first.
- [ ] **README and top-level docs** — `README.md`, any
      `CHEAT-SHEETS.md`, `GLOSSARY.md` in `.context/` if
      they reference the verb-first form.
- [ ] **CI drift gate** — add a check:
      `rg "ctx add (decision|task|learning|convention)\b"
      --type md --type go` returns zero hits *excluding*
      `docs/blog/` (where the editor's notes are intentional
      and the body is preserved). If there's an existing
      drift-test infrastructure (`internal/audit/`,
      `lint-drift`), wire this into it; otherwise add a
      standalone test.
- [ ] **Release notes** — breaking-change section with the
      one-line migration sed:
      `find . -type f \( -name '*.md' -o -name '*.go' \) -exec sed -i.bak -E 's/ctx add (decision|task|learning|convention)/ctx \1 add/g' {} \;`

### Phase E: Verification

- [ ] `make lint && make test` clean
- [ ] Manual smoke: `ctx task add "foo" --session-id ...
      --branch ... --commit ...` succeeds
- [ ] Manual smoke: same for decision, learning, convention
- [ ] Manual smoke: `ctx <noun> --help` lists `add` for every
      noun
- [ ] If D1: manual smoke that `ctx add task` returns "unknown
      command" with the "did you mean `ctx task add`?"
      suggestion
- [ ] If D2: manual smoke that `ctx add task` still works and
      prints the deprecation notice on stderr
