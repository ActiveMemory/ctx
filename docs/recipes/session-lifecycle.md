---
title: "The Complete Session"
icon: lucide/play
---

![ctx](../images/ctx-banner.png)

## Problem

"What does a full ctx session look like start to finish?"

You have ctx installed and your `.context/` directory initialized, but the 
individual commands and skills feel disconnected. How do they fit together 
into a coherent workflow? This recipe walks through a complete session -- from 
opening your editor to saving a snapshot before you close it -- so you can see 
how each piece connects.

## Commands and Skills Used

| Tool                   | Type        | Purpose                                          |
|------------------------|-------------|--------------------------------------------------|
| `ctx status`           | CLI command | Quick health check on context files              |
| `ctx agent`            | CLI command | Load token-budgeted context packet               |
| `ctx session save`     | CLI command | Save a session snapshot                          |
| `ctx session list`     | CLI command | List previous session snapshots                  |
| `ctx session load`     | CLI command | Load a previous session by index, date, or topic |
| `/ctx-remember`        | Skill       | Recall project context with structured readback  |
| `/ctx-agent`           | Skill       | Load full context packet inside the assistant    |
| `/ctx-status`          | Skill       | Show context summary with commentary             |
| `/ctx-next`            | Skill       | Suggest what to work on with rationale           |
| `/ctx-commit`          | Skill       | Commit code and prompt for context capture       |
| `/ctx-reflect`         | Skill       | Structured reflection checkpoint                 |
| `/ctx-save`            | Skill       | Save and enrich a session snapshot               |
| `/ctx-context-monitor` | Skill       | Automatic context capacity monitoring            |

## The Workflow

The session lifecycle has seven steps. You will not always use every step -- 
a quick bugfix might skip reflection, and a research session might skip 
committing -- but the full arc looks like this:

**Load context** > **Orient** > **Pick task** > **Work** > **Commit** > **Reflect** > **Save snapshot**

---

### Step 1: Load Context

Start every session by loading what you know. The fastest way is a single prompt:

```
Do you remember what we were working on?
```

This triggers the `/ctx-remember` skill. Behind the scenes, the assistant 
runs `ctx agent --budget 4000`, reads the files listed in the context packet 
(TASKS.md, DECISIONS.md, LEARNINGS.md, CONVENTIONS.md), checks 
`ctx session list --limit 3` for recent sessions, and then presents a 
structured readback:

- **Last session**: topic, date, what was accomplished
- **Active work**: pending and in-progress tasks
- **Recent context**: 1-2 decisions or learnings that matter now
- **Next step**: suggestion or question about what to focus on

The readback should feel like recall, not a file system tour. If you see 
"Let me check if there are files..." instead of a confident summary, the 
context system is not loaded properly.

**Alternative**: if you want raw data instead of a readback, run `ctx status` 
in your terminal or invoke `/ctx-status` for a summarized health check showing 
file counts, token usage, and recent activity.

---

### Step 2: Orient

After loading context, verify you understand the current state. This is where 
you read the room before acting.

```
/ctx-status
```

The status output shows which context files are populated, how many tokens they 
consume, and which files were recently modified. Look for:

- **Empty core files**: TASKS.md or CONVENTIONS.md with no content means context is sparse
- **High token count** (over 30k): context is bloated and might need `ctx compact`
- **No recent activity**: files may be stale and need updating

If the status looks healthy and the readback from Step 1 gave you enough 
context, skip ahead. If something seems off -- stale tasks, missing 
decisions -- spend a minute reading the relevant file before proceeding.

---

### Step 3: Pick What to Work On

With context loaded, choose a task. You can pick one yourself, or ask the 
assistant to recommend:

```
/ctx-next
```

The skill reads TASKS.md, checks recent sessions to avoid re-suggesting 
completed work, and presents 1-3 ranked recommendations with rationale. It 
prioritizes in-progress tasks over new starts (finishing is better than 
starting), respects explicit priority tags, and favors momentum -- continuing a 
thread from a recent session is cheaper than context-switching.

If you already know what you want to work on, state it directly:

```
Let's work on the session enrichment feature.
```

---

### Step 4: Do the Work

This is the main body of the session. Write code, fix bugs, refactor, 
research -- whatever the task requires. During this phase, a few ctx-specific 
patterns help:

**Check decisions before choosing**: when you face a design choice, check if a 
prior decision covers it.

```
Is this consistent with our decisions?
```

**Constrain scope**: keep the assistant focused on the task at hand.

```
Only change files in internal/cli/session/. Nothing else.
```

**Use `/ctx-implement` for multi-step plans**: if the task has a plan with 
multiple steps, this skill executes them one at a time with build/test 
verification between each step.

**Context monitoring runs automatically**: the `/ctx-context-monitor` skill 
is triggered by a hook at adaptive intervals. Early in a session it stays 
silent. After 16+ prompts it starts monitoring, and past 30 prompts it checks 
frequently. If context capacity is running high, it will proactively suggest 
saving unsaved work. You do not need to invoke this manually.

---

### Step 5: Commit with Context

When the work is ready, use the context-aware commit instead of raw `git commit`:

```
/ctx-commit
```

The skill runs a pre-commit build check (for Go projects, `go build`), reviews 
the staged changes, drafts a commit message focused on "why" rather than 
"what", and commits. After the commit succeeds, it prompts you:

> **Any context to capture?**
>
> - **Decision**: Did you make a design choice or trade-off?
> - **Learning**: Did you hit a gotcha or discover something?
> - **Neither**: No context to capture -- we are done.

If you made a decision, the skill records it with `ctx add decision`. If you 
learned something, it records it with `ctx add learning` including context, 
lesson, and application fields. This is the bridge between committing code and 
remembering *why* the code looks the way it does.

If source code changed in areas that affect documentation, the skill also 
offers to check for doc drift.

---

### Step 6: Reflect

At natural breakpoints -- after finishing a feature, resolving a complex bug, 
or before switching tasks -- pause to reflect:

```
/ctx-reflect
```

The skill works through a checklist: learnings discovered, decisions made, tasks 
completed or created, and whether the session is worth a full snapshot. It then 
presents a summary with specific items to persist, each with the exact 
command to run:

> I would suggest persisting:
> - **Learning**: `$PPID` in PreToolUse hooks resolves to the Claude Code PID
>   `ctx add learning --context "..." --lesson "..." --application "..."`
> - **Task**: mark "Add cooldown to ctx agent" as done
> - **Decision**: tombstone-based cooldown with 10m default
>   `ctx add decision "..."`
>
> Want me to persist any of these?

The key: the skill asks before persisting anything. You choose what to keep.

Not every commit needs reflection. A typo fix does not. But when you have 
been debugging for an hour and finally understand the root cause, that is 
worth a reflection checkpoint.

---

### Step 7: Save Snapshot

Before ending the session, save a snapshot that a future session can load:

```
/ctx-save auth-refactor
```

This runs `ctx session save` with the topic you provide, then enriches the 
saved file by replacing placeholder sections with real content:

- **Summary**: what was built or fixed, key decisions, problems encountered
- **Tasks for Next Session**: concrete items to pick up next time, in checkbox format
- **Files Modified**: list of changed files with brief descriptions

The snapshot is written to `.context/sessions/YYYY-MM-DD-HHMMSS-<topic>.md`. 
Future sessions can find it with `ctx session list` and load it with `ctx session load`.

To verify the save worked:

```bash
ctx session list --limit 3
```

---

## Putting It Together

Quick-reference checklist for a complete session:

- [ ] **Load**: "Do you remember?" or `/ctx-remember`
- [ ] **Orient**: `/ctx-status` -- check file health and token usage
- [ ] **Pick**: `/ctx-next` -- choose what to work on
- [ ] **Work**: implement, test, iterate (scope with "only change X")
- [ ] **Commit**: `/ctx-commit` -- commit and capture decisions/learnings
- [ ] **Reflect**: `/ctx-reflect` -- identify what to persist (at milestones)
- [ ] **Save**: `/ctx-save <topic>` -- snapshot for the next session

Short sessions (quick bugfix) might only use: Load, Work, Commit, Save.

Long sessions should Reflect after each major milestone and Save at least once before ending.

---

## Tips

**Save early if context is running low.** The `/ctx-context-monitor` skill will 
warn you when capacity is high, but do not wait for the warning. If you have 
been working for a while and have unsaved learnings, save proactively.

**Use descriptive topic names.** `ctx session save "auth-refactor"` is findable 
later. `ctx session save` defaults to "manual-save", which is harder to locate 
among many sessions.

**Load previous sessions by topic.** If you need context from a prior session, 
`ctx session load auth` will match by keyword. You do not need to remember the 
exact date or index number.

**Reflection is optional, saving is not.** You can skip `/ctx-reflect` for small 
changes, but always `/ctx-save` before ending a session where you did meaningful 
work. The snapshot is what the next session loads.

**Let the hook handle context loading.** The PreToolUse hook runs `ctx agent` 
automatically with a cooldown, so context loads on first tool use without you 
asking. The `/ctx-remember` prompt at session start is for *your* benefit -- 
to get a readback you can verify -- not because the assistant needs it.

## See Also

- [CLI Reference](../cli-reference.md) -- full documentation for all `ctx` commands
- [Prompting Guide](../prompting-guide.md) -- effective prompts for ctx-enabled projects
- [Tracking Work Across Sessions](task-management.md) -- deep dive on task management
- [Persisting Decisions, Learnings, and Conventions](knowledge-capture.md) -- deep dive on knowledge capture
- [Detecting and Fixing Drift](context-health.md) -- keeping context files accurate
