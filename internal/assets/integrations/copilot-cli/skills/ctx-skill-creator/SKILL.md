---
name: ctx-skill-creator
description: "Create, improve, test, and deploy skills. Full skill lifecycle from intent to working skill file."
tools: [bash, read, write, edit, glob, grep]
---

Create new skills or improve existing ones through a structured
workflow.

## When to Use

- Creating a new skill from scratch
- Improving an underperforming skill
- Porting a skill from one integration to another

## When NOT to Use

- Quick one-off automations (just script it)
- When the need is too vague (brainstorm first)

## Process

### 1. Intent capture

Gather:
- What should this skill do?
- When should it trigger?
- What tools does it need?
- What's the expected output?

### 2. Draft the SKILL.md

Use the standard structure:

```yaml
---
name: ctx-{name}
description: "..."
tools: [bash, read, write, ...]
---
```

Sections: When to Use, When NOT to Use, Process, Quality Checklist.

### 3. Validate

Check against skill audit dimensions:
- Positive framing
- Clear scope
- Good examples
- No phantom references
- Overtriggering guard

### 4. Test

If possible, do a dry run of the skill's workflow to verify
it works end-to-end.

### 5. Deploy

Write the file to the appropriate skills directory:
- Claude: `internal/assets/claude/skills/{name}/SKILL.md`
- Copilot CLI: `internal/assets/integrations/copilot-cli/skills/{name}/SKILL.md`

### 6. Build

Run `go build ./cmd/ctx/...` to verify the embed compiles.

## Quality Checklist

- [ ] Frontmatter is complete (name, description, tools)
- [ ] When to Use / When NOT to Use sections exist
- [ ] Process has numbered, actionable steps
- [ ] Quality Checklist at the end
- [ ] No phantom references
- [ ] Build passes with new skill embedded
