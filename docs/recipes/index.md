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

**Uses**: `ctx init`, `ctx setup`, `ctx agent`, `ctx completion`,
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

**Uses**: `ctx journal source/import`, `ctx journal site`,
`ctx journal obsidian`, `ctx serve`, `/ctx-history`,
`/ctx-journal-enrich`, `/ctx-journal-enrich-all`

---

### [Session Reminders](session-reminders.md)

Leave a **message for your next session**. Reminders surface
**automatically at session start** and repeat until dismissed.
Date-gate reminders to surface only after a specific date.

**Uses**: `ctx remind`, `ctx remind list`, `ctx remind dismiss`,
`ctx system check-reminders`

---

### [Reviewing Session Changes](session-changes.md)

See what moved since your last session: context file edits, code
commits, directories touched. Auto-detects session boundaries from
state markers.

**Uses**: `ctx change`, `ctx agent`, `ctx status`

---

### [Pausing Context Hooks](session-pause.md)

Silence all nudge hooks for a **quick task** that doesn't need ceremony
overhead. Session-scoped: Other sessions are unaffected. Security
hooks still fire.

**Uses**: `ctx hook pause`, `ctx hook resume`, `/ctx-pause`, `/ctx-resume`

---

## Knowledge & Tasks

### [Persisting Decisions, Learnings, and Conventions](knowledge-capture.md)

Record **architectural decisions** with **rationale**, capture **gotchas**
and lessons learned, and **codify** conventions so they
survive across sessions and team members.

**Uses**: `ctx add decision`, `ctx add learning`,
`ctx add convention`, `ctx decision reindex`,
`ctx learning reindex`, `/ctx-decision-add`,
`/ctx-learning-add`, `/ctx-convention-add`, `/ctx-reflect`

---

### [Tracking Work Across Sessions](task-management.md)

**Add**, **prioritize**, **complete**, **snapshot**, and **archive** tasks. Keep
`TASKS.md` focused as your project evolves across dozens of
sessions.

**Uses**: `ctx add task`, `ctx task complete`, `ctx task archive`,
`ctx task snapshot`, `/ctx-task-add`, `/ctx-archive`, `/ctx-next`

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

### [Bridging Claude Code Auto Memory](memory-bridge.md)

Mirror Claude Code's **auto memory** (MEMORY.md) into `.context/` for
**version control**, **portability**, and **drift detection**. Import
entries into structured context files with heuristic classification.

**Uses**: `ctx memory sync`, `ctx memory status`, `ctx memory diff`,
`ctx memory import`, `ctx memory publish`, `ctx system check-memory-drift`

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

**Uses**: `ctx hook message list`, `ctx hook message show`,
`ctx hook message edit`, `ctx hook message reset`

---

### [Hook Sequence Diagrams](hook-sequence-diagrams.md)

**Mermaid sequence diagrams** for every system hook: entry conditions,
state reads, output, throttling, and exit points. Includes throttling
summary table and state file reference.

**Uses**: All `ctx system` hooks

---

### [Auditing System Hooks](system-hooks-audit.md)

The 12 system hooks that run **invisibly** during every session: what each
one does, why it exists, and how to **verify** they're actually firing.
Covers webhook-based audit trails, log inspection, and detecting silent
hook failures.

**Uses**: `ctx system`, `ctx hook notify`, `.context/logs/`, `.ctxrc`
`notify.events`

---

### [Webhook Notifications](webhook-notifications.md)

Get **push notifications** when loops complete, hooks fire, or agents hit
milestones. Webhook URL is **encrypted**: never stored in plaintext.
Works with IFTTT, Slack, Discord, ntfy.sh, or any HTTP endpoint.

**Uses**: `ctx hook notify setup`, `ctx hook notify test`, `ctx hook notify --event`,
`.ctxrc` `notify.events`

---

### [Configuration Profiles](configuration-profiles.md)

Switch between **dev** and **base** runtime configurations without
editing `.ctxrc` by hand. Verbose logging and webhooks for debugging,
clean defaults for normal sessions.

**Uses**: `ctx config switch`, `ctx config status`, `/ctx-config`

---

## Maintenance

### [Detecting and Fixing Drift](context-health.md)

Keep context files accurate by detecting **structural drift**
(*stale paths, missing files, stale file ages*) and task
staleness.

**Uses**: `ctx drift`, `ctx sync`, `ctx compact`, `ctx status`,
`/ctx-drift`, `/ctx-status`, `/ctx-prompt-audit`

---

### [State Directory Maintenance](state-maintenance.md)

Clean up session tombstones from `.context/state/`. Prune old
per-session files, identify stale global markers, and keep the
state directory lean.

**Uses**: `ctx prune`

---

### [Troubleshooting](troubleshooting.md)

Diagnose hook failures, noisy nudges, stale context, and configuration
issues. Start with `ctx doctor` for a structural health check, then
use `/ctx-doctor` for agent-driven analysis of event patterns.

**Uses**: `ctx doctor`, `ctx hook event`, `/ctx-doctor`

---

### [Claude Code Permission Hygiene](claude-code-permissions.md)

Keep `.claude/settings.local.json` clean: recommended **safe defaults**,
what to **never** pre-approve, and a **maintenance workflow** for cleaning
up session debris.

**Uses**: `ctx init`, `/ctx-drift`, `/ctx-permission-sanitize`,
`ctx permission snapshot`, `ctx permission restore`

---

### [Permission Snapshots](permission-snapshots.md)

Capture a known-good permission **baseline** as a **golden image**, then restore
at session start to automatically drop session-accumulated permissions.

**Uses**: `ctx permission snapshot`, `ctx permission restore`,
`/ctx-permission-sanitize`

---

### [Turning Activity into Content](publishing.md)

Generate **blog posts** from project activity, write **changelog
posts** from commit ranges, and publish a browsable journal
site from your **session history**.

The output is generic Markdown, but the skills are tuned for the `ctx`-style
blog artifacts you see on this website.

**Uses**: `ctx journal site`, `ctx journal obsidian`, `ctx serve`,
`ctx journal import`, `/ctx-blog`, `/ctx-blog-changelog`,
`/ctx-journal-enrich`

---

### [Importing Claude Code Plans](import-plans.md)

Import Claude Code **plan files** (`~/.claude/plans/*.md`) into `specs/`
as permanent project specs. Filter by date, select interactively, and
optionally create tasks referencing each imported spec.

**Uses**: `/ctx-plan-import`, `/ctx-task-add`

---

### [Design Before Coding](design-before-coding.md)

Front-load design with a four-skill chain: **brainstorm** the approach,
**spec** the design, **task** the work, **implement** step-by-step.
Each step produces an artifact that feeds the next.

**Uses**: `/ctx-brainstorm`, `/ctx-spec`, `/ctx-task-add`,
`/ctx-implement`, `/ctx-decision-add`

---

## Agents & Automation

### [Building Project Skills](building-skills.md)

Encode repeating workflows into reusable **skills** the agent loads
automatically. Covers the full cycle: **identify** a pattern, **create**
the skill, **test** with realistic prompts, and **iterate** until it
triggers correctly.

**Uses**: `/ctx-skill-create`, `ctx init`

---

### [Running an Unattended AI Agent](autonomous-loops.md)

Set up a **loop** where an AI agent works through tasks overnight
without you at the keyboard, using `ctx` for **persistent memory**
between iterations.

This recipe shows how `ctx` supports long-running agent loops
without losing context or intent.

**Uses**: `ctx init`, `ctx loop`, `ctx watch`, `ctx load`,
`/ctx-loop`, `/ctx-implement`

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

---

### [Architecture Deep Dive](architecture-deep-dive.md)

Three-pass pipeline for understanding a codebase: **map** what
exists, **enrich** with code intelligence, then **hunt** for
where it will silently fail. Produces architecture docs,
quantified dependency data, and ranked failure hypotheses.

**Uses**: `/ctx-architecture`, `/ctx-architecture-enrich`,
`/ctx-architecture-failure-analysis`

---

### [Writing Steering Files](steering.md)

Tell your AI assistant **how to behave** — rule-based prompt
injection that fires automatically when prompts match a
description. Walks through scaffolding a steering file,
previewing matches, and syncing to each AI tool's native
format.

**Uses**: `ctx steering add`, `ctx steering preview`,
`ctx steering list`, `ctx steering sync`

---

### [Authoring Lifecycle Triggers](triggers.md)

Run **executable shell scripts** at session-start,
pre-tool-use, file-save, and other lifecycle events.
Script-based automation (complementary to steering's
rule-based prompts), with a security-first workflow: scaffold
disabled, test with mock input, enable only after review.

**Uses**: `ctx trigger add`, `ctx trigger test`,
`ctx trigger enable`, `ctx trigger disable`, `ctx trigger list`

---

## `ctx` Hub

### [`ctx` Hub: Getting Started](hub-getting-started.md)

Stand up a single-node hub on localhost, register two projects,
publish a decision from one, and watch it appear in the other.
End-to-end in under five minutes.

**Uses**: `ctx hub start`, `ctx connection register`,
`ctx connection subscribe`, `ctx connection sync`, `ctx connection listen`,
`ctx add --share`, `ctx agent --include-hub`

---

### [Personal cross-project brain](hub-personal.md)

**Story 1** day-to-day workflow: one developer, many
projects, one hub on localhost. Records a learning in
project A, watches it show up automatically in project B.
Walks through a realistic day of using the hub as passive
infrastructure — no manual `sync`, no `git push`, no
ceremony.

**Uses**: `ctx add --share`, `ctx connection subscribe`,
`ctx agent --include-hub`

---

### [Team knowledge bus](hub-team.md)

**Story 2** day-to-day workflow: a small trusted team
sharing decisions, learnings, and conventions via a hub on
an internal server. Covers the team publishing culture,
what belongs on the hub vs. local, token management, and
the social rules that make a shared knowledge stream
stay signal-rich.

**Uses**: `ctx add --share`, `ctx connection status`,
`ctx connection subscribe`, `ctx hub status`

---

### [`ctx` Hub: Multi-machine](hub-multi-machine.md)

Run the hub on a **LAN host** as a daemon and connect from project
directories on other workstations. Firewall guidance, TLS via a
reverse proxy, and safe daemon restart semantics.

**Uses**: `ctx hub start --daemon`, `ctx hub stop`,
`ctx connection register`, `ctx connection status`

---

### [`ctx` Hub: HA cluster](hub-cluster.md)

Raft-based leader election across three or more nodes for
redundancy. Covers bootstrap, runtime peer management, graceful
stepdown, and the Raft-lite durability caveat.

**Uses**: `ctx hub start --peers`, `ctx hub status`,
`ctx hub peer add/remove`, `ctx hub stepdown`

