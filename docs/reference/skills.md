---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Skills
icon: lucide/sparkles
---

![ctx](../images/ctx-banner.png)

## Skills

Skills are slash commands that run **inside your AI assistant** (*e.g.,
`/ctx-next`*), as opposed to CLI commands that run in your terminal
(*e.g., `ctx status`*). 

Skills give your agent structured workflows: It knows what to read, what to 
run, and when to ask. Most wrap one or more `ctx` CLI commands with 
opinionated behavior on top. 

!!! tip "Skills Are Best Used Conversationally"
    The beauty of `ctx` is that it's designed to be intuitive and 
    conversational, allowing you to interact with your AI assistant 
    naturally. That's why you don't have to memorize many of
    these skills.

    See the [**Prompting Guide**](../home/prompting-guide.md) for natural-language 
    triggers that invoke these skills conversationally.

    However, when you need a more precise control, you have the option
    to invoke the relevant skills directly.

<!-- drift-check: ls internal/assets/claude/skills/ | wc -l -->
<!-- drift-check: diff <(ls internal/assets/claude/skills/ | sort) <(sed -n '/^## All Skills/,/^---$/p' docs/reference/skills.md | grep -oP '\| \[`/\K[a-z-]+(?=`\])' | grep -v check-links | sort -u) -->
## All Skills

| Skill                                                | Description                                                     | Type             |
|------------------------------------------------------|-----------------------------------------------------------------|------------------|
| [`/ctx-remember`](#ctx-remember)                     | Recall project context and present structured readback          | user-invocable   |
| [`/ctx-wrap-up`](#ctx-wrap-up)                       | End-of-session context persistence ceremony                     | user-invocable   |
| [`/ctx-status`](#ctx-status)                         | Show context summary with interpretation                        | user-invocable   |
| [`/ctx-agent`](#ctx-agent)                           | Load full context packet for AI consumption                     | user-invocable   |
| [`/ctx-next`](#ctx-next)                             | Suggest 1-3 concrete next actions with rationale                | user-invocable   |
| [`/ctx-commit`](#ctx-commit)                         | Commit with integrated context persistence                      | user-invocable   |
| [`/ctx-reflect`](#ctx-reflect)                       | Pause and reflect on session progress                           | user-invocable   |
| [`/ctx-add-task`](#ctx-add-task)                     | Add actionable task to TASKS.md                                 | user-invocable   |
| [`/ctx-add-decision`](#ctx-add-decision)             | Record architectural decision with rationale                    | user-invocable   |
| [`/ctx-add-learning`](#ctx-add-learning)             | Record gotchas and lessons learned                              | user-invocable   |
| [`/ctx-add-convention`](#ctx-add-convention)         | Record coding convention for consistency                        | user-invocable   |
| [`/ctx-archive`](#ctx-archive)                       | Archive completed tasks from TASKS.md                           | user-invocable   |
| [`/ctx-pad`](#ctx-pad)                               | Manage encrypted scratchpad entries                             | user-invocable   |
| [`/ctx-recall`](#ctx-recall)                         | Browse and export AI session history                            | user-invocable   |
| [`/ctx-journal-enrich`](#ctx-journal-enrich)         | Enrich single journal entry with metadata                       | user-invocable   |
| [`/ctx-journal-enrich-all`](#ctx-journal-enrich-all) | Batch-enrich all unenriched journal entries                     | user-invocable   |
| [`/ctx-journal-normalize`](#ctx-journal-normalize)   | Normalize journal markdown for clean rendering                  | user-invocable   |
| [`/ctx-blog`](#ctx-blog)                             | Generate blog post draft from project activity                  | user-invocable   |
| [`/ctx-blog-changelog`](#ctx-blog-changelog)         | Generate themed blog post from a commit range                   | user-invocable   |
| [`/ctx-consolidate`](#ctx-consolidate)               | Consolidate redundant learnings or decisions                    | user-invocable   |
| [`/ctx-drift`](#ctx-drift)                           | Detect and fix context drift                                    | user-invocable   |
| [`/ctx-alignment-audit`](#ctx-alignment-audit)       | Audit docs claims against agent instructions                    | user-invocable   |
| [`/ctx-prompt-audit`](#ctx-prompt-audit)             | Analyze prompting patterns for improvement                      | user-invocable   |
| [`/check-links`](#check-links)                       | Audit docs for dead internal and external links                 | project-specific |
| [`/ctx-context-monitor`](#ctx-context-monitor)       | Respond to context checkpoint signals                           | automatic        |
| [`/ctx-import-plans`](#ctx-import-plans)             | Import Claude Code plan files into project specs                | user-invocable   |
| [`/ctx-implement`](#ctx-implement)                   | Execute a plan step-by-step with verification                   | user-invocable   |
| [`/ctx-loop`](#ctx-loop)                             | Generate autonomous loop script                                 | user-invocable   |
| [`/ctx-worktree`](#ctx-worktree)                     | Manage git worktrees for parallel agents                        | user-invocable   |
| [`/ctx-map`](#ctx-map)                               | Build and maintain architecture maps                            | user-invocable   |
| [`/ctx-remind`](#ctx-remind)                         | Manage session-scoped reminders                                 | user-invocable   |
| [`/ctx-doctor`](#ctx-doctor)                         | Troubleshoot ctx behavior with health checks and event analysis | user-invocable   |
| [`/ctx-pause`](#ctx-pause)                           | Pause context hooks for this session                            | user-invocable   |
| [`/ctx-resume`](#ctx-resume)                         | Resume context hooks after a pause                              | user-invocable   |

---

## Session Lifecycle

Skills for starting, running, and ending a productive session.

!!! note "Session Ceremonies"
    Two skills in this group are **ceremony skills**: `/ctx-remember` (session
    start) and `/ctx-wrap-up` (session end). Unlike other skills that work
    conversationally, these should be invoked as **explicit slash commands**
    for completeness. See [Session Ceremonies](../recipes/session-ceremonies.md).

### `/ctx-remember`

Recall project context and present a structured readback.
**Ceremony skill**: invoke explicitly at session start.

**Wraps**: `ctx agent --budget 4000`, `ctx recall list --limit 3`,
reads TASKS.md, DECISIONS.md, LEARNINGS.md

**See also**: [Session Ceremonies](../recipes/session-ceremonies.md),
[The Complete Session](../recipes/session-lifecycle.md)

---

### `/ctx-status`

Show context summary (*files, token budget, tasks, recent activity*)
with interpreted suggestions.

**Wraps**: `ctx status [--verbose] [--json]`

**See also**: [The Complete Session](../recipes/session-lifecycle.md),
[`ctx status` CLI](../cli/init-status.md#ctx-status)

---

### `/ctx-agent`

Load the full context packet optimized for AI consumption.
Also runs automatically via the PreToolUse hook with cooldown.

**Wraps**: `ctx agent [--budget] [--format] [--cooldown] [--session]`

**See also**: [The Complete Session](../recipes/session-lifecycle.md),
[`ctx agent` CLI](../cli/init-status.md#ctx-agent)

---

### `/ctx-next`

Suggest 1-3 concrete next actions ranked by priority, momentum,
and unblocked status.

**Wraps**: reads TASKS.md, `ctx recall list --limit 3`

**See also**: [The Complete Session](../recipes/session-lifecycle.md),
[Tracking Work Across Sessions](../recipes/task-management.md)

---

### `/ctx-commit`

Commit code with integrated context persistence: pre-commit checks,
staged files, Co-Authored-By trailer, and a post-commit prompt to
capture decisions and learnings.

**Wraps**: `git add`, `git commit`, optionally chains to
`/ctx-add-decision` and `/ctx-add-learning`

**See also**: [The Complete Session](../recipes/session-lifecycle.md)

---

### `/ctx-reflect`

Pause and reflect on session progress. Walks through a checklist of
learnings, decisions, task completions, and session notes to persist.

**Wraps**: chains to `ctx add learning`, `ctx add decision`,
manual TASKS.md updates

**See also**: [The Complete Session](../recipes/session-lifecycle.md),
[Persisting Decisions, Learnings, and Conventions](../recipes/knowledge-capture.md)

---

### `/ctx-wrap-up`

End-of-session context persistence ceremony. Gathers signal from
git diff, recent commits, and conversation themes. Proposes
candidates (learnings, decisions, conventions, tasks) with complete
structured fields for user approval, then persists via `ctx add`.
Offers `/ctx-commit` if uncommitted changes remain.
**Ceremony skill**: invoke explicitly at session end.

**Wraps**: `git diff --stat`, `git log`, `ctx add learning`,
`ctx add decision`, `ctx add convention`, `ctx add task`,
chains to `/ctx-commit`

**See also**: [Session Ceremonies](../recipes/session-ceremonies.md),
[The Complete Session](../recipes/session-lifecycle.md)

---

## Context Persistence

Skills for recording work artifacts: tasks, decisions, learnings,
conventions: into `.context/` files.

### `/ctx-add-task`

Add an actionable task with optional priority and phase section.

**Wraps**: `ctx add task "description" [--priority high|medium|low]`

**See also**: [Tracking Work Across Sessions](../recipes/task-management.md)

---

### `/ctx-add-decision`

Record an architectural decision with context, rationale, and
consequences. Supports Y-statement (lightweight) and full ADR formats.

**Wraps**: `ctx add decision "title" --context "..." --rationale "..."
--consequences "..."`

**See also**:
[Persisting Decisions, Learnings, and Conventions](../recipes/knowledge-capture.md)

---

### `/ctx-add-learning`

Record a project-specific gotcha, bug, or unexpected behavior.
Filters for insights that are searchable, project-specific, and
required real effort to discover.

**Wraps**: `ctx add learning "title" --context "..." --lesson "..."
--application "..."`

**See also**:
[Persisting Decisions, Learnings, and Conventions](../recipes/knowledge-capture.md)

---

### `/ctx-add-convention`

Record a coding convention that should be standardized across sessions.
Targets patterns seen 2-3+ times.

**Wraps**: `ctx add convention "rule" --section "Name"`

**See also**:
[Persisting Decisions, Learnings, and Conventions](../recipes/knowledge-capture.md)

---

### `/ctx-archive`

Archive completed tasks from TASKS.md to a timestamped file in
`.context/archive/`. Preserves phase headers for traceability.

**Wraps**: `ctx tasks archive [--dry-run]`

**See also**: [Tracking Work Across Sessions](../recipes/task-management.md)

---

## Scratchpad

### `/ctx-pad`

Manage the encrypted scratchpad: add, remove, edit, and reorder
one-liner notes. Encrypted at rest with AES-256-GCM.

**Wraps**: `ctx pad`, `ctx pad add`, `ctx pad rm`, `ctx pad edit`,
`ctx pad mv`, `ctx pad import`, `ctx pad export`, `ctx pad merge`

**See also**: [Scratchpad](scratchpad.md),
[Using the Scratchpad](../recipes/scratchpad-with-claude.md)

---

## Journal & History

Skills for browsing, exporting, and enriching your AI session history
into a structured journal.

### `/ctx-recall`

Browse, inspect, and export AI session history. List recent sessions,
show details by slug or ID, and export to `.context/journal/`.

**Wraps**: `ctx recall list`, `ctx recall show`, `ctx recall export`

**See also**:
[Browsing and Enriching Past Sessions](../recipes/session-archaeology.md)

---

### `/ctx-journal-enrich`

Enrich a single journal entry with YAML frontmatter: title, type,
outcome, topics, technologies, and summary. Shows diff before writing.

**Wraps**: reads and edits `.context/journal/*.md` files

**See also**:
[Browsing and Enriching Past Sessions](../recipes/session-archaeology.md),
[Turning Activity into Content](../recipes/publishing.md)

---

### `/ctx-journal-enrich-all`

Batch-enrich all unenriched journal entries. Filters out short sessions
and continuations. Can spawn subagents for large backlogs.

**Wraps**: iterates `/ctx-journal-enrich` across all entries

**See also**:
[Browsing and Enriching Past Sessions](../recipes/session-archaeology.md)

---

### `/ctx-journal-normalize`

Normalize journal markdown for clean rendering: fix fence nesting,
metadata formatting, list indentation, and collapse large tool outputs.

**Wraps**: reads and edits `.context/journal/*.md` files

**See also**:
[Browsing and Enriching Past Sessions](../recipes/session-archaeology.md),
[Turning Activity into Content](../recipes/publishing.md)

---

## Content Creation

Skills for turning project activity into publishable content.

### `/ctx-blog`

Generate a blog post draft from recent project activity: git history,
decisions, learnings, tasks, and journal entries. Requires a narrative
arc (problem, approach, outcome).

**Wraps**: reads `git log`, DECISIONS.md, LEARNINGS.md, TASKS.md,
journal entries; writes to `docs/blog/`

**See also**: [Turning Activity into Content](../recipes/publishing.md)

---

### `/ctx-blog-changelog`

Generate a themed blog post from a commit range. Takes a starting
commit and unifying theme, analyzes diffs and journal entries from
that period.

**Wraps**: `git log`, `git diff --stat`; writes to `docs/blog/`

**See also**: [Turning Activity into Content](../recipes/publishing.md)

---

## Auditing & Health

Skills for detecting drift, auditing alignment, and improving
prompt quality.

### `/ctx-consolidate`

Consolidate redundant entries in LEARNINGS.md or DECISIONS.md. Groups
overlapping entries by keyword similarity, presents candidates, and
(*with user approval*) merges groups into denser combined entries.
Originals are archived, not deleted.

**Wraps**: reads LEARNINGS.md and DECISIONS.md, writes consolidated
entries, archives originals, runs `ctx reindex`

**See also**:
[Detecting and Fixing Drift](../recipes/context-health.md)

---

### `/ctx-drift`

Detect and fix context drift: stale paths, missing files, file age
staleness, task accumulation, entry count warnings, and constitution
violations via `ctx drift`. Also detects skill drift against canonical
templates.

**Wraps**: `ctx drift [--fix]`

**See also**:
[Detecting and Fixing Drift](../recipes/context-health.md)

---

### `/ctx-alignment-audit`

Audit behavioral claims in docs and recipes against actual agent
instructions. Traces each claim to its backing instruction and reports
coverage as Covered, Partial, or Gap.

**Wraps**: reads AGENT_PLAYBOOK.md, plugin skill definitions, CLAUDE.md,
and docs/recipes

**See also**:
[Detecting and Fixing Drift](../recipes/context-health.md)

---

### `/ctx-prompt-audit`

Analyze recent prompting patterns to identify vague or ineffective
prompts. Reviews 3-5 journal entries and suggests rewrites with
positive observations.

**Wraps**: reads `.context/journal/` entries

**See also**:
[Detecting and Fixing Drift](../recipes/context-health.md)

---

### `/ctx-doctor`

Troubleshoot ctx behavior. Runs structural health checks via `ctx doctor`,
analyzes event log patterns via `ctx system events`, and presents findings
with suggested actions. The CLI provides the structural baseline; the agent
adds semantic analysis of event patterns and correlations.

**Wraps**: `ctx doctor --json`, `ctx system events --json --last 100`,
`ctx remind list`, `ctx system message list`, reads `.ctxrc`

**Trigger phrases**: "diagnose", "troubleshoot", "doctor", "health check",
"why didn't my hook fire?", "hooks seem broken", "something seems off"

**Graceful degradation**: If `event_log` is not enabled, the skill still
works but with reduced capability. It runs structural checks and notes:
"Enable `event_log: true` in `.ctxrc` for hook-level diagnostics."

**See also**: [Troubleshooting](../recipes/troubleshooting.md),
[`ctx doctor` CLI](../cli/doctor.md#ctx-doctor),
[`ctx system events` CLI](../cli/system.md#ctx-system-events)

---

### `/check-links`

Scan all markdown files under `docs/` for broken links. Two passes:
internal links (verify file targets exist on disk) and external links
(HTTP HEAD with timeout, report failures as warnings). Also checks
image references.

Invoked automatically as check #12 during `/audit`.

**Wraps**: Glob + Grep to scan, `curl` for external checks

**See also**:
[`/audit`](#audit-related-skills)

---

### `/ctx-context-monitor`

Respond to context checkpoint signals when usage hits high thresholds.
Fires at adaptive intervals and offers context persistence before
the session ends.

**Type**: Automatic: Triggered by the `check-context-size` hook,
not user-invocable

**Wraps**: hook-driven; suggests `/ctx-reflect`

**See also**:
[Running an Unattended AI Agent](../recipes/autonomous-loops.md)

---

## Planning & Execution

Skills for structured implementation and parallel agent workflows.

### `/ctx-import-plans`

Import Claude Code plan files (`~/.claude/plans/*.md`) into the project's
`specs/` directory. Lists plans with dates and H1 titles, supports
filtering (`--today`, `--since`, `--all`), slugifies headings for
filenames, and optionally creates tasks referencing each imported spec.

**Wraps**: reads `~/.claude/plans/*.md`, writes to `specs/`,
optionally chains to `/ctx-add-task`

**See also**:
[Importing Claude Code Plans](../recipes/import-plans.md),
[Tracking Work Across Sessions](../recipes/task-management.md)

---

### `/ctx-implement`

Execute a multi-step plan with build and test verification at each
step. Loads a plan from a file or conversation context, breaks it
into atomic steps, and checkpoints after every 3-5 steps.

**Wraps**: reads plan file, runs verification commands
(`go build`, `go test`, etc.)

**See also**:
[Running an Unattended AI Agent](../recipes/autonomous-loops.md)

---

### `/ctx-loop`

Generate a ready-to-run shell script for autonomous AI iteration.
Supports Claude Code, Aider, and generic tool templates with
configurable completion signals.

**Wraps**: `ctx loop [--tool] [--prompt] [--max-iterations]
[--completion] [--output]`

**See also**: [Autonomous Loops](../operations/autonomous-loop.md),
[Running an Unattended AI Agent](../recipes/autonomous-loops.md)

---

### `/ctx-worktree`

Manage git worktrees for parallel agent development. Create sibling
worktrees on dedicated branches, analyze task blast radius for
grouping, and tear down with merge.

**Wraps**: `git worktree add`, `git worktree list`,
`git worktree remove`, `git merge`

**See also**:
[Parallel Agent Development with Git Worktrees](../recipes/parallel-worktrees.md)

---

### `/ctx-map`

Build and maintain architecture maps incrementally. Creates or refreshes
`ARCHITECTURE.md` (*succinct project map, loaded at session start*) and
`DETAILED_DESIGN.md` (*deep per-module reference, consulted on-demand*).
Coverage is tracked in `map-tracking.json` so each run extends the map
rather than re-analyzing everything.

**Wraps**: `ctx status`, `git log`, reads source files; writes
`ARCHITECTURE.md`, `DETAILED_DESIGN.md`, `map-tracking.json`

**See also**:
[Detecting and Fixing Drift](../recipes/context-health.md)

---

### `/ctx-remind`

Manage session-scoped reminders via natural language. Translates user
intent (*"remind me to refactor swagger"*) into the corresponding
`ctx remind` command. Handles date conversion for `--after` flags.

**Wraps**: `ctx remind`, `ctx remind list`, `ctx remind dismiss`

**See also**:
[Session Reminders](../recipes/session-reminders.md)

---

## Session Control

Skills for controlling hook behavior during a session.

### `/ctx-pause`

Pause all context nudge and reminder hooks for the current session.
Security hooks still fire. Use for quick investigations or tasks that
don't need ceremony overhead.

**Wraps**: `ctx pause`

**Trigger phrases**: "pause ctx", "pause context", "stop the nudges",
"quiet mode"

**See also**:
[Pausing Context Hooks](../recipes/session-pause.md)

---

### `/ctx-resume`

Resume context hooks after a pause. Restores normal nudge, reminder,
and ceremony behavior. Silent no-op if not paused.

**Wraps**: `ctx resume`

**Trigger phrases**: "resume ctx", "resume context", "turn nudges back on",
"unpause"

**See also**:
[Pausing Context Hooks](../recipes/session-pause.md)

---

## Project-Specific Skills

The ctx plugin ships the skills listed above.
Teams can add their own project-specific skills to `.claude/skills/` in the
project root: These are separate from plugin-shipped skills and are scoped
to the project.

Project-specific skills follow the same format and are invoked the same way.

Custom skills are not covered in this reference.
