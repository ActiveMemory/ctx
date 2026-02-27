# Decisions Reference

Module-specific, already-shipped, and historical decisions moved from
DECISIONS.md to keep the main file within token budget. All entries
preserved verbatim.

---

## [2026-02-26-100000] Blog and content publishing architecture (consolidated)

**Status**: Accepted

**Consolidated from**: 3 decisions (2026-02-06 to 2026-02-17)

- Scattered themes that appear across multiple posts but were never the primary subject deserve standalone deep-dive posts. Publishing means updating dates, fixing paths, weaving cross-links, and adding an "Arc" section.
- Every blog post includes a "The Arc" section near the end that explicitly connects it to related posts in the series, making the blog a navigable web rather than a flat list.
- Drop ctx-journal-summarize skill because it duplicates ctx-blog. The blog skill serves both internal summary and public post use cases with a prompt tweak.

---

## [2026-02-26-100003] Documentation and navigation structure (consolidated)

**Status**: Accepted

**Consolidated from**: 4 decisions (2026-02-14 to 2026-02-21)

- Restructure docs nav sections with dedicated index pages (reference, operations, security) so all sections have icons on mobile and serve as lightweight landing pages.
- Add TL;DR admonitions to recipes longer than ~200 lines. A tip admonition after the intro surfaces the quick-start commands immediately; users who want depth still read the full page.
- Pair judgment recipes with mechanical recipes: mechanical recipes answer "how" but not "when" or "why." Index the judgment recipe before the mechanical one.
- Place Adopting ctx at nav position 3 (after Getting Started). Users with existing projects need the migration guide before reference material.

---

## [2026-02-27-002832] Naming and tool conventions (consolidated)

**Status**: Accepted

**Consolidated from**: 3 decisions (2026-02-21 to 2026-02-24)

- **Lowercase error strings**: Go convention is lowercase, no-punctuation error strings (ST1005). The CLI formats user-facing messages differently from returned errors. Spec wording adjusted for Go idiom.
- **Rename .contextrc to .ctxrc**: Tool identity should be consistent — .ctxrc follows the `.<tool>rc` convention (.npmrc, .bashrc). All source, tests, docs, and specs reference .ctxrc; historical records retain the old name.
- **Drop ctx- prefix on project-level skills**: Project-level skills (.claude/skills/) use plain names; only plugin skills (ctx:ctx-*) use the ctx- namespace. Renamed ctx-borrow to absorb as first instance.

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

## [2026-02-24-013829] Add token usage to journal frontmatter, not build a usage dashboard

**Status**: Accepted

**Context**: Discussed whether to build token analytics features or keep it minimal

**Decision**: Add token usage to journal frontmatter, not build a usage dashboard

**Rationale**: ccusage and cmonitor already solve the dashboard problem; ctx's core value is persistent context, not usage analytics. Token metadata in journal entries serves session archaeology without scope creep.

**Consequences**: Token/model fields are auto-populated at export time; enrichment skill is documented not to overwrite them. Users wanting dashboards are pointed to ccusage.

---

## [2026-02-22-120012] Recall system design (consolidated)

**Status**: Accepted

**Consolidated from**: 4 decisions (2026-01-28 to 2026-02-20)

**Context**: The recall system parses AI session history from JSONL files and exports enriched markdown to the journal.

**Decision**: Claude-first with tool-agnostic types; default export preserves enrichment (no `--update` flag needed); spec-driven development supersedes ad-hoc bug-fix tasks.

**Rationale**: Claude Code is primary audience; parser updates follow its releases. Tool-agnostic `SessionParser` interface enables future parsers. Default export already preserved frontmatter — the real fix was `--force` behavior. `specs/recall-export-safety.md` replaced 4 narrow tasks with 7 comprehensive spec-aligned tasks.

**Consequences**: Features assume Claude Code conventions. Parser registry auto-detects format. Export has safe defaults with `--regenerate` opt-in. Aider/Cursor parsers are community-contributed, best-effort.

---

## [2026-02-21-195839] Secure-by-default dev server: localhost bind with opt-in LAN targets

**Status**: Accepted

**Context**: dev_addr = 0.0.0.0:8000 was added to both zensical.toml files, binding the dev server to all interfaces — incompatible with ctx secure-by-default stance

**Decision**: Secure-by-default dev server: localhost bind with opt-in LAN targets

**Rationale**: Removed dev_addr from committed config (zensical defaults to localhost:8000). Added make site-serve-lan and make journal-serve-lan targets that pass -a 0.0.0.0:8000 via CLI flag. Avoids modifying config files at runtime and keeps the opt-in explicit

**Consequences**: make site-serve and make journal-serve are safe by default. LAN access requires deliberate make *-lan invocation. journal-serve-lan calls zensical directly (bypasses ctx journal site --serve) because the Go code does not pass through extra flags

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

## [2026-02-27-002833] Skill design philosophy (consolidated)

**Status**: Accepted

**Consolidated from**: 2 decisions (2026-02-04 to 2026-02-14)

- **Skills over CLI for judgment-heavy workflows**: Borrow-from-the-future (now /absorb) implemented as skill, not CLI command. Workflows requiring interactive judgment (conflict resolution, selective application, strategy selection) need agent adaptability; CLI flags cannot cover edge cases.
- **E/A/R classification standard**: Expert/Activation/Redundant taxonomy for skill evaluation. Good Skill = Expert Knowledge - What Claude Already Knows. Target: >70% Expert, <10% Redundant. All skills evaluated against this framework.

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
