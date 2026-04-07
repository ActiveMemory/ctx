---
title: "Building Project Skills"
icon: lucide/hammer
---

![ctx](../images/ctx-banner.png)

## The Problem

You have workflows your agent needs to repeat across sessions: a deploy
checklist, a review protocol, a release process. Each time, you re-explain
the steps. The agent gets it mostly right but forgets edge cases you
corrected last time.

Skills solve this by encoding domain knowledge into a reusable document
the agent loads automatically when triggered. A skill is not code - it is
a structured prompt that captures what took you sessions to learn.

## TL;DR

```text
/ctx-skill-create
```

The skill-creator walks you through: **identify** a repeating workflow,
**draft** a skill, **test** with realistic prompts, **iterate** until
it triggers correctly and produces good output.

## Commands and Skills Used

| Tool                 | Type    | Purpose                                                     |
|----------------------|---------|-------------------------------------------------------------|
| `/ctx-skill-create` | Skill   | Interactive skill creation and improvement workflow         |
| `ctx init`           | Command | Deploys template skills to `.claude/skills/` on first setup |

## The Workflow

### Step 1: Identify a Repeating Pattern

Good skill candidates:

- **Checklists** you repeat: deploy steps, release prep, code review
- **Decisions the agent gets wrong**: if you keep correcting the same
  behavior, encode the correction
- **Multi-step workflows**: anything with a sequence of commands and
  conditional branches
- **Domain knowledge**: project-specific terminology, architecture
  constraints, or conventions the agent cannot infer from code alone

**Not** good candidates: one-off instructions, things the platform
already handles (file editing, git operations), or tasks too narrow
to reuse.

### Step 2: Create the Skill

Invoke the skill-creator:

```text
You: "I want a skill for our deploy process"

Agent: [Asks about the workflow: what steps, what tools,
        what edge cases, what the output should look like]
```

Or capture a workflow you just did:

```text
You: "Turn what we just did into a skill"

Agent: [Extracts the steps from conversation history,
        confirms understanding, drafts the skill]
```

The skill-creator produces a `SKILL.md` file in `.claude/skills/your-skill/`.

### Step 3: Test with Realistic Prompts

The skill-creator proposes 2-3 test prompts - the kind of thing a real
user would say. It runs each one and shows the result alongside a
baseline (same prompt without the skill) so you can compare.

```text
Agent: "Here are test prompts I'd try:
        1. 'Deploy to staging'
        2. 'Ship the hotfix'
        3. 'Run the release checklist'
        Want to adjust these?"
```

### Step 4: Iterate on the Description

The `description` field in frontmatter determines when a skill triggers.
Claude tends to undertrigger - descriptions need to be specific and
slightly "pushy":

```yaml
# Weak - too vague, will undertrigger
description: "Use for deployments"

# Strong - covers situations and synonyms
description: >-
  Use when deploying to staging or production, running the release
  checklist, or when the user says 'ship it', 'deploy this', or
  'push to prod'. Also use after merging to main when a deploy
  is expected.
```

The skill-creator helps you tune this iteratively.

### Step 5: Deploy as Template (Optional)

If the skill should be available to all projects (not just this one),
place it in `internal/assets/claude/skills/` so `ctx init` deploys it
to new projects automatically.

Most project-specific skills stay in `.claude/skills/` and travel with
the repo.

## Skill Anatomy

```
my-skill/
  SKILL.md         # Required: frontmatter + instructions (<500 lines)
  scripts/         # Optional: deterministic code the skill can execute
  references/      # Optional: detail loaded on demand (not always)
  assets/          # Optional: output templates, not loaded into context
```

Key sections in `SKILL.md`:

| Section           | Purpose                            | Required?          |
|-------------------|------------------------------------|--------------------|
| Frontmatter       | Name, description (trigger)        | Yes                |
| When to Use       | Positive triggers                  | Yes                |
| When NOT to Use   | Prevents false activations         | Yes                |
| Process           | Steps and commands                 | Yes                |
| Examples          | Good/bad output pairs              | Recommended        |
| Quality Checklist | Verify before reporting completion | For complex skills |

## Tips

* **Description is everything.** A great skill with a vague description
  never fires. Spend time on trigger coverage - synonyms, concrete
  situations, edge cases.
* **Stay under 500 lines.** If your skill is growing past this, move
  detail into `references/` files and point to them from `SKILL.md`.
* **Do not duplicate the platform.** If the agent already knows how to
  do something (edit files, run git commands), do not restate it. Tag
  paragraphs as Expert/Activation/Redundant and delete Redundant ones.
* **Explain why, not just what.** "Sort by date because users want
  recent results first" beats "ALWAYS sort by date." The agent
  generalizes from reasoning better than from rigid rules.
* **Test negative triggers.** Make sure the skill does *not* fire on
  unrelated prompts. A skill that activates too broadly becomes noise.

## Next Up

**[Parallel Agent Development with Git Worktrees ->](parallel-worktrees.md)**:
Split work across multiple agents using git worktrees.

## See Also

* [Skills Reference](../reference/skills.md): full listing of all bundled and
  project-local skills
* [Guide Your Agent](guide-your-agent.md): how commands, skills, and
  conversational patterns work together
* [Design Before Coding](design-before-coding.md): the four-skill chain for
  front-loading design work
