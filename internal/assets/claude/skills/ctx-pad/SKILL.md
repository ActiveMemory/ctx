---
name: ctx-pad
description: "Manage encrypted scratchpad. Use for short, sensitive one-liners that travel with the project."
allowed-tools: Bash(ctx:*)
---

Manage the encrypted scratchpad via `ctx pad` commands using
natural language. Translate what the user says into the right
command.

## When to Use

- User wants to jot down a quick note, reminder, or sensitive value
- User asks to see, add, remove, edit, or reorder scratchpad entries
- User mentions "scratchpad", "pad", "notes", or "sticky notes"
- User says "jot down", "remember this", "note to self"

## When NOT to Use

- For structured tasks (use `ctx add task` instead)
- For architectural decisions (use `ctx add decision` instead)
- For lessons learned (use `ctx add learning` instead)

## Command Mapping

| User intent | Command |
|---|---|
| "show my scratchpad" / "what's on my pad" | `ctx pad` |
| "show me entry 3" / "what's in entry 3" | `ctx pad show 3` |
| "add a note: check DNS" / "jot down: check DNS" | `ctx pad add "check DNS"` |
| "delete the third one" / "remove entry 3" | `ctx pad rm 3` |
| "change entry 2 to ..." / "replace entry 2 with ..." | `ctx pad edit 2 "new text"` |
| "append '-- important' to entry 3" / "add to entry 3: ..." | `ctx pad edit 3 --append "-- important"` |
| "prepend 'URGENT:' to entry 1" | `ctx pad edit 1 --prepend "URGENT:"` |
| "move entry 4 to the top" / "prioritize entry 4" | `ctx pad mv 4 1` |
| "move entry 1 to the bottom" | `ctx pad mv 1 N` (where N = last position) |

## Execution

**List entries:**
```bash
ctx pad
```

**Show a single entry (raw text, pipe-friendly):**
```bash
ctx pad show 3
```

**Add an entry:**
```bash
ctx pad add "remember to check DNS config on staging"
```

**Remove an entry:**
```bash
ctx pad rm 2
```

**Replace an entry:**
```bash
ctx pad edit 1 "updated note text"
```

**Append to an entry:**
```bash
ctx pad edit 3 --append " -- this is important"
```

**Prepend to an entry:**
```bash
ctx pad edit 1 --prepend "URGENT: "
```

**Move an entry:**
```bash
ctx pad mv 3 1    # move entry 3 to position 1
```

**Compose entries (pipe show into edit):**
```bash
ctx pad edit 1 --append "$(ctx pad show 3)"
```

## Interpreting User Intent

When the user's intent is ambiguous:

- "update entry 2" with new text → **replace** (full rewrite)
- "add X to entry 2" → **append** (partial update)
- "put X before entry 2's text" → **prepend**
- "prioritize" / "bump up" / "move to top" → **mv N 1**
- "deprioritize" / "move to bottom" → **mv N last**

When the user says "add" — check context:
- "add a note" / "add to my pad" → `ctx pad add` (new entry)
- "add to entry 3" / "add this to the third one" → `ctx pad edit 3 --append` (modify existing)

## Important Notes

- The encryption key (`.context.key` in the context directory)
  must NEVER be printed to stdout or shown to the user
- Do not attempt to read `scratchpad.enc` in the context directory
  directly; always use `ctx pad`
- If the user gets a "no key" error, tell them to obtain the
  key file from a teammate
- Entries are one-liners; do not add multi-line content
- After modifying, show the updated scratchpad so the user can
  verify the change
