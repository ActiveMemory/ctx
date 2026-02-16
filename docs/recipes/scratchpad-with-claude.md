---
title: "Using the Scratchpad with Claude"
icon: lucide/sticky-note
---

![ctx](../images/ctx-banner.png)

## Problem

During a session you accumulate quick notes, reminders, intermediate values,
and sometimes sensitive tokens. They don't fit `TASKS.md` (not work items) or
`DECISIONS.md` (not decisions). They don't have the structured fields that
`LEARNINGS.md` requires.

Without somewhere to put them, they get lost between sessions.

**How do you capture working memory that persists across sessions without
polluting your structured context files?**

## Commands and Skills Used

| Tool | Type | Purpose |
|------|------|---------|
| `ctx pad` | CLI command | List all scratchpad entries |
| `ctx pad show N` | CLI command | Output raw text of entry N (pipe-friendly) |
| `ctx pad add` | CLI command | Add a new entry |
| `ctx pad edit` | CLI command | Replace, append to, or prepend to an entry |
| `ctx pad rm` | CLI command | Remove an entry |
| `ctx pad mv` | CLI command | Reorder entries |
| `/ctx-pad` | Skill | Natural language interface to all pad commands |

## The Workflow

### Step 1: Add a Note Naturally

You don't need to remember any syntax. Just tell your assistant what to
jot down:

```text
You: "jot down: check DNS propagation after deploy"

Agent: "Added to scratchpad:
  1. check DNS propagation after deploy"
```

Behind the scenes, the agent runs `ctx pad add "check DNS propagation after deploy"`.

### Step 2: View Your Scratchpad

```text
You: "show my scratchpad"

Agent: "Your scratchpad has 3 entries:
  1. check DNS propagation after deploy
  2. staging API endpoint: api.staging.example.com/v2
  3. retry limit should be 5, not 3"
```

### Step 3: Edit an Entry

```text
You: "update entry 2 to include the new endpoint"

Agent: "Updated entry 2:
  was: staging API endpoint: api.staging.example.com/v2
  now: staging API endpoint: api.prod.example.com/v2"
```

For partial edits, the agent can use `--append` or `--prepend`:

```text
You: "append the port number to entry 2"

Agent: "Updated entry 2:
  staging API endpoint: api.prod.example.com/v2 :8443"
```

### Step 4: Compose Entries

Use `ctx pad show N` to pipe entry content into other operations:

```text
You: "append the content of entry 3 to entry 1"

Agent: [runs ctx pad edit 1 --append "$(ctx pad show 3)"]
       "Done. Entry 1 is now:
        check DNS propagation after deploy retry limit should be 5, not 3"
```

### Step 5: Reorder by Priority

```text
You: "move entry 4 to the top"

Agent: [runs ctx pad mv 4 1]
       "Moved entry 4 to position 1. Scratchpad reordered."
```

### Step 6: Clean Up

```text
You: "remove entry 2, it's done"

Agent: [runs ctx pad rm 2]
       "Removed entry 2. 3 entries remaining."
```

## Conversational Approach

The `/ctx-pad` skill translates natural language into commands. You
describe intent; the agent handles the mechanics.

| You say | What the agent does |
|---------|---------------------|
| "jot down: check DNS after deploy" | `ctx pad add "check DNS after deploy"` |
| "remember this: retry limit is 5" | `ctx pad add "retry limit is 5"` |
| "show my scratchpad" / "what's on my pad" | `ctx pad` |
| "show me entry 3" / "what's in entry 3" | `ctx pad show 3` |
| "delete the third one" / "remove entry 3" | `ctx pad rm 3` |
| "change entry 2 to ..." / "replace entry 2 with ..." | `ctx pad edit 2 "new text"` |
| "append '-- important' to entry 3" / "add to entry 3: ..." | `ctx pad edit 3 --append " -- important"` |
| "prepend 'URGENT:' to entry 1" | `ctx pad edit 1 --prepend "URGENT: "` |
| "add the port to entry 2" | `ctx pad edit 2 --append ":8443"` |
| "prioritize entry 4" / "move entry 4 to the top" | `ctx pad mv 4 1` |
| "move entry 1 to the bottom" / "deprioritize entry 1" | `ctx pad mv 1 N` |
| "anything on my scratchpad?" | `ctx pad` |

The skill recognizes variations: "scratchpad", "pad", "notes", "sticky
notes", "jot down", "remember this", "note to self". You don't need to use
exact trigger phrases.

The key insight is **ambiguity around "add"**: "add a note" creates a new
entry, while "add to entry 3" appends to an existing one. The skill
distinguishes these from context.

## When to Use Scratchpad vs Context Files

| Situation | Use |
|-----------|-----|
| Temporary reminders ("check X after deploy") | **Scratchpad** |
| Working values during debugging (ports, endpoints, counts) | **Scratchpad** |
| Sensitive tokens or API keys (short-term storage) | **Scratchpad** |
| Quick notes that don't fit anywhere else | **Scratchpad** |
| Work items with completion tracking | **TASKS.md** |
| Trade-offs between alternatives with rationale | **DECISIONS.md** |
| Reusable lessons with context/lesson/application | **LEARNINGS.md** |
| Codified patterns and standards | **CONVENTIONS.md** |

**Decision guide:**

* If it has structured fields (context, rationale, lesson, application),
  it belongs in a context file.
* If it's a work item you'll mark done, it belongs in `TASKS.md`.
* If it's a quick note, reminder, or working value — especially if it's
  sensitive or ephemeral — it belongs on the scratchpad.

!!! tip "Scratchpad Is Not a Junk Drawer"
    The scratchpad is for working memory, not long-term storage.
    If a note is still relevant after several sessions, promote it:
    a persistent reminder becomes a task, a recurring value becomes a
    convention, a hard-won insight becomes a learning.

## Tips

* **Entries persist across sessions.** The scratchpad is committed
  (encrypted) to git, so entries survive session boundaries. Pick up
  where you left off.
* **Entries are numbered and reorderable.** Use `ctx pad mv` to put
  high-priority items at the top.
* **`ctx pad show N` enables unix piping.** Output raw entry text
  with no numbering prefix. Compose with `--append`, `--prepend`, or
  other shell tools.
* **Never mention the key file contents to the AI.** The agent knows
  how to use `ctx pad` commands but should never read or print
  `.context/.scratchpad.key` directly.
* **Encryption is transparent.** You interact with plaintext; the
  encryption/decryption happens automatically on every read/write.

## Next Up

**[Syncing Scratchpad Notes Across Machines →](scratchpad-sync.md)**: Distribute
encryption keys and scratchpad data across environments.

## See Also

* [Scratchpad](../scratchpad.md): feature overview, all 9 commands,
  encryption details, plaintext override
* [Persisting Decisions, Learnings, and Conventions](knowledge-capture.md):
  for structured knowledge that outlives the scratchpad
* [The Complete Session](session-lifecycle.md): full session lifecycle
  showing how the scratchpad fits into the broader workflow
