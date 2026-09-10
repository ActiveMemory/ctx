# Spec: OpenCode Mirror Sync — Adopt the Codex Full-Mirror Model

## Problem

The OpenCode skill sync introduced for issue #158
(`specs/opencode-skill-parity.md`) enrolls skills opt-in by directory
presence: a canonical skill syncs only if someone has already created
its directory under `internal/assets/integrations/opencode/skills/`.
The Codex integration (PR #163) shipped a stronger model, leaving the
OpenCode tree behind on three counts:

1. **Re-drift by default.** Every future canonical skill silently
   does not ship to OpenCode until someone remembers to `mkdir` —
   the same drift class issue #158 diagnosed. Canonical has 54
   skills; Codex mirrors 50; OpenCode shipped a curated 17.
2. **Orphans are blessed.** A deleted canonical skill leaves a stale
   OpenCode directory, and `TestSyncedSkillParity` exempts
   counterpart-less directories, so nothing flags it.
3. **`references/` are not shipped.** Skill bodies cite
   `references/` files; the OpenCode sync copies only `SKILL.md`.
   Currently latent (none of the 17 cite references), it becomes a
   live 404 the moment a references-bearing skill joins the tree —
   which the mirror model does immediately (`ctx-humanize`,
   `ctx-journal-enrich-all`, `ctx-skill-audit`).

Conversely, Codex's freshness check runs only via `make audit`,
which CI never runs — its tree can go stale without CI noticing.

## Fix

Adopt the mirror model for OpenCode and extend CI-level parity
enforcement to Codex — each integration inherits the other's
strength:

- Rewrite `hack/sync-opencode-skills.sh` on the
  `hack/sync-codex-skills.sh` model: default-include every canonical
  skill, exclusion list for Claude Code-only skills
  (`ctx-permission-sanitize`, `ctx-plan-import`, `ctx-dream`,
  `ctx-skill-create` — same four as Codex, same rationale), mirror
  each skill's `references/` directory, remove orphaned OpenCode
  directories.
- Embed `integrations/opencode/skills/*/references/*` and add
  `agent.OpenCodeSkillReferences()` plus reference deployment in
  `ctx setup opencode --write` (mirroring the Codex deploy path), so
  shipped skill bodies never cite files that do not exist.
- Extend `TestSyncedSkillParity`: add the Codex tree to the
  byte-parity check; for mirror trees (OpenCode, Codex) assert
  completeness (every canonical skill outside the exclusion list is
  present), no orphans, and reference parity (every embedded
  canonical `references/*.md` byte-matches the generated copy;
  non-`.md` references are outside the canonical embed glob and are
  covered by the sync scripts, not the test).
- Harden `check-opencode-skills` restore to full-replace
  (`rm -rf` + `cp -r`), since a mirror sync can add and remove
  directories.

## Decisions

- **`ctx-serendipity` ships despite its companion `ctx-dream` being
  excluded.** Its single `/ctx-dream` reference is a passing
  provenance mention ("the pass that produces the proposals you
  review here"), not an instruction to invoke it — and the
  cross-tool workflow is real: the dream runs headless under Claude
  Code, the dreams/ notebook lives in the repo, so the review walk
  can happen from any tool. The Codex mirror ships it under the
  same reasoning.

## Non-Goals

- Changing the Copilot CLI model (its tree carries tool-only wrapper
  skills; presence-based enrollment remains correct there).
- Widening the canonical Claude embed glob beyond
  `references/*.md`.
- Any transform beyond stripping `allowed-tools:` (the "no terse
  transform" decision from `specs/opencode-skill-parity.md` stands).

## Verification

- `make audit` green: OpenCode tree mirrors 50 skills, in sync.
- `go test ./internal/assets/read/skill/` passes with completeness,
  orphan, and reference assertions active for both mirror trees.
- `ctx setup opencode --write` deploys `references/` files alongside
  `SKILL.md` (verified via the deploy test suite).
