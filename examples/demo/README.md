![ctx](../../assets/ctx-banner.png)

# Demo Project

This is a sample project demonstrating Context (ctx) structure and best practices.

## Quick Start

```bash
# View context status
ctx status

# Get AI-ready context packet
ctx agent

# Add a new task
ctx add task "Implement feature X"

# Mark a task complete
ctx complete "feature X"

# Check for stale context
ctx drift
```

## Context Files

The `.context/` directory contains markdown files that provide persistent
context for AI coding assistants:

| File                 | Purpose                                           |
|----------------------|---------------------------------------------------|
| `AGENT_PLAYBOOK.md`  | **Read first** — How agents should use this system |
| `CONSTITUTION.md`    | Inviolable rules — NEVER violate these            |
| `TASKS.md`           | Current work items with phases and timestamps     |
| `CONVENTIONS.md`     | Coding standards and patterns                     |
| `ARCHITECTURE.md`    | System overview and component layout              |
| `DECISIONS.md`       | Technical decisions with rationale                |
| `LEARNINGS.md`       | Gotchas, tips, lessons learned                    |

## Key Concepts

### Agent Playbook

`AGENT_PLAYBOOK.md` is the bootstrap file for AI agents. It explains:
- The mental model (memory = files, not conversation)
- Read order for context files
- When and how to persist learnings/decisions
- How to avoid hallucinating memory

### Phase-Based Tasks

Tasks in `TASKS.md` stay in their phase permanently. Use inline labels
(`#in-progress`) instead of moving tasks between sections:

```markdown
## Phase 2: Authentication

- [x] Implement user registration #added:2026-01-04-100000 #done:2026-01-05-140000
- [ ] Implement OAuth2 login #added:2026-01-04-100000 #in-progress
- [ ] Add session management #added:2026-01-04-100000
```

### Structured Entries

Learnings and decisions follow structured formats with timestamps:

```markdown
## [2026-01-15-143022] Database connections need explicit timeouts

**Context**: What situation led to this learning

**Lesson**: What we learned

**Application**: How to apply it going forward
```

## Adding Context

```bash
# Add a learning with full structure
ctx add learning "Title" \
  --context "What happened" \
  --lesson "What we learned" \
  --application "How to apply it"

# Add a decision with rationale
ctx add decision "Title" \
  --context "What prompted this" \
  --rationale "Why this choice" \
  --consequences "What changes"

# Add a task
ctx add task "Implement feature X"
```

## Ralph Loop Integration

This demo includes Ralph Loop infrastructure for iterative AI development:

| File | Purpose |
|------|---------|
| `PROMPT.md` | Directive for AI agents — defines the work loop |
| `specs/` | Detailed specifications for features |

The Ralph Loop pattern:
1. AI reads `PROMPT.md` to understand the workflow
2. Picks ONE task from `.context/TASKS.md`
3. Reads relevant spec from `specs/` for requirements
4. Implements the task
5. Updates context files
6. Exits — the loop restarts with fresh context

This is separate from but complementary to ctx:
- **ctx** = context persistence (`.context/`)
- **Ralph Loop** = iterative AI workflow (`PROMPT.md` + `specs/`)

## Session Persistence

Session dumps are saved to `.context/sessions/` with timestamps:
- `2026-01-20-164600-feature-discussion.md` — Manual session notes
- Auto-saved transcripts (if Claude Code hooks are configured)

This allows future sessions to understand past context without relying on
conversation memory.
