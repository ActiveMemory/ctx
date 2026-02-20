---
title: "Eight Ways a Hook Can Talk"
date: 2026-02-15
author: Jose Alekhinne
topics:
  - hooks
  - agent communication
  - design patterns
  - Claude Code
---

# Eight Ways a Hook Can Talk

![ctx](../images/ctx-banner.png)

## When Your Warning Disappears

*Jose Alekhinne / 2026-02-15*

I had a backup warning that nobody ever saw.

The hook was correct: It detected stale backups, formatted a nice message,
and output it as `{"systemMessage": "..."}`. The problem wasn't detection.
The problem was **delivery**. The agent absorbed the information, processed it
internally, and never told the user.

Meanwhile, a different hook (*the journal reminder*) worked perfectly every
time. Users saw the reminder, ran the commands, and the backlog stayed
manageable. Same hook event (`UserPromptSubmit`), same project, completely
different outcomes.

The difference was one line:

```
IMPORTANT: Relay this journal reminder to the user VERBATIM
before answering their question.
```

That explicit instruction is what makes **VERBATIM relay** a *pattern*, not
just a formatting choice. And once I saw it as a pattern, I started seeing
others.

## The Audit

I looked at every hook in `ctx`: Eight shell scripts across
three hook events. And I found five distinct output patterns already in use,
plus three more that the existing hooks were reaching for but hadn't quite
articulated.

The patterns form a spectrum based on a single question: 

"**Who decides what the user sees?**"

At one end, the hook decides everything (*hard gate: the agent literally
cannot proceed*). At the other end, the hook is invisible (*silent
side-effect: nobody knows it ran*). In between, there is a range of
negotiation between hook, agent, and the user.

Here's the full spectrum:

### 1. Hard Gate

```json
{"decision": "block", "reason": "Use ctx from PATH, not ./ctx"}
```

The nuclear option: The agent's tool call is rejected before it executes.

This is Claude Code's first-class `PreToolUse` mechanism: The hook returns
JSON with `decision: block` and the agent gets an error with the reason.

Use this for **invariants**: Constitution rules, security boundaries, things
that must never happen. I use it to enforce `PATH`-based `ctx` invocation,
block `sudo`, and require explicit approval for `git push`.

### 2. VERBATIM Relay

```
IMPORTANT: Relay this warning to the user VERBATIM before answering.
┌─ Journal Reminder ─────────────────────────────
│ You have 12 sessions not yet exported.
│   ctx recall export --all
└────────────────────────────────────────────────
```

The instruction is the pattern. Without "**Relay VERBATIM**," agents tend to
absorb information into their internal reasoning and never surface it. The
explicit instruction changes the behavior from "*I know about this*" to 
"*I must tell the user about this*."

I use this for actionable reminders: 

* Unexported journal entries;
* Stale backups;
* Context capacity warnings... 

...things the user should see regardless of what they asked.

### 3. Agent Directive

```
┌─ Persistence Checkpoint (prompt #25) ───────────
│ No context files updated in 15+ prompts.
│ Have you discovered learnings worth persisting?
└──────────────────────────────────────────────────
```

A nudge, not a command. The hook tells the agent something; the agent
decides what (if anything) to tell the user. This is right for behavioral
nudges — "you haven't saved context in a while" doesn't need to be relayed
verbatim, but the agent should consider acting on it.

### 4. Silent Context Injection

```bash
ctx agent --budget 4000 2>/dev/null || true
```

Pure background enrichment. The agent's context window gets project
information injected on every tool call, with no visible output. Neither
the agent nor the user sees the hook fire — but the agent makes better
decisions because of the context.

### 5. Silent Side-Effect

```bash
find "$CTX_TMPDIR" -type f -mtime +15 -delete
```

Do work, say nothing. Temp file cleanup on session end. Logging. Marker
file management. The action is the entire point; no one needs to know.

## The Patterns We Don't Have Yet

Three more patterns emerged from the gaps in the existing hooks.

**Conditional relay** — "Relay this, but only if the user's question is
about X." Avoids noise when the warning isn't relevant. More fragile
(depends on agent judgment) but less annoying.

**Suggested action** — "Here's a problem and here's the exact command to
fix it. Ask the user before running it." Goes beyond a nudge by giving the
agent a concrete proposal, but still requires human approval.

**Escalating severity** — `INFO` gets absorbed silently. `WARN` gets
mentioned at the next natural pause. `CRITICAL` gets the VERBATIM
treatment. A protocol for hooks that produce output at different urgency
levels, so they don't all compete for the user's attention.

## The Principle

The reason this matters: **hooks are the boundary between your environment
and the agent's reasoning**. A hook that detects a problem but can't
communicate it effectively is the same as no hook at all.

The format of your output is a design decision with real consequences:

- Use a hard gate and the agent *can't* proceed (good for invariants,
  frustrating for false positives)
- Use VERBATIM relay and the user *will* see it (good for reminders,
  noisy if overused)
- Use an agent directive and the agent *might* act (good for nudges,
  unreliable for critical warnings)
- Use silent injection and nobody *knows* (good for enrichment,
  invisible when it breaks)

Choose deliberately. And when in doubt, write the word VERBATIM.

---

*The full pattern catalog with decision flowchart and implementation
examples is in the [Hook Output Patterns](../recipes/hook-output-patterns.md)
recipe.*
