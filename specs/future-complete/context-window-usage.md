# Context Window Token Usage in Checkpoint Nudges

## Problem

The check-context-size hook uses a prompt counter heuristic to nudge users
about long sessions. It says "this session is getting deep" but gives no
concrete numbers. Meanwhile, the session JSONL file contains real-time token
usage data — the actual context window fill level — that updates live with
every assistant turn.

Users benefit from seeing actual numbers in two scenarios:
1. **During existing checkpoint nudges** (prompt 20+) — adds concrete data
   to the "consider wrapping up" message
2. **When context window exceeds 80%** — urgent, fires independently of the
   prompt counter because a code-heavy session can blow past 80% in 15 turns

## Solution

A lightweight JSONL reader in the system package reads the last assistant
message's usage data from the current session's JSONL file. Called on every
`UserPromptSubmit` hook invocation (low cost: one glob + cached path + read
last ~32KB of file). Token usage displayed conditionally based on two rules.

### Why not use `recall/parser`?

It parses entire files and scans all projects — too expensive for a
per-prompt hook. The session_tokens reader finds one file and reads from
the end.

### Context window size

Default 200,000 tokens (Opus/Sonnet). Configurable via `.ctxrc` key
`context_window` for different models or future changes.

## Files Created

| File | Purpose |
|------|---------|
| `internal/cli/system/session_tokens.go` | Find JSONL file, read last usage data |
| `internal/cli/system/session_tokens_test.go` | Tests for path resolution and usage parsing |
| `internal/assets/hooks/messages/check-context-size/window.txt` | Template for >80% context window warning |

## Files Modified

| File | Change |
|------|--------|
| `internal/cli/system/check_context_size.go` | Token reading, >80% independent trigger, token line display |
| `internal/cli/system/check_context_size_test.go` | Tests for new display logic |
| `internal/assets/hooks/messages/registry.go` | `window` variant entry |
| `internal/rc/types.go` | `ContextWindow` field on `CtxRC` |
| `internal/rc/default.go` | `DefaultContextWindow = 200000` |
| `internal/rc/rc.go` | `ContextWindow()` accessor |
| `docs/recipes/customizing-hook-messages.md` | `window` variant in tables |
| `docs/home/configuration.md` | `context_window` key documented |

## Display Rules

| Condition | Token line? | Template |
|-----------|-------------|----------|
| Checkpoint fires, tokens available | Appended as info | `checkpoint` + token line |
| pct >= 80, no checkpoint due | Primary warning | `window` template |
| pct >= 80, checkpoint firing | Combined | `checkpoint` + token line (urgent) |
| No checkpoint, pct < 80 | No | silent |

### Token line format (inside checkpoint box)

- Under 80%: `⏱ Context window: ~52k tokens (~26% of 200k)`
- At/over 80%: `⚠ Context window: ~164k tokens (~82% of 200k) — running low`

## Non-Goals

- Replacing the prompt counter (it drives adaptive frequency)
- `ctx recall usage` subcommand (hook reads JSONL directly)
- Configurable 80% threshold (customize via template silence)
