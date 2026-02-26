# Archived Decisions (consolidated 2026-02-26)

Originals replaced by consolidated entries in DECISIONS.md.

## Group: Blog and content publishing architecture

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

## [2026-02-06-181708] Drop ctx-journal-summarize skill (duplicates ctx-blog)

**Status**: Accepted

**Context**: ctx-journal-summarize and ctx-blog both read journal entries over a time range and produce narrative summaries. The only difference was audience framing: internal summary vs public blog post.

**Decision**: Drop ctx-journal-summarize skill (duplicates ctx-blog)

**Rationale**: The blog skill can serve both use cases with a prompt tweak. One fewer skill to maintain, less surface area for drift.

**Consequences**: Removed skill dir, template, and references from integrations.md and two blog posts. Timeline narrative deferred item in TASKS.md marked as dropped. Users who want internal summaries use /ctx-blog instead.

---

## Group: Hook and notification design

## [2026-02-24-204550] Tone down proactive content suggestion claims rather than add more hooks

**Status**: Accepted

**Context**: publishing.md claims agents proactively suggest blog posts and journal rebuilds at natural moments, but no hook or playbook mechanism exists to trigger this.

**Decision**: Tone down proactive content suggestion claims rather than add more hooks

**Rationale**: Already have 9 UserPromptSubmit hooks — adding another nudge risks fatigue. The claim is aspirational, not functional. Conversational prompting (ask your agent) already works.

**Consequences**: Update docs to describe the conversational approach rather than claiming automatic behavior. Avoids over-promising. If demand emerges later, a hook can be added then.

---

## [2026-02-22-194444] Hook commands use structured JSON output instead of plain text

**Status**: Accepted

**Context**: qa-reminder and post-commit hooks were being ignored despite firing correctly

**Decision**: Hook commands use structured JSON output instead of plain text

**Rationale**: JSON with hookSpecificOutput.additionalContext is parsed as a directive by Claude Code, while plain text is treated as ambient context the agent can ignore

**Consequences**: Added HookResponse/HookSpecificOutput types and printHookContext() helper to internal/cli/system/input.go; converted qa_reminder.go and post_commit.go; future hooks should use printHookContext() for non-blocking directives

---

## [2026-02-12-005516] Drop prompt-coach hook

**Status**: Accepted

**Context**: Prompt-coach has been running since installation with zero useful tips fired. All counters across all state files are 0. The delivery mechanism is broken (stdout goes to AI not user, stderr is swallowed). Even if fixed with systemMessage, the coaching patterns are too narrow for experienced users and the prompting guide already covers best practices.

**Decision**: Drop prompt-coach hook

**Rationale**: Three layers of not-working: (1) patterns too narrow to match real prompts, (2) output channel invisible to user, (3) L-3 PID bug creates orphan temp files. Removing it eliminates the largest source of temp file accumulation, simplifies the hook stack, and removes dead code.

**Consequences**: One fewer hook in UserPromptSubmit (faster prompt submission). Eliminates prompt-coach temp file accumulation entirely — reduces cleanup burden. Need to remove: template script, config constant, script loader, hookScripts entry, settings.local.json reference, and active hook file.

---

## [2026-02-22-221724] De-emphasize /ctx-journal-normalize from default journal pipeline

**Status**: Accepted

**Context**: The journal pipeline previously prescribed normalize -> enrich as the default workflow. With improved programmatic normalization during export and simplified markdown generation (no code fences), the AI-based normalize skill is rarely needed.

**Decision**: De-emphasize /ctx-journal-normalize from default journal pipeline

**Rationale**: The normalize skill is expensive (reads entire journal files through LLM), nondeterministic on large inputs, and blows up subagent context windows on non-ctx projects with millions of lines of session JSON. Programmatic normalization handles most cases.

**Consequences**: Normalize removed from relay nudge, make journal, skill prerequisites, and all docs. Skill remains available for targeted per-file use via /ctx-journal-normalize when rendering issues occur.

---

## Group: ctx init and CLAUDE.md handling

## [2026-01-20-180000] Handle CLAUDE.md Creation/Merge in ctx init

**Status**: Accepted (to be implemented)

**Context**: Both `claude init` and `ctx init` want to create/modify CLAUDE.md.
Users of ctx will likely want ctx's context-aware version,
but may already have a CLAUDE.md from `claude init`.

**Decision**: `ctx init` handles CLAUDE.md intelligently:
- **No CLAUDE.md exists** -> Create it with ctx's context-loading template
- **CLAUDE.md exists** -> Don't overwrite. Instead:
  1. **Backup first** -> Copy to `CLAUDE.md.<unix_timestamp>.bak`
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

---

## Group: Documentation and navigation structure

## [2026-02-21-200038] Restructure docs nav sections with dedicated index pages

**Status**: Accepted

**Context**: Reference, Operations, and Security nav sections lacked icons in the mobile menu because they had no section index pages

**Decision**: Restructure docs nav sections with dedicated index pages

**Rationale**: Created reference/index.md, operations/index.md, security/index.md with linked summaries of sub-pages. Moved security.md to security/reporting.md to avoid file/directory name conflict. Renamed page titles to remove redundant ctx prefix (CLI, Skills, Tool Ecosystem).

**Consequences**: All nav sections now have icons on mobile. security.md URL changes to security/reporting/. Three internal links updated. Index pages serve as lightweight landing pages for each section.

---

## [2026-02-15-194828] Add TL;DR admonitions to recipes longer than ~200 lines

**Status**: Accepted

**Context**: Recipes bury the actionable pipeline at the bottom in Putting It All Together sections. Users must scroll past 300+ lines of explanation.

**Decision**: Add TL;DR admonitions to recipes longer than ~200 lines

**Rationale**: A tip admonition after the intro surfaces the quick-start commands immediately. Users who want depth still read the full page.

**Consequences**: 10 recipes now have TL;DRs. New recipes over ~200 lines should follow the pattern. Short recipes (permission-snapshots, scratchpad-with-claude) skip it.

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

## Group: Task and knowledge management

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

## [2026-02-18-071514] Knowledge scaling: archive path for decisions and learnings

**Status**: Accepted

**Context**: DECISIONS.md and LEARNINGS.md grow monotonically with no archival path. Tasks have ctx tasks archive but knowledge files accumulate forever. Long-lived projects will hit token budget pressure and signal-to-noise decay.

**Decision**: Knowledge scaling: archive path for decisions and learnings

**Rationale**: Follow the existing task archive pattern. Move old entries to .context/archive/ files. Extend ctx compact --archive to cover all three file types. Add superseded-entry convention for decisions.

**Consequences**: New spec at specs/knowledge-scaling.md. Phase 5 tasks (P5.1-P5.7) added to TASKS.md. New CLI commands: ctx decisions archive, ctx learnings archive. New .contextrc keys: archive_knowledge_after_days, archive_keep_recent.

---

## Group: Agent autonomy and separation of concerns

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

## Group: Security and permissions

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

## [2026-02-24-015709] Drop absolute-path-to-ctx regex from block-dangerous-commands

**Status**: Accepted

**Context**: Shell script had a regex blocking absolute paths to ctx binary. The block-non-path-ctx Go subcommand already covers this with better patterns (start-of-command, after-separator positions, test exception).

**Decision**: Drop absolute-path-to-ctx regex from block-dangerous-commands

**Rationale**: Duplicating it would create two sources of truth for the same rule.

**Consequences**: Only block-non-path-ctx owns absolute-path blocking; block-dangerous-commands focuses on sudo, git push, and bin-directory installs.
