# Runbook Migration: hack/runbooks/ to docs/operations/runbooks/

## Problem

Runbooks in `hack/runbooks/` are buried where humans don't look.
No agent has ever autonomously executed them since the project's
inception — they're effectively human-triggered "paste this prompt"
guides. Moving them to the docs site makes them discoverable,
searchable, and part of the site nav.

## Approach

- Move 3 existing runbooks from `hack/runbooks/` to `docs/operations/runbooks/`
- Create 5 new runbooks per TASKS.md runbook section
- Reframe narrative: "use this prompt with your agent"
- Update all cross-references (docs, skills, nav)
- Remove `hack/runbooks/`

## Scope

| Item | Action |
|------|--------|
| codebase-audit | Moved |
| docs-semantic-audit | Moved |
| sanitize-permissions | Moved |
| release-checklist | Created |
| breaking-migration | Created |
| hub-deployment | Created |
| new-contributor | Created |
| plugin-release | Created |

## Decision

Runbooks belong in the public docs because they are human-readable
procedures, not internal build scripts. The `hack/` convention is
for scripts the build system calls, not guides a person reads.
