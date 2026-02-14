---
name: ctx-pad
description: "Manage encrypted scratchpad. Use for short, sensitive one-liners that travel with the project."
allowed-tools: Bash(ctx:*)
---

Manage the encrypted scratchpad via `ctx pad` commands.

## When to Use

- User wants to jot down a quick note, reminder, or sensitive value
- User asks to see, add, remove, edit, or reorder scratchpad entries
- User mentions "scratchpad", "pad", "notes", or "sticky notes"

## When NOT to Use

- For structured tasks (use `ctx add task` instead)
- For architectural decisions (use `ctx add decision` instead)
- For lessons learned (use `ctx add learning` instead)

## Command Mapping

| User intent | Command |
|---|---|
| "show my scratchpad" / "what's on my pad" | `ctx pad` |
| "add a note: check DNS" | `ctx pad add "check DNS"` |
| "delete the third one" | `ctx pad rm 3` |
| "change entry 2 to ..." | `ctx pad edit 2 "new text"` |
| "move the last one to the top" | `ctx pad mv N 1` |

## Execution

**List entries:**
```bash
ctx pad
```

**Add an entry:**
```bash
ctx pad add "remember to check DNS config on staging"
```

**Remove an entry:**
```bash
ctx pad rm 2
```

**Edit an entry:**
```bash
ctx pad edit 1 "updated note text"
```

**Move an entry:**
```bash
ctx pad mv 3 1
```

## Important Notes

- The scratchpad key (.context/.scratchpad.key) must NEVER be printed to stdout
- Do not attempt to read .context/scratchpad.enc directly; always use `ctx pad`
- If the user gets a "no key" error, tell them to obtain the key file from a teammate
- Entries are one-liners; do not add multi-line content
