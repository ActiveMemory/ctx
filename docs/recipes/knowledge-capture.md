---
title: "Persisting Decisions, Learnings, and Conventions"
icon: lucide/brain
---

![ctx](../images/ctx-banner.png)

## Problem

You debug a subtle issue, discover the root cause, and move on. Three weeks
later, a different session hits the same issue. The knowledge existed briefly
in one session's memory but was never written down. Architectural decisions
suffer the same fate: you weigh trade-offs, pick an approach, and six sessions
later the AI suggests the alternative you already rejected.

How do you make sure important context survives across sessions?

## Commands and Skills Used

| Tool                    | Type    | Purpose                                       |
|-------------------------|---------|-----------------------------------------------|
| `ctx add decision`      | Command | Record an architectural decision              |
| `ctx add learning`      | Command | Record a gotcha, tip, or lesson               |
| `ctx add convention`    | Command | Record a coding pattern or standard           |
| `ctx decisions reindex` | Command | Rebuild the quick-reference index             |
| `ctx learnings reindex` | Command | Rebuild the quick-reference index             |
| `/ctx-add-decision`     | Skill   | AI-guided decision capture with validation    |
| `/ctx-add-learning`     | Skill   | AI-guided learning capture with validation    |
| `/ctx-reflect`          | Skill   | Surface items worth persisting at breakpoints |

## The Workflow

### Step 1: Understand What to Persist

Three context files serve different purposes:

**Decisions** (DECISIONS.md) answer "why is it this way?" They record
trade-offs between alternatives with structured fields: context, rationale,
and consequences. Decisions prevent re-debating settled questions.

**Learnings** (LEARNINGS.md) answer "what did we discover the hard way?"
They record gotchas and debugging insights specific to this project with
structured fields: context, lesson, and application. Learnings prevent
repeating past mistakes.

**Conventions** (CONVENTIONS.md) answer "how do we do things here?" They
record patterns and standards. No structured fields required -- just a name,
a rule, and an example. Conventions keep code consistent across sessions.

The decision point: if you chose between alternatives, it is a decision.
If you discovered something surprising, it is a learning. If you are
codifying a repeated pattern, it is a convention.

### Step 2: Record Decisions

Decisions have three structured fields: context, rationale, and
consequences. The easiest way to record one is to tell your agent:

```text
You: "Save that as a decision — we picked file-based cooldown
      tokens over env vars."

Agent: "Recorded in DECISIONS.md:
  - Decision: Use file-based cooldown tokens instead of env vars
  - Context: Hook subprocesses cannot persist env vars to parent shell
  - Rationale: File tokens survive across processes. Simpler than
    IPC. Cleanup is automatic via TTL.
  - Consequences: Tombstone files accumulate in /tmp. Cannot share
    state across machines."
```

The agent was present for the trade-off discussion. It fills in the
structured fields from what it observed — context it has but you would
struggle to reconstruct in a terminal one-liner. When multiple options
were considered, the agent includes rejected alternatives in the
rationale automatically.

The `/ctx-add-decision` skill also guides you through the fields
interactively. For quick decisions, it supports a Y-statement: "In the
context of [situation], facing [constraint], we decided for [choice]
and against [alternatives], to achieve [benefit], accepting that
[trade-off]."

??? note "CLI command for scripting and automation"
    When no agent is in the loop (CI pipelines, shell scripts, hooks),
    use the CLI directly with structured flags:

    ```bash
    ctx add decision "Use file-based cooldown tokens instead of env vars" \
      --context "Hook subprocesses cannot persist env vars to parent shell" \
      --rationale "File tokens survive across processes. Simpler than IPC. Cleanup is automatic via TTL." \
      --consequences "Tombstone files accumulate in /tmp. Cannot share state across machines."
    ```

### Step 3: Record Learnings

Learnings have three structured fields: context, lesson, and
application. Again, your agent is the best tool for this — it was there
when you discovered the gotcha and can express it more completely than
a backslash-escaped terminal command:

```text
You: "That subprocess thing was painful — save it as a learning."

Agent: "Added to LEARNINGS.md:
  - Learning: Claude Code hooks run in a subprocess
  - Context: Set env var in PreToolUse hook, but it was not visible
    in the main session
  - Lesson: Hook scripts execute in a child process. Env changes do
    not propagate to parent.
  - Application: Use tombstone files for hook-to-session
    communication. Never rely on hook env vars."
```

You said one sentence. The agent produced a structured entry with all
three fields filled from session context. That is the point — the agent
has the full conversation and can be far more expressive than what you
would type at a prompt.

The `/ctx-add-learning` skill applies three quality filters: (1) Could
someone Google this in 5 minutes? (2) Is it specific to this codebase?
(3) Did it take real effort to discover? All three must pass. Learnings
capture principles and heuristics, not code snippets.

??? note "CLI command for scripting and automation"
    When no agent is in the loop:

    ```bash
    ctx add learning "Claude Code hooks run in a subprocess" \
      --context "Set env var in PreToolUse hook, but it was not visible in the main session" \
      --lesson "Hook scripts execute in a child process. Env changes do not propagate to parent." \
      --application "Use tombstone files for hook-to-session communication. Never rely on hook env vars."
    ```

### Step 4: Record Conventions

Conventions are simpler — just a name, a rule, and optionally a section.
These are short enough that either approach works:

```text
You: "We've been using kebab-case for every CLI flag. Codify that."

Agent: "Added to CONVENTIONS.md under Naming:
  'Use kebab-case for all CLI flag names.'"
```

Or from the terminal:

```bash
ctx add convention "Use kebab-case for all CLI flag names" --section "Naming"
```

Conventions work best for rules that come up repeatedly. Codify a pattern
the third time you see it, not the first.

### Step 5: Reindex After Manual Edits

DECISIONS.md and LEARNINGS.md maintain a quick-reference index at the top --
a compact table of date and title for each entry. The index updates
automatically via `ctx add`, but falls out of sync after hand-edits.

```bash
ctx decisions reindex
ctx learnings reindex
```

Run reindex after any manual edit. The index lets AI tools scan all entries
without reading the full file, which matters when token budgets are tight.

### Step 6: Use /ctx-reflect to Surface What to Capture

At natural breakpoints -- after completing a feature, fixing a bug, or
before ending a session -- use `/ctx-reflect` to identify items worth
persisting.

```text
/ctx-reflect
```

The skill walks through learnings, decisions, tasks, and session notes,
skipping categories with nothing to report. The output includes specific
commands for each suggested persist:

> This session implemented file-based cooldown for `ctx agent` and
> discovered that hook subprocesses cannot set env vars in the parent.
>
> I'd suggest persisting:
> - **Learning**: Hook subprocesses cannot propagate env vars
>   `ctx add learning "..." --context "..." --lesson "..." --application "..."`
> - **Decision**: File-based cooldown tokens over env vars
>   `ctx add decision "..." --context "..." --rationale "..." --consequences "..."`
>
> Want me to persist any of these?

The skill always asks before persisting.

Note that `/ctx-reflect` is not the only way to trigger reflection. Agents
trained on the ctx playbook naturally surface persist-worthy items at
breakpoints -- even without invoking the skill explicitly. A conversational
prompt like "anything worth saving?" or "let's wrap up" can trigger the
same review. The skill provides a structured checklist, but the underlying
behavior is available through natural conversation.

### Step 7: The Conversational Approach

The commands shown in Steps 2-4 are precise, but you rarely need to type
them yourself. When working with an AI agent that has loaded the ctx
playbook, natural language triggers the same operations.

**You ask, the agent acts:**

| You say | What the agent does |
|---|---|
| "What have we learned this session?" | Reviews recent work, identifies learnings, offers to persist |
| "Save that as a decision" | Captures the trade-off you just discussed with structured fields |
| "Any gotchas worth remembering?" | Scans recent work for non-obvious lessons |
| "Record that convention" | Adds the pattern you just agreed on to CONVENTIONS.md |
| "Let's wrap up" | Runs through the reflect checklist conversationally |

Example exchange:

```text
You: "What did we learn from that debugging session?"

Agent: "We discovered three things worth persisting:
  1. Hook subprocesses can't propagate env vars to the parent -- we
     spent 20 minutes on this before switching to tombstone files.
  2. Tombstone files need a TTL or they accumulate in /tmp.
  3. The cooldown logic must check file mtime, not just existence.
  Want me to save these as learnings?"

You: "Yes, all of them."

Agent: "Done. Added all three to LEARNINGS.md with structured fields."
```

**The agent can also be proactive -- without being asked.** The playbook
encourages agents to persist context at natural milestones. You may see:

- After fixing a tricky bug: *"I noticed this was a non-obvious issue.
  Want me to save it as a learning?"*
- After choosing between approaches: *"We just made a trade-off between
  file-based tokens and IPC -- should I record this as a decision?"*
- At session milestones: the agent checks what is worth persisting and
  offers a brief summary.

**When agents run unattended (no human in the loop), they just do it:**

- *"I've added the subprocess env var gotcha to LEARNINGS.md so we
  don't hit it again."*
- *"Recorded the decision to use file-based tokens over env vars in
  DECISIONS.md."*
- *"Marked the cooldown task done and added a follow-up for TTL
  cleanup."*

This is by design. The playbook's self-check prompt -- "If this session
ended right now, would the next session know what happened?" -- drives
agents to persist early and often rather than waiting for explicit
instructions.

## Putting It Together

### Command-line approach (scripting and automation)

```bash
# Decision: record the trade-off
ctx add decision "Use PostgreSQL over SQLite" \
  --context "Need concurrent multi-user access" \
  --rationale "SQLite locks on writes; Postgres handles concurrency" \
  --consequences "Requires a database server; team needs Postgres training"

# Learning: record the gotcha
ctx add learning "SQL migrations must be idempotent" \
  --context "Deploy failed when migration ran twice after rollback" \
  --lesson "CREATE TABLE without IF NOT EXISTS fails on retry" \
  --application "Always use IF NOT EXISTS guards in migrations"

# Convention: record the pattern
ctx add convention "API handlers return structured errors" --section "API"

# Reindex after manual edits
ctx decisions reindex
ctx learnings reindex

# Reflect at breakpoints (from AI assistant)
# /ctx-reflect
```

### Conversational approach (interactive sessions)

```text
You: "We just picked Postgres over SQLite for the concurrency
      reasons. Save that?"

Agent: "Recorded. Added to DECISIONS.md:
  - Decision: Use PostgreSQL over SQLite
  - Context: Need concurrent multi-user access
  - Rationale: SQLite locks on writes; Postgres handles concurrency
  - Consequences: Requires a database server; team needs Postgres
    training"

You: "Also that migration thing -- the deploy blew up when the
      migration ran twice."

Agent: "Got it. Added to LEARNINGS.md:
  - Learning: SQL migrations must be idempotent
  - Lesson: CREATE TABLE without IF NOT EXISTS fails on retry
  - Application: Always use IF NOT EXISTS guards in migrations"

You: "Good. Anything else worth saving before we wrap up?"

Agent: "We've been using structured error returns in every API
  handler this session. Want me to codify that as a convention?"

You: "Do it."

Agent: "Done. Added 'API handlers return structured errors' to
  CONVENTIONS.md under the API section."
```

Both approaches produce the same structured entries in the same context
files. The conversational approach is the natural fit for interactive
sessions; the CLI commands are better suited for scripts, hooks, and
automation pipelines.

## Tips

- **Record decisions at the moment of choice.** The alternatives you
  considered and the reasons you rejected them fade quickly. Capture
  trade-offs while they are fresh.
- **Learnings should fail the Google test.** If someone could find it
  in a 5-minute search, it does not belong in LEARNINGS.md.
- **Conventions earn their place through repetition.** Add a convention
  the third time you see a pattern, not the first.
- **Use `/ctx-reflect` at every natural breakpoint.** The skill's
  checklist catches items you might otherwise lose.
- **Keep entries self-contained.** Each entry should make sense on its
  own. A future session may load only one due to token budget constraints.
- **Reindex after every hand-edit.** It takes less than a second. A stale
  index causes AI tools to miss entries.
- **Prefer the structured fields.** The verbosity forces clarity. A
  decision without rationale is just a fact; a learning without application
  is just a story.
- **Talk to your agent, don't type commands.** In interactive sessions,
  the conversational approach is the recommended way to capture knowledge.
  Say "save that as a learning" or "any decisions worth recording?" and let
  the agent handle the structured fields. Reserve the CLI commands for
  scripting, automation, and CI/CD pipelines where there is no agent in
  the loop.
- **Trust the agent's proactive instincts.** Agents trained on the ctx
  playbook will offer to persist context at milestones. Let them. A brief
  "want me to save this?" is cheaper than re-discovering the same lesson
  three sessions later.

## See Also

- [Tracking Work Across Sessions](task-management.md) -- managing the tasks
  that decisions and learnings support
- [The Complete Session](session-lifecycle.md) -- full session lifecycle
  including reflection and context persistence
- [Detecting and Fixing Drift](context-health.md) -- keeping knowledge
  files accurate as the codebase evolves
- [CLI Reference](../cli-reference.md) -- full documentation for `ctx add`,
  `ctx decisions`, `ctx learnings`
- [Context Files](../context-files.md) -- format and conventions for
  DECISIONS.md, LEARNINGS.md, and CONVENTIONS.md
