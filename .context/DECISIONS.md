# Decisions

<!-- INDEX:START -->
| Date | Decision |
|------|--------|
| 2026-02-24 | Notify events are opt-in, not opt-out |
| 2026-02-24 | Document worktree limitations rather than reroute paths |
| 2026-02-24 | RSS/Atom feed is infrastructure, not a user feature |
| 2026-02-24 | DETAILED_DESIGN.md lives outside FileReadOrder |
| 2026-02-24 | Use lowercase error strings in ctx remind for Go conventions |
| 2026-02-24 | Drop absolute-path-to-ctx regex from block-dangerous-commands |
| 2026-02-24 | Add token usage to journal frontmatter, not build a usage dashboard |
| 2026-02-22 | De-emphasize /ctx-journal-normalize from default journal pipeline |
| 2026-02-22 | Hook commands use structured JSON output instead of plain text |
| 2026-02-22 | Webhook URL encrypted with shared encryption key, not a dedicated key |
| 2026-02-22 | Journal site rendering architecture (consolidated) |
| 2026-02-22 | Plugin and skill distribution architecture (consolidated) |
| 2026-02-22 | Recall system design (consolidated) |
| 2026-02-21 | Restructure docs nav sections with dedicated index pages |
| 2026-02-21 | Rename .contextrc to .ctxrc for tool-name consistency |
| 2026-02-21 | Secure-by-default dev server: localhost bind with opt-in LAN targets |
| 2026-02-21 | Drop ctx- prefix on project-level skills |
| 2026-02-19 | Smart retrieval: budget-aware ctx agent |
| 2026-02-19 | Try-decrypt-first for pad merge format auto-detection |
| 2026-02-18 | Knowledge scaling: archive path for decisions and learnings |
| 2026-02-15 | allow_outside_cwd belongs in .contextrc, not just CLI |
| 2026-02-15 | Add TL;DR admonitions to recipes longer than ~200 lines |
| 2026-02-15 | Hook output patterns are a reference catalog, not an implementation backlog |
| 2026-02-15 | Pair judgment recipes with mechanical recipes |
| 2026-02-14 | Place Adopting ctx at nav position 3 |
| 2026-02-14 | Borrow-from-the-future implemented as skill, not CLI command |
| 2026-02-13 | Spec-first planning for non-trivial features |
| 2026-02-12 | Drop prompt-coach hook |
| 2026-02-06 | Drop ctx-journal-summarize skill (duplicates ctx-blog) |
| 2026-02-04 | E/A/R classification as the standard for skill evaluation |
| 2026-01-29 | Add quick reference index to DECISIONS.md |
| 2026-01-28 | No custom UI - IDE is the interface |
| 2026-01-28 | Tasks must include explicit deliverables, not just implementation steps |
| 2026-01-27 | Use reverse-chronological order (newest first) for DECISIONS.md and LEARNINGS.md |
| 2026-01-25 | Removed AGENTS.md from project root |
| 2026-01-25 | Keep CONSTITUTION Minimal |
| 2026-01-25 | Centralize Constants with Semantic Prefixes |
| 2026-01-21 | Separate Orchestrator Directive from Agent Tasks |
| 2026-01-21 | Hooks Use ctx from PATH, Not Hardcoded Paths |
| 2026-01-20 | Handle CLAUDE.md Creation/Merge in ctx init |
| 2026-01-20 | Always Generate Claude Hooks in Init (No Flag Needed) |
| 2026-01-20 | Generic Core with Optional Claude Code Enhancements |
<!-- INDEX:END -->

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

## [2026-02-24-015709] Drop absolute-path-to-ctx regex from block-dangerous-commands

**Status**: Accepted

**Context**: Shell script had a regex blocking absolute paths to ctx binary. The block-non-path-ctx Go subcommand already covers this with better patterns (start-of-command, after-separator positions, test exception).

**Decision**: Drop absolute-path-to-ctx regex from block-dangerous-commands

**Rationale**: Duplicating it would create two sources of truth for the same rule.

**Consequences**: Only block-non-path-ctx owns absolute-path blocking; block-dangerous-commands focuses on sudo, git push, and bin-directory installs.

---

## [2026-02-24-013829] Add token usage to journal frontmatter, not build a usage dashboard

**Status**: Accepted

**Context**: Discussed whether to build token analytics features or keep it minimal

**Decision**: Add token usage to journal frontmatter, not build a usage dashboard

**Rationale**: ccusage and cmonitor already solve the dashboard problem; ctx's core value is persistent context, not usage analytics. Token metadata in journal entries serves session archaeology without scope creep.

**Consequences**: Token/model fields are auto-populated at export time; enrichment skill is documented not to overwrite them. Users wanting dashboards are pointed to ccusage.

---

## [2026-02-22-221724] De-emphasize /ctx-journal-normalize from default journal pipeline

**Status**: Accepted

**Context**: The journal pipeline previously prescribed normalize → enrich as the default workflow. With improved programmatic normalization during export and simplified markdown generation (no code fences), the AI-based normalize skill is rarely needed.

**Decision**: De-emphasize /ctx-journal-normalize from default journal pipeline

**Rationale**: The normalize skill is expensive (reads entire journal files through LLM), nondeterministic on large inputs, and blows up subagent context windows on non-ctx projects with millions of lines of session JSON. Programmatic normalization handles most cases.

**Consequences**: Normalize removed from relay nudge, make journal, skill prerequisites, and all docs. Skill remains available for targeted per-file use via /ctx-journal-normalize when rendering issues occur.

---

## [2026-02-22-194444] Hook commands use structured JSON output instead of plain text

**Status**: Accepted

**Context**: qa-reminder and post-commit hooks were being ignored despite firing correctly

**Decision**: Hook commands use structured JSON output instead of plain text

**Rationale**: JSON with hookSpecificOutput.additionalContext is parsed as a directive by Claude Code, while plain text is treated as ambient context the agent can ignore

**Consequences**: Added HookResponse/HookSpecificOutput types and printHookContext() helper to internal/cli/system/input.go; converted qa_reminder.go and post_commit.go; future hooks should use printHookContext() for non-blocking directives

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

## [2026-02-21-200038] Restructure docs nav sections with dedicated index pages

**Status**: Accepted

**Context**: Reference, Operations, and Security nav sections lacked icons in the mobile menu because they had no section index pages

**Decision**: Restructure docs nav sections with dedicated index pages

**Rationale**: Created reference/index.md, operations/index.md, security/index.md with linked summaries of sub-pages. Moved security.md to security/reporting.md to avoid file/directory name conflict. Renamed page titles to remove redundant ctx prefix (CLI, Skills, Tool Ecosystem).

**Consequences**: All nav sections now have icons on mobile. security.md URL changes to security/reporting/. Three internal links updated. Index pages serve as lightweight landing pages for each section.

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

## [2026-02-18-071514] Knowledge scaling: archive path for decisions and learnings

**Status**: Accepted

**Context**: DECISIONS.md and LEARNINGS.md grow monotonically with no archival path. Tasks have ctx tasks archive but knowledge files accumulate forever. Long-lived projects will hit token budget pressure and signal-to-noise decay.

**Decision**: Knowledge scaling: archive path for decisions and learnings

**Rationale**: Follow the existing task archive pattern. Move old entries to .context/archive/ files. Extend ctx compact --archive to cover all three file types. Add superseded-entry convention for decisions.

**Consequences**: New spec at specs/knowledge-scaling.md. Phase 5 tasks (P5.1-P5.7) added to TASKS.md. New CLI commands: ctx decisions archive, ctx learnings archive. New .contextrc keys: archive_knowledge_after_days, archive_keep_recent.

---

## [2026-02-17] Scattered themes deserve standalone blog posts when they haven't been dissected

**Status**: Accepted

**Context**: The "context as infrastructure" theme appeared across 5+ posts but was never the main topic. Similarly, the 3:1 ratio was mentioned but never analyzed. "Code is cheap, judgment is not" was implicit throughout but never stated. User feedback existed as raw notes but not as a narrative.

**Decision**: When a theme is scattered across the blog but never dissected as the primary subject, it deserves a standalone deep-dive post. The ideas/ drafts serve as raw material; publishing means: updating dates, fixing paths, weaving cross-links, and adding an "Arc" section.

**Rationale**: Scattered mentions create implicit understanding. A standalone post creates explicit, linkable, searchable understanding. The cross-link web strengthens both the new post and every post that referenced the theme.

**Consequences**: Published 4 posts in one session (3:1 Ratio, Code Is Cheap, Context as Infrastructure, When a System Starts Explaining Itself). Each required cross-linking to/from 3-6 companion posts. The blog now has a coherent arc with explicit connections.

---

## [2026-02-17] Blog arc structure: each post has an "Arc" section connecting to the series

**Status**: Accepted

**Context**: The blog series grew to 18+ posts. Each post was standalone but the narrative connections were implicit. Readers landing on one post couldn't see where it fit in the larger argument.

**Decision**: Every blog post includes a "The Arc" section near the end that explicitly connects it to related posts in the series, framing where this post sits in the broader narrative.

**Rationale**: The Arc section serves two purposes: (1) it helps readers navigate the series, and (2) it forces the author to articulate how each post relates to the whole, which improves coherence and catches thematic gaps.

**Consequences**: All new posts must include an Arc section. Existing posts gain Arc sections and "See also" links as they are cross-linked from new posts. The blog becomes a web, not a list.

---

## [2026-02-15-231015] allow_outside_cwd belongs in .contextrc, not just CLI

**Status**: Accepted

**Context**: External context recipe claimed .contextrc could persist the boundary override, but the field didn't exist. Choice: fix the docs or make the promise true.

**Decision**: allow_outside_cwd belongs in .contextrc, not just CLI

**Rationale**: If a user already declared context_dir pointing outside the project, requiring --allow-outside-cwd on every command is redundant ceremony. .contextrc is configure-once-forget-about-it — the boundary flag should live there too.

**Consequences**: New allow_outside_cwd bool field in CtxRC. PersistentPreRun checks both the CLI flag and .contextrc. Shell aliases (Option C) become optional rather than necessary.

---

## [2026-02-15-194828] Add TL;DR admonitions to recipes longer than ~200 lines

**Status**: Accepted

**Context**: Recipes bury the actionable pipeline at the bottom in Putting It All Together sections. Users must scroll past 300+ lines of explanation.

**Decision**: Add TL;DR admonitions to recipes longer than ~200 lines

**Rationale**: A tip admonition after the intro surfaces the quick-start commands immediately. Users who want depth still read the full page.

**Consequences**: 10 recipes now have TL;DRs. New recipes over ~200 lines should follow the pattern. Short recipes (permission-snapshots, scratchpad-with-claude) skip it.

---

## [2026-02-15-170006] Hook output patterns are a reference catalog, not an implementation backlog

**Status**: Accepted

**Context**: Patterns 6-8 in hook-output-patterns.md (conditional relay, suggested action, escalating severity) were initially framed as 'not yet implemented' which implied planned work. Analysis showed all three are either already used in practice (suggested action appears in check-journal.sh, check-backup-age.sh, block-non-path-ctx.sh; conditional relay is just bash if-then-else already in check-persistence.sh and check-journal.sh) or not justified by current need (escalating severity would require agent-side protocol training for a three-tier system when the existing two-tier silent/VERBATIM split covers all use cases).

**Decision**: Hook output patterns are a reference catalog, not an implementation backlog

**Rationale**: The recipe documents hook patterns for anyone writing hooks — it is not scoped to ctx-only patterns. Removing them would lose legitimate reference material. But framing them as 'not yet implemented' violated the ctx manifesto: not written means nonexistent, and there were no backing tasks. The patterns stay as equal entries in the catalog without implementation promises.

**Consequences**: Patterns 6-8 are presented as first-class patterns alongside 1-5, without a 'not yet implemented' section. No tasks created. If a concrete need arises for any of these patterns in ctx hooks, a task gets created at that point — not before.

---

## [2026-02-15-105923] Pair judgment recipes with mechanical recipes

**Status**: Accepted

**Context**: Created 'When to Use Agent Teams' as a decision-framework companion to the existing 'Parallel Worktrees' how-to recipe

**Decision**: Pair judgment recipes with mechanical recipes

**Rationale**: Mechanical recipes answer 'how' but not 'when' or 'why'. Users need judgment guidance to avoid misapplying powerful features. The same pattern applies to permissions (recipe + runbook) and drift (skill + permission drift section).

**Consequences**: New advanced features should ship with both a how-to recipe and a when-to-use guide. Index the judgment recipe before the mechanical one so users encounter the thinking before the doing.

---

## [2026-02-14-164103] Place Adopting ctx at nav position 3

**Status**: Accepted

**Context**: Adding migration/adoption guide to the docs site navigation

**Decision**: Place Adopting ctx at nav position 3

**Rationale**: After 'how do I install?' (Getting Started) the immediate next question for most users is 'I already have stuff, how do I add this?' Context Files is reference material that comes after adoption.

**Consequences**: New users with existing projects find the guide early in the nav flow. Getting Started remains the entry point for greenfield projects.

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

**Consequences**: All non-trivial features now follow: write specs/feature.md → task out in TASKS.md with Phase header referencing the spec → first task includes bold read-the-spec instruction. AGENT_PLAYBOOK.md updated with 'Planning Non-Trivial Work' section.

---

## [2026-02-12-005516] Drop prompt-coach hook

**Status**: Accepted

**Context**: Prompt-coach has been running since installation with zero useful tips fired. All counters across all state files are 0. The delivery mechanism is broken (stdout goes to AI not user, stderr is swallowed). Even if fixed with systemMessage, the coaching patterns are too narrow for experienced users and the prompting guide already covers best practices.

**Decision**: Drop prompt-coach hook

**Rationale**: Three layers of not-working: (1) patterns too narrow to match real prompts, (2) output channel invisible to user, (3) L-3 PID bug creates orphan temp files. Removing it eliminates the largest source of temp file accumulation, simplifies the hook stack, and removes dead code.

**Consequences**: One fewer hook in UserPromptSubmit (faster prompt submission). Eliminates prompt-coach temp file accumulation entirely — reduces cleanup burden. Need to remove: template script, config constant, script loader, hookScripts entry, settings.local.json reference, and active hook file.

---

## [2026-02-11] Remove .context/sessions/ storage layer and ctx session command

**Status**: Accepted

**Context**: The session/recall/journal system had three overlapping storage layers: `~/.claude/projects/` (raw JSONL transcripts, owned by Claude Code), `.context/sessions/` (JSONL copies + context snapshots), and `.context/journal/` (enriched markdown from `ctx recall export`). The recall pipeline reads directly from `~/.claude/projects/`, making `.context/sessions/` a dead-end write sink that nothing reads from. The auto-save hook copied transcripts to a directory nobody consumed. The `ctx session save` command created context snapshots that git already provides through version history. This was ~15 Go source files, a shell hook, ~20 config constants, and 30+ doc references supporting infrastructure with no consumers.

**Decision**: Remove `.context/sessions/` entirely. Two stores remain: raw transcripts (global, tool-owned in `~/.claude/projects/`) and enriched journal (project-local in `.context/journal/`).

**Rationale**: Dead-end write sinks waste code surface, maintenance effort, and user attention. The recall pipeline already proved that reading directly from `~/.claude/projects/` is sufficient. Context snapshots are redundant with git history. Removing the middle layer simplifies the architecture from three stores to two, eliminates an entire CLI command tree (`ctx session`), and removes a shell hook that fired on every session end.

**Consequences**: Deleted `internal/cli/session/` (15 files), removed auto-save hook, removed `--auto-save` from watch, removed pre-compact auto-save from compact, removed `/ctx-save` skill, updated ~45 documentation files. Four earlier decisions superseded (SessionEnd hook, Auto-Save Before Compact, Session Filename Format, Two-Tier Persistence Model). Users who want session history use `ctx recall list/export` instead.

---

## [2026-02-06-181708] Drop ctx-journal-summarize skill (duplicates ctx-blog)

**Status**: Accepted

**Context**: ctx-journal-summarize and ctx-blog both read journal entries over a time range and produce narrative summaries. The only difference was audience framing: internal summary vs public blog post.

**Decision**: Drop ctx-journal-summarize skill (duplicates ctx-blog)

**Rationale**: The blog skill can serve both use cases with a prompt tweak. One fewer skill to maintain, less surface area for drift.

**Consequences**: Removed skill dir, template, and references from integrations.md and two blog posts. Timeline narrative deferred item in TASKS.md marked as dropped. Users who want internal summaries use /ctx-blog instead.

---

## [2026-02-04-230933] E/A/R classification as the standard for skill evaluation

**Status**: Accepted

**Context**: Reviewed ~30 external skill/prompt files; needed a systematic way to evaluate what to keep vs delete

**Decision**: E/A/R classification as the standard for skill evaluation

**Rationale**: Expert/Activation/Redundant taxonomy from judge.txt captures the key insight: Good Skill = Expert Knowledge - What Claude Already Knows. Gives a concrete target (>70% Expert, <10% Redundant)

**Consequences**: skill-creator SKILL.md updated with E/A/R as core principle. All future skills evaluated against this framework

---

## [2026-01-29-044515] Add quick reference index to DECISIONS.md

**Status**: Accepted

**Context**: AI agents need to locate decisions quickly without reading the
entire file when context budget is limited

**Decision**: Add quick reference index to DECISIONS.md

**Rationale**: Compact table at top allows scanning; agents can grep for full
timestamp to jump to entry

**Consequences**: Index auto-updated on ctx add decision; ctx decisions
reindex for manual edits

---

## [2026-01-28-051426] No custom UI - IDE is the interface

**Status**: Accepted

**Context**: Considering whether to build a web/desktop UI for browsing
sessions, editing journal entries, and analytics. Export feature creates
editable markdown files.

**Decision**: No custom UI - IDE is the interface

**Rationale**: UI is a liability: maintenance burden, security surface,
dependencies. IDEs already excel at what we'd build: file browsing,
full-text search, markdown editing, git integration. Any UI we build either
duplicates IDE features poorly or becomes an IDE itself.

**Consequences**:
1) No UI codebase to maintain.
2) Users use their preferred editor.
3) Focus CLI efforts on good markdown output.
4) Analytics stays CLI-based (ctx recall stats).
5) **Non-technical users learn VS Code**.

---

## [2026-01-28-041239] Tasks must include explicit deliverables, not just implementation steps

**Status**: Accepted

**Context**: AI prematurely marked parent task complete after finishing
subtasks (internal parser library) but missing the actual deliverable
(CLI command and slash command). The task description said 'create a CLI
command and slash command' but subtasks only covered implementation details.

**Decision**: Tasks must include explicit deliverables, not just implementation
steps

**Rationale**: Subtasks decompose HOW to build something. The parent task
defines WHAT the user gets. Without explicit deliverables, AI optimizes for
checking boxes rather than delivering value. Task descriptions are indirect
prompts to the agent.

**Consequences**: 1. Parent tasks should state deliverable explicitly
(e.g., 'Deliverable: ctx recall list command'). 2. Consider acceptance criteria
checkboxes. 3. Update prompting guide with task-writing best practices.

---

## [2026-01-27-065902] Use reverse-chronological order (newest first) for DECISIONS.md and LEARNINGS.md

**Status**: Accepted

**Context**: With chronological order, oldest items consume tokens first, and
newest (most relevant) items risk being truncated when budget is tight. The AI
reads files from line 1 by default and has no way of knowing to read the
tail first.

**Decision**: Use reverse-chronological order (newest first) for DECISIONS.md
and LEARNINGS.md. Prepending is slightly awkward but more robust than relying
on AI cleverness to read file tails.

**Rationale**: Ensures most recent/relevant items are read first regardless of
token budget or whether AI uses ctx agent.

**Consequences**:
- `ctx add` must prepend instead of append
- File structure is self-documenting (newest = first)
- Works correctly regardless of how file is consumed

---

## [2026-01-25-220800] Removed AGENTS.md from project root

**Status**: Accepted

**Context**: AGENTS.md was not auto-loaded by any AI tool and created confusion
with redundant content alongside CLAUDE.md and .context/AGENT_PLAYBOOK.md.

**Decision**: Consolidated on CLAUDE.md + .context/AGENT_PLAYBOOK.md as the
canonical agent instruction path.

**Rationale**: Single source of truth; CLAUDE.md is auto-loaded by Claude Code,
AGENT_PLAYBOOK.md provides ctx-specific instructions.

**Consequences**: Projects using ctx should not create AGENTS.md.

---

## [2026-01-25-180000] Keep CONSTITUTION Minimal

**Status**: Accepted

**Context**: When codifying lessons learned, temptation was to add all
conventions to CONSTITUTION.md as "invariants."

**Decision**: CONSTITUTION.md contains only truly inviolable rules:
- Security invariants (secrets, path traversal)
- Correctness invariants (tests pass)
- Process invariants (decision records)

Style preferences and best practices go in CONVENTIONS.md instead.

**Rationale**:
- Overly strict constitution creates friction and gets ignored
- "Crying wolf" effect — developers stop reading it
- Conventions can be bent; constitution cannot
- Security vs style are fundamentally different categories

**Consequences**:
- CONVENTIONS.md becomes the living style guide
- CONSTITUTION.md stays short and scary
- New rules must pass "is this truly inviolable?" test

---

## [2026-01-25-170000] Centralize Constants with Semantic Prefixes

**Status**: Accepted (implemented)

**Context**: YOLO-mode feature development scattered magic strings across the
codebase. Same literals (`"TASKS.md"`, `"task"`, `".context"`) appeared in
10+ files. Human-guided refactoring session consolidated them.

**Decision**: All repeated literals go in `internal/config/config.go` with
semantic prefixes:
- `Dir*` for directories (`DirContext`, `DirArchive`, `DirSessions`)
- `File*` for file paths (`FileSettings`, `FileClaudeMd`)
- `Filename*` for file names only (`FilenameTask`, `FilenameDecision`)
- `UpdateType*` for entry types (`UpdateTypeTask`, `UpdateTypeDecision`)

Maps must use constants as keys:
```go
var FileType = map[string]string{
    UpdateTypeTask: FilenameTask,  // not "task": "TASKS.md"
}
```

**Rationale**:
- Single source of truth for all identifiers
- Refactoring is find-replace on constant name
- IDE navigation works (go-to-definition)
- Typos caught at compile time, not runtime
- Self-documenting code (constants have godoc)

**Consequences**:
- All new literals must go through config package
- Existing code migrated to use constants
- Slightly more verbose but much more maintainable

---

## [2026-01-21-140000] Separate Orchestrator Directive from Agent Tasks

**Status**: Accepted

**Context**: Two task systems existed: `IMPLEMENTATION_PLAN.md`
(Ralph Loop orchestrator) and `.context/TASKS.md` (ctx's own context).
Ralph would find IMPLEMENTATION_PLAN.md complete and exit,
ignoring .context/TASKS.md.

**Decision**: Clean separation of concerns:
- **`.context/TASKS.md`** = Agent's mind. Tasks the agent decided need doing.
- **`IMPLEMENTATION_PLAN.md`** = Orchestrator's directive.
  A single meta-task: "Check your tasks."

The orchestrator doesn't maintain a parallel ledger — it just tells the
agent to check its own mind.

**Rationale**:
- Agent autonomy: the agent owns its task list
- Single source of truth for tasks
- Orchestrator is minimal, not a micromanager
- Fresh `ctx init` deployments can have one directive: "Check .context/TASKS.md"
- Prevents task list drift between two files

**Consequences**:
- `PROMPT.md` now references `.context/TASKS.md` for task selection
- `IMPLEMENTATION_PLAN.md` becomes a thin directive layer
- Historical milestones are archived, not active tasks
- North Star goals live in IMPLEMENTATION_PLAN.md (meta-level, not tasks)

---

## [2026-01-21-120000] Hooks Use ctx from PATH, Not Hardcoded Paths

**Status**: Accepted (implemented)

**Context**: Original implementation hardcoded absolute paths in hooks
(e.g., `/home/parallels/WORKSPACE/ActiveMemory/dist/ctx-linux-arm64`).
This breaks when:
- Sharing configs with other developers
- Moving projects
- Dogfooding in separate directories

**Decision**:
1. Hooks use `ctx` from PATH (e.g., `ctx agent --budget 4000`)
2. `ctx init` checks if `ctx` is in PATH before proceeding
3. If not in PATH, init fails with clear instructions to install

**Rationale**:
- Standard Unix practice — tools should be in PATH
- Portable across machines/users
- Dogfooding becomes realistic (tests the real user experience)
- No manual path editing required

**Consequences**:
- Users must run `sudo make install` or equivalent before `ctx init`
- Tests need `CTX_SKIP_PATH_CHECK=1` env var to bypass check
- README must document PATH installation requirement

---

## [2026-01-20-180000] Handle CLAUDE.md Creation/Merge in ctx init

**Status**: Accepted (to be implemented)

**Context**: Both `claude init` and `ctx init` want to create/modify CLAUDE.md.
Users of ctx will likely want ctx's context-aware version,
but may already have a CLAUDE.md from `claude init`.

**Decision**: `ctx init` handles CLAUDE.md intelligently:
- **No CLAUDE.md exists** → Create it with ctx's context-loading template
- **CLAUDE.md exists** → Don't overwrite. Instead:
  1. **Backup first** → Copy to `CLAUDE.md.<unix_timestamp>.bak`
     (e.g., `CLAUDE.md.1737399000.bak`)
  2. Check if it already has ctx content (idempotent check via marker comment)
  3. If not, output the snippet to append and offer to merge
  4. `ctx init --merge` flag to auto-append without prompting

**Rationale**:
- Timestamped backups preserve history across multiple runs
- Unix timestamp is fine for backups (rarely read by humans, easy to sort)
- Respects user's existing CLAUDE.md customizations
- Doesn't silently overwrite important config
- Idempotency prevents duplicate content on re-runs

**Consequences**:
- Need to detect existing ctx content (marker comment like `<!-- ctx:context -->`)
- Backup files accumulate: `CLAUDE.md.<timestamp>.bak` (may want cleanup command later)
- Init output must clearly show what was created vs what needs manual merge
- Should work gracefully even if user runs `ctx init` multiple times

---

## [2026-01-20-100000] Always Generate Claude Hooks in Init (No Flag Needed)

**Status**: Accepted (to be implemented)

**Context**: Setting up Claude Code hooks manually is error-prone.
Considered `--claude` flag but realized it's unnecessary.

**Decision**: `ctx init` ALWAYS creates `.claude/hooks/` alongside `.context/`:
```bash
ctx init    # Creates BOTH .context/ AND .claude/hooks/
```

**Rationale**:
- Other AI tools (Cursor, Aider, Copilot) don't know/care about `.claude/`
- No downside to creating hooks that sit unused
- Claude Code users get seamless experience with zero extra steps
- If user later switches to Claude Code, hooks are already there
- Simpler UX - no flags to remember

**Consequences**:
- `ctx init` creates both directories always
- Hook scripts are embedded in binary (like templates)
- Need to detect platform for binary path in hooks
- `.claude/` becomes part of ctx's standard output

---

## [2026-01-20-080000] Generic Core with Optional Claude Code Enhancements

**Status**: Accepted

**Context**: `ctx` should work with any AI tool, but Claude Code users could
benefit from deeper integration (auto-load, auto-save via hooks).

**Decision**: Keep `ctx` generic as the core tool, but provide optional
Claude Code-specific enhancements:
- `ctx hook claude-code` generates Claude-specific configs
- `.claude/hooks/` contains Claude Code hook scripts
- Features work without Claude Code, but are enhanced with it

**Rationale**:
- Maintains tool-agnostic philosophy from core-architecture.md
- Doesn't lock users into Claude Code
- Claude Code users get seamless experience without extra work
- Other AI tools can be supported similarly (`ctx hook cursor`, etc.)

**Consequences**:
- Need to maintain both generic and Claude-specific documentation
- Hook scripts are optional, not required
- Testing must cover both with and without Claude Code
