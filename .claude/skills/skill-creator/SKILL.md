---
name: skill-creator
description: "Guide for creating or improving skills. Use when building a new skill or evaluating an existing one."
---

Create or improve skills that extend the agent with
project-specific knowledge the platform does not already
provide.

## When to Use

- When building a new skill from scratch
- When evaluating or improving an existing skill
- When a skill is not triggering or producing poor results
- When the user asks to create a slash command

## When NOT to Use

- For one-off instructions (just tell the agent directly)
- When the platform already handles it (do not restate
  system prompt behavior)
- When the task is too narrow for reuse (a skill should apply
  to more than one situation)

## Usage Examples

```text
/skill-creator
/skill-creator (improve the qa skill)
/skill-creator (create a new deploy skill)
```

## Core Principles

**The agent is already smart.** Only add context it does not
have. Tag each paragraph as **Expert** (agent does not know
this: keep), **Activation** (agent knows but a brief reminder
helps: keep if short), or **Redundant** (agent definitely
knows: delete). Target: >70% Expert, <10% Redundant.

**Description is the trigger.** The `description` field in
frontmatter determines when a skill activates. Put "when to
use" there, not in the body. Be specific:

```yaml
# Bad
description: "Use when starting any conversation"

# Good
description: "Use after writing code, before commits, or when CI might fail"
```

**Match freedom to fragility.** Narrow bridge with cliffs
needs guardrails (exact commands). Open field allows many
routes (principles and heuristics).

## Skill Anatomy

A well-structured skill has these sections:

| Section           | Purpose                                    |
|-------------------|--------------------------------------------|
| Frontmatter       | Name, description (trigger), allowed-tools |
| Before X-ing      | Pre-flight checks before execution         |
| When to Use       | Positive triggers                          |
| When NOT to Use   | Negative triggers (prevent misuse)         |
| Usage Examples    | Invocation patterns                        |
| Process/Execution | What to do (commands, steps)               |
| Good/Bad Examples | Show desired vs undesired output           |
| Quality Checklist | Verify before reporting completion         |

Not every skill needs all sections. Short skills (< 30 lines)
can skip the checklist. Complex skills should have all of them.

## Skill Structure

```
skill-name/
├── SKILL.md      (required: frontmatter + instructions)
├── scripts/      (optional: deterministic, reusable code)
├── references/   (optional: loaded into context on demand)
└── assets/       (optional: used in output, not loaded)
```

Keep SKILL.md under 100 lines. Move detailed reference
material to separate files and describe when to read them.

## What NOT to Include

- Guidance the system prompt already provides
- Personality roleplay ("You are an expert...")
- Capability lists ("Masters X, Y, Z...")
- README, CHANGELOG, or auxiliary documentation
- References to files that do not exist
- Great content with a vague description: it will never
  trigger

## Workflow Patterns

**Sequential**: overview the steps upfront:

```markdown
1. Check formatting (`gofmt -l .`)
2. Run linter (`golangci-lint run`)
3. Run tests (`go test ./...`)
4. Report results
```

**Conditional**: guide through decision points:

```markdown
1. Determine the change type:
   **New feature?** -> Run full audit
   **Docs only?** -> Skip checks
```

## Output Patterns

**Template**: provide structure when consistency matters:

```markdown
After running checks, report:
1. **Result**: Pass or fail
2. **Failures**: What failed and how to fix
3. **Files touched**: List of modified files
```

**Examples**: show input/output pairs when style matters:

```markdown
Good:  ran `make audit` -> "All checks pass"
Bad:   "Should pass now" (without running anything)
```

Examples communicate style more clearly than descriptions.
Good/bad pairs are especially effective; they set boundaries
without being prescriptive.

## Litmus Test

Before finalizing a skill:

1. Does the platform already do this? Do not restate it.
2. Does it suppress AI judgment? It is a jailbreak, not a
   skill.
3. Do all referenced files exist? Fix or remove phantom
   references.
4. Is it under 100 lines? If not, split into references.
5. Is the description specific? Narrow the trigger.
6. Would an expert say "this captures knowledge that took
   years to learn"? If not, it is probably Redundant.
7. Does it have examples? Without them, execution quality
   degrades over time.

## Process

1. **Understand the need**: what problem does this skill
   solve? Is it recurring enough to justify a skill?
2. **Check for existing skills**: search `.claude/skills/`
   to avoid duplicates or find one to extend
3. **Draft the skill**: start with frontmatter and core
   content; tag paragraphs (Expert/Activation/Redundant)
4. **Add examples**: good/bad pairs for output quality
5. **Run the litmus test**: verify all 7 points above
6. **If template-worthy**: also create the version in
   `internal/tpl/claude/skills/` so it deploys with
   `ctx init`
