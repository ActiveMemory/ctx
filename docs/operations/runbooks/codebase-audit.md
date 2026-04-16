---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Codebase Audit
icon: lucide/search-code
---

![ctx](../../images/ctx-banner.png)

# Codebase Audit

A structured audit of the codebase: dead code, magic strings,
documentation drift, security surface, and roadmap opportunities.

**When to run**: Before a release, after a long YOLO sprint,
quarterly, or when planning the next phase of work.

**Time**: ~15-30 minutes with a team of agents.

---

## How to Use This Runbook

Start a Claude Code session with a clean git state
(`git stash` or commit first). Paste or adapt the prompt below.
The agent does the analysis; you triage the findings.

---

## Prompt

```
I want you to create an agent team to audit this codebase. Save each report as
a separate markdown file under `./ideas/` (or another directory if you prefer).

Use read-only agents (subagent_type: Explore) for all analyses. No code changes.

For each report, use this structure:
- Executive Summary (2-3 sentences + severity table)
- Findings (grouped, with file:line references)
- Ranked Recommendations (high/medium/low priority)
- Methodology (what was examined, how)

Keep reports actionable: every finding should suggest a concrete fix or next step.

## Analyses to Run

### 1. Extractable Patterns (session mining)
Search session JSONL files, journal entries, and task archives for repetitive
multi-step workflows. Count frequency of bash command sequences, slash command
usage, and recurring user prompts. Identify patterns that could become skills
or scripts. Cross-reference with existing skills to find coverage gaps.
Output: ranked list of automation opportunities with frequency data.

### 2. Documentation Drift (godoc + inline)
Compare every doc.go against its package's actual exports and behavior. Check
inline godoc comments on exported functions against their implementations.
Scan for stale TODO/FIXME/HACK comments. Check package-level comments match
package names. Output: drift items ranked by severity with exact file:line refs.

### 3. Maintainability
Look for: functions >80 lines that have logical split points; switch blocks
with >5 cases that could be table-driven or extracted; inline comments that
say "step 1", "step 2" or similar (sign the block wants to be a function);
files with >400 lines; packages with flat structure that could benefit from
sub-packages; functions that seem misplaced in their file. Do NOT flag
things that are fine as-is just because they could theoretically be different.
Output: concrete refactoring suggestions, not style nitpicks.

### 4. Security Review
This is a CLI app: focus on CLI-relevant attack surface, not web OWASP:
file path traversal (does user input flow into file paths unsanitized?),
command injection (does user input flow into exec calls?), symlink following
(does the tool follow symlinks when writing to .context/?), permission
handling (are file permissions set correctly?), sensitive data in outputs
(do any commands leak secrets or session content?). Output: findings with
severity ratings and exploit scenarios.

### 5. Blog Theme Discovery
Read existing blog posts for style and narrative voice. Analyze git log,
recent session discussions, and DECISIONS.md for story arcs worth writing
about. Suggest 3-5 blog post themes with: title, angle, target audience,
key commits/sessions to reference, and a 2-sentence pitch. Prioritize
themes that build a coherent narrative across posts.

### 6. Roadmap & Value Opportunities
Based on current features, recent momentum, and gaps found in other analyses:
what are the highest-value improvements? Consider: user-facing features,
developer experience, integration opportunities, and low-hanging fruit.
Output: prioritized list with effort/impact estimates (not time estimates).

### 7. User-Facing Documentation
Evaluate README, help text, and any user docs. Suggest improvements
structured as use-case pages: the problem, how ctx solves it, typical
workflow, gotchas. Identify gaps where a user would get stuck without
reading source code. Output: list of documentation gaps and suggested
page outlines.

### 8. Agent Team Strategies
Based on the codebase structure, suggest 2-3 agent team configurations for
upcoming work sessions. For each: team composition (roles, agent types),
task distribution strategy, coordination approach, and which types of work
it suits. Ground suggestions in actual project patterns, not generic advice.
```

---

## Tips

- **Clean state matters**: the prompt says "no code changes" but accidents
  happen. Start from a clean git state so you can `git checkout .` if needed.

- **Adjust scope**: drop analyses you don't need. Analyses 1-4 are the most
  actionable. Analyses 5-8 are planning/creative and can be skipped if you
  just want a technical audit.

- **Reports feed TASKS.md**: after the audit, read each report and create
  tasks in the appropriate Phase section. The reports are input, not output.

- **ideas/ is gitignored**: reports saved there won't be committed. Move
  specific findings to TASKS.md, DECISIONS.md, or LEARNINGS.md to persist them.

## History

- 2026-02-08: Original prompt created after a codebase audit sprint.
- 2026-02-17: Improved with read-only agents, report structure template,
  CLI-scoped security review, and maintainability thresholds.
- 2026-04-16: Moved from `hack/runbooks/` to `docs/operations/runbooks/`.
