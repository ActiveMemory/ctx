---
title: Recipes
icon: lucide/chef-hat
---

![ctx](../images/ctx-banner.png)

Workflow recipes *combining* `ctx` **commands** and **skills** to solve
*specific* problems.

---

## Getting Started

### [Guide Your Agent](guide-your-agent.md)

How commands, skills, and conversational patterns work together.
Train your agent to be proactive through **ask, guide, reinforce**.

---

### [Setup Across AI Tools](multi-tool-setup.md)

Initialize `ctx` and configure hooks for Claude Code, Cursor,
Aider, Copilot, or Windsurf. Includes **shell completion**,
**watch mode** for non-native tools, and **verification**.

**Uses**: `ctx init`, `ctx hook`, `ctx agent`, `ctx completion`,
`ctx watch`

---

### [Keeping Context in a Separate Repo](external-context.md)

Store context files **outside** the project tree: in a private repo,
shared directory, or anywhere else. Useful for open source projects
with private context or **multi-repo** setups.

**Uses**: `ctx init`, `--context-dir`, `--allow-outside-cwd`,
`.ctxrc`, `/ctx-status`

---

## Sessions

### [The Complete Session](session-lifecycle.md)

Walk through a full `ctx` session from **start to finish**:

* **Loading** context,
* **Picking** what to work on,
* **Committing** with context,
* **Capturing**, reflecting, and saving a snapshot.

**Uses**: `ctx status`, `ctx agent`,
`/ctx-remember`, `/ctx-next`, `/ctx-commit`, `/ctx-reflect`

---

### [Session Ceremonies](session-ceremonies.md)

The two bookend **rituals** for every session: `/ctx-remember` at the
start to load and confirm context, `/ctx-wrap-up` at the end to
review the session and persist **learnings**, **decisions**, and **tasks**.

**Uses**: `/ctx-remember`, `/ctx-wrap-up`, `/ctx-commit`, `ctx agent`,
`ctx add`

---

### [Browsing and Enriching Past Sessions](session-archaeology.md)

Export your AI session history to a **browsable journal site**.
**Enrich** entries with metadata and **search** across months of work.

**Uses**: `ctx recall list/show/export`, `ctx journal site`,
`ctx journal obsidian`, `ctx serve`, `/ctx-recall`,
`/ctx-journal-normalize`, `/ctx-journal-enrich`,
`/ctx-journal-enrich-all`

---

### [Session Reminders](session-reminders.md)

Leave a **message for your next session**. Reminders surface
**automatically at session start** and repeat until dismissed.
Date-gate reminders to surface only after a specific date.

**Uses**: `ctx remind`, `ctx remind list`, `ctx remind dismiss`,
`ctx system check-reminders`

---

## Knowledge & Tasks

### [Persisting Decisions, Learnings, and Conventions](knowledge-capture.md)

Record **architectural decisions** with **rationale**, capture **gotchas**
and lessons learned, and **codify** conventions so they
survive across sessions and team members.

**Uses**: `ctx add decision`, `ctx add learning`,
`ctx add convention`, `ctx decisions reindex`,
`ctx learnings reindex`, `/ctx-add-decision`,
`/ctx-add-learning`, `/ctx-add-convention`, `/ctx-reflect`

---

### [Tracking Work Across Sessions](task-management.md)

**Add**, **prioritize**, **complete**, **snapshot**, and **archive** tasks. Keep
`TASKS.md` focused as your project evolves across dozens of
sessions.

**Uses**: `ctx add task`, `ctx complete`, `ctx tasks archive`,
`ctx tasks snapshot`, `/ctx-add-task`, `/ctx-archive`, `/ctx-next`

---

### [Using the Scratchpad](scratchpad-with-claude.md)

Use the encrypted **scratchpad** for quick notes, working memory, and
sensitive values during AI sessions. Natural language in, encrypted
storage out.

**Uses**: `ctx pad`, `/ctx-pad`, `ctx pad show`, `ctx pad edit`

---

### [Syncing Scratchpad Notes Across Machines](scratchpad-sync.md)

Distribute your **scratchpad** encryption key, push and pull encrypted
notes via git, and resolve merge conflicts when two machines edit
simultaneously.

**Uses**: `ctx init`, `ctx pad`, `ctx pad resolve`, `scp`

---

## Hooks & Notifications

### [Hook Output Patterns](hook-output-patterns.md)

Choose the right output pattern for your Claude Code hooks: `VERBATIM`
relay for user-facing reminders, **hard gates** for invariants, agent
directives for nudges, and five more patterns across the spectrum.

**Uses**: ctx plugin hooks, `settings.local.json`

---

### [Customizing Hook Messages](customizing-hook-messages.md)

Customize what hooks **say** without changing what they **do**. Override
the QA gate for Python (`pytest` instead of `make lint`), silence noisy
ceremony nudges, or tailor post-commit instructions for your stack.

**Uses**: `ctx system message list`, `ctx system message show`,
`ctx system message edit`, `ctx system message reset`

---

### [Auditing System Hooks](system-hooks-audit.md)

The 12 system hooks that run **invisibly** during every session — what each
one does, why it exists, and how to **verify** they're actually firing.
Covers webhook-based audit trails, log inspection, and detecting silent
hook failures.

**Uses**: `ctx system`, `ctx notify`, `.context/logs/`, `.ctxrc`
`notify.events`

---

### [Webhook Notifications](webhook-notifications.md)

Get **push notifications** when loops complete, hooks fire, or agents hit
milestones. Webhook URL is **encrypted** — never stored in plaintext.
Works with IFTTT, Slack, Discord, ntfy.sh, or any HTTP endpoint.

**Uses**: `ctx notify setup`, `ctx notify test`, `ctx notify --event`,
`.ctxrc` `notify.events`

---

## Maintenance

### [Detecting and Fixing Drift](context-health.md)

Keep context files accurate by detecting **structural drift**
(*stale paths, missing files, stale file ages*) and task
staleness. Includes alignment audits to verify documentation
claims match agent instructions.

**Uses**: `ctx drift`, `ctx sync`, `ctx compact`, `ctx status`,
`/ctx-drift`, `/ctx-alignment-audit`, `/ctx-status`,
`/ctx-prompt-audit`

---

### [Claude Code Permission Hygiene](claude-code-permissions.md)

Keep `.claude/settings.local.json` clean: recommended **safe defaults**,
what to **never** pre-approve, and a **maintenance workflow** for cleaning
up session debris.

**Uses**: `ctx init`, `/ctx-drift`, `/sanitize-permissions`,
`ctx permissions snapshot`, `ctx permissions restore`

---

### [Permission Snapshots](permission-snapshots.md)

Capture a known-good permission **baseline** as a **golden image**, then restore
at session start to automatically drop session-accumulated permissions.

**Uses**: `ctx permissions snapshot`, `ctx permissions restore`,
`/sanitize-permissions`

---

### [Turning Activity into Content](publishing.md)

Generate **blog posts** from project activity, write **changelog
posts** from commit ranges, and publish a browsable journal
site from your **session history**.

The output is generic Markdown, but the skills are tuned for the `ctx`-style
blog artifacts you see on this website.

**Uses**: `ctx journal site`, `ctx journal obsidian`, `ctx serve`,
`ctx recall export`, `/ctx-blog`, `/ctx-blog-changelog`,
`/ctx-journal-enrich`, `/ctx-journal-normalize`

---

## Agents & Automation

### [Running an Unattended AI Agent](autonomous-loops.md)

Set up a **loop** where an AI agent works through tasks overnight
without you at the keyboard, using `ctx` for **persistent memory**
between iterations.

This recipe shows how `ctx` supports long-running agent loops
without losing context or intent.

**Uses**: `ctx init --ralph`, `ctx loop`, `ctx watch`, `ctx load`,
`/ctx-loop`, `/ctx-implement`, `/ctx-context-monitor`

---

### [When to Use a Team of Agents](when-to-use-agent-teams.md)

**Decision framework** for choosing between a single agent, parallel
worktrees, and a full agent team.

This recipe covers the file overlap test, when teams make things worse, and
what ctx provides at each level.

**Uses**: `/ctx-worktree`, `/ctx-next`, `ctx status`

---

### [Parallel Agent Development with Git Worktrees](parallel-worktrees.md)

Split a large backlog across 3-4 agents using **git worktrees**,
each on its own branch and working directory. Group tasks by
file overlap, work in parallel, merge back.

**Uses**: `/ctx-worktree`, `/ctx-next`, `git worktree`,
`git merge`
