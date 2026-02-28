---
name: _ctx-skill-creator
description: "Create, improve, and test skills. Use when building a new skill, evaluating or fixing an existing one, turning a conversation workflow into a reusable skill, or when a skill is not triggering or producing poor results. Also use when the user says 'turn this into a skill', 'make a slash command for this', or 'this should be a skill'."
---

Create or improve skills that extend the agent with
project-specific knowledge the platform does not already provide.

## When to Use

- Building a new skill from scratch
- Evaluating or improving an existing skill
- A skill is not triggering or producing poor results
- The user asks to create a slash command
- The conversation contains a workflow the user wants to capture
  ("turn this into a skill", "can we automate this?")

## When NOT to Use

- For one-off instructions (just tell the agent directly)
- When the platform already handles it (do not restate system
  prompt behavior)
- When the task is too narrow for reuse (a skill should apply
  to more than one situation)

## Core Principles

**The agent is already smart.** Only add context it does not
have. Tag each paragraph as **Expert** (agent does not know
this: keep), **Activation** (agent knows but a brief reminder
helps: keep if short), or **Redundant** (agent definitely
knows: delete). Target: >70% Expert, <10% Redundant.

**Description is the trigger.** The `description` field in
frontmatter determines when a skill activates. All "when to
use" context belongs there, not in the body. Claude tends to
undertrigger — make descriptions a little "pushy" by listing
concrete situations and synonyms the user might say:

```yaml
# Weak — too vague, will undertrigger
description: "Use when starting any conversation"

# Strong — specific triggers, covers synonyms
description: >-
  Use after writing code, before commits, or when CI might
  fail. Also use when the user says 'run checks', 'lint this',
  or 'does this pass tests?'
```

**Match freedom to fragility.** Narrow bridge with cliffs
needs guardrails (exact commands). Open field allows many
routes (principles and heuristics).

**Explain the why, not heavy-handed MUSTs.** Today's LLMs are
smart. They have good theory of mind and respond better to
reasoning than rigid directives. If you find yourself writing
ALWAYS or NEVER in all caps, reframe: explain *why* the thing
matters so the model understands and generalizes, rather than
memorizing a rule it might misapply.

**Match communication to user skill level.** Pay attention to
context cues. Some users are experienced developers; others are
new to terminals entirely. Briefly explain terms ("assertions —
automated checks that verify the output") when in doubt.

## Skill Anatomy

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

## Skill Structure and Progressive Disclosure

```
skill-name/
├── SKILL.md      (required: frontmatter + instructions)
├── scripts/      (optional: deterministic, reusable code)
├── references/   (optional: loaded into context on demand)
└── assets/       (optional: used in output, not loaded)
```

Skills use a three-level loading system:

1. **Metadata** (name + description) — always in context (~100 words)
2. **SKILL.md body** — loaded when skill triggers (<500 lines)
3. **Bundled resources** — loaded as needed (unlimited size;
   scripts can execute without loading into context)

Keep SKILL.md under 500 lines. If approaching that limit, move
detail to `references/` with clear pointers about when to read
each file. For large reference files (>300 lines), include a
table of contents.

For multi-domain skills, organize by variant so the agent reads
only the relevant reference:

```
cloud-deploy/
├── SKILL.md (workflow + selection logic)
└── references/
    ├── aws.md
    ├── gcp.md
    └── azure.md
```

## Writing Guide

### Output Formats

Provide structure when consistency matters:

```markdown
After running checks, report:
1. **Result**: Pass or fail
2. **Failures**: What failed and how to fix
3. **Files touched**: List of modified files
```

### Examples

Show input/output pairs when style matters. Good/bad pairs set
boundaries without being prescriptive:

```markdown
## Commit message format
**Example 1:**
Input: Added user authentication with JWT tokens
Output: feat(auth): implement JWT-based authentication

**Bad:** "Should pass now" (without running anything)
**Good:** ran `make audit` -> "All checks pass"
```

### Workflow Patterns

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
   **New feature?** -> Run full audit
   **Docs only?** -> Skip checks
```

## Process

### 1. Capture Intent

Start by understanding what the user wants. Two paths:

**From conversation**: The user says "turn this into a skill."
Extract answers from the conversation history — the tools used,
the sequence of steps, corrections made, input/output formats
observed. Confirm your understanding before proceeding.

**From scratch**: Ask these questions:
1. What should this skill enable the agent to do?
2. When should it trigger? (what user phrases or contexts)
3. What is the expected output format?
4. Are there edge cases or failure modes to handle?

### 2. Interview and Research

Proactively ask about edge cases, dependencies, and success
criteria *before* writing. Don't wait for the user to think of
everything — come prepared with questions based on what you
know about the domain.

Check for existing skills in `.claude/skills/` to avoid
duplicates or find one to extend.

### 3. Draft the SKILL.md

Write the skill following the anatomy above. As you write:

- Tag paragraphs (Expert/Activation/Redundant) and remove
  anything Redundant
- Write the description field to be "pushy" — cover synonyms,
  concrete situations, edge-case triggers
- Explain reasoning behind instructions instead of relying on
  rigid MUST/NEVER directives
- Add good/bad example pairs for output quality

### 4. Test

Propose 2-3 realistic test prompts — the kind of thing a real
user would actually say. Share them: "Here are a few test cases
I'd like to try. Do these look right, or do you want to add
more?"

Run each test by spawning a subagent with the skill loaded.
For new skills, also run a baseline (same prompt, no skill) so
the user can compare. For skill improvements, baseline against
the previous version.

### 5. Review with User

Present the results. Let the user evaluate quality. Their
feedback drives the next iteration. Empty feedback means "looks
good."

### 6. Improve

Apply feedback, but think carefully about *how*:

- **Generalize from feedback.** The skill will be used across
  many prompts. Don't overfit to test examples. If a fix only
  helps one test case, find the underlying principle and encode
  that instead.

- **Keep the prompt lean.** Read the test transcripts, not just
  final outputs. If the skill is making the agent waste time on
  unproductive steps, remove those instructions.

- **Bundle repeated work.** If all test runs independently
  wrote similar helper scripts or took the same multi-step
  approach, bundle that script into `scripts/` and tell the
  skill to use it. This saves every future invocation from
  reinventing the wheel.

- **Explain the why.** When adding a new instruction, include
  the reasoning. "Sort results by date because users typically
  want the most recent first" beats "ALWAYS sort by date."

### 7. Iterate

Repeat steps 4-6 until:
- The user is satisfied
- Feedback is all empty (everything looks good)
- No meaningful progress is being made

### 8. Optimize Description

After the skill content is stable, revisit the description
field. A great skill with a vague description will never
trigger. Consider:

- Does it cover common phrasings the user might say?
- Does it mention related concepts and synonyms?
- Is it specific enough to avoid false triggers?
- Would the agent understand *when* to reach for this skill
  just by reading the description?

### 9. Deploy Template (if applicable)

If the skill should ship with `ctx init`, also create the
version in `internal/assets/claude/skills/` so it deploys to
new projects.

## Litmus Test

Before finalizing a skill:

1. Does the platform already do this? Do not restate it.
2. Does it suppress AI judgment? That is a jailbreak, not a
   skill.
3. Do all referenced files exist? Fix or remove phantom refs.
4. Is it under 500 lines? If over, split into references.
5. Is the description specific and pushy? Narrow the trigger
   and cover synonyms.
6. Would an expert say "this captures knowledge that took
   years to learn"? If not, it is probably Redundant.
7. Does it have examples? Without them, execution quality
   degrades over time.

## What NOT to Include

- Guidance the system prompt already provides
- Personality roleplay ("You are an expert...")
- Capability lists ("Masters X, Y, Z...")
- README, CHANGELOG, or auxiliary documentation
- References to files that do not exist
- Great content with a vague description — it will never
  trigger
