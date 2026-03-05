---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: My AI Keeps Making the Same Mistakes
icon: lucide/repeat
---

![ctx](../images/ctx-banner.png)

## The Problem

You found a bug last Tuesday. You debugged it, understood the root cause,
and moved on. Today, a new session hits the exact same bug. The AI
rediscovers it from scratch, burning twenty minutes on something you
already solved.

Worse: you spent an hour last week evaluating two database migration
strategies, picked one, documented why in a comment somewhere, and now
the AI is cheerfully suggesting the approach you rejected. Again.

This is not a model problem. It is a **memory** problem. Without
persistent context, every session starts with amnesia.

## How `ctx` Stops the Loop

`ctx` gives your AI three files that directly prevent repeated mistakes,
each targeting a different failure mode.

### DECISIONS.md: Stop Relitigating Settled Choices

When you make an architectural decision, record it with rationale and
rejected alternatives. The AI reads this at session start and treats
it as settled.

```markdown
## [2026-02-12] Use JWT for Authentication

**Status**: Accepted

**Context**: Need stateless auth for the API layer.

**Decision**: JWT with short-lived access tokens and refresh rotation.

**Rationale**: Stateless, scales horizontally, team has prior experience.

**Alternatives Considered**:
- Session-based auth: Rejected. Requires sticky sessions or shared store.
- API keys only: Rejected. No user identity, no expiry rotation.
```

Next session, when the AI considers auth, it reads this entry and builds
on the decision instead of re-debating it. If someone asks "why not
sessions?", the rationale is already there.

### LEARNINGS.md: Capture Gotchas Once

Learnings are the bugs, quirks, and non-obvious behaviors that cost you
time the first time around. Write them down so they cost you zero time
the second time.

```markdown
## Build

### CGO Required for SQLite on Alpine

**Discovered**: 2026-01-20

**Context**: Docker build failed silently with "no such table" at runtime.

**Lesson**: The go-sqlite3 driver requires CGO_ENABLED=1 and gcc
installed in the build stage. Alpine needs apk add build-base.

**Application**: Always use the golang:alpine image with build-base
for SQLite builds. Never set CGO_ENABLED=0.
```

Without this entry, the next session that touches the Dockerfile will
hit the same wall. With it, the AI knows before it starts.

### CONSTITUTION.md: Draw Hard Lines

Some mistakes are not about forgetting -- they are about boundaries the
AI should never cross. CONSTITUTION.md sets inviolable rules.

```markdown
* [ ] Never commit secrets, tokens, API keys, or credentials
* [ ] Never disable security linters without a documented exception
* [ ] All database migrations must be reversible
```

The AI reads these as absolute constraints. It does not weigh them
against convenience. It refuses tasks that would violate them.

## The Accumulation Effect

Each of these files grows over time. Session one captures two decisions.
Session five adds a tricky learning about timezone handling. Session
twelve records a convention about error message formatting.

By session twenty, your AI has a knowledge base that no single person
carries in their head. New team members -- human or AI -- inherit it
instantly.

The key insight: **you are not just coding. You are building a knowledge
layer that makes every future session faster.**

`ctx` files version with your code in git. They survive branch switches,
team changes, and model upgrades. The context outlives any single session.

## Getting Started

Capture your first decision or learning right now:

```bash
ctx add decision "Use PostgreSQL" \
  --context "Need a relational database for the project" \
  --rationale "Team expertise, JSONB support, mature ecosystem"

ctx add learning "Vitest mock hoisting" \
  --context "Tests failing intermittently" \
  --lesson "vi.mock() must be at file top level" \
  --application "Use vi.doMock() for dynamic mocks"
```

## Further Reading

* [Knowledge Capture](../recipes/knowledge-capture.md): the full workflow
  for persisting decisions, learnings, and conventions
* [Context Files Reference](context-files.md): structure and format for
  every file in `.context/`
* [About ctx](about.md): the bigger picture -- why persistent context
  changes how you work with AI
