# Decisions

<!-- INDEX:START -->
| Date | Decision |
|------|--------|
| 2026-02-26 | Blog and content publishing architecture (consolidated) |
| 2026-02-26 | Hook and notification design (consolidated) |
| 2026-02-26 | ctx init and CLAUDE.md handling (consolidated) |
| 2026-02-26 | Documentation and navigation structure (consolidated) |
| 2026-02-26 | Task and knowledge management (consolidated) |
| 2026-02-26 | Agent autonomy and separation of concerns (consolidated) |
| 2026-02-26 | Security and permissions (consolidated) |
| 2026-02-26 | All webhook payloads must include session_id |
| 2026-02-26 | Context-load-gate v2: auto-injection replaces directive-based compliance |
| 2026-02-26 | .context/state/ directory for project-scoped runtime state |
| 2026-02-26 | Use imperative framing for context load gate hooks |
| 2026-02-24 | Notify events are opt-in, not opt-out |
| 2026-02-24 | Document worktree limitations rather than reroute paths |
| 2026-02-24 | RSS/Atom feed is infrastructure, not a user feature |
| 2026-02-24 | DETAILED_DESIGN.md lives outside FileReadOrder |
| 2026-02-24 | Use lowercase error strings in ctx remind for Go conventions |
| 2026-02-24 | Add token usage to journal frontmatter, not build a usage dashboard |
| 2026-02-22 | Webhook URL encrypted with shared encryption key, not a dedicated key |
| 2026-02-22 | Journal site rendering architecture (consolidated) |
| 2026-02-22 | Plugin and skill distribution architecture (consolidated) |
| 2026-02-22 | Recall system design (consolidated) |
| 2026-02-21 | Rename .contextrc to .ctxrc for tool-name consistency |
| 2026-02-21 | Secure-by-default dev server: localhost bind with opt-in LAN targets |
| 2026-02-21 | Drop ctx- prefix on project-level skills |
| 2026-02-19 | Smart retrieval: budget-aware ctx agent |
| 2026-02-19 | Try-decrypt-first for pad merge format auto-detection |
| 2026-02-15 | allow_outside_cwd belongs in .contextrc, not just CLI |
| 2026-02-15 | Hook output patterns are a reference catalog, not an implementation backlog |
| 2026-02-14 | Borrow-from-the-future implemented as skill, not CLI command |
| 2026-02-13 | Spec-first planning for non-trivial features |
| 2026-02-04 | E/A/R classification as the standard for skill evaluation |
<!-- INDEX:END -->

## [2026-02-26-200000] Context-load-gate v2: auto-injection replaces directive-based compliance

**Status**: Accepted

Soft instructions ("read these files") have a ~75-85% compliance ceiling because "don't apply judgment to this rule" is itself evaluated by judgment. Every instruction passes through the same attention/evaluation pipeline it's trying to override.

The v2 context-load-gate reads context files in the hook itself and injects content directly via `additionalContext`. The agent never chooses whether to comply — the content is already in its context window. This moves enforcement from the reasoning layer (subject to judgment) to the infrastructure layer (not subject to evaluation).

Injection strategy: CONSTITUTION, CONVENTIONS, ARCHITECTURE, AGENT_PLAYBOOK verbatim; DECISIONS, LEARNINGS index-only; TASKS mention-only; GLOSSARY skipped. Total ~7,700 tokens.

See: `specs/context-load-gate-v2.md`, `docs/blog/2026-02-25-the-homework-problem.md`.

---

## [2026-02-26-200001] .context/state/ directory for project-scoped runtime state

**Status**: Accepted

New gitignored directory under `context_dir` resolution for ephemeral project-scoped state. Follows `.context/logs/` precedent — added to `config.GitignoreEntries` and root `.gitignore`.

First use: injection oversize flag written by context-load-gate when injected tokens exceed the configurable `injection_token_warn` threshold (`.ctxrc`, default 15000). The check-context-size VERBATIM hook reads the flag and nudges the user to run `/ctx-consolidate`.

See: `specs/injection-oversize-nudge.md`.

---

## [2026-02-26-100000] Blog and content publishing architecture (consolidated)

**Status**: Accepted

**Consolidated from**: 3 decisions (2026-02-06 to 2026-02-17)

- Scattered themes that appear across multiple posts but were never the primary subject deserve standalone deep-dive posts. Publishing means updating dates, fixing paths, weaving cross-links, and adding an "Arc" section.
- Every blog post includes a "The Arc" section near the end that explicitly connects it to related posts in the series, making the blog a navigable web rather than a flat list.
- Drop ctx-journal-summarize skill because it duplicates ctx-blog. The blog skill serves both internal summary and public post use cases with a prompt tweak.

---

## [2026-02-26-100001] Hook and notification design (consolidated)

**Status**: Accepted

**Consolidated from**: 4 decisions (2026-02-12 to 2026-02-24)

- Tone down proactive content suggestion claims in docs rather than add more hooks. Already have 9 UserPromptSubmit hooks; adding another risks fatigue. Conversational prompting already works.
- Hook commands must use structured JSON output (hookSpecificOutput.additionalContext) instead of plain text, because Claude Code treats plain text as ignorable ambient context.
- Drop prompt-coach hook entirely: zero useful tips fired, output channel invisible to user, orphan temp file accumulation. The prompting guide already covers best practices.
- De-emphasize /ctx-journal-normalize from the default journal pipeline. The normalize skill is expensive and nondeterministic; programmatic normalization handles most cases. Skill remains available for targeted per-file use.

---

## [2026-02-26-100002] ctx init and CLAUDE.md handling (consolidated)

**Status**: Accepted

**Consolidated from**: 3 decisions (2026-01-20)

- `ctx init` handles CLAUDE.md intelligently: creates if missing, backs up and offers merge if existing, uses marker comment for idempotency. The `--merge` flag enables non-interactive append.
- `ctx init` always generates `.claude/hooks/` alongside `.context/` with no flag needed. Other AI tools ignore `.claude/`; Claude Code users get seamless zero-config experience.
- Core tool stays generic and tool-agnostic, with optional Claude Code enhancements via `.claude/hooks/`. Other AI tools can be supported similarly (`ctx hook cursor`, etc.).

---

## [2026-02-26-100003] Documentation and navigation structure (consolidated)

**Status**: Accepted

**Consolidated from**: 4 decisions (2026-02-14 to 2026-02-21)

- Restructure docs nav sections with dedicated index pages (reference, operations, security) so all sections have icons on mobile and serve as lightweight landing pages.
- Add TL;DR admonitions to recipes longer than ~200 lines. A tip admonition after the intro surfaces the quick-start commands immediately; users who want depth still read the full page.
- Pair judgment recipes with mechanical recipes: mechanical recipes answer "how" but not "when" or "why." Index the judgment recipe before the mechanical one.
- Place Adopting ctx at nav position 3 (after Getting Started). Users with existing projects need the migration guide before reference material.

---

## [2026-02-26-100004] Task and knowledge management (consolidated)

**Status**: Accepted

**Consolidated from**: 4 decisions (2026-01-27 to 2026-02-18)

- Tasks must include explicit deliverables, not just implementation steps. Parent tasks define WHAT the user gets; subtasks decompose HOW to build it. Without explicit deliverables, AI optimizes for checking boxes.
- Use reverse-chronological order (newest first) for DECISIONS.md and LEARNINGS.md. Ensures most recent items are read first regardless of token budget.
- Add quick reference index to DECISIONS.md: compact table at top allows scanning; agents can grep for full timestamp to jump to entry. Auto-updated on `ctx add decision`.
- Knowledge scaling via archive path for decisions and learnings: follow the task archive pattern, move old entries to `.context/archive/`, extend `ctx compact --archive` to cover all three file types.

---

## [2026-02-26-100005] Agent autonomy and separation of concerns (consolidated)

**Status**: Accepted

**Consolidated from**: 3 decisions (2026-01-21 to 2026-01-28)

- Removed AGENTS.md from project root. Consolidated on CLAUDE.md (auto-loaded) + .context/AGENT_PLAYBOOK.md as the canonical agent instruction path. Projects using ctx should not create AGENTS.md.
- Separate orchestrator directive from agent tasks: `.context/TASKS.md` is the agent's mind (tasks the agent owns); `IMPLEMENTATION_PLAN.md` is the orchestrator's thin directive layer ("check your tasks"). Prevents task list drift.
- No custom UI -- IDE is the interface. UI is a liability; IDEs already excel at file browsing, search, markdown editing, and git integration. Focus CLI efforts on good markdown output.

---

## [2026-02-26-100006] Security and permissions (consolidated)

**Status**: Accepted

**Consolidated from**: 4 decisions (2026-01-21 to 2026-02-24)

- Keep CONSTITUTION.md minimal: only truly inviolable rules (security, correctness, process invariants). Style preferences go in CONVENTIONS.md. Overly strict constitution gets ignored.
- Centralize constants with semantic prefixes in `internal/config/config.go`: `Dir*` for directories, `File*` for paths, `Filename*` for names, `UpdateType*` for entry types. Single source of truth, compile-time typo checks.
- Hooks use `ctx` from PATH, not hardcoded absolute paths. Standard Unix practice; portable across machines/users. `ctx init` checks PATH availability before proceeding.
- Drop absolute-path-to-ctx regex from block-dangerous-commands shell script. The block-non-path-ctx Go subcommand already covers this with better patterns; duplicating creates two sources of truth.

---

## [2026-02-26-000056] All webhook payloads must include session_id

**Status**: Accepted

**Context**: Multiple agents running concurrently produce identical webhook events with no way to attribute which session triggered which hook

**Decision**: All webhook payloads must include session_id

**Rationale**: Session ID is available on stdin for every hook invocation. Reading it even when the hook logic does not need it costs nothing and enables multi-agent diagnostics

**Consequences**: All run functions now take stdin parameter. Tests pass createTempStdin. Webhook payloads always have session_id populated when Claude Code provides one

---

## [2026-02-26-000052] Use imperative framing for context load gate hooks

**Status**: Accepted

**Context**: Advisory framing (Read your context files before proceeding) allowed agents to assess relevance and skip files, defeating the gate purpose

**Decision**: Use imperative framing for context load gate hooks

**Rationale**: Imperative framing (STOP. You must read...) with unconditional compliance checkpoint (Context Loaded block) removes the relevance assessment escape hatch. Verbatim relay becomes fallback safety net, not primary instruction

**Consequences**: Gate message is more assertive. Agents must always output Context Loaded block. Double failure (no block, no relay) is the only unobservable failure mode

---

## [2026-02-24-032946] Notify events are opt-in, not opt-out

**Status**: Accepted

**Context**: Bug found where relay messages fired without any .ctxrc configuration — EventAllowed treated nil/empty as allow all

**Decision**: Notify events are opt-in, not opt-out

**Rationale**: The correct default for notifications is silence. Users must explicitly list events in notify.events to receive them. This matches the principle that new features should not change behavior for existing users.

**Consequences**: EventAllowed returns false for nil/empty lists. All docs updated to reflect opt-in. ctx notify test bypasses the filter as a special case for connectivity checks.

---

## [2026-02-24-030254] Document worktree limitations rather than reroute paths

**Status**: Accepted

**Context**: Considered 3 engineering fixes (git toplevel, git-common-dir, worktree-aware split) to make untracked files resolve to the main checkout

**Decision**: Document worktree limitations rather than reroute paths

**Rationale**: All three add complexity for edge cases that workflow naturally avoids. Pad is key-gated already. Journal enrichment belongs on main after merge. Only notify is a real gap, and it deserves its own targeted fix.

**Consequences**: Worktree limitations documented in skill doc, parallel-worktrees recipe, scratchpad reference, and webhook-notifications recipe. Separate task filed for enabling notify in worktrees.

---

## [2026-02-24-025505] RSS/Atom feed is infrastructure, not a user feature

**Status**: Accepted

**Context**: Researched zensical RSS support (none found), discussed whether RSS matters for ctx.ist

**Decision**: RSS/Atom feed is infrastructure, not a user feature

**Rationale**: RSS serves as a replication protocol, zero-auth public API, and automation glue — targeting power users and future machine consumers, not casual readers

**Consequences**: Will implement as static Atom 1.0 feed generated at build time; spec captured at specs/rss-feed.md

---

## [2026-02-24-015939] DETAILED_DESIGN.md lives outside FileReadOrder

**Status**: Accepted

**Context**: Designing the /ctx-map skill output documents — needed to decide where DETAILED_DESIGN.md fits in the context loading pipeline

**Decision**: DETAILED_DESIGN.md lives outside FileReadOrder

**Rationale**: DETAILED_DESIGN.md is a deep per-module reference that can grow large. Loading it at session start would waste token budget. Agents consult specific sections on-demand when working on a module.

**Consequences**: ARCHITECTURE.md remains the session-start overview (~4000 tokens). DETAILED_DESIGN.md is never auto-loaded — agents must explicitly Read relevant sections. Two-tier documentation: succinct map vs. deep reference.

---

## [2026-02-24-015825] Use lowercase error strings in ctx remind for Go conventions

**Status**: Accepted

**Context**: The spec used 'No reminder with ID %d.' but golangci-lint flags capitalized error strings (ST1005) and trailing punctuation

**Decision**: Use lowercase error strings in ctx remind for Go conventions

**Rationale**: Go convention is lowercase, no-punctuation error strings; the CLI can format user-facing messages differently from returned errors

**Consequences**: Error messages match the rest of the codebase; the spec's exact wording is adjusted for Go idiom

---

## [2026-02-24-013829] Add token usage to journal frontmatter, not build a usage dashboard

**Status**: Accepted

**Context**: Discussed whether to build token analytics features or keep it minimal

**Decision**: Add token usage to journal frontmatter, not build a usage dashboard

**Rationale**: ccusage and cmonitor already solve the dashboard problem; ctx's core value is persistent context, not usage analytics. Token metadata in journal entries serves session archaeology without scope creep.

**Consequences**: Token/model fields are auto-populated at export time; enrichment skill is documented not to overwrite them. Users wanting dashboards are pointed to ccusage.

---

## [2026-02-22-101958] Webhook URL encrypted with shared encryption key, not a dedicated key

**Status**: Accepted

**Context**: ctx notify needs to encrypt webhook URLs. A new key per feature adds complexity.

**Decision**: Webhook URL encrypted with shared encryption key, not a dedicated key

**Rationale**: Reusing .context.key keeps the key management surface area minimal -- one key, one gitignore entry, one rotation cycle. The notify feature is a peer of the scratchpad (both store user secrets encrypted at rest).

**Consequences**: Key rename from .scratchpad.key to .context.key is now a follow-up task. Rotating the encryption key requires re-running ctx notify setup.

---

## [2026-02-22-120010] Journal site rendering architecture (consolidated)

**Status**: Accepted

**Consolidated from**: 5 decisions (2026-02-20)

**Context**: Journal site rendering required multiple architectural decisions to handle tool output, title formatting, and content normalization.

**Decision**: Journal site uses HTML-escaped `<pre><code>` blocks for tool output wrapping, code-level `normalizeContent` pipeline for rendering, CSS overflow for visual containment, and 75-char title limit.

**Rationale**: Fenced code blocks were tried first (survive blank lines, prevent markdown interpretation) but inner content conflicts remained. Switching to pre/code with HTML escaping ("defencify") eliminated all conflicts. pymdownx.highlight required `use_pygments=false`. CSS `max-height + overflow-y: auto` replaced `<details>` (which is a Type 6 HTML block incompatible with fenced code). AI normalization found zero issues across 290 files — code pipeline handles everything at build time.

**Consequences**: Journal site ships `docs/stylesheets/extra.css`. normalize.go is dramatically simpler. Title truncation at 75 chars (`RecallMaxTitleLen` in `config/limit.go`) applied in three places. AI normalization reserved for specific files only.

---

## [2026-02-22-120011] Plugin and skill distribution architecture (consolidated)

**Status**: Accepted

**Consolidated from**: 4 decisions (2026-02-16)

**Context**: ctx v0.6.0 converted from per-project shell hooks to a Go-based plugin model distributed via Claude Code marketplace.

**Decision**: Go subcommands (`ctx system *`) replace shell hooks; `internal/assets/claude/skills/` is the single source of truth for distributed skills; no symlinks for cross-directory sharing; permission docs match `DefaultClaudePermissions` exactly.

**Rationale**: Go subcommands eliminate jq dependency and enable `go test`. Single source prevents duplicate skill entries. Symlinks break on Windows without Developer Mode. Permission doc mismatches confuse users when `ctx init` seeds more than docs recommend.

**Consequences**: `ctx init` no longer creates `.claude/hooks/` or `.claude/skills/`. Existing projects need plugin installation. `.claude/skills/` holds only dev-only skills. Future skill additions must update both `config/file.go` and the recipe.

---

## [2026-02-22-120012] Recall system design (consolidated)

**Status**: Accepted

**Consolidated from**: 4 decisions (2026-01-28 to 2026-02-20)

**Context**: The recall system parses AI session history from JSONL files and exports enriched markdown to the journal.

**Decision**: Claude-first with tool-agnostic types; default export preserves enrichment (no `--update` flag needed); spec-driven development supersedes ad-hoc bug-fix tasks.

**Rationale**: Claude Code is primary audience; parser updates follow its releases. Tool-agnostic `SessionParser` interface enables future parsers. Default export already preserved frontmatter — the real fix was `--force` behavior. `specs/recall-export-safety.md` replaced 4 narrow tasks with 7 comprehensive spec-aligned tasks.

**Consequences**: Features assume Claude Code conventions. Parser registry auto-detects format. Export has safe defaults with `--regenerate` opt-in. Aider/Cursor parsers are community-contributed, best-effort.

---

## [2026-02-21-200037] Rename .contextrc to .ctxrc for tool-name consistency

**Status**: Accepted

**Context**: The RC file was called .contextrc but the CLI tool is ctx. Users saw the mismatch in docs and help text.

**Decision**: Rename .contextrc to .ctxrc for tool-name consistency

**Rationale**: Tool identity should be consistent — the file a user creates should match the tool they invoke. .ctxrc follows the .<tool>rc convention (.npmrc, .bashrc).

**Consequences**: All Go source, tests, docs, specs, and context files now reference .ctxrc. Historical records (blog posts, released specs, decision log) retain the old name as accurate history. A canonical .ctxrc template exists at project root. A new docs/configuration.md page provides dedicated config reference.

---

## [2026-02-21-195839] Secure-by-default dev server: localhost bind with opt-in LAN targets

**Status**: Accepted

**Context**: dev_addr = 0.0.0.0:8000 was added to both zensical.toml files, binding the dev server to all interfaces — incompatible with ctx secure-by-default stance

**Decision**: Secure-by-default dev server: localhost bind with opt-in LAN targets

**Rationale**: Removed dev_addr from committed config (zensical defaults to localhost:8000). Added make site-serve-lan and make journal-serve-lan targets that pass -a 0.0.0.0:8000 via CLI flag. Avoids modifying config files at runtime and keeps the opt-in explicit

**Consequences**: make site-serve and make journal-serve are safe by default. LAN access requires deliberate make *-lan invocation. journal-serve-lan calls zensical directly (bypasses ctx journal site --serve) because the Go code does not pass through extra flags

---

## [2026-02-21-195818] Drop ctx- prefix on project-level skills

**Status**: Accepted

**Context**: The ctx- prefix on .claude/skills/ctx-borrow made it look like a ctx plugin skill when it's a generic project-level utility

**Decision**: Drop ctx- prefix on project-level skills

**Rationale**: Project-level skills (.claude/skills/) should have plain names; only plugin skills (ctx:ctx-*) use the ctx- namespace

**Consequences**: Future project-level skills use plain names (e.g., absorb, audit, backup). Renamed ctx-borrow to absorb as first instance.

---

## [2026-02-19-192630] Smart retrieval: budget-aware ctx agent

**Status**: Accepted

**Context**: Issue #19 identified that ctx agent --budget is cosmetic — LEARNINGS.md excluded, decisions title-only, no relevance filtering, no graceful degradation

**Decision**: Smart retrieval: budget-aware ctx agent

**Rationale**: Phase 1 (smart retrieval) has the highest impact with no file format changes. Scoring entries by recency and task relevance, with tier-based budget allocation, solves the scaling problem at the presentation layer

**Consequences**: ctx agent output becomes richer (learnings, decision bodies) and budget-aware. Packet struct gains new fields (additive, backward compatible). New files: score.go, budget.go in internal/cli/agent/

---

## [2026-02-19-214858] Try-decrypt-first for pad merge format auto-detection

**Status**: Accepted

**Context**: Pad merge needs to handle both encrypted (.enc) and plaintext (.md) scratchpad files without requiring the user to specify format. Considered file extension matching, UTF-8 heuristics, and try-decrypt-first.

**Decision**: Try-decrypt-first for pad merge format auto-detection

**Rationale**: AES-256-GCM is self-authenticating — wrong key always fails cleanly. This makes try-decrypt a reliable discriminator with zero ambiguity. Fall back to plaintext on failure, with a UTF-8 validity warning to catch encrypted files mistakenly parsed as text.

**Consequences**: No --format flag needed. Users can mix encrypted and plaintext files in a single merge call. Foreign encrypted files with wrong key fall back gracefully instead of aborting.

---

## [2026-02-15-231015] allow_outside_cwd belongs in .contextrc, not just CLI

**Status**: Accepted

**Context**: External context recipe claimed .contextrc could persist the boundary override, but the field didn't exist. Choice: fix the docs or make the promise true.

**Decision**: allow_outside_cwd belongs in .contextrc, not just CLI

**Rationale**: If a user already declared context_dir pointing outside the project, requiring --allow-outside-cwd on every command is redundant ceremony. .contextrc is configure-once-forget-about-it — the boundary flag should live there too.

**Consequences**: New allow_outside_cwd bool field in CtxRC. PersistentPreRun checks both the CLI flag and .contextrc. Shell aliases (Option C) become optional rather than necessary.

---

## [2026-02-15-170006] Hook output patterns are a reference catalog, not an implementation backlog

**Status**: Accepted

**Context**: Patterns 6-8 in hook-output-patterns.md (conditional relay, suggested action, escalating severity) were initially framed as 'not yet implemented' which implied planned work. Analysis showed all three are either already used in practice (suggested action appears in check-journal.sh, check-backup-age.sh, block-non-path-ctx.sh; conditional relay is just bash if-then-else already in check-persistence.sh and check-journal.sh) or not justified by current need (escalating severity would require agent-side protocol training for a three-tier system when the existing two-tier silent/VERBATIM split covers all use cases).

**Decision**: Hook output patterns are a reference catalog, not an implementation backlog

**Rationale**: The recipe documents hook patterns for anyone writing hooks — it is not scoped to ctx-only patterns. Removing them would lose legitimate reference material. But framing them as 'not yet implemented' violated the ctx manifesto: not written means nonexistent, and there were no backing tasks. The patterns stay as equal entries in the catalog without implementation promises.

**Consequences**: Patterns 6-8 are presented as first-class patterns alongside 1-5, without a 'not yet implemented' section. No tasks created. If a concrete need arises for any of these patterns in ctx hooks, a task gets created at that point — not before.

---

## [2026-02-14-163859] Borrow-from-the-future implemented as skill, not CLI command

**Status**: Accepted

**Context**: Task proposed either /absorb skill or ctx borrow CLI command for merging deltas between two directories

**Decision**: Borrow-from-the-future implemented as skill, not CLI command

**Rationale**: The workflow requires interactive judgment: conflict resolution, selective file application, strategy selection between 3 tiers. An agent adapts to edge cases; CLI flags cannot.

**Consequences**: No ctx borrow subcommand. Users invoke /absorb in their AI tool. Non-AI users would need to manually run git diff/patch commands.

---

## [2026-02-13-133318] Spec-first planning for non-trivial features

**Status**: Accepted

**Context**: Designed ctx pad (encrypted scratchpad). Created spec, then tasks. Noticed the tasks alone wouldn't lead a future session to the spec.

**Decision**: Spec-first planning for non-trivial features

**Rationale**: Implementation sessions work from TASKS.md. If the spec isn't referenced there, the session builds from task summaries alone — incomplete context leads to design drift. Redundant references catch agents that skip ahead.

**Consequences**: All non-trivial features now follow: write specs/feature.md -> task out in TASKS.md with Phase header referencing the spec -> first task includes bold read-the-spec instruction. AGENT_PLAYBOOK.md updated with 'Planning Non-Trivial Work' section.

---

## [2026-02-11] Remove .context/sessions/ storage layer and ctx session command

**Status**: Accepted

**Context**: The session/recall/journal system had three overlapping storage layers: `~/.claude/projects/` (raw JSONL transcripts, owned by Claude Code), `.context/sessions/` (JSONL copies + context snapshots), and `.context/journal/` (enriched markdown from `ctx recall export`). The recall pipeline reads directly from `~/.claude/projects/`, making `.context/sessions/` a dead-end write sink that nothing reads from. The auto-save hook copied transcripts to a directory nobody consumed. The `ctx session save` command created context snapshots that git already provides through version history. This was ~15 Go source files, a shell hook, ~20 config constants, and 30+ doc references supporting infrastructure with no consumers.

**Decision**: Remove `.context/sessions/` entirely. Two stores remain: raw transcripts (global, tool-owned in `~/.claude/projects/`) and enriched journal (project-local in `.context/journal/`).

**Rationale**: Dead-end write sinks waste code surface, maintenance effort, and user attention. The recall pipeline already proved that reading directly from `~/.claude/projects/` is sufficient. Context snapshots are redundant with git history. Removing the middle layer simplifies the architecture from three stores to two, eliminates an entire CLI command tree (`ctx session`), and removes a shell hook that fired on every session end.

**Consequences**: Deleted `internal/cli/session/` (15 files), removed auto-save hook, removed `--auto-save` from watch, removed pre-compact auto-save from compact, removed `/ctx-save` skill, updated ~45 documentation files. Four earlier decisions superseded (SessionEnd hook, Auto-Save Before Compact, Session Filename Format, Two-Tier Persistence Model). Users who want session history use `ctx recall list/export` instead.

---

## [2026-02-04-230933] E/A/R classification as the standard for skill evaluation

**Status**: Accepted

**Context**: Reviewed ~30 external skill/prompt files; needed a systematic way to evaluate what to keep vs delete

**Decision**: E/A/R classification as the standard for skill evaluation

**Rationale**: Expert/Activation/Redundant taxonomy from judge.txt captures the key insight: Good Skill = Expert Knowledge - What Claude Already Knows. Gives a concrete target (>70% Expert, <10% Redundant)

**Consequences**: skill-creator SKILL.md updated with E/A/R as core principle. All future skills evaluated against this framework
