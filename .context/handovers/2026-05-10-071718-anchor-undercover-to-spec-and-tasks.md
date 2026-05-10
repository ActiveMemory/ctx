---
generated-at: 2026-05-10T07:17:18Z
---
# Handover [2026-05-10-071718] anchor undercover analysis → spec and task breakdown for kb-editorial-pipeline + require-git

**Provenance:** commit `f9582163` on branch `main`

## Summary

Did an undercover analysis of a sibling clean-room editorial-memory
tool under `~/Desktop/WORKSPACE/anchor/`, deliberately scrubbed of
company / project naming so its design ideas could travel into ctx
without inadvertent specifics. Produced three idea documents:
`ideas/001-sibling-project-undercover-analysis.md` (feature inventory
and ranked borrowable mechanisms), `ideas/002-editorial-pipeline-and-
skill-rigor.md` (deeper editorial-pipeline lift plan plus side-by-side
skill ceremony comparison highlighting the wishy-washy plan→spec
handoff in our skills vs the sibling's path-required spec contract),
`ideas/003-editorial-pipeline-debated-brief.md` (debated brief output
of a `/ctx-plan` adversarial-interview pass; carries the two
organizing principles).

Wrote two specs: `specs/kb-editorial-pipeline.md` (4 modes,
9 KB artifacts, closeout/fold mechanism, handover, render via existing
zensical) and `specs/require-git.md` (constitutional precondition;
breaking change). Appended three task phases to `.context/TASKS.md`:
**Phase SK** (skill surface polish — `MarkFlagRequired` on capture
flags, `--brief <path>` on `/ctx-spec`, authority boundary sections,
"light compression" wording standardization), **Phase RG** (require git
as architectural precondition; refuse-on-no-git with no auto-init),
**Phase KB** (editorial pipeline + handover; depends on Phase SK +
Phase RG; sub-phases KB-2 things-wtf port for validation and KB-3 docs).

Saved one feedback memory at `~/.claude/projects/-Users-volkan-Desktop-
WORKSPACE-ctx/memory/feedback_no_defer_unfamiliar_scope.md` capturing
the meta-correction: when lifting from a battle-tested external design,
default to lifting the whole shape; do not yeet on uncertainty. Two
organizing principles surfaced in the planning round and earned their
place at the top of the brief: **P1** — the LLM is the migration tool
(makes committing to specific schemas in v1 cheap rather than reckless);
**P2** — a KB of KBs is a KB (recursive composability collapses
federation, multi-team consolidation, and taxonomy-was-wrong recovery
into one mechanism). Captured 5 decisions and 5 learnings via the
ceremony skills before this handover. Working tree carries the new
specs/handover plus the unrelated pre-session
`.context/DECISIONS.md` and `.context/LEARNINGS.md` modifications;
nothing committed deliberately.

## Next session

Pick **Phase SK** (skill surface polish — smaller, unblocks
`/ctx-spec --brief` for cheaper future spec sessions) OR **Phase RG**
(require-git constitutional change — deserves its own PR with
prominent breaking-change notes in `dist/RELEASE_NOTES.md`) as the
first concrete implementation work. Both are independent and both
are hard prerequisites for Phase KB.

Specific first action: `cd ~/Desktop/WORKSPACE/ctx && ls
.context/TASKS.md` to find the chosen phase header, mark its first
task `#in-progress`, then read the corresponding source — for SK
that is `ideas/002-editorial-pipeline-and-skill-rigor.md` §3
("Reframing the wishy-washy skills"); for RG that is
`specs/require-git.md`. Slight argument for SK first because once
`--brief` ships, future spec sessions (including iterating on
`specs/require-git.md` itself) become cheaper and more disciplined.

Before writing any code: review the unrelated modifications already
sitting in `.context/DECISIONS.md` and `.context/LEARNINGS.md` from
before this session started — decide whether to commit them first,
stash, or fold into the SK or RG PR. The HARD GATE hook will block
any commit until lint + tests pass and the working tree is clean
of uncommitted incidentals.

## Highlights

- New (ideas/): `001-sibling-project-undercover-analysis.md`,
  `002-editorial-pipeline-and-skill-rigor.md`,
  `003-editorial-pipeline-debated-brief.md`.
- New (specs/): `kb-editorial-pipeline.md`, `require-git.md`.
- New (auto-memory):
  `feedback_no_defer_unfamiliar_scope.md` plus a row added to
  the auto-memory `MEMORY.md` index.
- Modified: `.context/TASKS.md` (Phase SK / Phase RG / Phase KB
  appended at end of file; structure mirrors sibling's named-
  phase grammar; references back to the two specs and the brief).
- Captured this session: 5 decisions (lift-pipeline, mandate-git,
  KB-ontology, pair-handover-and-editorial, KB-RULES-naming) +
  5 learnings (P1 LLM-migration-tool, P2 KB-of-KBs-is-a-KB,
  KB epistemology, lift-renames-with-features, workaround-tax-as-
  validation).
- Side-task crystallization: Phase 0 Grounding's existing item
  about "good 'phasing' mechanism for tasks" now has empirical
  material (this session used named phase codes SK/RG/KB, spec
  references inline, prerequisites called out) — useful when
  whoever picks that exploratory task up models what
  `ctx task add --phase` should look like.
- Dogfooded the future handover shape from
  `specs/kb-editorial-pipeline.md`: this file lives at
  `.context/handovers/<TS>-<slug>.md` with `generated-at` YAML
  frontmatter, Provenance line, and Summary / Next session /
  Highlights / Open questions H2 sections — the exact shape the
  Phase KB CLI will produce. `/ctx-remember` won't auto-pick-it-up
  until Phase KB ships, but a future cold-start can manually
  point to it as a forward-compatible artifact.

## Open questions

- **Phase ordering:** SK before RG, or RG before SK? Both are
  independent and Phase KB needs both. Slight argument for SK first
  (the `--brief` flag makes iterating on `specs/require-git.md`
  cheaper), but RG first is also defensible because it's a
  constitutional change that wants a clean focused PR.
- **Phase KB subcommand grammar:** `ctx kb ingest|ask|...`
  (kb-prefixed) vs `ctx ingest|ask|...` (top-level, sibling-style).
  Spec leans prefixed; pin during Phase KB implementation kickoff.
- **Phase RG opt-out list completeness:** beyond `--help` /
  `--version` / `ctx system bootstrap`, audit the command tree
  during Phase RG implementation; default new entries to
  git-required.
- **Pre-session unrelated `.context/DECISIONS.md` and
  `.context/LEARNINGS.md` modifications** carried into next session;
  the new captures from this session sit on top of them. Decide
  intent before committing.
