---
name: ctx-spec
description: "Scaffold a feature spec from the project template. Use when planning a new feature, writing a design document, or when a task references a missing spec."
---

Scaffold a new spec from `specs/tpl/spec-template.md` and walk through
each section with the user to produce a complete design document.

## When to Use

- Before implementing a non-trivial feature
- When a task says "Spec: `specs/X.md`" and the file does not exist
- When `/ctx-brainstorm` has produced a validated design that needs
  a written artifact
- When the user says "let's spec this out" or "write a spec for..."

## When NOT to Use

- Bug fixes or small changes (just do them)
- When a spec already exists (read it instead)
- When the design is still vague (use `/ctx-brainstorm` first)

## Usage Examples

```text
/ctx-spec
/ctx-spec (session checkpointing)
/ctx-spec (rss feed generation)
```

## Process

### 1. Gather the Feature Name

If not provided as an argument, ask:
> "What feature should this spec cover?"

Derive the filename: lowercase, hyphens, no spaces.
Target path: `specs/{feature-name}.md`

If the file already exists, warn and offer to review it instead.

### 2. Read the Template

Read `specs/tpl/spec-template.md` to get the current structure.

### 3. Walk Through Sections

Work through each section **one at a time**. For each section:

1. Explain what belongs there (one sentence)
2. Ask the user for input or propose content based on context
3. Write their answer into the section
4. Move to the next section

**Section order and prompts:**

| Section | Prompt |
|---------|--------|
| **Problem** | "What user-visible problem does this solve? Why now?" |
| **Approach** | "High-level: how does this work? Where does it fit?" |
| **Happy Path** | "Walk me through what happens when everything goes right." |
| **Edge Cases** | "What could go wrong? Think: empty input, partial failure, duplicates, concurrency, missing deps." |
| **Validation Rules** | "What input constraints are enforced? Where?" |
| **Error Handling** | "For each error condition: what message does the user see? How do they recover?" |
| **Interface** | "CLI command? Skill? Both? What flags?" |
| **Implementation** | "Which files change? Key functions? Existing helpers to reuse?" |
| **Configuration** | "Any .ctxrc keys, env vars, or settings?" |
| **Testing** | "Unit, integration, edge case tests?" |
| **Non-Goals** | "What does this intentionally NOT do?" |

**Spend extra time on Edge Cases and Error Handling.** These are
where specs earn their value. Push for at least 3 edge cases and
their expected behaviors. Do not accept "none" without challenge.

### 4. Open Questions

After all sections, ask:
> "Anything unresolved? If not, I'll remove the Open Questions
> section."

### 5. Write the Spec

Write the completed spec to `specs/{feature-name}.md`.

### 6. Cross-Reference

- If a Phase exists in TASKS.md referencing this spec, confirm
  the path matches
- If no tasks exist yet, offer to create them:
  > "Want me to break this into tasks in TASKS.md?"

## Skipping Sections

Not every spec needs every section. If a section clearly does not
apply (e.g., no CLI for an internal refactor), the user can say
"skip" and the section is omitted entirely — not left with
placeholder text.

## Quality Checklist

Before writing the file, verify:

- [ ] Problem section explains *why*, not just *what*
- [ ] At least 3 edge cases enumerated with expected behavior
- [ ] Error handling has user-facing messages and recovery steps
- [ ] Non-goals are explicit (prevents scope creep later)
- [ ] No placeholder `...` text remains
- [ ] Filename matches the convention: `specs/{feature-name}.md`
