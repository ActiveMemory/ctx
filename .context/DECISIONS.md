# Decisions

<!-- INDEX:START -->
| Date | Decision |
|----|--------|
| 2026-08-23 | Codex memories are out of scope for the ctx memory bridge |
| 2026-08-23 | ctx never parses Codex config.toml: it appends the [mcp_servers.ctx] table and scans header lines |
| 2026-08-23 | Codex plugin root lives at internal/assets/codex with a repo marketplace at .agents/plugins/marketplace.json |
| 2026-07-25 | Beyond a byte ceiling, knowledge content should become tooling, not more Markdown |
| 2026-07-25 | M5 knowledge health is two suggest-only signals: foldable root (staging count) and heavy page (bytes) |
| 2026-07-25 | Theme declaration via the Themes section keyword, on all three canonical kinds |
| 2026-07-25 | ctx convention add requires --section; no default; placeholders rejected (strict) |
| 2026-07-19 | M4 conventions digestion: curated ## -section taxonomy, unified into the entry-kind mover |
<!-- INDEX:END -->

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

## [2026-08-23-120839] Codex memories are out of scope for the ctx memory bridge

**Status**: Accepted

**Context**: ctx bridges Claude Code auto-memory (~/.claude/projects/<slug>/memory/MEMORY.md) into .context. Codex memories live under ~/.codex/memories as SQLite-backed generated state, are off by default (features.memories), and OpenAI documents them as not to be edited by hand.

**Decision**: Codex memories are out of scope for the ctx memory bridge

**Rationale**: There is no stable, documented file contract to mirror; mirroring opaque generated state would be fragile and the feature is opt-in and disabled by default.

**Consequence**: ctx setup codex delivers hooks, MCP, skills, AGENTS.md, and journal import, but not a memory bridge. Revisit when OpenAI documents the memory file format.

---

## [2026-08-23-120839] ctx never parses Codex config.toml: it appends the [mcp_servers.ctx] table and scans header lines

**Status**: Accepted

**Context**: ctx setup codex --write must register the ctx MCP server in .codex/config.toml and detect whether the ctx plugin is enabled in ~/.codex/config.toml. ctx has no TOML dependency; both files are user-owned and carry comments and ordering.

**Decision**: ctx never parses Codex config.toml: it appends the [mcp_servers.ctx] table and scans header lines

**Rationale**: A read-modify-write through a TOML library would drop comments and reorder tables in a file the user edits by hand. Appending a table header at EOF is always valid TOML, and skipping when the exact header line already exists is sufficient for idempotency. Detection only needs the plugin table header and its enabled key. This keeps go.mod free of a new dependency for a narrow need.

**Consequence**: ctx does not update an existing [mcp_servers.ctx] body (the user owns it). Detection is a line scan, so unusual TOML (the header inside a multi-line string) could misdetect; acceptable for config files Codex itself writes. If ctx ever needs to rewrite Codex config, revisit with a TOML library.

---

## [2026-08-23-120839] Codex plugin root lives at internal/assets/codex with a repo marketplace at .agents/plugins/marketplace.json

**Status**: Accepted

**Context**: Codex 0.148 ships plugins (.codex-plugin/plugin.json + hooks/hooks.json + skills/ + .mcp.json) and repo marketplaces (.agents/plugins/marketplace.json) as stable features, with a hook contract that mirrors Claude Code's. ctx already delivers Claude support as a plugin rooted at internal/assets/claude, referenced by .claude-plugin/marketplace.json.

**Decision**: Codex plugin root lives at internal/assets/codex with a repo marketplace at .agents/plugins/marketplace.json

**Rationale**: Mirroring the Claude layout (plugin root under internal/assets/<tool>, marketplace at the repo root) gives one-command install (codex plugin marketplace add ActiveMemory/ctx; codex plugin add ctx@activememory-ctx) and lets the same embedded hooks.json/skills serve the project-local route (ctx setup codex --write). Putting it under internal/assets/integrations/ like Copilot CLI would have broken the marketplace source path convention (./internal/assets/<plugin-root>) and split the Claude/Codex symmetry.

**Consequence**: Two plugin roots must stay version-synced (make sync-version / check-version-sync cover both plus the Codex marketplace). Codex skills are generated from the Claude skills by hack/sync-codex-skills.sh (allowed-tools stripped; Claude-only skills excluded) and guarded by make check-codex-skills, the same way Copilot CLI skills are.

---

## [2026-07-25-190410] Beyond a byte ceiling, knowledge content should become tooling, not more Markdown

**Status**: Accepted

**Context**: Designing M5's heavy-page signal; the failure case is a convention theme file growing toward ~1MB.

**Decision**: Beyond a byte ceiling, knowledge content should become tooling, not more Markdown

**Rationale**: An LLM is a poor linter: it cannot reliably apply a large ruleset expressed as prose. Past a ceiling, recursive folding (more Markdown) is the wrong fix; the right move is extracting the content to actual tooling (a linter), code, or docs. So the heavy-page remedy is an advisory 'split OR extract to tooling', decided by the human, not tier-2 auto-fold.

**Consequence**: The heavy-page nudge frames the ceiling as a smell ('a context file this heavy is a linter in prose'), prompting extraction rather than endless subdivision. Tier-2 theme recursion remains deferred, not precluded.

---

## [2026-07-25-190410] M5 knowledge health is two suggest-only signals: foldable root (staging count) and heavy page (bytes)

**Status**: Accepted

**Context**: M5 wires the growth nudge, /ctx-remember, and /ctx-wrap-up to progressive disclosure. The 2026-07-16 foundational decision already fixed staging-as-watermark, the three surfaces, and suggest-only. Left open: the signal/threshold, and that the original root-only measure was blind to theme-file bloat that folding itself creates.

**Decision**: M5 knowledge health is two suggest-only signals: foldable root (staging count) and heavy page (bytes)

**Rationale**: Two questions need two measures. 'Should I fold?' is a COUNT of staged entries/sections via disclosure.StagedEntries (suggest /ctx-digest). 'Is this page too heavy to be useful context?' is a BYTE count, scanned over the root AND every theme file under .context/<noun>/ (suggest split or extract). Folding relocates bulk into theme files, so a root-only signal is blind to the bloat it creates. One knowledge.Health function feeds all three surfaces so hook and skills cannot drift. No state file: staging plus on-disk bytes are self-describing.

**Consequence**: convention_line_count (200 lines) is replaced by convention_section_count (foldability) plus a shared theme_page_byte_ceiling (weight); a deliberate behavior change, framed as user education. Theme files are now scanned. Reuses disclosure.StagedEntries and the check-knowledge daily-throttle / log-first plumbing.

---

## [2026-07-25-132759] Theme declaration via the Themes section keyword, on all three canonical kinds

**Status**: Accepted

**Context**: pd-m4 needed a way to NAME a theme up front (ctx <kind> add --section Themes), distinct from folding entries into it. The section-vs-entry distinction in conventions surfaced the need.

**Decision**: Theme declaration via the Themes section keyword, on all three canonical kinds

**Rationale**: Content is '<name> — <gist>', split on the same em-dash separator the theme parser reads back, so it round-trips. It writes BOTH the gist bullet and the theme file, because a bullet without its file fails the gist-to-file pairing invariant and would leave a root ctx disclosure refuses to touch. Generalized to learnings/decisions since they also have a ## Themes region.

**Consequence**: New internal/write/theme package + disclosure.AddTheme; theme adds bypass the per-kind body-flag gate (a theme has only a name and gist); a new user-facing CLI surface across all three canonical kinds.

---

## [2026-07-25-132719] ctx convention add requires --section; no default; placeholders rejected (strict)

**Status**: Accepted

**Context**: The convention add-path (pd-m4). A default or catch-all section is where an undecided caller dumps every convention, defeating the H2-section grouping the digest pass folds on.

**Decision**: ctx convention add requires --section; no default; placeholders rejected (strict)

**Rationale**: Choosing the section IS the thinking; the CLI refuses to do it for the caller. Enforced strictly via validate.RejectPlaceholder, which rejects empty/whitespace AND the shipped placeholder set (TBD, n/a, none, pending, …). Deviates deliberately from the plan's original T16 ('just move the anchor').

**Consequence**: ctx convention add errors without a concrete --section; agents and scripts must name the target H2 section. Pinned by TestConventionAddRequiresSection (7 cases); a refused add leaves the root byte-identical.

---

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
