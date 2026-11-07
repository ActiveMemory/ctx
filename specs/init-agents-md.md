# `ctx init` writes AGENTS.md; CLAUDE.md imports it

## Problem

`ctx init` writes the agent bootstrap instructions (run
`ctx system bootstrap`, read `AGENT_PLAYBOOK.md`, the "do you
remember?" readback, hook authority, handovers, the KB workflow) only
into `CLAUDE.md`. `AGENTS.md` is the vendor-neutral instruction file:
stewarded by the Agentic AI Foundation under the Linux Foundation
since 2025-12-09 and read natively by Codex, GitHub Copilot, Cursor,
Gemini CLI, Jules, OpenCode, Pi, Zed, Aider and others. Claude Code
2.1.277 (2026-09-18) reads it too, but only when no `CLAUDE.md`
(or `.claude/CLAUDE.md`, `CLAUDE.local.md`) exists in the working
directory or above it; when one exists, Claude reads `CLAUDE.md` only,
unless that file imports `AGENTS.md` with an `@AGENTS.md` line.

Two consequences today:

- A Codex, Cursor or Copilot user who runs `ctx init` gets no file
  their agent reads.
- `ctx setup agents|codex|opencode|pi` deploy `AGENTS.md` from a
  second, older template (`internal/assets/integrations/agents.md`):
  no bootstrap step, no handovers, no KB workflow, and it still points
  at `.context/sessions/`. A team that mixes tools gets two diverging
  instruction texts.

The 2026-01-25 decision "Removed AGENTS.md from project root" rested
on "no tool auto-loads AGENTS.md"; that premise no longer holds.
Upstream issue: ActiveMemory/ctx#197.

## Approach

One template, loaded once by every agent:

- `AGENTS.md` carries the ctx block between `<!-- ctx:agents -->` and
  `<!-- ctx:agents:end -->`, created or merged by `ctx init` with the
  same `merge.OrCreate` semantics `CLAUDE.md` has today: create when
  absent, prompt (or `--merge`) to insert after the H1 of an existing
  file, `--reset` replaces the block. The embedded template is
  `internal/assets/integrations/agents.md`, rewritten as the
  tool-neutral version of the current CLAUDE.md content.
  `ctx setup agents|codex|opencode|pi` keep deploying that same file
  through `agents.Deploy`, which already skips a file that carries the
  `ctx:agents` marker, so `ctx setup codex` after `ctx init` is a
  no-op for `AGENTS.md`.
- `CLAUDE.md` is still created or merged, but its ctx block
  (`<!-- ctx:context -->` ... `<!-- ctx:end -->`) shrinks to a bridge:
  an explanatory comment and the line `@AGENTS.md`. Claude Code
  expands the import at launch, relative to the importing file. The
  bridge is written unconditionally rather than relying on the native
  fallback because the fallback needs Claude Code 2.1.277 or later, is
  unavailable in some sessions (the Claude Code docs recommend the
  import there), and is switched off by any hand-written `CLAUDE.md`
  in the tree, which is exactly the `--merge` case.

Alternatives considered:

- `AGENTS.md` only, no `CLAUDE.md`: silently loses the instructions on
  older Claude Code and the moment the user adds their own
  `CLAUDE.md`. Rejected.
- Full content in both files: two copies to keep in sync, loaded
  twice. Rejected.
- Keep the full `CLAUDE.md` and symlink `AGENTS.md` to it (the
  agents.md migration tip): symlinks do not survive Windows checkouts
  and `agents.Deploy` rejects them. Rejected.

## Behavior

### Happy Path

1. `ctx init` in a project with neither file writes `AGENTS.md`
   (title, ctx block) and `CLAUDE.md` (title, bridge block). Both are
   reported as created.
2. Claude Code loads `CLAUDE.md` and expands `@AGENTS.md`; Codex,
   Copilot, Cursor, OpenCode and Pi load `AGENTS.md` directly. Every
   agent sees one copy of the instructions.
3. A later `ctx setup codex` finds the `ctx:agents` marker in
   `AGENTS.md` and skips it (unchanged behavior).

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| `AGENTS.md` exists without ctx markers (hand-written) | Default: prompt "Would you like to append ctx context management instructions?"; `--merge`: no prompt. On yes: backup, block inserted after the H1 (at the top when there is no H1), rest untouched. On no: skipped. |
| `AGENTS.md` exists with a ctx block (earlier `ctx init`, or `ctx setup codex` with the old template) | Left alone: "ctx content exists". `ctx init --reset` replaces the block with the current template after a backup. |
| `CLAUDE.md` exists with the old full ctx block (project initialized by an earlier ctx) | Left alone by default, so the old block and the new `AGENTS.md` coexist: the instructions load twice but stay consistent. `ctx init --reset --merge`, the documented upgrade step, replaces the block with the bridge. |
| `CLAUDE.md` exists, hand-written, no markers | Prompt or `--merge`, then the bridge block is inserted after the H1; the user's text is untouched. Without the bridge Claude Code would ignore `AGENTS.md` because a `CLAUDE.md` is present. |
| `AGENTS.md` is a directory or unreadable | The read fails, the write fails, init prints the existing non-fatal warning labelled `AGENTS.md` and continues, as for `CLAUDE.md` today. |
| `--minimal` | Unchanged: both root files are infrastructure, not knowledge files, and are still written. |
| Re-running `ctx init` | Idempotent: both files report "ctx content exists". |

### Validation Rules

None beyond the existing marker handling: a template must contain both
markers or `merge.UpdateMarkedSection` returns `TemplateMissingMarkers`.
The asset tests guard the markers and the bootstrap line in both
templates.

### Error Handling

| Error condition | User-facing message | Recovery |
|-----------------|---------------------|----------|
| Embedded template unreadable | existing `ReadTemplate` error as a non-fatal warning labelled `AGENTS.md` | rebuild the binary |
| `AGENTS.md` write fails | existing `FileWrite` error as a non-fatal warning | fix permissions, re-run `ctx init` |
| Merge declined at the prompt | "skipped AGENTS.md" | re-run with `--merge` |

## Interface

### CLI

```
ctx init [--merge] [--reset]
```

No new flags. `--merge` also auto-merges into an existing `AGENTS.md`;
`--reset` also refreshes its ctx block. The `--merge` help text and
the init long description name both files.

## Implementation

### Files to Create/Modify

| File | Change |
|------|--------|
| `internal/config/marker/marker.go` | add `AgentsEnd = "<!-- ctx:agents:end -->"` |
| `internal/cli/initialize/core/agents/{agents.go,doc.go}` | new: `HandleMd(cmd, force, autoMerge)` for `AGENTS.md`; mirrors `core/claude.HandleMd` with `agent.AgentsMd()`, `cfgHook.FileAgentsMd`, `marker.AgentsStart`/`AgentsEnd` |
| `internal/cli/initialize/cmd/root/run.go` | call `coreAgents.HandleMd` before `coreClaude.HandleMd` |
| `internal/cli/initialize/core/merge/create.go`, `section.go` | `OrCreate` inserts only the marker-bounded block into an existing file (the template's own `# Project Context` title used to come along as a second H1); the block extraction moves to a `markedSection` helper in its own file, shared with `UpdateMarkedSection` |
| `internal/assets/integrations/agents.md` | tool-neutral version of the CLAUDE.md content between `ctx:agents` markers |
| `internal/assets/claude/CLAUDE.md` | bridge: comment plus `@AGENTS.md` inside the `ctx:context` markers |
| `internal/cli/initialize/core/claude/doc.go`, `core/doc.go`, `cmd/root/doc.go`, `cmd/root/cmd.go` | doc comments |
| `internal/assets/commands/commands.yaml`, `flags.yaml` | init long help and `--merge` text |
| `internal/assets/README.md` | directory map rows |
| `docs/cli/init-status.md`, `docs/cli/setup.md`, `docs/operations/migration.md`, `docs/operations/integrations.md`, `docs/operations/upgrading.md`, `docs/operations/runbooks/breaking-migration.md`, `docs/home/common-workflows.md`, `docs/home/configuration.md` | describe `AGENTS.md` plus the bridge |
| `.context/DECISIONS.md` | decision superseding the 2026-01-25 entry |

### Helpers to Reuse

`merge.OrCreate` and `merge.UpdateMarkedSection` (create, prompt,
backup, H1 insertion, `--reset`), `agent.AgentsMd()` (template
reader), `agents.Deploy` (setup path, untouched),
`write/initialize.Created`.

## Testing

- `internal/cli/initialize/init_test.go`: a fresh init creates
  `AGENTS.md` with both markers and the bootstrap line, and `CLAUDE.md`
  with the ctx markers, an `@AGENTS.md` line and no bootstrap text;
  `--merge` into a hand-written `AGENTS.md` inserts after the H1 and
  keeps the content; an `AGENTS.md` with an existing ctx block is left
  alone. The existing CLAUDE.md merge tests keep passing against the
  bridge, and the merge tests for both files assert that the
  template heading does not leak into the merged file.
- `internal/assets/project_test.go`: both templates carry their
  markers; the CLAUDE.md template has the import and not the
  instructions.
- `internal/cli/setup/core/{agents,codex,pi}` tests: unchanged, and
  pass against the new template because they key on the marker only.

## Non-Goals

- Changing the `ctx setup agents|codex|opencode|pi` deployment logic.
- Flipping the ctx repository's own root `CLAUDE.md` / `AGENTS.md`
  (follow-up).
- Auto-migrating an old full CLAUDE.md block without `--reset`, or
  printing a hint about it (follow-up).
- Nested or per-directory `AGENTS.md` files.
- Copilot's `.github/copilot-instructions.md` and the other per-tool
  files `ctx setup` writes.
