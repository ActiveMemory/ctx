---
name: ctx-alignment-audit
description: "Audit alignment between docs and agent instructions. Use when recipes or docs make claims about agent behavior that may not be backed by the playbook or skills."
allowed-tools: Bash(ctx:*), Read, Grep, Glob
---

Audit whether behavioral claims in documentation are backed by
actual agent instructions in the playbook and skills. Identify
gaps between what the docs promise and what the agent is taught
to do.

## When to Use

- After writing or updating recipe documentation
- After modifying the Agent Playbook or skills
- When a recipe makes claims about proactive agent behavior
- Periodically to catch drift between docs and instructions
- When the user asks "is this real or wishful thinking?"

## When NOT to Use

- For code-level drift (use `/ctx-drift` instead)
- For context file staleness (use `/ctx-status`)
- When reviewing docs for prose quality (not behavioral claims)

## What to Audit

Behavioral claims are statements in docs that describe what an
agent **will do**, **may do**, or **offers to do**. They appear in:

- **"Agent note" blockquotes** in recipes
- **"Conversational Approach"** sections
- **"Tips"** sections describing proactive behavior
- Any paragraph using: "the agent will", "the agent offers",
  "proactively", "without being asked", "automatically"

## Process

### Step 1: Collect Claims

Read the target documentation file(s). Extract every behavioral
claim — each statement that describes agent behavior the user
should expect. Record:

- File and line number
- The exact claim
- What behavior it describes

### Step 2: Trace Each Claim

For each claim, search for matching instructions in:

1. **AGENT_PLAYBOOK.md** (in the context directory) — the primary
   behavioral source; the agent reads this at session start
2. **`.claude/skills/*/SKILL.md`** — skill-specific instructions
   loaded when a skill is invoked
3. **`CLAUDE.md`** — project-level instructions always in context

For each claim, determine:

- **Covered**: a matching instruction exists that would produce
  the described behavior
- **Partial**: related instructions exist but don't fully cover
  the claim (e.g., the playbook says "persist at milestones" but
  the recipe claims the agent identifies specific follow-up tasks)
- **Gap**: no instruction exists that would produce the behavior

### Step 3: Report

Present findings as a table:

| Claim (file:line) | Status | Backing instruction | Gap description |
|--------------------|--------|---------------------|-----------------|
| "agent creates follow-up tasks" (task-mgmt:57) | Gap | None | Playbook doesn't teach follow-up task identification |
| "agent offers to save learnings" (knowledge:237) | Covered | Playbook: Proactive Behavior During Work | — |

### Step 4: Fix (if requested)

For each gap, propose a fix:

- **Playbook addition**: if the behavior should apply broadly
  across all sessions (add to AGENT_PLAYBOOK.md)
- **Skill addition**: if the behavior is specific to one skill
  (add to the relevant SKILL.md)
- **Doc correction**: if the claim overpromises and should be
  toned down

Always update both the live copy (AGENT_PLAYBOOK.md in the
context directory, or `.claude/skills/`) AND the template copy
(`internal/assets/AGENT_PLAYBOOK.md` or
`internal/assets/claude/skills/`) to keep them in sync.

## What Makes a Claim "Covered"

A claim is covered when the agent instructions would **reliably
produce** the described behavior. Check:

- Is the instruction **specific enough**? "Persist at milestones"
  is vague. "After completing a task, offer to add follow-up
  tasks" is specific.
- Is the **trigger** clear? The instruction must say WHEN to act,
  not just WHAT to do.
- Is the **mechanism** clear? If the agent needs to run a command
  or check a file, is that stated?
- Does the instruction include **example phrasing**? Agents
  follow examples more reliably than abstract rules.

## Common Gap Patterns

- **"The agent proactively does X"** but the playbook only says
  to do X when asked → add proactive trigger
- **"Natural language triggers Y"** but no conversational trigger
  mapping exists → add to triggers table
- **"The agent chains A then B then C"** but each step is taught
  independently → add chained flow instruction
- **Recipe shows an example dialogue** but the underlying skill
  doesn't teach that flow → update the skill
- **"The agent notices X during work"** but no detection mechanism
  is described → add HOW to detect, not just WHAT to notice

## Instruction File Health

After auditing alignment, check that instruction files are still
a healthy size for agent consumption. Bloated instruction files
get ignored — agents stop following rules buried deep in long
documents.

### Size Heuristics

| File | Healthy | Warning | Danger |
|------|---------|---------|--------|
| AGENT_PLAYBOOK.md | < 5k tokens | 5-8k tokens | > 8k tokens |
| Individual SKILL.md | < 2k tokens | 2-3k tokens | > 3k tokens |
| CLAUDE.md | < 2k tokens | 2-3k tokens | > 3k tokens |
| Total context dir  | < 25k tokens | 25-40k tokens | > 40k tokens |

### How to Measure

```bash
# Total context token estimate
ctx status | grep "Token Estimate"

# Rough per-file estimate (~4 chars per token)
wc -c AGENT_PLAYBOOK.md  # divide by 4 (run from context dir)
wc -c .claude/skills/*/SKILL.md   # divide by 4
```

### When Files Are Too Large

- **Playbook too large**: Extract project-specific sections
  (Pre-Flight Checklist, Go Documentation Standard) into
  CONVENTIONS.md where they belong. The playbook should teach
  behavior, not coding standards.
- **Skill too large**: Split into a focused skill and a
  reference doc. The skill teaches WHEN and HOW; reference
  docs provide details the agent can look up if needed.
- **Total context too large**: Run `ctx compact` to archive
  completed tasks and deduplicate learnings.

### Signs of Instruction Bloat

- Agent ignores rules from the bottom half of a file
- Same instruction appears in multiple places (playbook AND
  skill AND CLAUDE.md)
- Instruction sections contain long examples that could be
  shorter
- File includes general knowledge the agent already has

## Quality Checklist

After completing the audit:
- [ ] Every behavioral claim in the target file was traced
- [ ] Each claim has a clear status (Covered / Partial / Gap)
- [ ] Gaps have proposed fixes with specific locations
- [ ] Fixes were applied to BOTH live and template copies
- [ ] No new claims were introduced without backing instructions
