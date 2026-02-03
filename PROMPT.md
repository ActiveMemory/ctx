# Context CLI — Session Prompt

You are working on **ctx**, a CLI tool for persistent AI context management.

## On Session Start

1. Read this file (you're doing it now)
2. Run `ctx status` to see current state
3. Check `.context/TASKS.md` for what to work on

## Context Files (Read Order)

| File                         | Purpose                                  |
|------------------------------|------------------------------------------|
| `.context/CONSTITUTION.md`   | Hard rules — NEVER violate               |
| `.context/TASKS.md`          | Current work items by phase              |
| `.context/DECISIONS.md`      | Architectural decisions with rationale   |
| `.context/LEARNINGS.md`      | Gotchas and lessons learned              |
| `.context/CONVENTIONS.md`    | Code patterns and standards              |
| `.context/AGENT_PLAYBOOK.md` | How to persist context, session patterns |

## Specs

Design documents live in `specs/`. Key specs for current work:

| Spec                            | Topic                                  |
|---------------------------------|----------------------------------------|
| `specs/monitor-architecture.md` | Overall ctx monitor design             |
| `specs/active-sessions.md`      | Session tracking via tombstones        |
| `specs/context-health.md`       | Token estimation, repetition detection |
| `specs/signals.md`              | Inter-session communication            |
| `specs/auditors.md`             | Programmatic + semantic audit system   |

## Build & Test

```bash
CGO_ENABLED=0 go build -o ctx ./cmd/ctx
CGO_ENABLED=0 go test ./...
```

## Working Style

- **Ask questions** when requirements are unclear
- **Persist context** as you work (don't wait for session end)
- **Use `ctx add`** for learnings, decisions, tasks
- **Check existing patterns** before writing new code

For detailed operational patterns (session files, timestamps, Go docs standards),
see `.context/AGENT_PLAYBOOK.md`.

## Current Focus: Phase 2 — Cross-Session Monitoring

The `ctx monitor` feature enables one process to inform active Claude Code
sessions about context health issues. See `.context/TASKS.md` Phase 2 for
the full breakdown.

### Architecture Summary

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│  Claude Code    │     │   ctx monitor    │     │  Claude Code    │
│  Session A      │     │   (watcher)      │     │  Session B      │
└────────┬────────┘     └────────┬─────────┘     └────────┬────────┘
         │                       │                        │
         │ SessionStart hook     │                        │
         │ writes tombstone      │                        │
         ▼                       │                        │
┌─────────────────┐              │               ┌────────▼────────┐
│ .context/       │◄─────────────┤               │ .context/       │
│ active-sessions/│   polls      │               │ active-sessions/│
└─────────────────┘   sessions   │               └─────────────────┘
                                 │
                      ┌──────────▼──────────┐
                      │ Analyze transcripts │
                      │ - Token usage       │
                      │ - Repetition        │
                      └──────────┬──────────┘
                                 │
                      ┌──────────▼──────────┐
                      │ Write signals to    │
                      │ .context/signals/   │
                      └──────────┬──────────┘
                                 │
         ┌───────────────────────┴───────────────────────┐
         ▼                                               ▼
┌─────────────────┐                             ┌─────────────────┐
│ UserPromptSubmit│                             │ UserPromptSubmit│
│ hook injects    │                             │ hook injects    │
│ signal content  │                             │ signal content  │
└─────────────────┘                             └─────────────────┘
```

## Session End

Before finishing:
1. Mark completed tasks in TASKS.md
2. Add learnings: `ctx add learning "..."`
3. Add decisions: `ctx add decision "..."`
4. For significant sessions, save to `.context/sessions/`
