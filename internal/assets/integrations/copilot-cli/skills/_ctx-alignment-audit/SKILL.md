---
name: _ctx-alignment-audit
description: "Audit alignment between docs and agent instructions. Use when docs make claims about agent behavior that may not be backed by the playbook or skills."
tools: [bash, read, glob, grep]
---

Audit whether behavioral claims in documentation are backed by
actual agent instructions.

## When to Use

- After writing or updating documentation
- After modifying the Agent Playbook or skills
- When a doc makes claims about proactive agent behavior
- Periodically to catch drift between docs and instructions

## When NOT to Use

- For code-level drift (use `ctx-drift` instead)
- For context file staleness (use `ctx-status`)
- When reviewing docs for prose quality (not behavioral claims)

## Process

### Step 1: Collect Claims

Read target docs. Extract every behavioral claim: statements
describing what an agent "will do", "may do", or "offers to do".

### Step 2: Trace Each Claim

Search for matching instructions in:
1. **AGENT_PLAYBOOK.md**: primary behavioral source
2. **skills/*/SKILL.md**: skill-specific instructions
3. **INSTRUCTIONS.md**: project-level instructions

For each claim, determine:
- **Covered**: matching instruction exists
- **Partial**: related but incomplete
- **Gap**: no instruction exists

### Step 3: Report

| Claim (file:line) | Status | Backing instruction | Gap |
|---|---|---|---|
| "agent creates tasks" | Gap | None | Not taught |
| "agent saves learnings" | Covered | Playbook | n/a |

### Step 4: Fix (if requested)

For each gap, propose:
- **Playbook addition**: if behavior applies broadly
- **Skill addition**: if specific to one skill
- **Doc correction**: if the claim overpromises

## Quality Checklist

- [ ] Every behavioral claim was traced
- [ ] Each claim has clear status (Covered/Partial/Gap)
- [ ] Gaps have proposed fixes
- [ ] No new claims introduced without backing
