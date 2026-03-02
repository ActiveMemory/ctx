---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Common Workflows
icon: lucide/repeat
---

![ctx](../images/ctx-banner.png)

The commands below cover what you'll use most often: 

* recording context, 
* checking health, 
* browsing history, 
* and running loops.

Each section is a self-contained snippet you can copy into your terminal.

For deeper, step-by-step guides, see [Recipes](../recipes/index.md).

## Track Context

```bash
# Add a task
ctx add task "Implement user authentication"

# Record a decision (full ADR fields required)
ctx add decision "Use PostgreSQL for primary database" \
  --context "Need a reliable database for production" \
  --rationale "PostgreSQL offers ACID compliance and JSON support" \
  --consequences "Team needs PostgreSQL training"

# Note a learning
ctx add learning "Mock functions must be hoisted in Jest" \
  --context "Tests failed with undefined mock errors" \
  --lesson "Jest hoists mock calls to top of file" \
  --application "Place jest.mock() before imports"

# Mark task complete
ctx complete "user auth"
```

## Leave a Reminder for Next Session

Drop a note that surfaces automatically at the start of your next session:

```bash
# Leave a reminder
ctx remind "refactor the swagger definitions"

# Date-gated: don't surface until a specific date
ctx remind "check CI after the deploy" --after 2026-02-25

# List pending reminders
ctx remind list

# Dismiss a reminder by ID
ctx remind dismiss 1
```

Reminders are relayed verbatim at session start by the `check-reminders` hook
and repeat every session until you dismiss them.

See [Session Reminders](../recipes/session-reminders.md) for the full recipe.

## Check Context Health

```bash
# Detect stale paths, missing files, potential secrets
ctx drift

# See full context summary
ctx status
```

## Browse Session History

List and search past AI sessions from the terminal:

```bash
ctx recall list --limit 5
```

### Journal Site

Export session transcripts to a browsable static site with search,
navigation, and topic indices.

!!! info ""
    The `ctx journal` command requires
    [zensical](https://pypi.org/project/zensical/) (**Python >= 3.10**).

    `zensical` is a Python-based static site generator from the
    *Material* for *MkDocs* team.

    (*[why zensical?](../blog/2026-02-15-why-zensical.md)*).

If you don't have it on your system,
install `zensical` once with [pipx](https://pipx.pypa.io/):

```bash
# One-time setup
pipx install zensical
```

!!! warning "Avoid `pip install zensical`"
    `pip install` often fails: For example, on macOS, system Python installs a
    non-functional stub (*`zensical` requires `Python >= 3.10`*), and
    Homebrew Python blocks system-wide installs (`PEP 668`).

    `pipx` creates an **isolated environment** with the
    **correct Python version** automatically.

### Export and Serve

Then, **export and serve**:

```bash
# Export all sessions to .context/journal/ (only new files)
ctx recall export --all

# Generate and serve the journal site
ctx journal site --serve
```

Open [http://localhost:8000](http://localhost:8000) to browse.

To update after new sessions, run the same two commands again.

### Safe By Default

`ctx recall export --all` is **safe by default**: 

* It only exports new sessions and **skips existing files**.
* Locked entries (*via `ctx recall lock`*) are **always skipped** by
  both export and enrichment skills.
* If you add `locked: true` to frontmatter during enrichment, run
  `ctx recall sync` to propagate the lock state to `.state.json`.

### Re-Exporting Existing Files

Here is how you regenerate existing files. 

**Backup your `.context` folder** before regeneration, as this is a 
potentially destructive action.

To re-export journal files, you need to explicitly opt-in using the 
`--regenerate` flag:


| Flag combination                        | Frontmatter     | Body                        |
|-----------------------------------------|-----------------|-----------------------------|
| `--regenerate`                          | Preserved       | **Overwritten** from source |
| `--regenerate --keep-frontmatter=false` | **Overwritten** | **Overwritten**             |

!!! danger "Regeneration Overwrites Body Edits"
    `--regenerate` preserves your YAML frontmatter (*tags, summary,
    enrichment metadata*) but it **replaces the Markdown body** with a
    fresh export. 

    **Any manual edits you made to the transcript will be lost**.

    **Lock entries you want to protect first**: `ctx recall lock <session-id>`.

See [Session Journal](../reference/session-journal.md) for the full pipeline
including **normalization** and **enrichment**.

## Scratchpad

Store short, sensitive one-liners in an encrypted scratchpad
that travels with the project:

```bash
# Write a note
ctx pad set db-password "postgres://user:pass@localhost/mydb"

# Read it back
ctx pad get db-password

# List all keys
ctx pad list
```

The scratchpad is encrypted with a key stored in `.context/.ctx.key`
(*`.gitignore`d by default*). 

See [Scratchpad](../reference/scratchpad.md) for details.

## Run an Autonomous Loop

Generate a script that iterates an AI agent until a completion
signal is detected:

```bash
ctx loop
chmod +x loop.sh
./loop.sh
```

See [Autonomous Loops](../operations/autonomous-loop.md) for configuration
and advanced usage.

## Agent Session Start

The first thing an AI agent should do at session start is discover where
context lives:

```bash
ctx system bootstrap
```

This prints the resolved context directory, the files in it, and the
operating rules. The `CLAUDE.md` template instructs the agent to run this
automatically. See [CLI Reference: bootstrap](../cli/system.md#ctx-system-bootstrap).

## The Two Skills You Should Always Use

Using **`/ctx-remember`** at session start and **`/ctx-wrap-up`** at
session end are the **highest-value skills** in the entire catalog:

```bash
# session begins:
/ctx-remember
... do work ...
# before closing the session:
/ctx-wrap-up
```

Let's provide some **context**, because this is **important**:

Although the agent *will* **eventually** discover your context through
`CLAUDE.md → AGENT_PLAYBOOK.md`, `/ctx-remember`
**hydrates the full context up front** (*tasks, decisions,
recent sessions*) so the agent **starts informed** rather than
piecing things together over several turns.

`/ctx-wrap-up` is the other half: A structured review that
captures learnings, decisions, and tasks before you close the
window.

Hooks like `check-persistence` remind *you* (*the user*) mid-session
that context hasn't been saved in a while, but they don't
trigger persistence automatically: You still have to act.
Also, a `CTRL+C` can end things at any moment with no reliable
"*before session end*" event. 

In short, `/ctx-wrap-up` is the **deliberate checkpoint** that makes 
sure **nothing slips through**. And `/ctx-remember` it its mirror skill
to be used at session start.

See [Session Ceremonies](../recipes/session-ceremonies.md) for
the full workflow.

## CLI Commands vs. AI Skills

Most `ctx` operations come in two flavors: a **CLI command** you run
in your terminal and an **AI skill** (*slash command*) you invoke
inside your coding assistant.

Commands and skills are **not interchangeable**: Each has a distinct role.

|                | ctx CLI command                    | ctx AI skill                                      |
|----------------|------------------------------------|---------------------------------------------------|
| **Runs where** | Your terminal                      | Inside the AI assistant                           |
| **Speed**      | Fast (*milliseconds*)              | Slower (*LLM round-trip*)                         |
| **Cost**       | Free                               | Consumes tokens and context                       |                                                   
| **Analysis**   | Deterministic heuristics           | Semantic / judgment-based                         |
| **Best for**   | Quick checks, scripting, CI        | Deep analysis, generation, workflow orchestration |

<!-- drift-check: diff <(ls internal/assets/claude/skills/ | sort) <(sed -n '/Paired Commands/,/CLI-Only Commands/p' docs/home/common-workflows.md | grep -oP 'ctx-[a-z-]+' | sort -u) -->
<!-- drift-check: diff <(ctx --help 2>&1 | sed -n '/Available Commands/,/Flags/p' | grep -oP '^\s+\K\w+' | sort) <(sed -n '/CLI-Only Commands/,/Rule of Thumb/p' docs/home/common-workflows.md | grep -oP 'ctx \K[a-z]+' | sort -u) -->

### Paired Commands

These have both a CLI and a skill counterpart. Use the CLI for
quick, deterministic checks; use the skill when you need the
agent's judgment.

| CLI                  | Skill                 | When to prefer the skill                                   |
|----------------------|-----------------------|------------------------------------------------------------|
| `ctx drift`          | `/ctx-drift`          | Semantic analysis: catches meaning drift the CLI misses    |
| `ctx status`         | `/ctx-status`         | Interpreted summary with recommendations                   |
| `ctx add task`       | `/ctx-add-task`       | Agent decomposes vague goals into concrete tasks           |
| `ctx add decision`   | `/ctx-add-decision`   | Agent drafts rationale and consequences from discussion    |
| `ctx add learning`   | `/ctx-add-learning`   | Agent extracts the lesson from a debugging session         |
| `ctx add convention` | `/ctx-add-convention` | Agent observes a repeated pattern and codifies it          |
| `ctx tasks archive`  | `/ctx-archive`        | Agent reviews which tasks are truly done                   |
| `ctx pad`            | `/ctx-pad`            | Agent reads/writes scratchpad entries in conversation flow |
| `ctx recall`         | `/ctx-recall`         | Agent searches session history with semantic understanding |
| `ctx agent`          | `/ctx-agent`          | Agent loads and acts on the context packet                 |
| `ctx loop`           | `/ctx-loop`           | Agent tailors the loop script to your project              |
| `ctx doctor`         | `/ctx-doctor`         | Agent adds semantic analysis to structural checks          |
| `ctx pause`          | `/ctx-pause`          | Agent pauses hooks with session-aware reasoning            |
| `ctx resume`         | `/ctx-resume`         | Agent resumes hooks after a pause                          |
| `ctx remind`         | `/ctx-remind`         | Agent manages reminders in conversation flow               |

### AI-Only Skills

These have no CLI equivalent. They require the agent's reasoning.

| Skill                     | Purpose                                                                                 |
|---------------------------|-----------------------------------------------------------------------------------------|
| `/ctx-remember`           | Load context and present structured readback at session start                           |
| `/ctx-wrap-up`            | End-of-session ceremony: persist learnings, decisions, tasks                            |
| `/ctx-next`               | Suggest 1-3 concrete next actions from context                                          |
| `/ctx-commit`             | Commit with integrated context capture                                                  |
| `/ctx-reflect`            | Pause and assess session progress                                                       |
| `/ctx-consolidate`        | Merge overlapping learnings or decisions                                                |
| `/ctx-alignment-audit`    | Verify docs claims match agent instructions                                             |
| `/ctx-prompt-audit`       | Analyze prompting patterns for improvement                                              |
| `/ctx-import-plans`       | Import Claude Code plan files into project specs                                        |
| `/ctx-implement`          | Execute a plan step-by-step with verification                                           |
| `/ctx-worktree`           | Manage parallel agent worktrees                                                         |
| `/ctx-journal-normalize`  | Fix markdown rendering issues in journal entries                                        |
| `/ctx-journal-enrich`     | Add metadata, tags, and summaries to journal entries                                    |
| `/ctx-journal-enrich-all` | Full journal pipeline: export if needed, then batch-enrich                               |
| `/ctx-blog`               | Generate a blog post ([zensical](https://pypi.org/project/zensical/)-flavored Markdown) |
| `/ctx-blog-changelog`     | Generate themed blog post from commits between releases                                 |
| `/ctx-map`                | Build and maintain architecture maps (ARCHITECTURE.md, DETAILED_DESIGN.md)              |

### CLI-Only Commands

These are infrastructure: used in scripts, CI, or one-time setup.

| Command                    | Purpose                                         |
|----------------------------|-------------------------------------------------|
| `ctx init`                 | Initialize `.context/` directory                |
| `ctx load`                 | Output assembled context for piping             |
| `ctx complete`             | Mark a task done by substring match             |
| `ctx sync`                 | Reconcile context with codebase state           |
| `ctx compact`              | Consolidate and clean up context files          |
| `ctx hook`                 | Generate AI tool integration config             |
| `ctx watch`                | Watch AI output and auto-apply context updates  |
| `ctx serve`                | Serve any zensical directory (default: journal) |
| `ctx permissions snapshot` | Save settings as a golden image                 |
| `ctx permissions restore`  | Restore settings from golden image              |
| `ctx journal site`         | Generate browsable journal from exports         |
| `ctx notify setup`         | Configure webhook notifications                 |
| `ctx decisions`            | List and filter decisions                       |
| `ctx learnings`            | List and filter learnings                       |
| `ctx tasks`                | List tasks, manage archival and snapshots       |
| `ctx why`                  | Read the philosophy behind ctx                  |
| `ctx guide`                | Quick-reference cheat sheet                     |
| `ctx site`                 | Site management commands                        |
| `ctx config`               | Manage runtime configuration profiles           |
| `ctx system`               | System diagnostics and hook commands            |
| `ctx system backup`        | Back up context and Claude data to tar.gz / SMB |
| `ctx completion`           | Generate shell autocompletion scripts           |

!!! tip "Rule of Thumb"
    **Quick check?** Use the CLI. 

    **Need judgment?** Use the skill.

    When in doubt, start with the CLI: It's free and instant.

    Escalate to the skill when heuristics aren't enough.

----

**Next Up**: [Context Files →](context-files.md): what each `.context/` file does and how to use it

**See Also**:

* [Recipes](../recipes/index.md): targeted how-to guides for specific tasks
* [Knowledge Capture](../recipes/knowledge-capture.md): patterns for recording decisions, learnings, and conventions
* [Context Health](../recipes/context-health.md): keeping your `.context/` accurate and drift-free
* [Session Archaeology](../recipes/session-archaeology.md): digging into past sessions
* [Task Management](../recipes/task-management.md): tracking and completing work items
