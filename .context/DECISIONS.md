# Decisions

<!-- INDEX:START -->
| Date | Decision |
|------|--------|
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
| 2026-01-28 | ctx recall is Claude-first |
| 2026-01-28 | Tasks must include explicit deliverables, not just implementation steps |
| 2026-01-28 | Use tool-agnostic Session type with tool-specific parsers for recall system |
| 2026-01-27 | Use reverse-chronological order (newest first) for DECISIONS.md and LEARNINGS.md |
| 2026-01-25 | Removed AGENTS.md from project root |
| 2026-01-25 | Keep CONSTITUTION Minimal |
| 2026-01-25 | Centralize Constants with Semantic Prefixes |
| 2026-01-21 | Separate Orchestrator Directive from Agent Tasks |
| 2026-01-21 | Hooks Use ctx from PATH, Not Hardcoded Paths |
| 2026-01-20 | Use SessionEnd Hook for Auto-Save |
| 2026-01-20 | Handle CLAUDE.md Creation/Merge in ctx init |
| 2026-01-20 | Auto-Save Before Compact |
| 2026-01-20 | Session Filename Format: YYYY-MM-DD-HHMMSS-topic.md |
| 2026-01-20 | Two-Tier Context Persistence Model |
| 2026-01-20 | Always Generate Claude Hooks in Init (No Flag Needed) |
| 2026-01-20 | Generic Core with Optional Claude Code Enhancements |
<!-- INDEX:END -->

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

**Context**: Task proposed either /ctx-borrow skill or ctx borrow CLI command for merging deltas between two directories

**Decision**: Borrow-from-the-future implemented as skill, not CLI command

**Rationale**: The workflow requires interactive judgment: conflict resolution, selective file application, strategy selection between 3 tiers. An agent adapts to edge cases; CLI flags cannot.

**Consequences**: No ctx borrow subcommand. Users invoke /ctx-borrow in their AI tool. Non-AI users would need to manually run git diff/patch commands.

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
