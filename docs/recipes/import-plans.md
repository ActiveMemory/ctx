---
title: "Importing Claude Code Plans"
icon: lucide/file-input
---

![ctx](../images/ctx-banner.png)

## The Problem

Claude Code plan files (`~/.claude/plans/*.md`) are ephemeral: They have
structured context, approach, and file lists, but they're orphaned after
the session ends. The filenames are UUIDs, so you can't tell what's in
them without opening each one.

**How do you turn a useful plan into a permanent project spec?**

## TL;DR

```text
You: /ctx-import-plans
Agent: [lists plans with dates and titles]
       1. 2026-02-28  Add authentication middleware
       2. 2026-02-27  Refactor database connection pool
You: "import 1"
Agent: [copies to specs/add-authentication-middleware.md]
```

Plans are copied (*not moved*) to `specs/`, slugified by their H1 heading.

## Commands and Skills Used

| Tool                 | Type  | Purpose                                     |
|----------------------|-------|---------------------------------------------|
| `/ctx-import-plans`  | Skill | List, filter, and import plan files to specs |
| `/ctx-add-task`      | Skill | Optionally add a task referencing the spec   |

## The Workflow

### Step 1: List Available Plans

Invoke the skill and it lists plans with modification dates and titles:

```text
You: /ctx-import-plans

Agent: Found 3 plan files:
         1. 2026-02-28  Add authentication middleware
         2. 2026-02-27  Refactor database connection pool
         3. 2026-02-25  Import plans skill
       Which plans would you like to import?
```

### Step 2: Filter (*Optional*)

You can narrow the list with arguments:

| Argument              | Effect                                   |
|-----------------------|------------------------------------------|
| `--today`             | Only plans modified today                |
| `--since YYYY-MM-DD`  | Only plans modified on or after the date |
| `--all`               | Import everything without prompting      |
| *(none)*              | Interactive selection                    |

```text
You: /ctx-import-plans --today
You: /ctx-import-plans --since 2026-02-27
You: /ctx-import-plans --all
```

### Step 3: Select and Import

Pick one or more plans by number:

```text
You: "import 1 and 3"

Agent: Imported 2 plan(s):
         ~/.claude/plans/abc123.md -> specs/add-authentication-middleware.md
         ~/.claude/plans/ghi789.md -> specs/import-plans-skill.md
       Want me to add tasks referencing these specs?
```

The agent reads the H1 heading from each plan and slugifies it for the
filename. If a plan has no H1 heading, the original filename (*minus
extension*) is used as the slug.

### Step 4: Add Follow-Up Tasks (*Optional*)

If you say yes, the agent creates tasks in `TASKS.md` that reference the
imported specs:

```text
You: "yes, add tasks"

Agent: [runs /ctx-add-task for each spec]
       Added:
         - [ ] Implement authentication middleware (spec: specs/add-authentication-middleware.md)
         - [ ] Import plans skill (spec: specs/import-plans-skill.md)
```

## Conversational Approach

You don't need to remember the exact skill name:

| You say                                     | What happens                    |
|---------------------------------------------|---------------------------------|
| "import my plans"                           | `/ctx-import-plans` (interactive) |
| "save today's plans as specs"               | `/ctx-import-plans --today`     |
| "import all plans from this week"           | `/ctx-import-plans --since ...` |
| "turn that plan into a spec"                | `/ctx-import-plans` (filtered)  |

## Tips

* **Plans are copied, not moved**: The originals stay in `~/.claude/plans/`.
  Claude Code manages that directory; `ctx` doesn't delete from it.
* **Conflict handling**: If `specs/{slug}.md` already exists, the agent
  asks whether to overwrite or pick a different name.
* **Specs are project memory**: Once imported, specs are tracked in git
  and available to future sessions. Reference them from `TASKS.md` phase
  headers with `Spec: specs/slug.md`.
* **Pair with `/ctx-implement`**: After importing a plan as a spec, use
  `/ctx-implement` to execute it step-by-step with verification.

## See Also

* [Skills Reference: /ctx-import-plans](../reference/skills.md#ctx-import-plans):
  full skill description
* [The Complete Session](session-lifecycle.md): where plan import fits
  in the session flow
* [Tracking Work Across Sessions](task-management.md): managing tasks
  that reference imported specs
