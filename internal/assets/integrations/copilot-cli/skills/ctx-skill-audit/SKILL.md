---
name: ctx-skill-audit
description: "Audit skills against prompting best practices. Check for quality, consistency, and common anti-patterns."
tools: [bash, read, glob, grep]
---

Audit skill files for quality, consistency, and adherence to
prompting best practices.

## When to Use

- After creating or modifying skills
- During periodic quality reviews
- When skills seem to underperform

## When NOT to Use

- No skills exist yet
- Just after a fresh skill creation (let it settle first)

## Audit Dimensions

### 1. Positive framing
Instructions should say what TO do, not just what NOT to do.

### 2. Motivation over mandates
Explain WHY a rule exists, not just the rule.

### 3. Structure
Uses clear sections: When to Use, When NOT to Use, Process,
Quality Checklist.

### 4. Examples
Includes good and bad examples for clarity.

### 5. Scope
Skill is focused on one task, not a catch-all.

### 6. Description quality
Frontmatter description is clear and actionable.

### 7. Overtriggering guard
"When NOT to Use" section prevents false activations.

### 8. Phantom references
No references to tools, files, or commands that don't exist.

### 9. Tool declarations
Tools listed in frontmatter match what the skill actually uses.

## Process

1. Glob all skill files: `internal/assets/**/skills/*/SKILL.md`
2. Read each skill
3. Score against the 9 dimensions (pass/fail/partial)
4. Report findings per skill with actionable fixes

## Output Format

```
## Skill Audit Report

| Skill | Score | Issues |
|-------|-------|--------|
| ctx-implement | 8/9 | Missing bad example |
| ctx-commit | 9/9 | Clean |
| ctx-reflect | 7/9 | Phantom ref to /ctx-update-docs |

### Details
...
```

## Quality Checklist

- [ ] All skill files scanned
- [ ] Each dimension checked per skill
- [ ] Actionable fixes provided for failures
- [ ] No false positives (verify references exist)
