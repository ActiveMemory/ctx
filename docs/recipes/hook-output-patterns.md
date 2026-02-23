---
title: "Hook Output Patterns"
icon: lucide/megaphone
---

![ctx](../images/ctx-banner.png)

## The Problem

Claude Code hooks can output text, JSON, or nothing at all. But the *format*
of that output determines **who sees it** and **who acts on it**. 

Choose the wrong pattern, and your carefully crafted warning gets silently 
absorbed by the agent, or your agent-directed nudge gets dumped on the user 
as noise.

This recipe catalogs the known hook output patterns and explains when to 
use each one.

## The Spectrum

These patterns form a spectrum based on **who decides** what the user sees:

| Pattern              | Who decides?                          |
|----------------------|---------------------------------------|
| Hard gate            | Hook decides (agent can't proceed)    |
| VERBATIM relay       | Hook decides (agent must show)        |
| Escalating severity  | Hook suggests, agent judges urgency   |
| Conditional relay    | Hook sets criteria, agent evaluates   |
| Suggested action     | Hook proposes, agent + user decide    |
| Agent directive      | Agent decides entirely                |
| Silent injection     | Nobody — invisible background context |
| Silent side-effect   | Nobody — invisible background work    |

The spectrum runs from **full hook control** (*hard gate*) to **full
invisibility** (*silent side effect*). Most hooks belong somewhere in the
middle.

---

## Pattern 1: Hard Gate

**Block the tool call entirely**. The agent cannot proceed: it must find
another approach or tell the user.

```bash
echo '{"decision": "block", "reason": "Use ctx from PATH, not ./ctx"}'
```

**When to use**: Enforcing invariants that must never be violated:
Constitution rules, security boundaries, destructive command prevention.

**Hook type**: `PreToolUse` only (*Claude Code first-class mechanism*).

**Examples in `ctx`**:

* `ctx system block-non-path-ctx`: Enforces the PATH invocation rule
* `block-git-push.sh`: Requires explicit user approval for pushes (*project-local*)
* `block-dangerous-commands.sh`: Prevents `sudo`, copies to `~/.local/bin` (*project-local*)

**Trade-off**: The agent gets a block response with a reason. Good reasons
help the agent recover ("*use X instead*"); bad reasons leave it stuck.

---

## Pattern 2: VERBATIM Relay

**Force the agent to show this to the user as-is**. The explicit instruction
overcomes the agent's tendency to silently absorb context.

```bash
echo "IMPORTANT: Relay this warning to the user VERBATIM before answering their question."
echo ""
echo "┌─ Journal Reminder ─────────────────────────────"
echo "│ You have 12 sessions not yet exported."
echo "└────────────────────────────────────────────────"
```

**When to use**: Actionable reminders the user needs to see regardless of
what they asked: Stale backups, unexported sessions, resource warnings.

**Hook type**: `UserPromptSubmit` (*runs before the agent sees the prompt*).

**Examples in `ctx`**:

- `ctx system check-journal`: Unexported sessions and unenriched entries
- `ctx system check-context-size`: Context capacity warning
- `ctx system check-resources`: Resource pressure (memory, swap, disk, load) — DANGER only
- `check-backup-age.sh`: Stale backup warning (*project-local*)

**Trade-off**: Noisy if overused. Every VERBATIM relay adds a preamble
before the agent's actual answer. Throttle with once-per-day markers or
adaptive frequency.

**Key detail**: The phrase `IMPORTANT: Relay this ... VERBATIM` is what
makes this work. Without it, agents tend to process the information
internally and never surface it. The explicit instruction is the pattern:
the box-drawing is just fancy formatting.

---

## Pattern 3: Agent Directive

**Tell the agent to do something, not the user**. The agent decides whether
and how to involve the user.

```bash
echo "┌─ Persistence Checkpoint (prompt #25) ───────────"
echo "│ No context files updated in 15+ prompts."
echo "│ Have you discovered learnings, decisions,"
echo "│ or completed tasks worth persisting?"
echo "└──────────────────────────────────────────────────"
```

**When to use**: Behavioral nudges. The hook detects a condition and
asks the agent to consider an action. The user may never need to know.

**Hook type**: `UserPromptSubmit`.

**Examples in `ctx`**:

* `ctx system check-persistence`: Nudges the agent to persist context

**Trade-off:** No guarantee the agent acts. The nudge is one signal among
many in the context window. Strong phrasing helps ("Have you...?" is better
than "Consider..."), but ultimately the agent decides.

---

## Pattern 4: Silent Context Injection

**Load context with no visible output**. The agent gets enriched without
either party noticing.

```bash
ctx agent --budget 4000 >/dev/null || true
```

**When to use:** Background context loading that should be invisible.
The agent benefits from the information, but neither it, nor the user needs
to know it happened.

**Hook type:** `PreToolUse` with `.*` matcher (runs on every tool call).

**Examples in `ctx`**:

* The `ctx agent` `PreToolUse` hook: injects project context silently

**Trade-off**: Adds latency to every tool call. Keep the injected content
small and fast to generate.

---

## Pattern 5: Silent Side-Effect

**Do work, produce no output**: Housekeeping that needs no acknowledgment.

```bash
find "$CTX_TMPDIR" -type f -mtime +15 -delete
```

**When to use**: Cleanup, log rotation, temp file management. Anything
where the action is the point and nobody needs to know it happened.

**Hook type**: `SessionEnd`, or any hook where output is irrelevant.

**Examples in `ctx`**:

* `ctx system cleanup-tmp`: Removes stale temp files on session end

**Trade-off**: None, if the action is truly invisible. If it can fail in
a way that matters, consider logging.

### Pattern 6: Conditional Relay

**Tell the agent to relay only if a condition holds in context.**

```bash
echo "If the user's question involves modifying .context/ files,"
echo "relay this warning VERBATIM:"
echo ""
echo "┌─ Context Integrity ─────────────────────────────"
echo "│ CONSTITUTION.md has not been verified in 7 days."
echo "└────────────────────────────────────────────────"
echo ""
echo "Otherwise, proceed normally."
```

**When to use**: Warnings that only matter in certain contexts. Avoids
noise when the user is doing unrelated work.

**Trade-off**: Depends on the agent's judgment about when the condition
holds. More fragile than VERBATIM relay, but less noisy.

### Pattern 7: Suggested Action

**Give the agent a specific command to propose to the user.**

```bash
echo "┌─ Stale Dependencies ──────────────────────────"
echo "│ go.sum is 30+ days newer than go.mod."
echo "│ Suggested: run \`go mod tidy\`"
echo "│ Ask the user before proceeding."
echo "└───────────────────────────────────────────────"
```

**When to use:** The hook detects a fixable condition and knows the fix.
Goes beyond a nudge: Gives the agent **a concrete next step**. The agent
still asks for permission but knows exactly what to propose.

**Trade-off:** The suggestion might be wrong or outdated. The "*ask the
user before proceeding*" part is critical.

### Pattern 8: Escalating Severity

**Different urgency tiers with different relay expectations.**

```bash
# INFO: agent processes silently, mentions if relevant
echo "INFO: Last test run was 3 days ago."

# WARN: agent should mention to user at next natural pause
echo "WARN: 12 uncommitted changes across 3 branches."

# CRITICAL: agent must relay immediately, before any other work
echo "CRITICAL: Relay VERBATIM before answering. Disk usage at 95%."
```

**When to use:** When you have multiple hooks producing output and need
to avoid overwhelming the user. INFO gets absorbed, WARN gets mentioned,
CRITICAL interrupts.

**Examples in `ctx`**:

* `ctx system check-resources`: Uses two tiers (WARNING/DANGER) internally
  but only fires the VERBATIM relay at DANGER level — WARNING is silent.
  See `ctx system` for the user-facing command that shows both tiers.

**Trade-off:** Requires agent training or convention to recognize the
tiers. Without a shared protocol, the prefixes are just text.

---

## Choosing a Pattern

```
Is the agent about to do something forbidden?
  └─ Yes → Hard gate

Does the user need to see this regardless of what they asked?
  └─ Yes → VERBATIM relay
  └─ Sometimes → Conditional relay

Should the agent consider an action?
  └─ Yes, with a specific fix → Suggested action
  └─ Yes, open-ended → Agent directive

Is this background context the agent should have?
  └─ Yes → Silent injection

Is this housekeeping?
  └─ Yes → Silent side-effect
```

## Design Tips

**Throttle aggressively**. VERBATIM relays that fire every prompt will be
ignored or resented. Use once-per-day markers (`touch $REMINDED`), adaptive
frequency (every Nth prompt), or staleness checks (only fire if condition
persists).

**Include actionable commands**. "You have 12 unexported sessions" is less
useful than "You have 12 unexported sessions. Run: `ctx recall export --all`."
Give the user (or agent) the exact next step.

**Use box-drawing for visual structure.** The `┌─ ─┐ │ └─ ─┘` pattern
makes hook output visually distinct from agent prose. It also signals
"this is machine-generated, not agent opinion."

**Test the silence path**. Most hook runs should produce no output (*the
condition isn't met*). Make sure the common case is fast and silent.

## Common Pitfalls

Lessons from 19 days of hook debugging in `ctx`. Every one of these was
encountered, debugged, and fixed in production.

### Silent misfire: wrong key name

```json
{ "PreToolUseHooks": [ ... ] }
```

The key is `PreToolUse`, not `PreToolUseHooks`. Claude Code validates
silently -- a misspelled key means the hook is ignored with no error.
**Always test with a debug `echo` first** to confirm the hook fires
before adding real logic.

### JSON escaping breaks shell commands

Go's `json.Marshal` escapes `>`, `<`, and `&` as Unicode sequences
(`\u003e`) by default. This breaks shell commands in generated config:

```json
"command": "ctx agent 2\u003e/dev/null"
```

Fix: use `json.Encoder` with `SetEscapeHTML(false)` when generating
hook configuration.

### Stdin, not environment variables

Hook input arrives as **JSON via stdin**, not environment variables:

```bash
# Wrong:
COMMAND="$CLAUDE_TOOL_INPUT"

# Right:
HOOK_INPUT=$(cat)
COMMAND=$(echo "$HOOK_INPUT" | jq -r '.tool_input.command // empty')
```

### Regex overfitting

A regex meant to catch `ctx` as a binary will also match `ctx` as a
directory component:

```bash
# Too broad -- blocks: git -C /home/jose/WORKSPACE/ctx status
(/home/|/tmp/|/var/)[^ ]*ctx[^ ]*

# Narrow to binary only:
(/home/|/tmp/|/var/)[^ ]*/ctx( |$)
```

Test hook regexes against paths that contain the target string as a
*substring*, not just as the final component.

### Repetition fatigue

Injecting context on every tool call sounds safe. In practice, after
seeing the same context injection fifteen times, the agent treats it as
background noise -- conventions stated in the injected context get
violated because salience has been destroyed by repetition.

Fix: **cooldowns**. `ctx agent --session $PPID --cooldown 10m` injects
at most once per ten minutes per session using a tombstone file in
`/tmp/`. This is not an optimization; it is a correction for a design
flaw. Every injection consumes attention budget -- 50 tool calls at
4,000 tokens each means 200,000 tokens of repeated context, most of
it wasted.

### Hardcoded paths

A username rename (`parallels` to `jose`) broke every hook at once.
Use `$CLAUDE_PROJECT_DIR` instead of absolute paths:

```json
"command": "\"$CLAUDE_PROJECT_DIR\"/.claude/hooks/block-git-push.sh"
```

If the platform provides a runtime variable for paths, always use it.

---

## Next Up

**[Auditing System Hooks →](system-hooks-audit.md)**: Verify
system hooks are firing and audit their behavior.

## See Also

- [Claude Code Permission Hygiene](claude-code-permissions.md): how
  permissions and hooks work together
- [Defense in Depth](../blog/2026-02-09-defense-in-depth-securing-ai-agents.md):
  why hooks matter for agent security
