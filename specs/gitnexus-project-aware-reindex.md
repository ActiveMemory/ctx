# Spec: project-aware GitNexus reindex guidance in skills

## Problem

Several shipped skill assets tell the agent to fix a stale GitNexus
index by running `gitnexus analyze` (with a generic "or your Docker
wrapper if the npm binary isn't viable" aside). On hosts where the npm
binary cannot build — GitNexus pins tree-sitter 0.21.1, whose native
addon has no arm64 prebuilt for recent Node ABIs — `gitnexus analyze` is
a silent no-op, and the skill's suggestion sends the agent down a dead
end even when a working repo-local indexer sits one directory over
(observed: `/ctx-remember` in the `os` project suggested the broken npx
path while `make gitnexus-index` existed).

The generic wording predates the repo-local indexing convention. It
names the fallback command first and never tells the agent to look for
the repository's own indexing entry point.

## Design

Invert the guidance: prefer the repository's own indexing entry point,
fall back to bare `gitnexus analyze` only when none exists. The skills
ship to arbitrary user projects, so the wording stays generic — it
names the *kinds* of repo-local entry point (a `make gitnexus-index`
target, an indexing script, or the steps in the repo's `GITNEXUS.md`)
rather than hard-coding any one project's command.

Sites (canonical Claude sources; the Copilot CLI copies are generated
from these by `hack/sync-copilot-skills.sh` and must be regenerated, not
hand-edited):

- `internal/assets/claude/skills/ctx-remember/SKILL.md` — the stale-index
  suggestion in the companion-tool check.
- `internal/assets/claude/skills/ctx-architecture-enrich/SKILL.md` — the
  indexing-precondition note, the "no MCP connected" block, and the
  ">5 commits behind" hard-stop remedy.

No behavior in ctx's Go code changes; this is skill-prompt wording only.

## Implementation

- Reword each site to: reindex with the repo's own entry point (a
  `make gitnexus-index` target, an indexing script, or its
  `GITNEXUS.md`) if it has one, else `gitnexus analyze`.
- Run `make sync-copilot-skills` so the Copilot CLI ctx-remember copy
  matches; `make check-copilot-skills` (in `make audit`) must pass.

## Acceptance

- No shipped skill suggests `gitnexus analyze` as the *first* remedy for
  a stale index; each points at the repo-local entry point first.
- `make audit` passes (`check-copilot-skills` confirms the generated
  Copilot copies are in sync).

Closes the TASKS item "Skill assets hard-code `npx gitnexus analyze` as
the stale-index remedy" (the project-aware follow-up to the earlier
de-npx pass).
