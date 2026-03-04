---
name: ctx-skill-audit
description: "Audit skills against Anthropic prompting best practices. Use when reviewing skill quality, after creating or modifying a skill, before releasing skills, or when a skill produces inconsistent results. Also use when the user says 'audit this skill', 'check skill quality', 'review the skills', or 'are our skills any good?'"
---

Audit one or more skills against Anthropic's prompting best
practices. The goal is to find patterns that degrade skill
effectiveness with current Claude models and suggest concrete
improvements.

## When to Use

- After creating or modifying a skill (quality gate)
- Reviewing all skills before a release (batch audit)
- When a skill produces inconsistent or poor results
- When skills were written for older models and may need
  calibration for Claude 4.5/4.6

## Before Auditing

1. Read `references/anthropic-best-practices.md` from this
   skill's directory — it contains the condensed audit criteria.
2. Identify which skill(s) to audit. If the user names a
   specific skill, audit that one. If they say "audit all
   skills," plan a batch pass.
3. For bundled skills, read from
   `internal/assets/claude/skills/*/SKILL.md`.
   For live skills, read from `.claude/skills/*/SKILL.md`.

## Audit Dimensions

Apply these checks to each skill. Each dimension maps to a
section in the best practices reference.

### 1. Positive Framing

Scan for negative instructions ("don't", "never", "avoid",
"do not") that lack a positive counterpart. Every negative
should be paired with what the agent *should* do instead.

**Pass:** negative instructions are supplements to clear
positive guidance.
**Fail:** primary instructions are negative, leaving the
agent to guess the desired behavior.

<example>
<fail>
Do not create new files. Do not modify tests. Do not add
comments.
</fail>
<pass>
Edit only the files specified in the task. Preserve existing
tests and comments — add new ones only when the user requests
them.
</pass>
</example>

### 2. Motivation Over Mandates

Check for MUST, NEVER, ALWAYS, CRITICAL used as emphasis
without explaining *why* the rule matters. Claude 4.5/4.6
responds better to reasoning than rigid directives.

**Pass:** important instructions include motivation ("because
X" or "so that Y") that lets the model generalize.
**Fail:** instructions rely on emphasis alone to convey
importance.

<example>
<fail>
You MUST ALWAYS run tests before reporting completion.
</fail>
<pass>
Run tests before reporting completion — untested changes
create silent regressions that compound across sessions.
</pass>
</example>

### 3. XML Tag Structure

Check whether the skill mixes instructions with variable
content (file paths, user input, injected code) without
clear delimiters. XML tags prevent the model from confusing
injected content with skill instructions.

**Pass:** variable content is wrapped in descriptive tags,
or the skill doesn't inject variable content.
**Fail:** the skill templates in external content alongside
instructions without delimiters.

### 4. Few-Shot Examples

Check whether non-trivial behaviors (output formats, decision
logic, style requirements) are demonstrated with examples.
Skills that describe complex output without showing it drift
over time.

**Pass:** key behaviors have at least one good/bad example
pair, or the behavior is simple enough that examples would
be redundant.
**Fail:** the skill describes a specific output format or
decision process but provides no examples.

### 5. Subagent Guard

If the skill spawns or encourages spawning subagents (via the
Agent tool), check that it states when subagents are and
aren't warranted. Claude Opus 4.6 over-delegates to subagents
when a direct tool call would be faster.

**Pass:** subagent usage has explicit scope (when to use,
when not to), or the skill doesn't involve subagents.
**Fail:** the skill defaults to subagent delegation without
stating when direct execution is preferable.

### 6. Overtriggering Calibration

Check for language written to combat undertriggering in older
models that may cause overtriggering in Claude 4.5/4.6:
excessive caps emphasis (CRITICAL, MUST), redundant capability
statements ("You are an expert"), or aggressive always/never
framing.

**Pass:** instructions use natural language with emphasis
reserved for genuinely critical points.
**Fail:** the skill reads like it was written for a less
capable model that needed constant nudging.

### 7. Phantom References

Every file path, tool name, and command referenced in the
skill must exist. Broken references are a form of hallucination
in the skill itself.

**Pass:** all references resolve to real files/tools.
**Fail:** the skill mentions files or commands that don't
exist.

### 8. Scope Discipline

Check whether the skill encourages work beyond what's
requested — "while you're in there" improvements, unsolicited
refactoring, or scope creep. Skills should state the minimum
viable outcome.

**Pass:** the skill's scope matches its stated purpose.
**Fail:** the skill encourages additional work beyond its
core task.

### 9. Description Trigger Quality

The `description` field determines when the skill activates.
Check that it:
- Covers concrete trigger situations and user phrases
- Includes synonyms and related concepts
- Is specific enough to avoid false triggers
- Is "pushy" enough to avoid undertriggering

**Pass:** reading the description alone, you'd know exactly
when to use this skill.
**Fail:** the description is vague ("use for general tasks")
or too narrow (misses common phrasings).

## Process

### Single Skill Audit

1. Read the skill's SKILL.md.
2. Apply all 9 audit dimensions.
3. Report findings using the output format below.
4. Suggest specific rewrites for any failures — show the
   current text and the proposed replacement.

### Batch Audit

1. List all skills to audit (bundled, live, or both).
2. Audit each skill, but report concisely — only dimensions
   that fail or have notable findings.
3. Summarize with a scorecard at the end.

## Output Format

For each audited skill, report:

```
### /skill-name

**Overall:** X/9 pass

| # | Dimension              | Result | Notes                    |
|---|------------------------|--------|--------------------------|
| 1 | Positive framing       | pass   |                          |
| 2 | Motivation over mandates | fail | 3 bare MUST/NEVER found  |
| 3 | XML tag structure      | pass   |                          |
| 4 | Few-shot examples      | fail   | No output format example |
| 5 | Subagent guard         | n/a    | No subagent usage        |
| 6 | Overtriggering         | pass   |                          |
| 7 | Phantom references     | pass   |                          |
| 8 | Scope discipline       | pass   |                          |
| 9 | Description quality    | warn   | Missing synonym coverage |

**Suggested fixes:**
- [Dimension 2] Line "You MUST ALWAYS run tests" →
  "Run tests before completion — untested changes create
  silent regressions."
- [Dimension 4] Add example showing expected output format
  after the "Report results" section.
```

For batch audits, end with a summary:

```
## Batch Summary

| Skill              | Score | Top Issue                |
|--------------------|-------|--------------------------|
| ctx-commit         | 8/9   | Missing example          |
| ctx-drift          | 7/9   | 2 bare mandates          |
| ctx-verify         | 9/9   | —                        |
```

## Quality Checklist

Before reporting audit results:

- [ ] Read the best practices reference before starting
- [ ] Applied all 9 dimensions (mark n/a where inapplicable)
- [ ] Every "fail" has a specific suggested rewrite, not just
      a description of the problem
- [ ] Phantom reference check actually verified file existence
      (used Glob/Read, not assumption)
- [ ] Description quality check considered real user phrases,
      not hypothetical ones
