# Archived Decisions (consolidated 2026-02-22)


Originals replaced by consolidated entries in DECISIONS.md.


## Group: Superseded session decisions (archived, not consolidated)



## [2026-01-20-200000] Use SessionEnd Hook for Auto-Save

**Status**: Superseded (removed in v0.4.0 — `.context/sessions/` eliminated; recall pipeline reads from `~/.claude/projects/` directly)

**Context**: Need to save context even when user exits abruptly (Ctrl+C).

**Decision**: Use Claude Code's `SessionEnd` hook to auto-save transcript:
- Hook fires on all exits including Ctrl+C
- Copies `transcript_path` to `.context/sessions/`
- Creates both .jsonl (raw) and .md (summary) files

**Rationale**:
- Catches all exit scenarios
- Transcript contains full conversation
- No user action required
- Graceful degradation (just doesn't save if hook fails)

**Consequences**:
- Only works with Claude Code (other tools need different approach)
- Requires jq for JSON parsing in hook script
- Session files are .jsonl format (need tooling to read)

---




## [2026-01-20-160000] Auto-Save Before Compact

**Status**: Superseded (removed in v0.4.0 — `.context/sessions/` eliminated; compact no longer writes session snapshots)

**Context**: `ctx compact` archives old tasks. Information could be
lost if not captured.

**Decision**: `ctx compact` should auto-save a session dump before archiving:
1. Save current state to `.context/sessions/YYYY-MM-DD-HHMM-pre-compact.md`
2. Then perform the compaction

**Rationale**:
- Safety net before destructive-ish operation
- User can always recover pre-compact state
- No extra user action required

**Consequences**:
- Compact command becomes slightly slower
- Sessions directory grows with each compact
- May want `--no-save` flag for automation

---




## [2026-01-20-140000] Session Filename Format: YYYY-MM-DD-HHMMSS-topic.md

**Status**: Superseded (removed in v0.4.0 — `.context/sessions/` eliminated; journal entries use `ctx recall export` naming)

**Context**: Multiple sessions per day would overwrite each other.
Also, multiple compacts in the same minute could collide.

**Decision**: Use `YYYY-MM-DD-HHMMSS-<topic>.md` format for session files.
Two file types:
- **Manual session files**: `HHMMSS-<topic>.md` - updated throughout session
- **Auto-snapshots**: `HHMMSS-<event>.jsonl` - immutable once created

**Rationale**:
- Human-readable (unlike unix timestamps)
- Naturally sorts chronologically
- Seconds precision prevents collision even with rapid compacts
- Clear distinction between manual notes and raw snapshots

**Consequences**:
- Slightly longer filenames
- Must ensure consistent format in all session-saving code
- Manual files keep getting updated; snapshots are write-once

---




## [2026-01-20-120000] Two-Tier Context Persistence Model

**Status**: Superseded (v0.4.0 — two tiers remain but `.context/sessions/` eliminated; full-dump tier is now `~/.claude/projects/` JSONL + `.context/journal/` enriched markdown)

**Context**: Need to persist context across sessions. Token budgets limit
what can be loaded. But nothing should be truly lost.

**Decision**: Implement two tiers of persistence:

| Tier          | Purpose                 | Location                 | Token Cost             |
|---------------|-------------------------|--------------------------|------------------------|
| **Curated**   | Quick context reload    | `.context/*.md`          | Low (budgeted)         |
| **Full dump** | Safety net, archaeology | `.context/sessions/*.md` | Zero (not auto-loaded) |

**Rationale**:
- Curated context is token-efficient for daily use
- Full dumps ensure nothing is ever truly lost
- Users can dive into sessions/ when they need deep context
- Separation prevents context bloat

**Consequences**:
- Need both manual and automatic ways to populate both tiers
- Session files grow over time (may need archival strategy)
- `ctx agent` only loads curated tier by default

---



## Group: Journal site rendering



## [2026-02-20-121937] Defencify journal site: pre/code replaces fenced code blocks

**Status**: Accepted

**Context**: Fenced code blocks for Tool Output caused cascading bugs: nesting conflicts, fence depth calculation, stray fences in User turns swallowing subsequent content. Multiple fixes (isAlreadyFenced, fenceForContent, fence tracking) added complexity without fully solving the problem.

**Decision**: Defencify journal site: pre/code replaces fenced code blocks

**Rationale**: HTML-escaped content in <pre><code> eliminates ALL inner content conflicts in one shot — no fence depth, no nesting, no tracking. Requires use_pygments=false in zensical config to prevent pymdownx.highlight from hijacking the blocks. Trade-off: content is plain preformatted text (no rich markdown rendering inside tool output or user messages), which is acceptable.

**Consequences**: Removed isAlreadyFenced, fenceForContent, codeFence. Line-by-line transforms track <pre> blocks instead of fences. Generated zensical.toml now includes full markdown_extensions config. Collapsible tool output (<details>) remains incompatible — needs CSS/JS approach.

---




## [2026-02-20-062043] Code-level normalize replaces AI source normalization

**Status**: Accepted

**Context**: Attempted AI normalization of 290 journal files (~1M lines). 10 agents found zero source edits needed — files were already well-formatted. Only finding: 12 files with broken fence language hints.

**Decision**: Code-level normalize replaces AI source normalization

**Rationale**: normalizeContent pipeline handles all rendering concerns at build time (fence stripping, tool output wrapping, heading demotion, list spacing). Source-level AI normalization adds no value beyond what the code pipeline provides.

**Consequences**: Mark all entries as normalized+fences_verified programmatically. Reserve AI normalization for specific files with known issues. Enrichment proceeds without blocking on source normalization.

---




## [2026-02-20-044438] Title length limit of 75 characters for journal entries

**Status**: Accepted

**Context**: H1 headings and link text wrapping to second line at around 80 chars caused rendering issues — the second line did not render as part of the heading.

**Decision**: Title length limit of 75 characters for journal entries

**Rationale**: 75 chars keeps headings on a single line below the typical wrap width. Truncation happens on word boundary (last space before limit). Applied in three places: cleanTitle (recall/slug.go) for export frontmatter, normalizeContent (journal/normalize.go) for H1 headings in site copy, and parseJournalEntry (journal/parse.go) inherits from frontmatter title. RecallMaxTitleLen = 75 in config/limit.go is the single source of truth.

**Consequences**: Titles may lose trailing words. FirstUserMsg from parser (100 chars) gets further truncated. The three-character ellipsis suffix from parser is stripped before truncation.

---




## [2026-02-20-044430] CSS overflow instead of details collapsibility for journal tool outputs

**Status**: Accepted

**Context**: Long tool outputs (hundreds of lines) need visual containment. Previously used details/summary for collapsibility, but details is a Type 6 HTML block incompatible with fenced code blocks.

**Decision**: CSS overflow instead of details collapsibility for journal tool outputs

**Rationale**: CSS max-height with overflow-y: auto on pre elements provides scroll-based containment without any HTML block interaction. Works with any content inside pre, including fenced code blocks rendered by the markdown parser. No class names needed — targets .md-typeset pre globally. extra_css in zensical.toml must appear under project section (after nav, before project.theme).

**Consequences**: Journal site ships docs/stylesheets/extra.css (generated by run.go). Tool outputs scroll at 30em instead of click-to-expand. The stylesheet is written fresh on every ctx journal site run.

---




## [2026-02-20-044421] Use fenced code blocks for tool output wrapping in journal site

**Status**: Accepted

**Context**: Tool output sections in exported session transcripts contain arbitrary content including markdown syntax (headings, thematic breaks, lists) and HTML fragments that break rendering when interpreted by Python-Markdown.

**Decision**: Use fenced code blocks for tool output wrapping in journal site

**Rationale**: Tried three approaches: 1) details/pre wrappers — Type 6 HTML block ends at blank lines in content. 2) pre/code wrappers — Python-Markdown does not implement CommonMark Type 1 blocks, so pre also ends at blank lines. 3) Fenced code blocks — correctly survive blank lines and prevent all markdown/HTML interpretation. Fences are safe because stripFences runs first and removes all fence lines from content before wrapToolOutputs adds new fence wrappers.

**Consequences**: Tool outputs render as monospace code blocks. No HTML escaping needed (content emitted verbatim). Lost details/summary collapsibility — replaced by CSS max-height + overflow-y: auto via extra.css. The normalize.go code is dramatically simpler (no escapeToolLine, no blank line neutralization, no HTML entity juggling).

---



## Group: Plugin and skill distribution



## [2026-02-16-164550] Permission docs match DefaultClaudePermissions exactly

**Status**: Accepted

**Context**: claude-code-permissions.md showed a curated subset of 13 skill permissions while DefaultClaudePermissions in config/file.go had 26 entries

**Decision**: Permission docs match DefaultClaudePermissions exactly

**Rationale**: A curated subset causes confusion when ctx init seeds more permissions than the docs recommend — users wonder where the extra entries came from and whether they are safe

**Consequences**: The recommended defaults section now lists all 26 skill entries from DefaultClaudePermissions; future skill additions must update both config/file.go and the recipe

---




## [2026-02-16-164512] No symlinks for cross-directory skill sharing

**Status**: Accepted

**Context**: Considered symlinking .claude/skills/ctx-* to internal/assets/claude/skills/ to avoid duplication

**Decision**: No symlinks for cross-directory skill sharing

**Rationale**: Git on Windows checks out symlinks as plain text files containing the target path unless Developer Mode is enabled and core.symlinks=true. Most contributors won't have this configured, breaking the build/embed on forked repos

**Consequences**: We keep two physical directories with distinct purposes instead of linking them. Contributors install the plugin from their local clone directory to get the ctx-* skills.

---




## [2026-02-16-164509] Single source of truth for distributed skills

**Status**: Accepted

**Context**: Discovered ctx-* skills were duplicated between .claude/skills/ and internal/assets/claude/skills/, causing double entries in Claude Code's skill list

**Decision**: Single source of truth for distributed skills

**Rationale**: internal/assets/claude/skills/ is what gets embedded in the binary and served by the marketplace plugin — it is the distribution path. .claude/skills/ now holds only the 12 dev-only skills (release, qa, backup, etc.) that are never distributed

**Consequences**: Skill edits for user-facing ctx-* skills happen in internal/assets/claude/skills/. Dev-only skills live in .claude/skills/. No sync script needed.

---




## [2026-02-16-100447] ctx v0.6.0: Plugin conversion — shell hooks to Go subcommands

**Status**: Accepted

**Context**: ctx v0.4.0 deployed 6 shell scripts to .claude/hooks/ via ctx init, requiring jq and coupling ctx to per-project scaffolding

**Decision**: ctx v0.6.0: Plugin conversion — shell hooks to Go subcommands

**Rationale**: Go subcommands (ctx system *) eliminate jq dependency, ship as compiled binary, enable testing with go test, and leverage Claude Code's plugin system for distribution

**Consequences**: ctx init no longer creates .claude/hooks/ or .claude/skills/. Users install the ctx plugin separately. Existing projects need to remove old .sh hooks and install the plugin. Version jumps from 0.4.0 to 0.6.0 to signal the magnitude.

---



## Group: Recall system design



## [2026-02-20-224112] Spec supersedes old Phase 2 export-preservation tasks

**Status**: Accepted

**Context**: The recall-export-safety spec covers a broader scope (locks, --keep-frontmatter, --dry-run, ergonomics) than the original 4 tasks (T2.1.1-T2.1.4) which only addressed --force behavior

**Decision**: Spec supersedes old Phase 2 export-preservation tasks

**Rationale**: Replaced narrow bug-fix tasks with 7 spec-aligned tasks (T2.1-T2.7) that cover the full design

**Consequences**: Old specs/export-update-mode.md is superseded by specs/recall-export-safety.md; TASKS.md Phase 2 section rewritten with 7 tasks mapping to spec phases

---




## [2026-02-20-142444] No --update flag needed for export — default is the update mode

**Status**: Accepted

**Context**: T2.1 requested ctx recall export --update to preserve enrichments during re-export

**Decision**: No --update flag needed for export — default is the update mode

**Rationale**: The default behavior already preserves enriched YAML frontmatter. Adding a redundant flag increases API surface without value. The real fix is making --force work correctly.

**Consequences**: T2.1 broken into 4 subtasks focused on the --force bug, state cleanup, tests, and docs. No new CLI flag added.

---




## [2026-01-28-045840] ctx recall is Claude-first

**Status**: Accepted

**Context**: Building recall feature to parse AI session history. JSONL formats
differ across tools (Claude Code, Aider, Cursor). Need to decide scope and
compatibility strategy.

**Decision**: ctx recall is Claude-first

**Rationale**: Claude Code is primary target audience. Most users auto-upgrade,
so supporting only recent versions avoids maintenance burden. Other tools can
add parsers but are secondary - not worth same polish.

**Consequences**: 1) Parser updates follow Claude Code releases,
no legacy schema support. 2) Aider/Cursor parsers are community-contributed,
best-effort. 3) Features can assume Claude Code conventions
(slugs, session IDs, tool result format).

---




## [2026-01-28-040251] Use tool-agnostic Session type with tool-specific parsers for recall system

**Status**: Accepted

**Context**: JSONL session formats are not standardized across AI coding
assistants. Claude Code, Cursor, Aider each have different formats or may not
export sessions at all. Need to support multiple tools eventually.

**Decision**: Use tool-agnostic Session type with tool-specific parsers for
recall system

**Rationale**: Separating the output type (Session) from the parsing logic
allows adding new tool support without changing downstream code. Starting with
Claude Code only, but the interface abstraction makes it easy to add
AiderParser, CursorParser, etc. later.

**Consequences**: 1. Session struct is tool-agnostic (common fields
only). 2. SessionParser interface defines ParseFile, ParseLine,
CanParse. 3. ClaudeCodeParser is first implementation. 4. Parser
registry/factory can auto-detect format from file content.

---

