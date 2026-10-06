# Private Mode

> **Status: stub.** This captures the problem and the candidate
> scope recorded in TASKS.md on 2026-09-23 (session 28b10323). No
> design decisions have been made; the Approach section lists
> options, not choices. Flesh out with `/ctx-plan` or `/ctx-spec`
> before implementing.

## Problem

`ctx` can't be used in a project without committing to that
project. Open-source projects can't make `ctx` a contributor
dependency, and today using `ctx` there either leaks into tracked
files or breaks things for other contributors. Observed in
spike-sdk-go:

- **`ctx init` edits tracked files.** It appends to `.gitignore`,
  adds `-include Makefile.ctx` to `Makefile`, and writes
  `CLAUDE.md`.
- **A committed `CLAUDE.md` without `.context/` blocks other
  agents.** A contributor who has `ctx` installed gets
  `Error: no .context here` from `ctx system bootstrap`, and the
  CLAUDE.md rule "installed but returns an error -> relay and STOP"
  halts their agent. Reproduced in a fresh clone containing only
  `CLAUDE.md`.
- **Ignoring `.context/` removes the undo layer.** The constitution
  requires git as the safety net for agent-driven edits ("persistent
  memory is dishonest without git reflog"). With `.context/`
  untracked, recoveries like
  `git show <sha>:.context/LEARNINGS.md` after a clobber become
  impossible.

## Approach

Candidate scope, from the task. Each item is an option to evaluate,
not a decision:

1. **Init without touching tracked files.** `ctx init --private`
   (or equivalent) writes ignore rules to `.git/info/exclude`
   instead of `.gitignore`, and adds no Makefile include (for
   example an untracked `GNUmakefile`, or no make targets at all).
2. **Untracked agent instructions.** The instructions live in an
   untracked file (for example `CLAUDE.local.md`), or the CLAUDE.md
   template treats a missing `.context/` as "not a ctx project"
   rather than as an error to STOP on.
3. **Undo for an untracked `.context/`.** A versioning story so
   recovery still works, for example a nested git repo inside
   `.context/` or snapshots.
4. **Leak detection.** `ctx drift` / `ctx doctor` flag private-mode
   state that has leaked into tracked files.

## Open Questions

- Is private mode a flag on `init`, a persisted mode in `.ctxrc`,
  or detected from the environment (`.context/` excluded via
  `.git/info/exclude`)?
- Item 2 has two different fixes: an untracked instructions file,
  or a CLAUDE.md template change that benefits non-private projects
  too. Pick one or both.
- How does a nested repo in item 3 interact with the "Git is
  required" invariant (`specs/require-git.md`) and with tools that
  walk the outer repo?
- Which existing `ctx` commands assume `.gitignore` or `Makefile`
  ownership and need a private-mode branch?
