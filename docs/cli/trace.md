---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Commit Context Tracing
icon: lucide/git-commit-horizontal
---

### `ctx trace`

Show the context behind git commits. Links commits back to the
decisions, tasks, learnings, and sessions that motivated them.

`git log` shows *what* changed, `git blame` shows *who* —
`ctx trace` shows *why*.

```bash
ctx trace [commit] [flags]
```

**Flags**:

| Flag       | Description                        |
|------------|------------------------------------|
| `--last N` | Show context for last N commits    |
| `--json`   | Output as JSON for scripting       |

**Examples**:

```bash
# Show context for a specific commit
ctx trace abc123

# Show context for last 10 commits
ctx trace --last 10

# JSON output
ctx trace abc123 --json
```

**Output**:

```
Commit: abc123 "Fix auth token expiry"
Date:   2026-03-14 10:00:00 -0700
Context:
  [Decision] #12: Use short-lived tokens with server-side refresh
    Date: 2026-03-10

  [Task] #8: Implement token rotation for compliance
    Status: completed
```

When listing recent commits with `--last`:

```
abc123  Fix auth token expiry         decision:12, task:8
def456  Add rate limiting             decision:15, learning:7
789abc  Update dependencies           (none)
```

---

### `ctx trace file`

Show the context trail for a file. Combines `git log` with
context resolution.

```bash
ctx trace file <path[:line-range]> [flags]
```

**Flags**:

| Flag       | Description                              |
|------------|------------------------------------------|
| `--last N` | Maximum commits to show (default: 20)    |

**Examples**:

```bash
# Show context trail for a file
ctx trace file src/auth.go

# Show context for specific line range
ctx trace file src/auth.go:42-60
```

---

### `ctx trace tag`

Manually tag a commit with context. For commits made without the
hook, or to add extra context after the fact.

Tags are stored in `.context/trace/overrides.jsonl` since git
trailers cannot be added to existing commits without rewriting
history.

```bash
ctx trace tag <commit> --note "<text>"
```

**Examples**:

```bash
ctx trace tag HEAD --note "Hotfix for production outage"
ctx trace tag abc123 --note "Part of Q1 compliance initiative"
```

---

### `ctx trace hook`

Enable or disable the prepare-commit-msg hook for automatic
context tracing. When enabled, commits automatically receive a
`ctx-context` trailer with references to relevant decisions,
tasks, learnings, and sessions.

```bash
ctx trace hook <enable|disable>
```

**What the hook does**:

1. Before each commit, collects context from three sources:
   - **Pending context** accumulated during work (`ctx add`, `ctx task complete`)
   - **Staged file changes** to `.context/` files
   - **Working state** (in-progress tasks, active AI session)
2. Injects a `ctx-context` trailer into the commit message
3. After commit, records the mapping in `.context/trace/history.jsonl`

**Examples**:

```bash
# Install the hook
ctx trace hook enable

# Remove the hook
ctx trace hook disable
```

**Resulting commit message**:

```
Fix auth token expiry handling

Refactored token refresh logic to handle edge case
where refresh token expires during request.

ctx-context: decision:12, task:8, session:abc123
```

---

### Reference Types

The `ctx-context` trailer supports these reference types:

| Prefix           | Points to                  | Example                             |
|------------------|----------------------------|-------------------------------------|
| `decision:<n>`   | Entry #n in DECISIONS.md   | `decision:12`                       |
| `learning:<n>`   | Entry #n in LEARNINGS.md   | `learning:5`                        |
| `task:<n>`       | Task #n in TASKS.md        | `task:8`                            |
| `convention:<n>` | Entry #n in CONVENTIONS.md | `convention:3`                      |
| `session:<id>`   | AI session by ID           | `session:abc123`                    |
| `"<text>"`       | Free-form context note     | `"Performance fix for P1 incident"` |

---

### Storage

Context trace data is stored in the `.context/` directory:

| File                            | Purpose                          | Lifecycle                    |
|---------------------------------|----------------------------------|------------------------------|
| `state/pending-context.jsonl`   | Accumulates refs during work     | Truncated after each commit  |
| `trace/history.jsonl`           | Permanent commit-to-context map  | Append-only, never truncated |
| `trace/overrides.jsonl`         | Manual tags for existing commits | Append-only                  |
