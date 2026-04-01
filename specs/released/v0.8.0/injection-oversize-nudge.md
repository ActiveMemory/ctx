o Injection Oversize Nudge

## Problem

The context-load-gate (v2) auto-injects project context on the first tool
use of every session. The injection size depends entirely on the project's
context files — CONSTITUTION, CONVENTIONS, ARCHITECTURE, AGENT_PLAYBOOK
verbatim, plus DECISIONS and LEARNINGS indexes.

For most projects, this is ~7-8k tokens. But context files grow:
CONVENTIONS accumulates patterns, ARCHITECTURE grows with the system,
DECISIONS and LEARNINGS index tables expand. A mature project can easily
reach 15-20k tokens of injected content.

There is no feedback loop today. The user doesn't know their injection
budget is growing. The agent won't warn — it receives the content as a
fait accompli. By the time context pressure becomes noticeable (shorter
effective sessions, more "lost in the middle" misses), the cause is
invisible.

The fix is cheap: the context-load-gate already computes the exact token
count. That number just needs to reach the user.

## Approach

Two-phase design: **detect** in context-load-gate, **notify** in
check-context-size.

1. **context-load-gate** (PreToolUse, first tool use per session):
   after computing `totalTokens`, compare against the configured
   threshold. If over, write a human-readable flag file to
   `.context/state/injection-oversize`.

2. **check-context-size** (UserPromptSubmit, adaptive frequency):
   on checkpoint passes, check for the flag file. If present, append
   a consolidation nudge to the VERBATIM relay message, then delete
   the flag (one-shot per session).

This separation matters:
- Detection happens exactly once per session (when we have the data).
- Notification targets the **user** via VERBATIM relay — not the agent.
  The agent can't rationalize skipping a message it never evaluates.
- The nudge arrives at a natural checkpoint, not buried in hook output.

### `.context/state/` Directory

New directory for project-scoped runtime state. Follows the precedent
of `.context/logs/` — gitignored, lives under `context_dir` resolution
rules, extensible for future ephemeral state.

Gitignore pattern: `.context/state/` (added to root `.gitignore` and
`config.GitignoreEntries` for `ctx init`).

### Configurable Threshold

The warning threshold must be configurable because context windows vary
by orders of magnitude:

| Model class | Context window | Reasonable threshold |
|-------------|---------------|---------------------|
| SLM (5k)   | 5,000 tokens  | 3,000 |
| Standard    | 128k-200k     | 15,000 (default) |
| Large       | 1M+           | 0 (disabled) |

New `.ctxrc` key: `injection_token_warn` (int, default 15000, 0 = disabled).

## Behavior

### Happy Path

1. User starts a new session
2. Agent uses first tool → context-load-gate fires
3. Gate reads context files, assembles injection, computes `totalTokens`
4. `totalTokens` (e.g., 18,200) exceeds `injection_token_warn` (15,000)
5. Gate writes `.context/state/injection-oversize` with:
   - Timestamp
   - Token count
   - Per-file breakdown (which files contribute most)
   - Instructions: "Run `/ctx-consolidate` to distill context files"
6. Gate proceeds normally — injects content, creates session marker
7. Later: check-context-size fires at a checkpoint (prompt #20, #25, etc.)
8. Hook reads `.context/state/injection-oversize` — file exists
9. Appends to the VERBATIM checkpoint box:
   ```
   │ ⚠ Context injection is large (~18,200 tokens).
   │ Run /ctx-consolidate to distill your context files.
   ```
10. Deletes the flag file (one-shot — don't nag every checkpoint)

### Flag File Format

Human-readable text at `.context/state/injection-oversize`:

```
Context injection oversize warning
===================================
Timestamp: 2026-02-26T14:30:00Z
Injected:  18,200 tokens (threshold: 15,000)

Per-file breakdown:
  CONSTITUTION.md    1,200 tokens
  CONVENTIONS.md     4,800 tokens  ← largest
  ARCHITECTURE.md    3,100 tokens
  AGENT_PLAYBOOK.md  2,400 tokens
  DECISIONS.md (idx)   900 tokens
  LEARNINGS.md (idx)   800 tokens

Action: Run /ctx-consolidate to distill context files.
Files with the most growth are the best candidates.
```

This format serves three audiences:
- **check-context-size hook**: reads the file, extracts token count
  for the VERBATIM nudge, then deletes it
- **User inspecting the file directly**: understands what it is,
  what to do, and which files to focus on
- **Debugging**: clear snapshot of injection state at flag-write time

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| `injection_token_warn: 0` | Disabled — never write flag file |
| Token count under threshold | No flag written, no nudge |
| Flag already exists (parallel sessions) | Overwrite — latest data wins |
| `state/` directory doesn't exist | Create it (MkdirAll) |
| check-context-size fires before gate | No flag file → no nudge (correct) |
| Flag file but no checkpoint triggers | File persists until next session's checkpoint |
| `.context/` not initialized | Both hooks bail early (existing behavior) |

### Error Handling

| Error condition | Behavior | Recovery |
|-----------------|----------|----------|
| Cannot create `state/` dir | Skip flag write, log warning | Injection proceeds normally |
| Cannot write flag file | Skip — nudge won't fire | Non-critical, silent |
| Cannot read flag in checkpoint | Skip nudge line | Checkpoint fires normally |
| Cannot delete flag after nudge | Nudge fires again on next checkpoint | Noisy but harmless |

## Interface

### CLI

No new commands. Changes are internal to existing hooks.

### Configuration

New `.ctxrc` key:

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `injection_token_warn` | int | 15000 | Token threshold for oversize warning. 0 = disabled. |

Example `.ctxrc`:
```yaml
# Warn when auto-injected context exceeds this token count.
# Set to 0 to disable. Adjust based on your model's context window.
injection_token_warn: 12000
```

## Implementation

### Files to Create/Modify

| File | Change |
|------|--------|
| `internal/config/dir.go` | Add `DirState = "state"` constant |
| `internal/config/dir.go` | Add `.context/state/` to `GitignoreEntries` |
| `internal/rc/types.go` | Add `InjectionTokenWarn int` field |
| `internal/rc/default.go` | Add `DefaultInjectionTokenWarn = 15000` |
| `internal/rc/rc.go` | Add default to `Default()` |
| `internal/cli/system/context_load_gate.go` | Write flag file when over threshold |
| `internal/cli/system/check_context_size.go` | Read flag, append nudge, delete flag |
| `internal/cli/system/context_load_gate_test.go` | Test flag write behavior |
| `internal/cli/system/check_context_size_test.go` | Test nudge append behavior |
| `.gitignore` | Add `.context/state/` entry |

### Key Implementation

#### context_load_gate.go — Flag Writer

After the existing `printHookContext` and webhook call, add:

```go
// Oversize nudge: write flag for check-context-size to pick up
warnThreshold := rc.RC().InjectionTokenWarn
if warnThreshold > 0 && totalTokens > warnThreshold {
    stateDir := filepath.Join(dir, config.DirState)
    _ = os.MkdirAll(stateDir, 0o750)

    var flag strings.Builder
    flag.WriteString("Context injection oversize warning\n")
    flag.WriteString(strings.Repeat("=", 35) + "\n")
    flag.WriteString(fmt.Sprintf("Timestamp: %s\n", time.Now().UTC().Format(time.RFC3339)))
    flag.WriteString(fmt.Sprintf("Injected:  %d tokens (threshold: %d)\n\n", totalTokens, warnThreshold))
    flag.WriteString("Per-file breakdown:\n")
    for _, entry := range perFileTokens {
        flag.WriteString(fmt.Sprintf("  %-20s %5d tokens\n", entry.name, entry.tokens))
    }
    flag.WriteString("\nAction: Run /ctx-consolidate to distill context files.\n")
    flag.WriteString("Files with the most growth are the best candidates.\n")

    _ = os.WriteFile(
        filepath.Join(stateDir, "injection-oversize"),
        []byte(flag.String()), 0o600)
}
```

This requires tracking per-file token counts during the injection loop
(a small `[]struct{ name string; tokens int }` accumulator).

#### check_context_size.go — Flag Reader

Inside the `if shouldCheck` block, before emitting the message:

```go
// Check for injection oversize flag
oversizeFile := filepath.Join(rc.ContextDir(), config.DirState, "injection-oversize")
if data, err := os.ReadFile(oversizeFile); err == nil {
    // Extract token count from flag file (line 4: "Injected:  NNNNN tokens ...")
    tokenCount := extractOversizeTokens(data)
    msg += fmt.Sprintf(
        "│ ⚠ Context injection is large (~%d tokens).\n"+
        "│ Run /ctx-consolidate to distill your context files.\n",
        tokenCount)
    _ = os.Remove(oversizeFile) // one-shot
}
```

### Helpers to Reuse

- `rc.RC()` — reads `.ctxrc` configuration (existing)
- `rc.ContextDir()` — resolves context directory (existing)
- `config.DirState` — new constant, follows `DirArchive`/`DirJournal` pattern
- `os.MkdirAll` — create state dir on first use
- `context.EstimateTokens` / `EstimateTokensString` — already used in gate

## Testing

### Unit Tests — context_load_gate_test.go

- **Under threshold**: totalTokens < warn → no flag file written
- **Over threshold**: totalTokens > warn → flag file exists with correct content
- **Threshold disabled (0)**: no flag file regardless of size
- **Per-file breakdown in flag**: verify file names and token counts appear
- **State dir auto-created**: works even when `state/` doesn't exist yet

### Unit Tests — check_context_size_test.go

- **Flag present at checkpoint**: nudge includes oversize warning line
- **Flag absent at checkpoint**: normal checkpoint, no oversize line
- **Flag deleted after nudge**: file removed after first emit
- **Malformed flag file**: nudge still fires with fallback token count (0)

### Integration

- Write context files that exceed threshold
- Run context-load-gate → verify flag file created
- Simulate prompts to checkpoint → verify VERBATIM includes nudge
- Run context-load-gate again (same session) → verify no second flag
  (marker prevents re-fire)

## Non-Goals

- **Blocking injection when over threshold**: The value of having context
  present outweighs the token cost. The nudge is advisory, not a gate.
- **Auto-consolidation**: The hook detects and notifies. The user decides
  when and how to consolidate. `/ctx-consolidate` is a separate workflow.
- **Per-file thresholds**: A single total threshold is sufficient. The
  per-file breakdown in the flag file tells the user where to look.
- **Modifying the injection strategy**: This spec adds a feedback loop
  to the existing v2 injection. It does not change what gets injected.
