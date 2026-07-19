# Decisions

<!-- DECISION FORMATS

## Quick Format (Y-Statement)

For lightweight decisions, a single statement suffices:

> "In the context of [situation], facing [constraint], we decided for [choice]
> and against [alternatives], to achieve [benefit], accepting that [trade-off]."

## Full Format

For significant decisions:

## [YYYY-MM-DD] Decision Title

**Status**: Accepted | Superseded | Deprecated

**Context**: What situation prompted this decision? What constraints exist?

**Alternatives Considered**:
- Option A: [Pros] / [Cons]
- Option B: [Pros] / [Cons]

**Decision**: What was decided?

**Rationale**: Why this choice over the alternatives?

**Consequence**: What are the implications? (Include both positive and negative)

**Related**: See also [other decision] | Supersedes [old decision]

## When to Record a Decision

✓ Trade-offs between alternatives
✓ Non-obvious design choices
✓ Choices that affect architecture
✓ "Why" that needs preservation

✗ Minor implementation details
✗ Routine maintenance
✗ Configuration changes
✗ No real alternatives existed

-->

## [2026-07-19-100259] M4 conventions digestion: curated ## -section taxonomy, unified into the entry-kind mover

**Status**: Accepted

**Context**: CONVENTIONS.md is the last unbounded canonical root (390 lines, 18 stable curated ## sections). The M3 mover refuses the convention kind (ErrApplyNotEntryKind). The spec assumed conventions accrete as ### entries under a ## Recent staging zone, but the real file is a curated ## -section taxonomy with no ## Recent and no timestamps — the blocking TBD for M4.

**Decision**: M4 conventions digestion: curated ## -section taxonomy, unified into the entry-kind mover

**Rationale**: Approach A (unify) over a separate convention path or a one-time script: the file already fits the entry-kind preamble|staging|## Themes shape, so parametrizing the mover by a per-kind entry prefix (## for conventions vs ## [) removes special-case code instead of adding a parallel path. Identity = section title (no timestamp). Keep the scanner dumb — no fence detection (a rabbit hole in a destructive parser); byte-conservation plus the human-gated plan review make a dumb scanner safe against false ## boundaries. Curated-taxonomy matches what conventions are: stable, not accreting.

**Consequence**: Retire parseConvention, ConventionLinePrefix and HeadingRecent (## Recent); the three-self-similar-tiers claim becomes literally true. New sentinel ErrDuplicateStagedTitle (fail-loud on duplicate ## titles). heading.ParseEntryBlocks stays untouched (zero regression to the live LEARNINGS/DECISIONS path). ctx convention add joins the unified prepend anchor (insert before the first ## section, skipping the structural ## Themes; fallback AfterHeader), landing new conventions newest-first in staging above ## Themes — correcting the old AppendAtEnd inconsistency (decisions/learnings already prepend). This add-path change is IN M4 scope. specs/progressive-disclosure.md revised accordingly.

---

## Themes

- package-structure-and-quality-gates — Package taxonomy & quality gates: write/ output, internal/err, config/ constants, doc.go floor, AST audit tests, log split, GraphBuilder interface → [package-structure-and-quality-gates](decisions/package-structure-and-quality-gates.md)
- cli-command-surface — CLI surface: singular command names, flag naming (--json-file/--consequence), flags-not-subcommands, hook->setup, bootstrap placement, backup/recall removal → [cli-command-surface](decisions/cli-command-surface.md)
- skills-and-agent-architecture — Skill & agent architecture: ctx-dream design, architecture skill triad, analysis/enrichment split, prompt-templates removed, skill promotion, agent autonomy → [skills-and-agent-architecture](decisions/skills-and-agent-architecture.md)
- context-model-and-state — Context model & on-disk state: CWD-anchored resolution, encryption-key resolution, pad snapshot, gitignore handovers/memory, server-authoritative Author, init guard → [context-model-and-state](decisions/context-model-and-state.md)
- hooks-session-and-telemetry — Hooks, session & telemetry: hook-relay provenance, memory-pressure signals, billing piggyback, heartbeat telemetry, context-window detection, notifications → [hooks-session-and-telemetry](decisions/hooks-session-and-telemetry.md)
- journal-and-knowledge-lifecycle — Journal & knowledge lifecycle: verbatim journal render, journal-local vs shareable LEARNINGS, #done removal, write-once-then-consolidate, task/knowledge mgmt → [journal-and-knowledge-lifecycle](decisions/journal-and-knowledge-lifecycle.md)
- kb-and-vocabulary — KB pipeline & vocabulary: Phase-KB editorial design, localizable i18n primitives, config-driven per-file freshness checks → [kb-and-vocabulary](decisions/kb-and-vocabulary.md)
- integrations-and-assets — Integrations & assets: YAML text externalization, ctxctl audit binary, companion-tool peer-MCP (no gateway), editor-integration harnesses, Desktop shell-out → [integrations-and-assets](decisions/integrations-and-assets.md)
- product-community-and-deps — Product, community & deps: progressive disclosure, ceremony-credit throttle, statusline informs-not-gates, sonnet 200k default, IRC->Discord, drop fatih/color → [product-community-and-deps](decisions/product-community-and-deps.md)
- security-and-permissions — Security & permissions: system-path deny-list as safety net not boundary; permission-model decisions → [security-and-permissions](decisions/security-and-permissions.md)
