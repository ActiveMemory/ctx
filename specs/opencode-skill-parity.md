# Spec: OpenCode Skill Parity — Generate the Tree from Canonical Claude Skills

Issue: https://github.com/ActiveMemory/ctx/issues/158

## Problem

The OpenCode integration ships 10 hand-written skills while Claude has
54 and Copilot CLI has 49. The entire planning arc documented in the
[Design Before Coding](https://ctx.ist/recipes/design-before-coding/)
recipe — `/ctx-brainstorm`, `/ctx-spec`, `/ctx-task-out`,
`/ctx-implement` — is absent, so OpenCode users cannot follow the
project's own recommended design workflow.

Diffing the 10 hand-written skills against their Claude counterparts
shows the terseness is mostly truncated reference material (dropped
flag tables, output descriptions), not OpenCode-specific adaptation.
Two conventions coexisting in one tree is unintentional divergence.

## Approach

Align OpenCode with the Copilot CLI model: generate the tree from the
canonical Claude skills at build time. `hack/sync-copilot-skills.sh`
already proves the shape in production (see the closed issue #61):
derive each enrolled skill from
`internal/assets/claude/skills/<name>/SKILL.md` with the Claude-specific
`allowed-tools:` frontmatter key stripped, opt-in by directory
presence, wired into `make build`, gated by a `check-*` target in
`make audit`.

No Go changes: the embed glob
(`integrations/opencode/skills/*/SKILL.md`), `agent.OpenCodeSkills()`,
and `deploySkills()` all walk whatever directories exist.

## Deliverables

1. **`hack/sync-opencode-skills.sh`** — sibling of
   `hack/sync-copilot-skills.sh`, same contract: iterate
   `internal/assets/integrations/opencode/skills/*/`, overwrite each
   `SKILL.md` from the Claude source with `allowed-tools:` stripped;
   skills with no Claude counterpart are left untouched.
2. **Makefile wiring** — `sync-opencode-skills` runs as part of
   `make build`; `check-opencode-skills` (freshness gate, mirrors
   `check-copilot-skills`) runs as part of `make audit`.
3. **Enrollment** — the existing 10 skills (ctx-agent, ctx-handover,
   ctx-kb-ask, ctx-kb-ground, ctx-kb-ingest, ctx-kb-note,
   ctx-kb-site-review, ctx-remember, ctx-status, ctx-wrap-up) become
   synced; 7 new skills enroll: the planning arc (ctx-brainstorm,
   ctx-spec, ctx-task-out, ctx-implement, ctx-plan) plus the capture
   pair the arc's workflow leans on (ctx-task-add, ctx-decision-add).
   17 total, canonical Claude names throughout.
4. **Docs** — `docs/home/opencode.md` slash-command section updated
   from the hand-listed 4 to the full synced set.

## Decisions

- **Canonical bodies replace terse variants.** The issue's open
  question (a `terse` transform for OpenCode's context budget) is
  resolved as: no terse transform. The truncation was unintentional
  divergence. The few genuinely OpenCode-specific lines in the
  hand-written bodies (e.g. ctx-status's "the slash command takes no
  arguments" note) are dropped with them; if OpenCode-specific
  adaptation is ever needed, the honest fix is a transform in the sync
  script, not hand-edits that the next sync overwrites.
- **The five planning-arc skills carry no `allowed-tools:` key**, so
  their sync transform is a byte-identical copy. Their `/ctx-*`
  cross-references are closed within the enrolled set.
- **`ctx-learning-add` and `ctx-convention-add` are enrolled** (review
  finding): synced `ctx-wrap-up`, `ctx-handover`, and `ctx-kb-note`
  route capture through them, and the capture-pair rationale that
  enrolled `ctx-task-add`/`ctx-decision-add` applies verbatim.
- **Five dangling `/ctx-*` references are an accepted gap**: the
  canonical bodies of `ctx-wrap-up` and `ctx-remember` mention
  `/ctx-commit`, `/ctx-reflect`, `/ctx-digest`, `/ctx-history`, and
  `/ctx-journal-enrich-all`, which are not enrolled here. Accepted
  because the follow-up mirror-sync branch
  (`specs/opencode-mirror-sync.md`) ships the full canonical tree and
  closes all of them; enrolling piecemeal now would churn the tree
  twice.

## Acceptance Criteria

- [ ] `hack/sync-opencode-skills.sh` exists and mirrors the Copilot
      script's contract
- [ ] `make build` syncs OpenCode skills; `make check-opencode-skills`
      fails on staleness
- [ ] The Design Before Coding arc is available in OpenCode
- [ ] Enrolled OpenCode skills are byte-identical to their Claude
      source minus `allowed-tools:`
- [ ] Skill names align 1:1 with the Claude tree
- [ ] Build, lint, and compliance tests pass
      (`TestSkillFrontmatter` covers the new dirs via the embed glob)

## Non-Goals

- Renaming or enrolling the 14 legacy-named unsynced Copilot skills
  (separate concern; OpenCode already uses canonical names)
- A terse/summarizing transform in the sync script
- Any change to skill deployment (`ctx setup opencode --write`) or the
  embed layer
