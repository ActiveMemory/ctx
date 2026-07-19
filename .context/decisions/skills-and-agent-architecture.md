# skills-and-agent-architecture

## [2026-06-07-112203] ctx-dream executor is a documented contract, not a hardcoded cron/claude assumption

**Status**: Accepted

**Context**: Settling ctx-dream v1 open questions. The executor runs the
out-of-band dream pass (read ideas/, classify+ground, write proposals). Question
was cron 'claude -p' vs a raw Anthropic-API scheduled loop.

**Decision**: ctx-dream executor is a documented contract, not a hardcoded
cron/claude assumption

**Rationale**: cron 'claude -p' is the reference executor (reuses Claude Code
auth, tool-calling, and PreToolUse hooks so the three guards are structural for
free; matches the existing skill draft and the cheap-validation goal). But we
must NOT assume it is the only executor: other harnesses (different AI CLI, raw
API loop, CI runner) must be able to run the same dream. So ctx owns an
executor-agnostic Go core (dreams/ layout, state record, ledger, proposal
schema, the three guards as callable logic) and the executor is a documented
contract: run one bounded pass, enforce the three guards STRUCTURALLY (Claude
Code via PreToolUse hooks; API loop via in-loop tool executor), fail loud, write
proposals-only into dreams/. Dream is opt-in, not enabled by default.

**Consequence**: Guards live as reusable Go logic in internal/dream/, not only
as a hook script. Two user-facing docs are required: a Claude Code enablement
guide and an executor-contract reference for other harnesses. The serendipity
review skill is split into its own spec (specs/ctx-serendipity.md). v1 ships the
cron/claude-p reference path but the data contract + guards stay
executor-portable.

---

## [2026-06-06-133805] ctx-dream: standalone proposing memory consolidator (Option B), human-gated via serendipity

**Status**: Accepted

**Context**: We explored whether ctx should grow a scheduled, background 'dream'
(a sleep-time memory process) and how it should relate to canonical memory. Felt
pain: the author's ideas/ folder is too overwhelming to triage, and canonical
files bloat over time (109 decisions, 151 learnings, 154 unimported sessions).
The risk to avoid: a background LLM job autonomously rewriting authoritative
memory and silently corrupting it (the research shows continuous LLM
consolidation is lossy and non-monotonic). Full debate:
.context/briefs/20260606T203414Z-ctx-dream-disciplined-consolidator.md

**Decision**: ctx-dream: standalone proposing memory consolidator (Option B),
human-gated via serendipity

**Rationale**: Chose a NEW, standalone, PROPOSING consolidator (Option B): it
writes only to its own sidecar + proposals queue + ledger + per-dream archive,
never autonomously to the five canonical files; a human 'serendipity' review
session is the sole bridge (accept/reject/amend) into canonical. One skill, two
modes: discipline (default; grounded, structured, provenanced proposals) and
creative/exploration (a safe relaxation: resurface + chance, reader-only).
Principle: decouple the cognition, reuse the plumbing (own the consolidation
logic; reuse import/enrich/kb-ingest via the enriched-journal data contract).
Standalone so mechanics evolve independently and changes to existing curation
skills can't break it, and for creative freedom (don't assume existing verbs
suffice). Discipline-first because it is the hard load-bearing substrate and
creative is a strict, safer relaxation of it. Grounded in
ideas/ctx-dreams/research: Auto-Dreamer (2605.20616) for the architecture,
'Useful Memories Become Faulty When Continuously Updated by LLMs' (2605.12978)
for the threat model, and the deep-research eval cluster for the finding that a
single agreeable LLM is not an adversarial gate (it silently repairs the missing
justification), which is why the gate must be human. Rejected: Option A (dream
owns a parallel canonical store, which does not fix bloat and creates two
divergent substrates); autonomous mutation / auto-approve (violates 'each memory
entry needs dedicated human attention'); pure-garden-only (under-serves
engineering's need for grounding and actionability); coupling to existing
skills' internals; garden-first build order.

**Consequence**: Positive: nothing autonomous touches canonical, so the system
is reversible by construction; the dream's mechanics can evolve freely; v1
(disciplined ideas/ triage, validated via a ctx-remind-nagged ~15-minute review
round) is low-stakes and validates the mechanism and author engagement cheaply.
Negative / trade-off: no human serendipity session = no consolidation, so the
dream's entire value is gated behind human review cadence, and the author
historically under-runs curation; mitigated only by ctx-remind nags + targeting
felt pain (ideas/) + a pleasure-not-chore framing. Validation of the full
product thesis (disciplined consolidation of canonical memory for engineering
teams) is deferred to a later test on a project where bloat actually bites. Spec
work proceeds via /ctx-spec --brief on the brief above; key mechanics remain
open (executor, proposal schema, ledger schema, .context/ layout).

---

## [2026-04-09-001332] Architecture skill pipeline is a triad not a quartet

**Status**: Accepted

**Context**: Had a proposed ctx-architecture-extend for extension point mapping,
making four skills

**Decision**: Architecture skill pipeline is a triad not a quartet

**Rationale**: Extension points already covered per-module in DETAILED_DESIGN
and by registration site discovery in enrich. Fourth skill fragments pipeline
without distinct value

**Consequence**: Pipeline is map enrich hunt. Three skills three questions: how
does it work, how well does it connect, where will it break

---

## [2026-03-25-233646] Architecture analysis and enrichment are separate skills — constraint is the feature

**Status**: Accepted

**Context**: Observed that agents take shortcuts when code intelligence tools
are available during architecture analysis. A 5.2x depth reduction was measured
(5866 vs 1124 lines) when GitNexus was available during reading. Mentioning
unavailable tools by name in a skill plants the idea for the agent to use them.

**Decision**: Architecture analysis and enrichment are separate skills —
constraint is the feature

**Rationale**: Discovery requires forced reading without shortcuts. Validation
and quantification are a separate pass. Two-pass compiler analogy: semantic
parsing (human-style reading) then static analysis (graph enrichment). Never
mention tools you want the agent to avoid — absence is the only reliable
constraint.

**Consequence**: ctx-architecture deliberately excludes code intelligence tools
from allowed-tools and never mentions them. ctx-architecture-enrich is a
separate skill that runs after, using the deep artifacts as baseline. Gemini is
allowed in both for upstream/external lookups only.

---


## [2026-03-25-173336] Prompt templates removed — skills are the single agent instruction mechanism

**Status**: Accepted

**Context**: Prompt templates (.context/prompts/) overlapped with skills but had
no discoverability — even the project creator didn't know they existed

**Decision**: Prompt templates removed — skills are the single agent
instruction mechanism

**Rationale**: Adding metadata to prompts to fix discoverability would recreate
the skill system. One concept is better than two.

**Consequence**: code-review, explain, refactor promoted to proper skills. ctx
prompt CLI removed. loop.md retained as ctx loop config file at
.context/loop.md.

---

## [2026-03-13-223111] Delete ctx-context-monitor skill — hook output is self-sufficient

**Status**: Accepted

**Context**: The skill documented how to relay context window warnings, but the
hook message already includes IMPORTANT: Relay this context window warning to
the user VERBATIM which agents follow without the skill.

**Decision**: Delete ctx-context-monitor skill — hook output is
self-sufficient

**Rationale**: No mechanism exists for hooks to trigger skills. The skill was
never loaded during sessions. Adding enforcement elsewhere would either be too
far back in context (playbook) or dilute the already-crisp hook message.

**Consequence**: One fewer skill to maintain. No behavioral change — agents
continue relaying warnings as before.

---



## [2026-03-12-133007] Rename ctx-map skill to ctx-architecture

**Status**: Accepted

**Context**: The name 'map' didn't convey the iterative, architectural nature of
the ritual

**Decision**: Rename ctx-map skill to ctx-architecture

**Rationale**: 'architecture' better describes surveying and evolving project
structure across sessions

**Consequence**: All cross-references updated across skills, docs, .context
files, and settings

---

---

## [2026-03-01-090124] Promote 6 private skills to bundled plugin skills; keep 7 project-local

**Status**: Accepted

**Context**: Reviewed all 13 _ctx-* private skills to determine which are
universally useful for any ctx user vs specific to the ctx codebase or personal
infra.

**Decision**: Promote 6 private skills to bundled plugin skills; keep 7
project-local

**Rationale**: Promote if the skill benefits any ctx-powered project without
project-specific hardcoding. Keep private if it references this repo's Go
internals, personal infra, or language-specific tooling. Promote list: _ctx-spec
(generic scaffolding), _ctx-brainstorm (design facilitation), _ctx-verify (claim
verification), _ctx-skill-create (skill authoring), _ctx-link-check (doc link
audit), _ctx-permission-sanitize (Claude Code permissions audit). Keep list:
_ctx-audit (Go/ctx checks), _ctx-qa (Go Makefile), _ctx-backup (SMB infra),
_ctx-release/_ctx-release-notes (ctx release workflow), _ctx-update-docs (ctx
package mapping), _ctx-absorb (borderline, revisit later).

**Consequence**: Six skills move from .claude/skills/ to
internal/assets/claude/skills/ and become available to all ctx users via ctx
init. Cross-references between skills need updating (e.g., /_ctx-brainstorm
becomes /ctx-brainstorm). The seven remaining private skills stay project-local.

---

## [2026-02-26-100005] Agent autonomy and separation of concerns (consolidated)

**Status**: Accepted

**Consolidated from**: 3 decisions (2026-01-21 to 2026-01-28)

- Removed AGENTS.md from project root. Consolidated on CLAUDE.md (auto-loaded) +
  .context/AGENT_PLAYBOOK.md as the canonical agent instruction path. Projects
  using ctx should not create AGENTS.md.
- ~~Separate orchestrator directive from agent tasks~~ (superseded 2026-03-25:
  IMPLEMENTATION_PLAN.md removed — TASKS.md is the single source of truth for
  work items, AGENT_PLAYBOOK.md covers agent behavior).
- No custom UI -- IDE is the interface. UI is a liability; IDEs already excel at
  file browsing, search, markdown editing, and git integration. Focus CLI
  efforts on good markdown output.

---

