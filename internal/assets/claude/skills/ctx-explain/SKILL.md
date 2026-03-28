---
name: ctx-explain
description: "Explain code for someone new to the project. Use when the user asks 'what does this do', 'explain this', or wants to understand unfamiliar code."
allowed-tools: Read, Grep, Glob
---

Explain the specified code for someone new to the project. Tailor
depth to the user's expertise if known from context.

## When to Use

- User says "explain this code", "explain this", "what does this do"
- User is onboarding to an unfamiliar area of the codebase
- User says "walk me through this" or "how does this work"

## When NOT to Use

- User wants deep architectural analysis (use `/ctx-architecture`)
- User wants to trace a bug (use `/gitnexus-debugging`)
- User wants execution flow tracing (use `/gitnexus-exploring`)

## Explanation Structure

Cover each dimension in order. Skip any that don't apply.

1. **What it does**: Describe the purpose and behavior in plain
   language.
2. **Why it exists**: What problem does it solve? What would break
   without it?
3. **How it connects**: Which modules call it, and which modules
   does it depend on?
4. **Key design decisions**: Why was this approach chosen over
   alternatives?
5. **Non-obvious details**: Anything surprising, subtle, or easy
   to misunderstand.

## Execution

1. Read the target code
2. Read `.context/ARCHITECTURE.md` for system-level context
3. Trace callers and callees if connections matter
4. Present the explanation following the structure above

Keep it concise. Lead with the "what": the reader wants to orient
before diving into "why" and "how."
