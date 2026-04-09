---
name: ctx-spec
description: "Scaffold a feature spec from the project template. Use when planning a new feature or when a task references a missing spec."
tools: [bash, read, write]
---

Scaffold a new spec from `specs/tpl/spec-template.md` and walk
through each section with the user.

## When to Use

- Before implementing a non-trivial feature
- When a task says "Spec: `specs/X.md`" and the file doesn't exist
- When `ctx-brainstorm` produced a validated design that needs
  a written artifact
- When the user says "let's spec this out"

## When NOT to Use

- Bug fixes or small changes
- When a spec already exists (read it instead)
- When the design is still vague (use `ctx-brainstorm` first)

## Process

### 1. Gather the Feature Name

If not provided, ask. Derive filename: lowercase, hyphens.
Target: `specs/{feature-name}.md`

### 2. Read the Template

Read `specs/tpl/spec-template.md`.

### 3. Walk Through Sections

Work through each section **one at a time**:

| Section              | Prompt                                                        |
|----------------------|---------------------------------------------------------------|
| **Problem**          | "What user-visible problem does this solve? Why now?"         |
| **Approach**         | "How does this work? Where does it fit?"                      |
| **Happy Path**       | "Walk through what happens when everything goes right."       |
| **Edge Cases**       | "What could go wrong? (empty input, failures, duplicates)"    |
| **Validation Rules** | "What input constraints are enforced?"                        |
| **Error Handling**   | "For each error: user message and recovery?"                  |
| **Interface**        | "CLI command? Skill? Both? Flags?"                            |
| **Implementation**   | "Which files change? Key functions? Helpers to reuse?"        |
| **Configuration**    | "Any .ctxrc keys, env vars, or settings?"                     |
| **Testing**          | "Unit, integration, edge case tests?"                         |
| **Non-Goals**        | "What does this intentionally NOT do?"                        |

**Spend extra time on Edge Cases and Error Handling.**

### 4. Open Questions

After all sections:
> "Anything unresolved? If not, I'll remove the Open Questions
> section."

### 5. Write the Spec

Write to `specs/{feature-name}.md`.

### 6. Cross-Reference

- If a Phase exists in TASKS.md, confirm the path matches
- If no tasks exist, offer to create them

## Quality Checklist

- [ ] Problem section explains *why*, not just *what*
- [ ] At least 3 edge cases with expected behavior
- [ ] Error handling has user messages and recovery
- [ ] Non-goals are explicit
- [ ] No placeholder text remains
