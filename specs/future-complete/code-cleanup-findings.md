# Spec: Code Cleanup Findings

## Problem

Accumulated code cleanup tasks from prior review sessions,
tracked in TASKS.md under "Code Cleanup Findings". Covers:
exec consolidation, MCP rename, Cmd() boilerplate extraction,
ctxrc constant audit, doc.go drift detection, companion check
wiring, SHALLOW_DOC fixes, spec nudge feature, skill discovery
hook, cmd/ purity enforcement, and schema fixes.

## Scope

See TASKS.md "Code Cleanup Findings" section for the full list.
Individual sub-specs exist for larger items:
- `specs/exec-package.md` — exec.Command consolidation
- `specs/ctxrc-constant-promotion.md` — ctxrc promotion audit

## Non-Goals

- Rewriting lint scripts in Go (blocked on ctxctl)
- Fixing 28 grandfathered cmd/ purity violations (separate task)
