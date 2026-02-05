---
name: skill-creator
description: "Guide for creating or improving skills. Use when building a new skill or evaluating an existing one."
---

Create skills that extend the agent with project-specific knowledge the
platform doesn't already provide.

## Core Principles

**Claude is already smart.** Only add context it doesn't have. Tag each
paragraph as **Expert** (Claude doesn't know this — keep), **Activation**
(Claude knows but a brief reminder helps — keep if short), or **Redundant**
(Claude definitely knows — delete). Target: >70% Expert, <10% Redundant.

**Description is the trigger.** The `description` field in frontmatter is
what determines when a skill activates. Put "when to use" there, not in
the body. Be specific:

```yaml
# Bad
description: "Use when starting any conversation"

# Good
description: "Use after writing code, before commits, or when CI might fail"
```

**Match freedom to fragility.** Narrow bridge with cliffs needs guardrails
(exact commands). Open field allows many routes (principles and heuristics).

## Skill Structure

```
skill-name/
├── SKILL.md          (required — frontmatter + instructions)
├── scripts/          (optional — deterministic, reusable code)
├── references/       (optional — loaded into context on demand)
└── assets/           (optional — used in output, not loaded into context)
```

Keep SKILL.md under 100 lines. Move detailed reference material to
separate files and describe when to read them.

## What NOT to Include

- Guidance the system prompt already provides (see "Skills That Fight
  the Platform" blog post)
- Personality roleplay ("You are an expert...")
- Capability lists ("Masters X, Y, Z...")
- README, CHANGELOG, or auxiliary documentation
- References to files that don't exist
- Checkbox procedures (Step 1, Step 2...) — teach thinking, not steps
- Great content with a vague description — it'll never trigger

## Workflow Patterns

**Sequential** — overview the steps upfront:

```markdown
1. Check formatting (`gofmt -l .`)
2. Run linter (`golangci-lint run`)
3. Run tests (`go test ./...`)
4. Report results
```

**Conditional** — guide through decision points:

```markdown
1. Determine the change type:
   **New feature?** → Run full audit
   **Docs only?** → Skip checks
```

## Output Patterns

**Template** — provide structure when consistency matters:

```markdown
After running checks, report:
1. **Result**: Pass or fail
2. **Failures**: What failed and how to fix
3. **Files touched**: List of modified files
```

**Examples** — show input/output pairs when style matters:

```markdown
Good:  ran `make audit` → "All checks pass"
Avoid: "Should pass now" (without running anything)
```

Examples communicate style more clearly than descriptions.

## Litmus Test

Before finalizing a skill:

1. Does the platform already do this? → Don't restate it
2. Does it suppress AI judgment? → It's a jailbreak, not a skill
3. Do all referenced files exist? → Fix or remove phantom references
4. Is it under 100 lines? → If not, split into references
5. Is the description specific? → Narrow the trigger
6. Would an expert say "this captures knowledge that took me years to learn"? → If not, it's probably Redundant
