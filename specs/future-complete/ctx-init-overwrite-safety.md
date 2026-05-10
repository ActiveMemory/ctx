---
title: ctx init Overwrite Safety
status: proposed
date: 2026-05-03
owner: jose
scope: behavioral — `ctx init` UX, hook scaffolding side effects, docs
related:
  - specs/single-source-context-anchor.md
---

# Spec: ctx init Overwrite Safety

## Problem

On 2026-04-25 01:44 PDT, commit `423016c3` ("feat: enforce explicit
CTX_DIR with single-source anchor resolution") replaced the live
contents of four `.context/` files with the embedded scaffolding
templates from `internal/assets/context/`:

| File           | Before | After (template) |
|----------------|-------:|-----------------:|
| TASKS.md       |   2138 |               36 |
| DECISIONS.md   |   1801 |               49 |
| LEARNINGS.md   |   1667 |               20 |
| CONVENTIONS.md |    272 |               62 |

Total: 5878 lines of curated context destroyed in a single commit
with no warning and no backup. Recovery was possible only because
the prior versions were preserved in git — had this happened to a
user without git (or had it happened between commits), the loss
would have been irreversible.

## Root cause

Two compounding defects in `ctx init`:

### Defect 1 — silent overwrite on "y"

`internal/cli/initialize/cmd/root/run.go` (lines 117-136) prompts
`"Overwrite existing context?" [y/N]` when essential files exist.
The prompt:

- Does not enumerate the files that will be overwritten.
- Does not show their current line counts or sizes.
- Does not warn that templates will replace curated content.
- Does not create a backup before overwriting.
- With `--force`, skips the prompt entirely.

The word "overwrite" is doing all the rhetorical work. A user
treating `init` as idempotent (a reasonable mental model — most
`init` commands are no-ops on initialized state) sees a binary
prompt and answers `y` without realizing they are about to delete
thousands of lines of decisions, learnings, and tasks.

### Defect 2 — phantom `.context/` from misbehaving hooks

User-reported, hypothesis to be verified during implementation:
some hook in an empty project fires before context exists, fails
to see `.context/`, and instead of bailing **creates** one and
writes a session-anchor JSONL there (the `jsonl-path-<sessionID>`
file already noted in LEARNINGS as the
`filepath.Join('', rel)` trap).

Symptom: a fresh project with no real context shows up with a
`.context/` containing only a `jsonl-path-*` file. If the
overwrite-safety fix from Defect 1 ships as "refuse
unconditionally when `.context/` is populated," then `ctx init`
on a fresh project will fail because a phantom `.context/`
already exists.

The two defects are coupled: Defect 2 must be fixed first or the
Defect 1 fix backfires.

## Goals (in order)

1. Make it impossible for `ctx init` to silently destroy curated
   context. Refusing-then-erroring is preferable to
   prompting-then-overwriting.
2. Eliminate the phantom-`.context/` hook bug so the refuse
   policy doesn't punish first-time users.
3. Preserve a recovery escape hatch for the rare case where a
   user genuinely wants to nuke their context (`--reset` or
   similar, with on-disk backup).

## Non-goals

- Restoring already-lost context (handled out-of-band via git
  recovery for this incident; users without git are outside this
  spec).
- Reworking the templates themselves.
- Changing the structure of `.context/` files.

## Design

### A. Make overwrite the impossible default

Change `ctx init` to follow this decision tree:

```
.context/ does not exist          → create + scaffold (current behavior)
.context/ exists, no populated    → scaffold missing files only
files (only state/ or empty)
.context/ exists, files populated → REFUSE with explanatory error
                                     pointing at `ctx init --reset`
```

"Populated" = any of TASKS.md / DECISIONS.md / LEARNINGS.md /
CONVENTIONS.md / CONSTITUTION.md exceeds the embedded template's
line count by more than a small tolerance, OR contains content
not present in the template (cheaper check: SHA mismatch with the
embedded asset).

Remove `--force` as an overwrite mechanism. Replace with
`--reset`, which:

- Requires interactive confirmation that names every populated
  file by basename and shows its line count.
- Writes a timestamped backup to `.context/.backup-init-<ISO>/`
  containing every file it is about to replace, and prints the
  backup path before proceeding.
- In non-interactive mode (`--caller`, no TTY, CI), refuses
  unconditionally — there is no scripted use case for "reset my
  context".

The `--reset` name is intentionally not `--force`: `force` reads
as "ignore my safety rails," `reset` reads as "I want a clean
slate." The semantic difference matters at the call site.

### B. Stop hooks from materializing `.context/`

Audit every hook in `internal/assets/hooks/` and every code path
that writes to `state.Dir()` / `rc.ContextDir()`. Any path that
creates a directory or writes a file under `.context/` must:

1. First call `state.Initialized()` (or equivalent) and bail
   silently if false.
2. Never call `os.MkdirAll(stateDir, ...)` as a side effect of
   reading. Read paths use `os.Stat` and return cleanly when the
   directory does not exist.

Specifically: trace the `jsonl-path-<sessionID>` writer (likely
in the session-anchor relay or journal-import hook) and confirm
it gates on `state.Initialized()`. The April LEARNINGS entry on
`filepath.Join('', rel)` closed the CWD-relative trap but did
not prove the writer is gated against bootstrapping a
`.context/` from scratch.

### C. Update user-facing surface

Every place that documents or invokes `ctx init` needs the new
contract:

- `ctx init --help` text and the embedded usage example.
- `docs/cli/init.md` (or wherever the canonical init docs live).
- `docs/recipes/activating-context.md` and any recipe that walks
  a user through first-run setup.
- `https://ctx.ist/recipes/activating-context/` (mirror in the
  docs site).
- `CLAUDE.md` template snippets that mention `ctx init`.
- AGENT_PLAYBOOK.md if it references `--force` anywhere.
- Any `_ctx-release` / release-notes scaffolding that lists
  flags.
- The error message in `rc.ErrDirNotDeclared` and the
  "Overwrite existing context?" prompt removal.

### D. Tests

- Unit: `init` against a populated `.context/` returns
  `ErrContextPopulated` (new sentinel), exit code != 0, no
  files modified.
- Unit: `init --reset` in non-interactive mode returns
  `ErrResetRequiresTTY`, exit code != 0, no files modified.
- Unit: `init --reset` in interactive mode with `n` answer:
  no files modified, no backup directory created.
- Unit: `init --reset` with `y` answer: backup directory
  exists with originals, target files match templates.
- Integration: hook-fired `ctx ...` invocation in a
  non-initialized project does not create `.context/`. (Add a
  test that runs every relay hook against an empty CWD and
  asserts no `.context/` materializes.)

## Risks

- **Refuse-too-aggressive**: edge case where `.context/` exists
  with one populated file (e.g., user manually created
  CONSTITUTION.md before running init). Mitigation: scaffold
  *missing* files in this case rather than refusing wholesale.
- **`--reset` rediscovers the original footgun**: if the
  confirmation prompt is too quiet, we recreate Defect 1 under a
  new flag name. Mitigation: file-by-file enumeration with line
  counts in the prompt, plus mandatory backup.
- **Hook audit misses a writer**: a future hook adds a
  `MkdirAll(stateDir)` and reintroduces phantom `.context/`.
  Mitigation: the integration test in section D catches this on
  CI; consider an AST audit (`internal/audit/`) that flags
  `MkdirAll` calls outside of `init`/`activate`.

## Tasks

Implementation order matters: B before A so the refuse policy
doesn't strand fresh-project users.

### Phase B: Stop phantom .context/ creation

- [ ] Audit `internal/assets/hooks/` — list every hook that
      reads or writes under `.context/` and its current
      `state.Initialized()` gate status
- [ ] Locate the writer that produces `jsonl-path-<sessionID>`
      files; confirm it bails when `.context/` does not exist
      (not just when `CTX_DIR` is unset)
- [ ] Add an integration test: run each relay hook in a temp
      CWD with no `.context/`; assert no directory or file is
      created
- [ ] If the AST audit route is taken: add
      `internal/audit/no_mkdir_in_hooks_test.go` that flags
      `os.MkdirAll` / `os.Mkdir` calls in
      `internal/assets/hooks/` and `internal/cli/*/hook/`

### Phase A: Refuse-by-default in ctx init

- [ ] Add `ErrContextPopulated` sentinel in
      `internal/cli/initialize/...` (or `internal/err/`)
- [ ] Replace the `[y/N]` overwrite prompt in
      `internal/cli/initialize/cmd/root/run.go` (lines 117-136)
      with the populated-check + refuse-with-error path
- [ ] Implement "scaffold missing files only" branch for the
      partial-init case
- [ ] Remove `--force` flag (lines 40, 63-71 in
      `internal/cli/initialize/cmd/root/cmd.go`); add
      `--reset` flag with the contract from section A
- [ ] `--reset` writes timestamped backup to
      `.context/.backup-init-<ISO>/` before any destructive op;
      print the backup path to stderr
- [ ] `--reset` refuses in non-interactive mode
      (`!isatty(stdin)` or `--caller` set) with
      `ErrResetRequiresTTY`
- [ ] All unit tests from section D pass

### Phase C: Docs and recipes

- [ ] Update `ctx init --help` text and embedded examples
- [ ] Update `docs/cli/init.md`
- [ ] Update `docs/recipes/activating-context.md` and any other
      recipe under `docs/recipes/` that mentions `ctx init`
- [ ] Cross-check the docs-site source (zensical.toml entries,
      mkdocs nav) — same change must reach
      https://ctx.ist/recipes/activating-context/
- [ ] Audit `internal/assets/claude/CLAUDE*.md` and
      `AGENT_PLAYBOOK*.md` for mentions of `ctx init --force`;
      replace with `ctx init --reset` (and update the mental
      model accordingly: this is now a destructive op, not a
      shortcut)
- [ ] Release-notes entry: breaking change, `--force` removed,
      replaced by `--reset` with stricter semantics
- [ ] Add a LEARNINGS.md entry post-implementation describing
      the incident, the fix, and the AST-audit guard if added

### Phase D: Verification

- [ ] `make lint && make test` clean
- [ ] Manual smoke: fresh empty project → `ctx init` succeeds
- [ ] Manual smoke: populated project → `ctx init` refuses with
      a useful error and exit code != 0
- [ ] Manual smoke: populated project → `ctx init --reset`,
      decline → no changes; accept → backup exists, files
      reset
- [ ] Manual smoke: `ctx init --reset < /dev/null` → refused,
      exit code != 0
