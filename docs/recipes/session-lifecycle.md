---
title: "The Complete Session"
icon: lucide/play
---

![ctx](../images/ctx-banner.png)

## The Problem

"**What does a full ctx session look like start to finish?**"

You have `ctx` installed and your `.context/` directory initialized, but the
individual commands and skills feel disconnected.

How do they fit together into a coherent workflow?

This recipe walks through a complete session, from opening your editor to
persisting context before you close it, so you can see how each piece connects.

!!! tip "TL;DR"
    1. **Load**: `/ctx-remember` — load context, get structured readback
    2. **Orient**: `/ctx-status` — check file health and token usage
    3. **Pick**: `/ctx-next` — choose what to work on
    4. **Work**: implement, test, iterate
    5. **Commit**: `/ctx-commit` — commit and capture decisions/learnings
    6. **Reflect**: `/ctx-reflect` — identify what to persist (at milestones)
    7. **Wrap up**: `/ctx-wrap-up` — end-of-session ceremony

    Read on for the full walkthrough with examples.

## Commands and Skills Used

| Tool                   | Type        | Purpose                                          |
|------------------------|-------------|--------------------------------------------------|
| `ctx status`           | CLI command | Quick health check on context files              |
| `ctx agent`            | CLI command | Load token-budgeted context packet               |
| `ctx recall list`      | CLI command | List previous sessions                           |
| `ctx recall show`      | CLI command | Inspect a specific session in detail             |
| `/ctx-remember`        | Skill       | Recall project context with structured readback  |
| `/ctx-agent`           | Skill       | Load full context packet inside the assistant    |
| `/ctx-status`          | Skill       | Show context summary with commentary             |
| `/ctx-next`            | Skill       | Suggest what to work on with rationale           |
| `/ctx-commit`          | Skill       | Commit code and prompt for context capture       |
| `/ctx-reflect`         | Skill       | Structured reflection checkpoint                 |
| `/ctx-recall`          | Skill       | Browse session history inside your AI assistant  |
| `/ctx-context-monitor` | Skill       | Automatic context capacity monitoring            |

## The Workflow

The session lifecycle has seven steps. You will not always use every step
(*for example, a quick bugfix might skip reflection, and a research session
might skip committing*) but the full arc looks like this:

**Load context** > **Orient** > **Pick a Task** > **Work** > **Commit** >
**Reflect**

---

### Step 1: Load Context

Start every session by loading what you know. The fastest way is a single prompt:

```

Do you remember what we were working on?

```

This triggers the `/ctx-remember` skill. Behind the scenes, the assistant
runs `ctx agent --budget 4000`, reads the files listed in the context packet
(`TASKS.md`, `DECISIONS.md`, `LEARNINGS.md`, `CONVENTIONS.md`), checks
`ctx recall list --limit 3` for recent sessions, and then presents a
structured **readback**:

!!! note "What is a Readback?"
    A **readback** is a **structured summary** where the agent plays back what
    it knows:

    * last session
    * active tasks
    * recent decisions

    This way, you can confirm it loaded the right context.

    The term "*readback*" comes from aviation, where pilots repeat instructions
    back to air traffic control to confirm they heard correctly.

    Same idea in `ctx`: The agent tells you what it
    "*thinks*" is going on, and you correct anything that's off before the work
    begins.

    * **Last session**: topic, date, what was accomplished
    * **Active work**: pending and in-progress tasks
    * **Recent context**: 1-2 decisions or learnings that matter now
    * **Next step**: suggestion or question about what to focus on

The **readback** should feel like a **recall**, not a file system tour.
If you see "*Let me check if there are files...*" instead of a confident
summary, the context system is not loaded properly.

As an alternative, if you want raw data instead of a *readback*, run
`ctx status` in your terminal or invoke `/ctx-status` for a summarized
health check showing file counts, token usage, and recent activity.

---

### Step 2: Orient

After loading context, verify you understand the current state.

```

/ctx-status

```

The status output shows which context files are populated, how many tokens they
consume, and which files were recently modified. Look for:

* **Empty core files**: `TASKS.md` or `CONVENTIONS.md` with no content means
  the context is sparse
* **High token count** (*over 30k*): the context is **bloated** and might
  need `ctx compact`
* **No recent activity**: files may be **stale** and need updating

If the status looks healthy and the readback from Step 1 gave you enough
context, skip ahead.

If something seems off (*stale tasks, missing decisions...*), spend a minute
reading the relevant file before proceeding.

---

### Step 3: Pick What to Work On

With context loaded, choose a task. You can pick one yourself, or ask the
assistant to recommend:

```
/ctx-next
```

The skill reads `TASKS.md`, checks recent sessions to avoid re-suggesting
completed work, and presents 1-3 ranked recommendations **with rationale**.

It prioritizes in-progress tasks over new starts (*finishing is better than
starting*), respects explicit priority tags, and favors momentum:
continuing a thread from a recent session is cheaper than context-switching.

If you already know what you want to work on, state it directly:

```
Let's work on the session enrichment feature.
```

---

### Step 4: Do the Work

This is the main body of the session: write code, fix bugs, refactor,
research—whatever the task requires.

During this phase, a few `ctx`-specific patterns help:

**Check decisions before choosing**: when you face a design choice, check if a
prior decision covers it.

```
Is this consistent with our decisions?
```

**Constrain scope**: keep the assistant focused on the task at hand.

```
Only change files in internal/cli/session/. Nothing else.
```

**Use `/ctx-implement` for multistep plans**: if the task has multiple steps,
this skill executes them one at a time with build/test verification between
each step.

**Context monitoring runs automatically**: the `/ctx-context-monitor` skill
is triggered by a hook at adaptive intervals. Early in a session it stays
silent. After 16+ prompts it starts monitoring, and past 30 prompts it checks
frequently. If context capacity is running high, it will suggest saving
unsaved work. You typically do not need to invoke this skill manually.

---

### Step 5: Commit with Context

When the work is ready, use the context-aware commit instead of raw `git commit`:

```
/ctx-commit
```

As an alternative, you can ask the assistant "can we commit this?" and it will
pick up the `/ctx-commit` skill for you.

The skill runs a pre-commit build check (*for Go projects, `go build`*),
reviews the staged changes, drafts a commit message focused on "*why*"
rather than "*what*", and then commits.

After the commit succeeds, it prompts you:

```text
**Any context to capture?**

- **Decision**: Did you make a design choice or trade-off?
- **Learning**: Did you hit a gotcha or discover something?
- **Neither**: No context to capture; we are done.
```

If you made a decision, the skill records it with `ctx add decision`. If you
learned something, it records it with `ctx add learning` including context,
lesson, and application fields. This is the bridge between committing code and
remembering why the code looks the way it does.

If source code changed in areas that affect documentation, the skill also
offers to check for doc drift.

!!! tip "The Agent May Recommend Committing"
    You do not always need to invoke `/ctx-commit` explicitly.

    After a commit, the agent may proactively offer to capture context:

    "*We just made a trade-off there. Want me to record it as a decision?*"

    This is normal: The Agent Playbook encourages persisting at milestones,
    and a commit is a natural milestone.

---

### Step 6: Reflect

At natural breakpoints (*after finishing a feature, resolving a complex bug,
or before switching tasks*) pause to reflect:

```
/ctx-reflect
```

The skill works through a checklist: learnings discovered, decisions made, tasks
completed or created, and whether there are items worth persisting. It then
presents a summary with specific items to persist, each with the exact command
to run:

> I would suggest persisting:
>
> - **Learning**: `$PPID` in PreToolUse hooks resolves to the Claude Code PID
>   `ctx add learning --context "..." --lesson "..." --application "..."`
> - **Task**: mark "Add cooldown to ctx agent" as done
> - **Decision**: tombstone-based cooldown with 10m default
>   `ctx add decision "..."`
>
> Want me to persist any of these?

The skill asks before persisting anything. You choose what to keep.

Not every commit needs reflection. A typo fix does not.
But when you have been debugging for an hour and finally understand the
root cause, that is worth a reflection checkpoint.

!!! info "Agents Reflect at Milestones"
    Agents often reflect without explicit invocation.

    After completing a significant piece of work, the agent may naturally
    surface items worth persisting:

    "*We discovered that `$PPID` resolves differently inside hooks. Should I
    save that as a learning?*"

    This is the agent following the **Work-Reflect-Persist** cycle from the
    **Agent Playbook**.

    You do not need to say `/ctx-reflect` for this to happen; the agent treats
    milestones as reflection triggers on its own.

---

### Step 7: Persist Before Ending

Before ending the session, run the wrap-up ceremony to capture outstanding
learnings, decisions, conventions, and tasks:

```
/ctx-wrap-up
```

The skill gathers signal from `git diff`, recent commits, and the
conversation, then proposes structured candidates for your approval. After
you select which to keep, it persists them via `ctx add` and offers
`/ctx-commit` if uncommitted changes remain.

Session transcripts are automatically captured by Claude Code and can be
browsed later with `ctx recall list` and `ctx recall show`.

!!! tip "Ceremony Skills"
    `/ctx-remember` at session start and `/ctx-wrap-up` at session end are
    **ceremony skills** — invoke them as explicit slash commands for
    completeness, not conversationally.

    See [Session Ceremonies](session-ceremonies.md) for details on why
    these two skills are the exception to the conversational rule.

---

## Conversational Session Flow

The steps above describe the mechanics. In practice, a session with a
`ctx`-aware agent feels more like a conversation than a command sequence.

The agent knows the lifecycle and moves through it naturally, often without
explicit invocation. Here is a realistic example:

---

**You**: Do you remember what we were working on?

**Agent** *(loads context silently, presents structured readback)*:

> Last session (Feb 7, caching-layer): we implemented the in-memory cache with
> TTL expiration and wrote tests for the eviction policy. 
>
> You have 3 pending tasks; the biggest one is adding Redis as a backend.
> 
> Recent decision:<br>
> we chose LRU over LFU because the access pattern is mostly recency-based.
> 
> Want to continue with the Redis backend, or tackle something else?

**You**: Let's work on the Redis caching backend.

**Agent** *(picks up the task, begins implementation)*:

> Got it. I see the task in TASKS.md: "Add Redis backend for cache layer."
> 
> I will start with the connection pool setup. Only touching
> `internal/cache/redis/` for now.

*(The agent works... writes code, runs tests, iterates.)*

**Agent** *(after implementing the connection pool and choosing a retry strategy)*:

> The connection pool is working. I went with exponential backoff with jitter
> for retries instead of fixed intervals; it handles thundering herd better
> under reconnection storms. That is a design choice worth recording. 
> 
> Want me to save it as a decision?

**You**: Yes, go ahead.

**Agent**:

> Done. Added to DECISIONS.md: "Exponential backoff with jitter for Redis
> reconnection, to avoid thundering herd on pool recovery."

*(More work happens... the agent implements the cache read/write paths.)*

**You**: What else did we learn during this session?

**Agent**:

> A few things worth capturing:
>
> * **Learning**: Redis `SET` with `NX` and `PX` flags handles atomic
>   set-if-not-exists with TTL in a single call. No need for a separate
>   `EXPIRE` command. We initially did it in two calls and hit a race condition.
> * **Learning**: The `go-redis` library's `PoolSize` default is 10 per CPU.
>   For our workload, explicit sizing at 20 was more predictable.
> * **Task update**: The connection pool subtask is done; the read/write path
>   is implemented but needs integration tests.
>
> Want me to persist these?

**You**: Yes. And let's wrap up.

**Agent** *(reflects and persists... all in one flow)*:

> Persisted:
>
> * 2 learnings added to LEARNINGS.md
> * Marked "connection pool setup" done in TASKS.md
> * Added "write integration tests for Redis cache" as a new task
>
> Context is up to date. You are good to close.

!!! tip "Notice What Happened"
    In the above workflow, the user never typed `/ctx-reflect` or
    `ctx add learning`.

    The agent moved through **Load**, **Orient**, **Pick**, **Work**,
    **Commit**, and **Reflect** driven by **natural conversation**.

    "*Let's wrap up*" was enough to trigger the full reflect-and-persist flow.

    The agent surfaced persist-worthy items at milestones: after a design
    choice, after discovering a gotcha—without waiting to be asked.

This is the intended experience. 

The commands and skills still exist for when you want precise control, but 
**the agent is a proactive partner** in the lifecycle, not a passive executor 
of slash commands.

---

## Putting It All Together

Quick-reference checklist for a complete session:

* [ ] **Load**: `/ctx-remember` — load context and confirm readback
* [ ] **Orient**: `/ctx-status` — check file health and token usage
* [ ] **Pick**: `/ctx-next` — choose what to work on
* [ ] **Work**: implement, test, iterate (scope with "only change X")
* [ ] **Commit**: `/ctx-commit` — commit and capture decisions/learnings
* [ ] **Reflect**: `/ctx-reflect` — identify what to persist (at milestones)
* [ ] **Wrap up**: `/ctx-wrap-up` — end-of-session ceremony

Conversational equivalents: you can drive the same lifecycle with plain language:

| Step    | Slash command       | Natural language                                        |
|---------|---------------------|---------------------------------------------------------|
| Load    | `/ctx-remember`     | "Do you remember?" / "What were we working on?"         |
| Orient  | `/ctx-status`       | "How's our context looking?"                            |
| Pick    | `/ctx-next`         | "What should we work on?" / "Let's do the caching task" |
| Work    | --                  | "Only change files in internal/cache/"                  |
| Commit  | `/ctx-commit`       | "Commit this" / "Ship it"                               |
| Reflect | `/ctx-reflect`      | "What did we learn?" / *(agent offers at milestones)*   |
| Wrap up | `/ctx-wrap-up`      | *(use the slash command for completeness)*               |

The agent understands both columns.

In practice, most sessions use a mix:

* explicit commands when you want precision
* natural language when you want flow and agentic autonomy

The agent will also initiate steps on its own (*particularly "Reflect"*)
when it recognizes a milestone.

Short sessions (*quick bugfix*) might only use: **Load**, **Work**, **Commit**.

Long sessions should **Reflect** after each major milestone and persist
learnings and decisions before ending.

---

## Tips

**Persist early if context is running low.** The `/ctx-context-monitor` skill will
warn you when capacity is high, but do not wait for the warning. If you have
been working for a while and have unpersisted learnings, persist proactively.

**Browse previous sessions by topic.** If you need context from a prior session,
`ctx recall show auth` will match by keyword. You do not need to remember the
exact date or slug.

**Reflection is optional but valuable.** You can skip `/ctx-reflect` for small
changes, but always persist learnings and decisions before ending a session
where you did meaningful work. These are what the next session loads.

**Let the hook handle context loading.** The `PreToolUse` hook runs `ctx agent`
automatically with a cooldown, so context loads on first tool use without you
asking. The `/ctx-remember` prompt at session start is for your benefit (to get
a readback you can verify), not because the assistant needs it.

**The agent is a proactive partner, not a passive tool.** A `ctx`-aware agent
follows the Agent Playbook: it watches for milestones (completed tasks, design
decisions, discovered gotchas) and offers to persist them without being asked.
If you finish a tricky debugging session, it may say "That root cause is worth
saving as a learning. Want me to record it?" before you think to ask. This is
by design.

## Next Up

**[Session Ceremonies →](session-ceremonies.md)**: The two bookend
rituals for every session — `/ctx-remember` at the start,
`/ctx-wrap-up` at the end.

## See Also

* [Session Ceremonies](session-ceremonies.md): why `/ctx-remember` and
  `/ctx-wrap-up` are explicit slash commands, not conversational
* [CLI Reference](../reference/cli-reference.md): full documentation for all `ctx` commands
* [Prompting Guide](../home/prompting-guide.md): effective prompts for ctx-enabled projects
* [Tracking Work Across Sessions](task-management.md): deep dive on task management
* [Persisting Decisions, Learnings, and Conventions](knowledge-capture.md):
  deep dive on knowledge capture
* [Detecting and Fixing Drift](context-health.md): keeping context files accurate
