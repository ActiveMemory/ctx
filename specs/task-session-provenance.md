# Spec: Task Provenance via Hook Relay

Link tasks to the session, commit, and branch that created them,
enabling full traceability from TASKS.md entries back to
conversations, code state, and git branches.

## Problem

Tasks added via `ctx add task` get a `#added:` timestamp but no
session, commit, or branch identifier. When reviewing tasks later,
there's no way to trace back to the conversation where the task was
discussed or the code state when it was created.

Journal entries already carry session IDs (from Claude Code JSONL
parsing). Tasks do not. The link is missing.

## Design: Hook Relay as Context Injection

The session ID, commit hash, and branch name are injected into the
agent's context window via the existing `UserPromptSubmit` hook
relay mechanism. No JSONL parsing, no env vars, no state files.

### How It Works

1. The `UserPromptSubmit` hook already receives `session_id` in
   its stdin JSON payload from Claude Code.
2. The hook script runs `git rev-parse --short HEAD` and
   `git branch --show-current` (instant, no I/O).
3. The hook emits a context block that includes all three values.
4. Claude Code relays this to the agent as part of the conversation.
5. The agent sees the block on every user prompt, absorbing the
   values through natural repetition.
6. When calling `ctx add task`, the agent passes the values via
   flags: `--session`, `--branch`, `--commit`.

### Relay Block Format

```
┌─ Context ────────────────────────────────────────────
│  Session: a92cadca
│  Branch: main @ 68fbc00a
│  Context: .context
│  [3] relevant pr to catalog: ...
└──────────────────────────────────────────────────
```

The session ID is the first 8 characters of the Claude session UUID
(matching `journal.ShortIDLen` used in journal filenames). The
commit hash is `git rev-parse --short HEAD` (7-8 chars). The branch
is `git branch --show-current`.

### Why Not JSONL Parsing or Env Vars?

- **JSONL parsing** assumes a single active session per project
  directory. With agent teams and parallel worktrees, multiple
  sessions write to the same project dir simultaneously. "Most
  recently modified" picks the wrong session.
- **Env vars** (`CTX_SESSION_ID`) require Claude Code to export
  the session ID to child processes, which it does not do.
- **Per-session state files** require counters, cleanup on resume,
  and handling of stale state.
- **Hook relay** is zero-state: the hook receives the session ID
  from Claude Code on every prompt, emits it, and forgets. No
  counters, no cleanup, no resume edge cases. The repetition *is*
  the persistence mechanism --- the agent's context window serves
  as the "store."

### Why Every Prompt (Not Just First)?

Emitting the session ID only on the first prompt would require
per-session counters, which introduces state management:
- How to detect session resume vs. new session?
- How to clean up stale counters?
- What if the first relay is lost to context compression?

Emitting on every prompt avoids all of this. A single short line
(`Session: a92cadca`) is no more distracting than the existing
`Context: .context` line. The agent absorbs it naturally.

## Task Format Change

Before:
```
- [ ] Fix the auth bug #priority:high #added:2026-04-06-143000
```

After:
```
- [ ] Fix the auth bug #priority:high #session:a92cadca #branch:main #commit:68fbc00a #added:2026-04-06-143000
```

### Flags on `ctx add task`

| Flag | Short | Required | Default |
|------|-------|----------|---------|
| `--session` | `-s` | No | `unknown` |
| `--branch` | `-b` | No | `unknown` |
| `--commit` | | No | `unknown` |

All three default to `"unknown"` when not provided. The tags are
always present so you know whether provenance was captured.

## Implementation

### 1. Hook Script Changes

In the `UserPromptSubmit` hook (e.g., `ctx system check_reminder`
or equivalent):

- Parse `session_id` from stdin JSON (already available)
- Run `git rev-parse --short HEAD` and
  `git branch --show-current`
- If git is not installed, not a repo, or the commands fail for
  any reason, use `"unknown"` for branch and commit. The hook
  must never block or error on git unavailability.
- Include all three in the relay output block
- Truncate session ID to 8 characters

### 2. Task Template Changes

`internal/assets/tpl/tpl_entry.go`:
```go
Task        = "- [ ] %s%s%s%s%s #added:%s\n"
TaskSession = " #session:%s"
TaskBranch  = " #branch:%s"
TaskCommit  = " #commit:%s"
```

### 3. Format Function Changes

`internal/cli/add/core/format/fmt.go`:
- Accept session, branch, commit parameters
- Format tags (use `"unknown"` when empty)
- Insert into template

### 4. CLI Flag Changes

`internal/cli/add/cmd/root/cmd.go`:
- Add `--session`, `--branch`, `--commit` flags
- Pass values to format function

### Correlation Flow

Given a task with `#session:a92cadca #branch:main #commit:68fbc00a`:

1. `grep a92cadca .context/journal/` finds the journal entry
   (journal filenames include the 8-char short ID)
2. `git log 68fbc00a` shows the code state when the task was
   created
3. `git log main` shows the branch history
4. Full session context is available for review

### Scope

- `ctx add task` only --- decisions and learnings already have
  timestamps that correlate well enough via journal date ranges
- No migration of existing tasks (they stay as-is, no tags)
- Tags are informational only --- no tooling depends on them yet
  (future: `ctx task trace` could use them)

## Test Plan

- [ ] `--session a92cadca` produces `#session:a92cadca` in task
- [ ] `--branch main` produces `#branch:main` in task
- [ ] `--commit 68fbc00a` produces `#commit:68fbc00a` in task
- [ ] Missing flags default to `unknown`
- [ ] All three tags always present (never omitted)
- [ ] Hook relay block includes Session and Branch lines
- [ ] Existing task format tests still pass
- [ ] Tasks with provenance tags are parseable by existing grep/
  filter tooling
