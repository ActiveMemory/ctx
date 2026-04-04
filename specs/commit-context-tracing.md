# Spec: Commit Context Tracing (`ctx trace`)

Link every git commit back to the decisions, tasks, learnings, and
sessions that motivated it. Today `git log` shows *what* changed and
`git blame` shows *who* — `ctx trace` shows *why*.

---

## Problem

Code changes lose their reasoning over time. A developer looking at
a six-month-old commit sees the diff but not the discussion, decision,
or lesson that drove it. The context exists — in `.context/` files,
in session histories, in task descriptions — but there is no link
from the commit to that context.

Questions like "why did we implement it this way?" require archaeology:
reading old decisions, guessing which session produced the code,
asking teammates. The answers are often lost entirely.

## Solution

Embed **context pointers** in git commit trailers. A prepare-commit-msg
hook automatically detects which context is relevant — from three
sources: accumulated pending context, staged file changes, and current
working state — then injects trailers into the commit message.
`ctx trace` resolves those pointers back to the original reasoning.

```
Fix auth token expiry handling

Refactored token refresh logic to handle edge case
where refresh token expires during request.

ctx-context: decision:12, task:8, session:abc123
```

---

## Core Concepts

### ctx-context Trailer

A standard git trailer added to commit messages. Contains one or more
comma-separated references:

```
ctx-context: decision:12, task:8, session:abc123
ctx-context: learning:5
ctx-context: "Manual note: legal compliance requirement"
```

**Reference types:**

| Prefix | Points to | Example |
|--------|-----------|---------|
| `decision:<n>` | Entry #n in DECISIONS.md | `decision:12` |
| `learning:<n>` | Entry #n in LEARNINGS.md | `learning:5` |
| `task:<n>` | Task #n in TASKS.md | `task:8` |
| `convention:<n>` | Entry #n in CONVENTIONS.md | `convention:3` |
| `session:<id>` | AI session by ID | `session:abc123` |
| `"<text>"` | Free-form context note | `"Performance fix for P1 incident"` |

Multiple `ctx-context` trailers per commit are allowed.

### Three-Source Detection

The hook collects context from **three sources** — the user does not
specify it manually.

**Source 1: Pending Context (accumulated during work)**

As ctx commands run, they append references to an accumulator file:

```
.context/state/pending-context.jsonl
```

| Event | What gets recorded |
|-------|-------------------|
| `ctx add decision "..."` | `decision:N` appended |
| `ctx add learning "..."` | `learning:N` appended |
| `ctx add convention "..."` | `convention:N` appended |
| `ctx complete N` | `task:N` appended |
| Task marked in-progress | `task:N` appended |
| AI session starts | `session:<id>` appended |

Format:
```jsonl
{"ref":"decision:12","timestamp":"2026-03-14T10:00:00Z"}
{"ref":"task:8","timestamp":"2026-03-14T10:05:00Z"}
{"ref":"session:abc123","timestamp":"2026-03-14T10:00:00Z"}
```

This captures context that happened *before* the commit — decisions
made earlier in the session, tasks completed along the way, etc.

**Source 2: Staged File Analysis (what's being committed right now)**

The hook inspects staged `.context/` files at commit time:

- If DECISIONS.md is staged → diff for added `##` headers → extract
  entry numbers
- If LEARNINGS.md is staged → diff for added `##` headers → extract
  entry numbers
- If CONVENTIONS.md is staged → diff for added `##` headers → extract
  entry numbers
- If TASKS.md is staged → diff for newly completed tasks (`- [x]`
  lines added)

This catches context changes that are part of *this* commit itself.

**Source 3: Current Working State (active context)**

- In-progress tasks in TASKS.md → `task:N`
- Active AI session via `CTX_SESSION_ID` environment variable →
  `session:<id>`

This captures the broader working context even when `.context/` files
didn't change.

**Merge & Deduplicate**

All three sources feed into a single refs list. Duplicates are removed
(same ref from pending + staged = one trailer entry). If the merged
list is empty, the hook exits silently.

### Recording Mechanism

Existing ctx commands gain a single side effect — appending one line
to `pending-context.jsonl`. This is a one-line call added to each
command's run function:

```go
// internal/trace/pending.go
func Record(ref string) error {
    f, err := os.OpenFile(pendingPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil // silent fail — tracing is best-effort
    }
    defer f.Close()
    entry := PendingEntry{Ref: ref, Timestamp: time.Now().UTC()}
    return json.NewEncoder(f).Encode(entry)
}
```

Called from:
- `ctx add` → `trace.Record("decision:N")` after successful write
- `ctx complete` → `trace.Record("task:N")` after marking done
- Session start → `trace.Record("session:<id>")` if env var is set

Recording is **best-effort** — if it fails, the command still
succeeds. Context tracing never blocks normal ctx operations.

### Optional Hook

The hook is **not installed by default**. Users opt in:

```bash
ctx hook prepare-commit-msg enable
```

This registers the prepare-commit-msg hook in the project's ctx hook
configuration. Users who don't want automatic context tagging simply
don't enable the hook.

---

## CLI Commands

### `ctx trace` — Query commit context

```bash
# Show context for a specific commit
ctx trace abc123
# Commit: abc123 "Fix auth token expiry"
# Date:   2026-03-14
#
# Context:
#   Decision #12: Use short-lived tokens with server-side refresh
#     Status: Accepted | Date: 2026-03-10
#     Rationale: Short-lived tokens reduce blast radius of token theft...
#
#   Task #8: Implement token rotation for compliance
#     Status: completed
#
#   Session: 2026-03-14-abc123 (47 messages, 12 tool calls)

# Show context for last N commits
ctx trace --last 5
# abc123  Fix auth token expiry       → decision:12, task:8
# def456  Add rate limiting           → decision:15, learning:7
# 789abc  Refactor middleware         → task:12
# ...

# Show context trail for a file (combines git log + trailer resolution)
ctx trace file src/auth.go
# abc123  Fix auth token expiry       → decision:12, task:8
# older1  Initial auth implementation → decision:3
# older2  Add OAuth2 support          → decision:7, task:2

# Show context trail for a file at a specific line range
ctx trace file src/auth.go:42-60
# abc123  Fix auth token expiry       → decision:12

# Raw output (for scripting)
ctx trace abc123 --json
```

### `ctx trace tag` — Manually tag a commit

For commits made without the hook, or to add extra context:

```bash
# Tag HEAD with context
ctx trace tag HEAD --note "Hotfix for production outage"

# Tag a specific commit
ctx trace tag abc123 --note "Part of Q1 compliance initiative"
```

Manual tags are stored in `.context/trace/overrides.jsonl` since
git trailers cannot be added to existing commits without rewriting
history.

```json
{"commit":"abc123","refs":["\"Hotfix for production outage\""],"timestamp":"2026-03-14T10:00:00Z"}
```

`ctx trace` checks the history file, commit trailer, and overrides
file when resolving context.

---

## Local Storage

### Two-Layer Storage

| File | Purpose | Lifecycle |
|------|---------|-----------|
| `state/pending-context.jsonl` | Accumulates refs during work | Truncated after each commit |
| `trace/history.jsonl` | Permanent commit→context map | Append-only, never truncated |
| `trace/overrides.jsonl` | Manual tags for existing commits | Append-only |

```
.context/
├── state/
│   └── pending-context.jsonl  ← accumulates refs between commits
├── trace/
│   ├── history.jsonl          ← permanent record of all commits
│   └── overrides.jsonl        ← manual tags for existing commits
```

### Pending Context Format

```jsonl
{"ref":"decision:12","timestamp":"2026-03-14T10:00:00Z"}
{"ref":"task:8","timestamp":"2026-03-14T10:05:00Z"}
{"ref":"session:abc123","timestamp":"2026-03-14T10:00:00Z"}
```

Append-only during work. Truncated after each commit by the
prepare-commit-msg hook.

### History Format (Permanent Record)

```jsonl
{"commit":"abc123","refs":["decision:12","task:8","session:abc123"],"message":"Fix auth token expiry","timestamp":"2026-03-14T10:00:00Z"}
{"commit":"def456","refs":["decision:15","learning:7"],"message":"Add rate limiting","timestamp":"2026-03-14T11:30:00Z"}
{"commit":"789abc","refs":["task:12"],"message":"Refactor middleware","timestamp":"2026-03-14T14:00:00Z"}
```

Written by the prepare-commit-msg hook after injecting the trailer.
This is the **primary source** for `ctx trace` — it survives even
if commits are squashed, rebased, or cherry-picked and lose their
trailers.

### Override Format

```jsonl
{"commit":"abc123","refs":["\"Hotfix for production outage\""],"timestamp":"2026-03-14T10:00:00Z"}
{"commit":"def456","refs":["decision:15"],"timestamp":"2026-03-14T11:00:00Z"}
```

Written by `ctx trace tag`. Append-only.

---

## Reference Resolution

`ctx trace` resolves context from **three sources**, merged and
deduplicated:

```
ctx trace abc123
  1. trace/history.jsonl   → refs for commit abc123 (primary)
  2. git trailer           → ctx-context from commit message (portable)
  3. trace/overrides.jsonl → manual tags (supplemental)
  4. merge all, deduplicate
  5. resolve each ref
```

The history file is the **primary source** — always available locally.
The git trailer is the **portable copy** — travels with the commit
across forks and cherry-picks. Overrides are **supplemental** — added
after the fact by `ctx trace tag`.

### Resolving Individual References

When resolving a ref, `ctx trace` reads the current state of context
files:

- **decision:12** → reads entry #12 from current DECISIONS.md
- **task:8** → reads task #8 from current TASKS.md (may be completed)
- **session:abc123** → looks up session in recall history

If an entry has been archived (via `ctx compact`), `ctx trace` falls
back to the archive directory:

```
.context/archive/YYYY-MM-DD-DECISIONS.md
```

If the reference cannot be resolved (entry deleted, session purged),
`ctx trace` shows the raw reference with a `[not found]` marker:

```
Decision #12: [not found — may have been archived]
```


## Examples

### Developer asks "why was this implemented this way?"

```bash
$ git blame src/auth/token.go | head -5
abc123 (dev1 2026-03-14) func refreshToken(ctx context.Context) {

$ ctx trace abc123
Commit: abc123 "Fix auth token expiry"
Date:   2026-03-14

Context:
  Decision #12: Use short-lived tokens with server-side refresh
    Status: Accepted | Date: 2026-03-10
    Rationale: Short-lived tokens reduce blast radius of token
    theft. Server-side refresh avoids client-side complexity.

  Task #8: Implement token rotation for compliance
    Status: completed

  Session: 2026-03-14-abc123
    47 messages, 12 tool calls
    Run `ctx recall show abc123` for full session
```

### Team reviews recent changes with context

```bash
$ ctx trace --last 10
abc123  Fix auth token expiry         → decision:12, task:8
def456  Add rate limiting middleware  → decision:15, learning:7
789abc  Refactor error handling       → task:12, convention:3
aaa111  Update dependencies           → (no context)
bbb222  Fix typo in README            → (no context)
ccc333  Implement audit logging       → decision:18, task:15, session:xyz
```

### Trace a file's history

```bash
$ ctx trace file src/auth/token.go
abc123  2026-03-14  Fix auth token expiry         → decision:12, task:8
older1  2026-02-28  Initial auth implementation   → decision:3
older2  2026-01-15  Add OAuth2 support            → decision:7, task:2
```
