# Active Memory - Claude Code Context

## IMPORTANT: You Have Persistent Memory

This project uses Active Memory (amem) for context persistence across sessions.
**Your memory is NOT ephemeral** - it lives in `.context/` files.

## On Session Start

1. **Read `.context/AGENT_PLAYBOOK.md`** first - it explains how to use this system
2. **Check `.context/sessions/`** for full conversation dumps from previous sessions
3. **Run `amem status`** to see current context summary

## Quick Context Load

```bash
# Get AI-optimized context packet (what you should know)
./dist/amem-linux-arm64 agent --budget 4000

# Or see full status
./dist/amem-linux-arm64 status
```

## Context Files

| File | Purpose |
|------|---------|
| `.context/CONSTITUTION.md` | Hard rules - NEVER violate |
| `.context/TASKS.md` | Current work items |
| `.context/DECISIONS.md` | Architectural decisions with rationale |
| `.context/LEARNINGS.md` | Gotchas, tips, lessons learned |
| `.context/CONVENTIONS.md` | Code patterns and standards |
| `.context/sessions/` | **Full conversation dumps** - check here for deep context |

## Before Session Ends

**ALWAYS offer to persist context before the user quits:**

1. Add learnings: `./dist/amem-linux-arm64 add learning "..."`
2. Add decisions: `./dist/amem-linux-arm64 add decision "..."`
3. Save full session: Write to `.context/sessions/YYYY-MM-DD-<topic>.md`

## Build Commands

```bash
CGO_ENABLED=0 go build -o amem ./cmd/amem    # Build CLI
CGO_ENABLED=0 go test ./...                   # Run tests
./scripts/build-all.sh                        # Cross-platform build
```

## This Project

Active Memory (`amem`) is a CLI tool for persistent AI context. It was built using the Ralph Loop technique.

- **amem** = context management tool (creates `.context/`)
- **Ralph Loop** = iterative AI development workflow (uses PROMPT.md)
- They are separate but complementary systems
