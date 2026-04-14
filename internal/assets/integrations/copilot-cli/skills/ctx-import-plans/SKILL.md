---
name: ctx-import-plans
description: "Import plan files into project specs directory. Use to convert external plans into project-tracked specs."
tools: [bash, read, write]
---

Import plan files into the project's `specs/` directory.

## When to Use

- When plan files exist outside the project (e.g., from AI
  tool plan modes)
- When converting external design docs to project specs
- When the user says "import that plan"

## When NOT to Use

- Plan is already in `specs/`
- Plan is too vague to be a spec (brainstorm first)

## Process

### 1. Locate the plan

If path provided, read it. Otherwise, check common locations:
- Current conversation context
- Session workspace files

### 2. Convert to spec format

Map plan sections to the spec template structure:
- Problem → Problem
- Steps/Tasks → Implementation
- Goals → Happy Path
- Risks → Edge Cases

### 3. Handle conflicts

If `specs/{name}.md` already exists:
- Compare contents
- Offer to merge, replace, or rename

### 4. Write the spec

Write to `specs/{name}.md`.

### 5. Create tasks (optional)

Offer to break the spec into tasks in TASKS.md.

## Quality Checklist

- [ ] Spec follows project template structure
- [ ] No conflicts with existing specs
- [ ] File written to correct location
- [ ] Tasks offered if applicable
