# Archived Decisions (consolidated 2026-02-27)

Originals replaced by consolidated entries in DECISIONS.md.

## Group: Context injection architecture (v2)

## [2026-02-26-182752] Extract diagrams from ARCHITECTURE.md into linked files to halve context injection budget

**Status**: Accepted

**Context**: context-load-gate was injecting ~20K tokens per session start. ARCHITECTURE.md contained 10 ASCII/Mermaid diagrams totaling ~600 lines. AGENT_PLAYBOOK.md had 3 sections duplicated in CLAUDE.md and CONVENTIONS.md, plus 3 overlapping persistence sections.

**Decision**: Extract diagrams from ARCHITECTURE.md into linked files to halve context injection budget

**Rationale**: Diagrams are reference material rarely needed at session start — agents need the structural understanding (verbal summary) not the visual representation. Duplicated content wastes tokens for zero incremental value. Linked files preserve diagrams for on-demand access.

**Consequences**: 5 new architecture-dia-*.md files in .context/ (not in FileReadOrder). ARCHITECTURE.md and AGENT_PLAYBOOK.md both slimmed. Embedded AGENT_PLAYBOOK.md template updated to match. Total injection dropped 53% (20K→9.5K tokens).

---

## [2026-02-26-200000] Context-load-gate v2: auto-injection replaces directive-based compliance

**Status**: Accepted

Soft instructions ("read these files") have a ~75-85% compliance ceiling because "don't apply judgment to this rule" is itself evaluated by judgment. Every instruction passes through the same attention/evaluation pipeline it's trying to override.

The v2 context-load-gate reads context files in the hook itself and injects content directly via `additionalContext`. The agent never chooses whether to comply — the content is already in its context window. This moves enforcement from the reasoning layer (subject to judgment) to the infrastructure layer (not subject to evaluation).

Injection strategy: CONSTITUTION, CONVENTIONS, ARCHITECTURE, AGENT_PLAYBOOK verbatim; DECISIONS, LEARNINGS index-only; TASKS mention-only; GLOSSARY skipped. Total ~7,700 tokens.

See: `specs/context-load-gate-v2.md`, `docs/blog/2026-02-25-the-homework-problem.md`.

---

## [2026-02-26-000052] Use imperative framing for context load gate hooks

**Status**: Accepted

**Context**: Advisory framing (Read your context files before proceeding) allowed agents to assess relevance and skip files, defeating the gate purpose

**Decision**: Use imperative framing for context load gate hooks

**Rationale**: Imperative framing (STOP. You must read...) with unconditional compliance checkpoint (Context Loaded block) removes the relevance assessment escape hatch. Verbatim relay becomes fallback safety net, not primary instruction

**Consequences**: Gate message is more assertive. Agents must always output Context Loaded block. Double failure (no block, no relay) is the only unobservable failure mode

---

## Group: Webhook and notification design

## [2026-02-26-000056] All webhook payloads must include session_id

**Status**: Accepted

**Context**: Multiple agents running concurrently produce identical webhook events with no way to attribute which session triggered which hook

**Decision**: All webhook payloads must include session_id

**Rationale**: Session ID is available on stdin for every hook invocation. Reading it even when the hook logic does not need it costs nothing and enables multi-agent diagnostics

**Consequences**: All run functions now take stdin parameter. Tests pass createTempStdin. Webhook payloads always have session_id populated when Claude Code provides one

---

## [2026-02-24-032946] Notify events are opt-in, not opt-out

**Status**: Accepted

**Context**: Bug found where relay messages fired without any .ctxrc configuration — EventAllowed treated nil/empty as allow all

**Decision**: Notify events are opt-in, not opt-out

**Rationale**: The correct default for notifications is silence. Users must explicitly list events in notify.events to receive them. This matches the principle that new features should not change behavior for existing users.

**Consequences**: EventAllowed returns false for nil/empty lists. All docs updated to reflect opt-in. ctx notify test bypasses the filter as a special case for connectivity checks.

---

## [2026-02-22-101958] Webhook URL encrypted with shared encryption key, not a dedicated key

**Status**: Accepted

**Context**: ctx notify needs to encrypt webhook URLs. A new key per feature adds complexity.

**Decision**: Webhook URL encrypted with shared encryption key, not a dedicated key

**Rationale**: Reusing .context.key keeps the key management surface area minimal -- one key, one gitignore entry, one rotation cycle. The notify feature is a peer of the scratchpad (both store user secrets encrypted at rest).

**Consequences**: Key rename from .scratchpad.key to .context.key is now a follow-up task. Rotating the encryption key requires re-running ctx notify setup.

---

## Group: Naming and tool conventions

## [2026-02-24-015825] Use lowercase error strings in ctx remind for Go conventions

**Status**: Accepted

**Context**: The spec used 'No reminder with ID %d.' but golangci-lint flags capitalized error strings (ST1005) and trailing punctuation

**Decision**: Use lowercase error strings in ctx remind for Go conventions

**Rationale**: Go convention is lowercase, no-punctuation error strings; the CLI can format user-facing messages differently from returned errors

**Consequences**: Error messages match the rest of the codebase; the spec's exact wording is adjusted for Go idiom

---

## [2026-02-21-200037] Rename .contextrc to .ctxrc for tool-name consistency

**Status**: Accepted

**Context**: The RC file was called .contextrc but the CLI tool is ctx. Users saw the mismatch in docs and help text.

**Decision**: Rename .contextrc to .ctxrc for tool-name consistency

**Rationale**: Tool identity should be consistent — the file a user creates should match the tool they invoke. .ctxrc follows the .<tool>rc convention (.npmrc, .bashrc).

**Consequences**: All Go source, tests, docs, specs, and context files now reference .ctxrc. Historical records (blog posts, released specs, decision log) retain the old name as accurate history. A canonical .ctxrc template exists at project root. A new docs/configuration.md page provides dedicated config reference.

---

## [2026-02-21-195818] Drop ctx- prefix on project-level skills

**Status**: Accepted

**Context**: The ctx- prefix on .claude/skills/ctx-borrow made it look like a ctx plugin skill when it's a generic project-level utility

**Decision**: Drop ctx- prefix on project-level skills

**Rationale**: Project-level skills (.claude/skills/) should have plain names; only plugin skills (ctx:ctx-*) use the ctx- namespace

**Consequences**: Future project-level skills use plain names (e.g., absorb, audit, backup). Renamed ctx-borrow to absorb as first instance.

---

## Group: Skill design philosophy

## [2026-02-14-163859] Borrow-from-the-future implemented as skill, not CLI command

**Status**: Accepted

**Context**: Task proposed either /absorb skill or ctx borrow CLI command for merging deltas between two directories

**Decision**: Borrow-from-the-future implemented as skill, not CLI command

**Rationale**: The workflow requires interactive judgment: conflict resolution, selective file application, strategy selection between 3 tiers. An agent adapts to edge cases; CLI flags cannot.

**Consequences**: No ctx borrow subcommand. Users invoke /absorb in their AI tool. Non-AI users would need to manually run git diff/patch commands.

---

## [2026-02-04-230933] E/A/R classification as the standard for skill evaluation

**Status**: Accepted

**Context**: Reviewed ~30 external skill/prompt files; needed a systematic way to evaluate what to keep vs delete

**Decision**: E/A/R classification as the standard for skill evaluation

**Rationale**: Expert/Activation/Redundant taxonomy from judge.txt captures the key insight: Good Skill = Expert Knowledge - What Claude Already Knows. Gives a concrete target (>70% Expert, <10% Redundant)

**Consequences**: skill-creator SKILL.md updated with E/A/R as core principle. All future skills evaluated against this framework

---
