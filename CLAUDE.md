# Context - Claude Code Context

## IMPORTANT: You Have Persistent Memory

This project uses Context (ctx) for context persistence across sessions.
**Your memory is NOT ephemeral** - it lives in `.context/` files.

## On Session Start

1. **Read `.context/AGENT_PLAYBOOK.md`** first - it explains how to use this system
2. **Run `ctx status`** to see current context summary

## When Asked "Do You Remember?"

When the user asks "Do you remember?", "What were we working on?", or any
memory-related question:

**Do this FIRST (silently):**
- Read `.context/TASKS.md`
- Read `.context/DECISIONS.md` and `.context/LEARNINGS.md`
- Run `ctx recall list --limit 5` for recent session history

**Then respond with a structured readback:**

1. **Last session**: cite the most recent session topic and date
2. **Active work**: list pending or in-progress tasks
3. **Recent context**: mention 1-2 recent decisions or learnings
4. **Next step**: offer to continue or ask what to focus on

**Never** lead with "I don't have memory", "Let me check if there are files",
or narrate your discovery process. The `.context/` files are your memory.
Read them silently, then present what you found as recall, not as a search.

## Quick Context Load

```bash
# Get AI-optimized context packet (what you should know)
ctx agent --budget 4000

# Or see full status
ctx status
```

## Context Files

| File                       | Purpose                                                   |
|----------------------------|-----------------------------------------------------------|
| `.context/CONSTITUTION.md` | Hard rules - NEVER violate                                |
| `.context/TASKS.md`        | Current work items                                        |
| `.context/DECISIONS.md`    | Architectural decisions with rationale                    |
| `.context/LEARNINGS.md`    | Gotchas, tips, lessons learned                            |
| `.context/CONVENTIONS.md`  | Code patterns and standards                               |

## Before Session Ends

**ALWAYS offer to persist context before the user quits:**

1. Add learnings: `ctx add learning "..."`
2. Add decisions: `ctx add decision "..."`

## Build Commands

```bash
CGO_ENABLED=0 go build -o ctx ./cmd/ctx    # Build CLI
CGO_ENABLED=0 go test ./...                   # Run tests
./hack/build-all.sh                        # Cross-platform build
```

## This Project

Context (`ctx`) is a CLI tool for persistent AI context. It was built using 
the Ralph Loop technique.

- **ctx** = context management tool (creates `.context/`)
- **Ralph Loop** = iterative AI development workflow (uses PROMPT.md)
- They are separate but complementary systems
