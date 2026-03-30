---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Keeping AI Honest
icon: lucide/shield-check
---

![ctx](../images/ctx-banner.png)

## The Problem

AI agents confabulate. They invent history that never happened, claim
familiarity with decisions that were never made, and sometimes declare
a task complete when it is not. This is not malice - it is the default
behavior of a system optimizing for plausible-sounding responses.

When your AI says "we decided to use Redis for caching last week," can
you verify that? When it says "the auth module is complete," can you
confirm it? Without grounded, persistent context, the answer is no.
You are trusting vibes.

`ctx` replaces vibes with verifiable artifacts.

## Grounded Memory

Every entry in `ctx` context files has a timestamp and structured
fields. When the AI cites a decision, you can check it.

```markdown
## [2026-01-28-143022] Use Event Sourcing for Audit Trail

**Status**: Accepted

**Context**: Compliance requires full mutation history.

**Decision**: Event sourcing for the audit subsystem only.

**Rationale**: Append-only log meets compliance requirements
without imposing event sourcing on the entire domain model.
```

The timestamp `2026-01-28-143022` is not decoration. It is a verifiable
anchor. If the AI references this decision, you can open DECISIONS.md,
find the entry, and confirm it says what the AI claims. If the entry
does not exist, the AI is hallucinating - and you know immediately.

This is **grounded memory**: claims that trace back to artifacts you
control and can audit.

## CONSTITUTION.md: Hard Guardrails

CONSTITUTION.md defines rules the AI must treat as inviolable. These
are not suggestions or best practices - they are constraints that
override task requirements.

```markdown
# Constitution

These rules are INVIOLABLE. If a task requires violating these,
the task is wrong.

* [ ] Never commit secrets, tokens, API keys, or credentials
* [ ] All public API changes require a decision record
* [ ] Never delete context files without explicit user approval
```

The AI reads these at session start, before anything else. A well-
integrated agent will refuse a task that conflicts with a constitutional
rule, citing the specific rule it would violate.

## The Agent Playbook's Anti-Hallucination Rules

The AGENT_PLAYBOOK.md file includes a section called **"How to Avoid
Hallucinating Memory"** with five explicit rules:

1. **Never assume.** If it is not in the context files, you do not
   know it.
2. **Never invent history.** Do not claim "we discussed" something
   without a file reference.
3. **Verify before referencing.** Search files before citing them.
4. **When uncertain, say so.** "I don't see a decision on this" is
   always better than a fabricated one.
5. **Trust files over intuition.** If the files say PostgreSQL but
   your training data suggests MySQL, the files win.

These rules create a behavioral contract. The AI is not left to guess
how confident it should be - it has explicit instructions to ground
every claim in the context directory.

## Drift Detection

Context files can go stale. You rename a package, delete a module, or
finish a sprint, and suddenly ARCHITECTURE.md references paths that
no longer exist. Stale context is almost as dangerous as no context:
the AI treats outdated information as current truth.

`ctx drift` detects this divergence:

```bash
ctx drift
```

It scans context files for references to files, paths, and symbols
that no longer exist in the codebase. Stale references get flagged
so you can update or remove them before they mislead the next session.

Regular drift checks - weekly, or after major refactors - keep your
context files honest the same way tests keep your code honest.

## The Verification Loop

The `/ctx-commit` skill includes a built-in verification step: before
staging, it maps claims to evidence and runs self-audit questions to
surface gaps. This catches inconsistencies at the point where they
matter most — right before code is committed.

This closes the loop. You write context. The AI reads context. The
verification step confirms that context still matches reality. When
it does not, you fix it - and the next session starts from truth,
not from drift.

## Trust Through Structure

The common thread across all of these mechanisms is **structure over
prose**. Timestamps make claims verifiable. Constitutional rules make
boundaries explicit. Drift detection makes staleness visible. The
playbook makes behavioral expectations concrete.

You do not need to trust the AI. You need to trust the system --
and verify when it matters.

## Further Reading

* [Detecting and Fixing Drift](../recipes/context-health.md): the
  full workflow for keeping context files accurate
* [Invariants](../reference/design-invariants.md): the properties
  that must hold for any valid `ctx` implementation
* [Agent Security](../security/agent-security.md): threat model and
  mitigations for AI agents operating with persistent context
