# ctx monitor — Architecture Overview

## Problem

Claude Code sessions run in isolation. When context degrades (high token usage,
repetition loops), the user must notice manually. There's no way for an external
process to inform a running session about issues.

## Solution

A monitoring system that:
1. Tracks active Claude Code sessions via tombstone files
2. Analyzes transcript health (tokens, repetition, turn count)
3. Runs pluggable auditors (context health, doc drift, etc.)
4. Writes signals that get injected into the next user prompt

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         Claude Code Session                             │
├─────────────────────────────────────────────────────────────────────────┤
│  SessionStart hook                                                      │
│    └─► creates .context/active-sessions/{id}--{encoded_path}            │
│                                                                         │
│  UserPromptSubmit hook                                                  │
│    └─► reads .context/signals/* → injects warnings into context         │
│                                                                         │
│  SessionEnd hook                                                        │
│    └─► removes .context/active-sessions/{id}--*                         │
└─────────────────────────────────────────────────────────────────────────┘
                              ▲
                              │ writes signals
                              │
┌─────────────────────────────┴───────────────────────────────────────────┐
│                      ctx monitor (separate terminal)                    │
├─────────────────────────────────────────────────────────────────────────┤
│  Loop:                                                                  │
│    1. List .context/active-sessions/ (stat only, no reads)              │
│    2. For each session:                                                 │
│       a. Parse transcript JSONL                                         │
│       b. Compute health metrics                                         │
│       c. Run auditors                                                   │
│    3. Write alerts to .context/signals/                                 │
│    4. Sleep(interval)                                                   │
└─────────────────────────────────────────────────────────────────────────┘
```

## Directory Structure

```
.context/
├── active-sessions/           # Tombstone files (created by hooks)
│   ├── abc123--...path...     # Session 1
│   └── def456--...path...     # Session 2
├── signals/                   # Alerts for injection
│   ├── context-health.md      # "Context at 85%"
│   └── doc-drift.md           # "DECISIONS.md stale"
├── TASKS.md
├── DECISIONS.md
└── ...
```

## Key Design Decisions

### Tombstone Files (No Content)

Session tracking uses empty tombstone files where the filename encodes all
needed information:

```
{session_id}--{transcript_path_with_slashes_as_underscores}
```

Benefits:
- `ls` is sufficient — no file reads needed
- Multiple sessions supported naturally
- Easy cleanup on session end
- Orphan detection via transcript file existence check

### Hook-Based Communication

The only way to "inform" a running Claude session is via hooks:
- **SessionStart**: Register session existence
- **UserPromptSubmit**: Inject signals (closest to "interrupt")
- **SessionEnd**: Cleanup

Signals aren't real-time — they're delivered on the *next* user message.

### Pluggable Auditors

Auditors are independent checkers that can be enabled/disabled:

```go
type Auditor interface {
    Name() string
    Check(sessions []Session, projectDir string) []Alert
}
```

Built-in auditors:
- **ContextHealthAuditor**: Token usage, repetition detection
- **DocDriftAuditor**: Stale .context/ files vs code changes

Future auditors (user-extensible):
- SecurityAuditor: Scan for secrets
- CodeQualityAuditor: Run linters

## CLI

```bash
# Watch all active sessions
ctx monitor

# Custom interval
ctx monitor --interval 30s

# Select auditors
ctx monitor --auditors health,drift

# Run once (for cron/scripts)
ctx monitor --once
```

## Success Criteria

- [ ] Multiple concurrent sessions tracked correctly
- [ ] Signals delivered on next user prompt
- [ ] Context health warnings at 70%, 85%, 95% thresholds
- [ ] Repetition detection catches obvious loops
- [ ] Orphan sessions cleaned up or flagged
- [ ] Monitor runs without impacting session performance
