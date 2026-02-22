---
#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Context Files
icon: lucide/files
---

![ctx](images/ctx-banner.png)

## `.context/`

Each context file in `.context/` serves a specific purpose. 

Files are designed to be human-readable, AI-parseable, and token-efficient.

## File Overview

| File                | Purpose                                    | Priority    |
|---------------------|--------------------------------------------|-------------|
| `CONSTITUTION.md`   | Hard rules that must **NEVER** be violated | 1 (highest) |
| `TASKS.md`          | Current and planned work                   | 2           |
| `CONVENTIONS.md`    | Project patterns and standards             | 3           |
| `ARCHITECTURE.md`   | System overview and components             | 4           |
| `DECISIONS.md`      | Architectural decisions with rationale     | 5           |
| `LEARNINGS.md`      | Lessons learned, gotchas, tips             | 6           |
| `GLOSSARY.md`       | Domain terms and abbreviations             | 7           |
| `AGENT_PLAYBOOK.md` | Instructions for AI tools                  | 8 (lowest)  |

## Read Order Rationale

The priority order follows a logical progression for AI tools:

1. `CONSTITUTION.md`: Inviolable rules first. The AI tool must know what it
   *cannot* do before attempting anything.
2. `TASKS.md`: Current work items. What the AI tool should focus on.
3. `CONVENTIONS.md`: How to write code. Patterns and standards to follow
   when implementing tasks.
4. `ARCHITECTURE.md`: System structure. Understanding of components and
   boundaries before making changes.
5. `DECISIONS.md`: Historical context. Why things are the way they are,
   to avoid re-debating settled decisions.
6. `LEARNINGS.md`: Gotchas and tips. Lessons from past work that inform
   the current implementation.
7. `GLOSSARY.md`: Reference material. Domain terms and abbreviations for
   lookup as needed.
8. `AGENT_PLAYBOOK.md`: Meta instructions last. How to use this context
   system itself. Loaded last because the agent should understand the
   *content* (*rules, tasks, patterns*) before the *operating manual*.

---

## `CONSTITUTION.md`

**Purpose:** Define hard invariants—rules that must **NEVER** be violated, 
regardless of the task.

AI tools read this first and should refuse tasks that violate these rules.

### Structure

```markdown
# Constitution

These rules are INVIOLABLE. If a task requires violating these, the task 
is wrong.

## Security Invariants

- [ ] Never commit secrets, tokens, API keys, or credentials
- [ ] Never store customer/user data in context files
- [ ] Never disable security linters without documented exception

## Quality Invariants

- [ ] All code must pass tests before commit
- [ ] No `any` types in TypeScript without documented reason
- [ ] No TODO comments in main branch (move to TASKS.md)

## Process Invariants

- [ ] All architectural changes require a decision record
- [ ] Breaking changes require version bump
- [ ] Generated files are never committed
```

### Guidelines

* Keep rules minimal and absolute
* Each rule should be enforceable (can verify compliance)
* Use checkbox format for clarity
* Never compromise on these rules

---

## `TASKS.md`

**Purpose:** Track current work, planned work, and blockers.

### Structure

Tasks are organized by **Phase**: logical groupings that preserve order and
enable replay. 

Tasks stay in their Phase permanently; status is tracked via checkboxes and 
inline tags.

```markdown
# Tasks

## Phase 1: Initial Setup

- [x] Set up project structure
- [x] Configure linting and formatting
- [ ] Add CI/CD pipeline `#in-progress`

## Phase 2: Core Features

- [ ] Implement user authentication `#priority:high`
- [ ] Add API rate limiting `#priority:medium`
  - Blocked by: Need to finalize auth first

## Backlog

- [ ] Performance optimization `#priority:low`
- [ ] Add metrics dashboard `#priority:deferred`
```

**Key principles:**

* Tasks never move between sections — mark as `[x]` or `[-]` in place
* Use `#in-progress` inline tag to indicate current work
* Phase headers provide structure and replay order
* Backlog section for unscheduled work

### Tags

Use inline backtick-wrapped tags for metadata:

| Tag            | Values                         | Purpose                   |
|----------------|--------------------------------|---------------------------|
| `#priority`    | `high`, `medium`, `low`        | Task urgency              |
| `#area`        | `core`, `cli`, `docs`, `tests` | Codebase area             |
| `#estimate`    | `1h`, `4h`, `1d`               | Time estimate (optional)  |
| `#in-progress` | (none)                         | Currently being worked on |

**Lifecycle tags** (for session correlation):

| Tag        | Format               | When to add                        |
|------------|----------------------|------------------------------------|
| `#added`   | `YYYY-MM-DD-HHMMSS`  | Auto-added by `ctx add task`       |
| `#started` | `YYYY-MM-DD-HHMMSS`  | When beginning work on the task    |
| `#done`    | `YYYY-MM-DD-HHMMSS`  | When marking the task `[x]`        |

These timestamps help correlate tasks with session files and track which
session started vs completed work.

### Status Markers

| Marker | Meaning                  |
|--------|--------------------------|
| `[ ]`  | Pending                  |
| `[x]`  | Completed                |
| `[-]`  | Skipped (include reason) |

### Guidelines

* Never delete tasks; mark as `[x]` completed or `[-]` skipped
* Never move tasks between sections; use inline tags for status
* Use `ctx tasks archive` periodically to move completed tasks to archive
* Mark current work with `#in-progress` inline tag

---

## `DECISIONS.md`

**Purpose:** Record architectural decisions with rationale so they don't
get re-debated.

### Structure

```markdown
# Decisions

## [YYYY-MM-DD] Decision Title

**Status**: Accepted | Superseded | Deprecated

**Context**: What situation prompted this decision?

**Decision**: What was decided?

**Rationale**: Why was this the right choice?

**Consequences**: What are the implications?

**Alternatives Considered**:
- Alternative A: Why rejected
- Alternative B: Why rejected
```

### Example

```markdown
## [2025-01-15] Use TypeScript Strict Mode

**Status**: Accepted

**Context**: Starting a new project, need to choose the type-checking level.

**Decision**: Enable TypeScript strict mode with all strict flags.

**Rationale**: Catches more bugs at compile time. Team has experience
with strict mode. Upfront cost pays off in reduced runtime errors.

**Consequences**: More verbose type annotations required. Some
third-party libraries need type assertions.

**Alternatives Considered**:
- Basic TypeScript: Rejected because it misses null checks
- JavaScript with JSDoc: Rejected because tooling support is weaker
```

### Status Values

| Status     | Meaning                                 |
|------------|-----------------------------------------|
| Accepted   | Current, active decision                |
| Superseded | Replaced by newer decision (link to it) |
| Deprecated | No longer relevant                      |

---

## `LEARNINGS.md`

**Purpose:** Capture lessons learned, gotchas, and tips that shouldn't
be forgotten.

### Structure

```markdown
# Learnings

## Category Name

### Learning Title

**Discovered**: YYYY-MM-DD

**Context**: When/how was this learned?

**Lesson**: What's the takeaway?

**Application**: How should this inform future work?
```

### Example

```markdown
## Testing

### Vitest Mocks Must Be Hoisted

**Discovered**: 2025-01-15

**Context**: Tests were failing intermittently when mocking fs module.

**Lesson**: Vitest requires `vi.mock()` calls to be hoisted to the
top of the file. Dynamic mocks need `vi.doMock()` instead.

**Application**: Always use `vi.mock()` at file top. Use `vi.doMock()`
only when mock needs runtime values.
```

### Categories

Organize learnings by topic:

* Testing
* Build & Deploy
* Performance
* Security
* Third-Party Libraries
* Git & Workflow

---

## `CONVENTIONS.md`

**Purpose**: Document project patterns, naming conventions, and standards.

### Structure

```markdown
# Conventions

## Naming

* **Files**: kebab-case for all source files
* **Components**: PascalCase for React components
* **Functions**: camelCase, verb-first (getUser, parseConfig)
* **Constants**: SCREAMING_SNAKE_CASE

## Patterns

### Pattern Name

**When to use**: Situation description

**Implementation**:
// in triple backticks
// Example code

**Why**: Rationale for this pattern
```

### Guidelines

* Include concrete examples
* Explain the "why" not just the "what"
* Keep patterns minimal—only document what's non-obvious

---

## `ARCHITECTURE.md`

**Purpose**: Provide system overview and component relationships.

### Structure

```markdown
# Architecture

## Overview

Brief description of what the system does and how it's organized.

## Components

### Component Name

**Responsibility**: What this component does

**Dependencies**: What it depends on

**Dependents**: What depends on it

**Key Files**:
- path/to/file.ts — Description

## Data Flow

Description or diagram of how data moves through the system.

## Boundaries

What's in scope vs out of scope for this codebase.
```

### Guidelines

* Keep diagrams simple (Mermaid works well)
* Focus on boundaries and interfaces
* Update when major structural changes occur

---

## `GLOSSARY.md`

**Purpose**: Define domain terms, abbreviations, and project vocabulary.

### Structure

```markdown
# Glossary

## Domain Terms

### Term Name

**Definition**: What it means in this project's context

**Not to be confused with**: Similar terms that mean different things

**Example**: How it's used

## Abbreviations

| Abbrev | Expansion                     | Context                |
|--------|-------------------------------|------------------------|
| ADR    | Architectural Decision Record | Decision documentation |
| SUT    | System Under Test             | Testing                |
```

### Guidelines

* Define project-specific meanings
* Clarify potentially ambiguous terms
* Include abbreviations used in code or docs

---

## `AGENT_PLAYBOOK.md`

**Purpose**: Explicit instructions for how AI tools should read, apply, 
and update context.

### Key Sections

**Read Order**: Priority order for loading context files

**When to Update**: Events that trigger context updates

**How to Avoid Hallucinating Memory:** Critical rules:

1. Never assume—if not in files, you don't know it
2. Never invent history—don't claim "we discussed" without evidence
3. Verify before referencing—search files before citing
4. When uncertain, say so
5. Trust files over intuition

**Context Update Commands:** Format for automated updates via `ctx watch`:

```xml
<context-update type="task">Implement rate limiting</context-update>
<context-update type="complete">user auth</context-update>
<context-update type="learning"
  context="Debugging hooks"
  lesson="Hooks receive JSON via stdin"
  application="Parse JSON stdin with the host language"
>Hook Input Format</context-update>
<context-update type="decision"
  context="Need a caching layer"
  rationale="Redis is fast and team has experience"
  consequences="Must provision Redis infrastructure"
>Use Redis for caching</context-update>
```

See [Integrations](integrations.md#context-update-commands) for full documentation.

---

## Parsing Rules

All context files follow these conventions:

1. **Headers define structure**: `#` for title, `##` for sections, `###` for 
   items
2. **Bold keys for fields**: `**Key**:` followed by value
3. **Code blocks are literal**: Never parse code block content as structure
4. **Lists are ordered**: Items appear in priority/chronological order
5. **Tags are inline**: Backtick-wrapped tags like `#priority:high`

## Further Reading

* [Refactoring with Intent](blog/2026-02-01-refactoring-with-intent.md):
  How persistent context prevents drift during refactoring sessions.

## Token Efficiency

Keep context files **concise**:

* Use abbreviations in tags, not prose;
* Omit obvious words ("*The*," "*This*");
* Prefer bullet points over paragraphs;
* Keep examples minimal but illustrative;
* Archive old completed items periodically.

----

**Next Up**: [Prompting Guide →](prompting-guide.md) — effective prompts for AI sessions with ctx.
